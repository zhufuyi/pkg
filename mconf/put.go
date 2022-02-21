package mconf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tomwright/dasel"
)

// 通用修改设置
type genericPutOptions struct {
	File                string
	Out                 string
	Parser              string
	ReadParser          string
	WriteParser         string
	Selector            string
	Value               string
	ValueType           string
	Init                func(genericPutOptions) genericPutOptions
	Reader              io.Reader
	Writer              io.Writer
	Multi               bool
	Compact             bool
	MergeInputDocuments bool
	EscapeHTML          bool
}

func shouldReadFromStdin(fileFlag string) bool {
	return fileFlag == "" || fileFlag == "stdin" || fileFlag == "-"
}

func getReadParser(fileFlag string, readParserFlag string, parserFlag string) (ReadParser, error) {
	useStdin := shouldReadFromStdin(fileFlag)

	if readParserFlag == "" {
		readParserFlag = parserFlag
	}

	if useStdin && readParserFlag == "" {
		return nil, fmt.Errorf("read parser flag required when reading from stdin")
	}

	if readParserFlag == "" {
		parser, err := NewReadParserFromFilename(fileFlag)
		if err != nil {
			return nil, fmt.Errorf("could not get read parser from filename: %w", err)
		}
		return parser, nil
	}
	parser, err := NewReadParserFromString(readParserFlag)
	if err != nil {
		return nil, fmt.Errorf("could not get read parser: %w", err)
	}
	return parser, nil
}

type getRootNodeOpts struct {
	File                string
	Reader              io.Reader
	Parser              ReadParser
	MergeInputDocuments bool
}

func getRootNode(opts getRootNodeOpts) (*dasel.Node, error) {
	if opts.Reader == nil {
		f, err := os.Open(opts.File)
		if err != nil {
			return nil, fmt.Errorf("could not open input file: %w", err)
		}
		defer f.Close()
		opts.Reader = f
	}

	value, err := Load(opts.Parser, opts.Reader)
	if err != nil {
		return nil, fmt.Errorf("could not load input: %w", err)
	}

	if opts.MergeInputDocuments {
		switch val := value.(type) {
		case SingleDocument:
			value = &BasicSingleDocument{Value: []interface{}{val.Document()}}
		case MultiDocument:
			value = &BasicSingleDocument{Value: val.Documents()}
		}
	}

	return dasel.New(value), nil
}

func parseValue(value string, valueType string) (interface{}, error) {
	switch strings.ToLower(valueType) {
	case "string", "str":
		return value, nil
	case "int", "integer":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse int [%s]: %w", value, err)
		}
		return val, nil
	case "bool", "boolean":
		switch strings.ToLower(value) {
		case "true", "t", "yes", "y", "1":
			return true, nil
		case "false", "f", "no", "n", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("could not parse bool [%s]: unhandled value", value)
		}
	default:
		return nil, fmt.Errorf("unhandled type: %s", valueType)
	}
}

func shouldWriteToStdout(fileFlag string, outFlag string) bool {
	return (outFlag == "stdout" || outFlag == "-") || outFlag == "" && shouldReadFromStdin(fileFlag)
}

func getWriteParser(readParser ReadParser, writeParserFlag string, parserFlag string,
	outFlag string, fileFlag string, formatTemplateFlag string) (WriteParser, error) {
	if formatTemplateFlag != "" {
		writeParserFlag = "plain"
	}

	if writeParserFlag == "" {
		writeParserFlag = parserFlag
	}

	if writeParserFlag != "" {
		parser, err := NewWriteParserFromString(writeParserFlag)
		if err != nil {
			return nil, fmt.Errorf("could not get write parser: %w", err)
		}
		return parser, nil
	}

	if !shouldWriteToStdout(fileFlag, outFlag) {
		p, err := NewWriteParserFromFilename(fileFlag)
		if err != nil {
			return nil, fmt.Errorf("could not get write parser from filename: %w", err)
		}
		return p, nil
	}

	if p, ok := readParser.(WriteParser); ok {
		return p, nil
	}
	return nil, fmt.Errorf("read parser cannot be used to write. please specify a write parser")
}

type writeNodeToOutputOpts struct {
	Node           *dasel.Node
	Parser         WriteParser
	File           string
	Out            string
	Writer         io.Writer
	FormatTemplate string
}

func getOutputWriter(in io.Writer, file string, out string) (io.Writer, func(), error) {
	if in == nil {
		switch {
		case out == "":
			// No out flag... write to the file we read from.
			f, err := os.Create(file)
			if err != nil {
				return nil, nil, fmt.Errorf("could not open output file: %w", err)
			}
			return f, func() {
				_ = f.Close()
			}, nil

		case out != "":
			// Out flag was set.
			f, err := os.Create(out)
			if err != nil {
				return nil, nil, fmt.Errorf("could not open output file: %w", err)
			}
			return f, func() {
				_ = f.Close()
			}, nil
		}
	}
	return in, func() {}, nil
}

func writeNodeToOutput(opts writeNodeToOutputOpts, options ...ReadWriteOption) error {
	writer, writerCleanUp, err := getOutputWriter(opts.Writer, opts.File, opts.Out)
	if err != nil {
		return err
	}
	opts.Writer = writer
	defer writerCleanUp()

	var value, originalValue interface{}
	if opts.FormatTemplate == "" {
		value = opts.Node.InterfaceValue()
		originalValue = opts.Node.OriginalValue
	} else {
		result, err := dasel.FormatNode(opts.Node, opts.FormatTemplate)
		if err != nil {
			return fmt.Errorf("could not format node: %w", err)
		}
		value = result.String()
		originalValue = value
	}

	if err := Write(opts.Parser, value, originalValue, opts.Writer, options...); err != nil {
		return fmt.Errorf("could not write to output file: %w", err)
	}

	return nil
}

func runGenericPutCommand(opts genericPutOptions) error {
	if opts.Init != nil {
		opts = opts.Init(opts)
	}
	readParser, err := getReadParser(opts.File, opts.ReadParser, opts.Parser)
	if err != nil {
		return err
	}
	rootNode, err := getRootNode(getRootNodeOpts{
		File:                opts.File,
		Parser:              readParser,
		Reader:              opts.Reader,
		MergeInputDocuments: opts.MergeInputDocuments,
	})
	if err != nil {
		return err
	}

	updateValue, err := parseValue(opts.Value, opts.ValueType)
	if err != nil {
		return err
	}

	if opts.Multi {
		if err := rootNode.PutMultiple(opts.Selector, updateValue); err != nil {
			return fmt.Errorf("could not put multi value: %w", err)
		}
	} else {
		if err := rootNode.Put(opts.Selector, updateValue); err != nil {
			return fmt.Errorf("could not put value: %w", err)
		}
	}

	writeParser, err := getWriteParser(readParser, opts.WriteParser, opts.Parser, opts.Out, opts.File, "")
	if err != nil {
		return err
	}

	writeOptions := []ReadWriteOption{
		EscapeHTMLOption(opts.EscapeHTML),
	}

	if opts.Compact {
		writeOptions = append(writeOptions, PrettyPrintOption(false))
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   rootNode,
		Parser: writeParser,
		File:   opts.File,
		Out:    opts.Out,
		Writer: opts.Writer,
	}, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}

// Put 修改或增加一个字段值，inFormat输入数据解析格式yaml、json 、toml，outFormat输出数据解析格式yaml、json 、toml，输出格式可以和输入格式不一致
func Put(in []byte, selector string, valueType string, value string, inFormat string, outFormat string) ([]byte, error) {
	outputBuffer := bytes.NewBuffer([]byte{})

	opts := genericPutOptions{
		Parser:      inFormat,
		WriteParser: outFormat,
		Selector:    selector,
		Value:       value,
		ValueType:   valueType,
		Reader:      bytes.NewReader(in),
		Writer:      outputBuffer,
	}

	err := runGenericPutCommand(opts)
	if err != nil {
		return nil, err
	}

	output, err := io.ReadAll(outputBuffer)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// PutYaml 修改或增加字段，针对yaml文档
func PutYaml(in []byte, selector string, valueType string, value string) ([]byte, error) {
	return Put(in, selector, valueType, value, YamlFormat, YamlFormat)
}

// PutYaml 修改或增加字段，针对json文档
func PutJson(in []byte, selector string, valueType string, value string) ([]byte, error) {
	return Put(in, selector, valueType, value, JsonFormat, JsonFormat)
}

// ---------------------------------------------------------------------------------------

type putObjectOpts struct {
	File                string
	Out                 string
	ReadParser          string
	WriteParser         string
	Parser              string
	Selector            string
	InputTypes          []string
	InputValues         []string
	Reader              io.Reader
	Writer              io.Writer
	Multi               bool
	Compact             bool
	MergeInputDocuments bool
	EscapeHTML          bool
}

func getMapFromTypesValues(inputTypes []string, inputValues []string) (map[string]interface{}, error) {
	if len(inputTypes) != len(inputValues) {
		return nil, fmt.Errorf("exactly %d types are required, want %d", len(inputValues), len(inputTypes))
	}

	updateValue := map[string]interface{}{}

	for k, arg := range inputValues {
		splitArg := strings.Split(arg, "=")
		name := splitArg[0]
		value := strings.Join(splitArg[1:], "=")
		parsedValue, err := parseValue(value, inputTypes[k])
		if err != nil {
			return nil, fmt.Errorf("could not parse value [%s]: %w", name, err)
		}
		updateValue[name] = parsedValue
	}

	return updateValue, nil
}

func runPutObjectCommand(opts putObjectOpts) error {
	readParser, err := getReadParser(opts.File, opts.ReadParser, opts.Parser)
	if err != nil {
		return err
	}
	rootNode, err := getRootNode(getRootNodeOpts{
		File:                opts.File,
		Parser:              readParser,
		Reader:              opts.Reader,
		MergeInputDocuments: opts.MergeInputDocuments,
	})
	if err != nil {
		return err
	}

	updateValue, err := getMapFromTypesValues(opts.InputTypes, opts.InputValues)
	if err != nil {
		return err
	}

	if opts.Multi {
		if err := rootNode.PutMultiple(opts.Selector, updateValue); err != nil {
			return fmt.Errorf("could not put object multi value: %w", err)
		}
	} else {
		if err := rootNode.Put(opts.Selector, updateValue); err != nil {
			return fmt.Errorf("could not put object value: %w", err)
		}
	}

	writeParser, err := getWriteParser(readParser, opts.WriteParser, opts.Parser, opts.Out, opts.File, "")
	if err != nil {
		return err
	}

	writeOptions := []ReadWriteOption{
		EscapeHTMLOption(opts.EscapeHTML),
	}

	if opts.Compact {
		writeOptions = append(writeOptions, PrettyPrintOption(false))
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   rootNode,
		Parser: writeParser,
		File:   opts.File,
		Out:    opts.Out,
		Writer: opts.Writer,
	}, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}

// PutObjectYaml 修改或增加多个字段值，参数valueTypes是值的类型，kvs是值类型对应kv键值对，例如valueTypes=["int","string"],values=["a=1","b=hello"]
// 参数inFormat输是入数据解析格式yaml、json 、toml，outFormat是输出数据解析格式yaml、json 、toml，输出格式可以和输入格式不一致
func PutObject(in []byte, selector string, valueTypes []string, kvs []string, inFormat string, outFormat string) ([]byte, error) {
	outputBuffer := bytes.NewBuffer([]byte{})

	err := runPutObjectCommand(putObjectOpts{
		Parser:      inFormat,
		WriteParser: outFormat,
		Selector:    selector,
		InputValues: kvs,
		InputTypes:  valueTypes,
		Reader:      bytes.NewReader(in),
		Writer:      outputBuffer,
	})
	if err != nil {
		return nil, err
	}

	output, err := io.ReadAll(outputBuffer)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// PutObjectYaml 修改或增加多个字段值，针对yaml文档
func PutObjectYaml(in []byte, selector string, valueTypes []string, kvs []string) ([]byte, error) {
	return PutObject(in, selector, valueTypes, kvs, YamlFormat, YamlFormat)
}

// PutObjectJson 修改或增加多个字段值，针对json文档
func PutObjectJson(in []byte, selector string, valueTypes []string, kvs []string) ([]byte, error) {
	return PutObject(in, selector, valueTypes, kvs, JsonFormat, JsonFormat)
}

// ---------------------------------------------------------------------------------------

type putDocumentOpts struct {
	File                string
	Out                 string
	ReadParser          string
	WriteParser         string
	Parser              string
	Selector            string
	DocumentString      string
	DocumentParser      string
	Reader              io.Reader
	Writer              io.Writer
	Multi               bool
	Compact             bool
	MergeInputDocuments bool
	EscapeHTML          bool
}

func runPutDocumentCommand(opts putDocumentOpts) error {
	readParser, err := getReadParser(opts.File, opts.ReadParser, opts.Parser)
	if err != nil {
		return err
	}
	rootNode, err := getRootNode(getRootNodeOpts{
		File:                opts.File,
		Parser:              readParser,
		Reader:              opts.Reader,
		MergeInputDocuments: opts.MergeInputDocuments,
	})
	if err != nil {
		return err
	}

	documentParser, err := getPutDocumentParser(readParser, opts.DocumentParser)
	if err != nil {
		return err
	}

	documentValue, err := documentParser.FromBytes([]byte(opts.DocumentString))
	if err != nil {
		return fmt.Errorf("could not parse document: %w", err)
	}

	if opts.Multi {
		if err := rootNode.PutMultiple(opts.Selector, documentValue); err != nil {
			return fmt.Errorf("could not put document multi value: %w", err)
		}
	} else {
		if err := rootNode.Put(opts.Selector, documentValue); err != nil {
			return fmt.Errorf("could not put document value: %w", err)
		}
	}

	writeParser, err := getWriteParser(readParser, opts.WriteParser, opts.Parser, opts.Out, opts.File, "")
	if err != nil {
		return err
	}

	writeOptions := []ReadWriteOption{
		EscapeHTMLOption(opts.EscapeHTML),
	}

	if opts.Compact {
		writeOptions = append(writeOptions, PrettyPrintOption(false))
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   rootNode,
		Parser: writeParser,
		File:   opts.File,
		Out:    opts.Out,
		Writer: opts.Writer,
	}, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}

func getPutDocumentParser(readParser ReadParser, documentParserFlag string) (ReadParser, error) {
	if documentParserFlag == "" {
		return readParser, nil
	}

	parser, err := NewReadParserFromString(documentParserFlag)
	if err != nil {
		return nil, fmt.Errorf("could not get document parser: %w", err)
	}
	return parser, nil
}

// PutDocumentYaml 修改或增加多个字段值，参数docs是yaml或json文档，例如`["foo","bar"]
// 参数inFormat输是入数据解析格式yaml、json 、toml，outFormat是输出数据解析格式yaml、json 、toml，输出格式可以和输入格式不一致
func PutDocument(in []byte, selector string, docs string, inFormat string, outFormat string) ([]byte, error) {
	outputBuffer := bytes.NewBuffer([]byte{})

	err := runPutDocumentCommand(putDocumentOpts{
		Parser:         inFormat,
		WriteParser:    outFormat,
		Selector:       selector,
		DocumentString: docs,
		Reader:         bytes.NewReader(in),
		Writer:         outputBuffer,
	})
	if err != nil {
		return nil, err
	}

	output, err := io.ReadAll(outputBuffer)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// PutObjectYaml 修改或增加文档，针对yaml文档
func PutDocumentYaml(in []byte, selector string, docs string) ([]byte, error) {
	return PutDocument(in, selector, docs, YamlFormat, YamlFormat)
}

// PutObjectJson 修改或增加文档，针对json文档
func PutDocumentJson(in []byte, selector string, docs string) ([]byte, error) {
	return PutDocument(in, selector, docs, JsonFormat, JsonFormat)
}

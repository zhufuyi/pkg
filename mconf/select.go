package mconf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/tomwright/dasel"
)

type selectOptions struct {
	File                string
	Parser              string
	ReadParser          string
	WriteParser         string
	Selector            string
	Reader              io.Reader
	Writer              io.Writer
	Multi               bool
	NullValueNotFound   bool
	Compact             bool
	DisplayLength       bool
	MergeInputDocuments bool
	FormatTemplate      string
	Colourise           bool
	EscapeHTML          bool
}

type customErrorHandlingOpts struct {
	File     string
	Out      string
	Writer   io.Writer
	Err      error
	NullFlag bool
}

func writeStringToOutput(value string, file string, out string, writer io.Writer) error {
	writer, writerCleanUp, err := getOutputWriter(writer, file, out)
	if err != nil {
		return err
	}
	defer writerCleanUp()

	if _, err := writer.Write([]byte(value)); err != nil {
		return fmt.Errorf("could not write to output file: %w", err)
	}

	return nil
}

func customErrorHandling(opts customErrorHandlingOpts) (bool, error) {
	if opts.Err == nil {
		return false, nil
	}

	if !opts.NullFlag {
		return false, opts.Err
	}

	var valNotFound *dasel.ValueNotFound
	if !errors.As(opts.Err, &valNotFound) {
		return false, opts.Err
	}

	if err := writeStringToOutput("null\n", opts.File, opts.Out, opts.Writer); err != nil {
		return false, fmt.Errorf("could not write string to output: %w", err)
	}

	return true, nil
}

func outputNodeLength(writer io.Writer, nodes ...*dasel.Node) error {
	for _, n := range nodes {
		val := n.Value
		if val.Kind() == reflect.Interface {
			val = val.Elem()
		}
		length := 0
		switch val.Kind() {
		case reflect.Map, reflect.Slice, reflect.String:
			length = val.Len()
		default:
			length = len(fmt.Sprint(val.Interface()))
		}
		if _, err := fmt.Fprintf(writer, "%d\n", length); err != nil {
			return err
		}
	}
	return nil
}

type writeNodesToOutputOpts struct {
	Nodes          []*dasel.Node
	Parser         WriteParser
	File           string
	Out            string
	Writer         io.Writer
	FormatTemplate string
}

func writeNodesToOutput(opts writeNodesToOutputOpts, options ...ReadWriteOption) error {
	writer, writerCleanUp, err := getOutputWriter( /*cmd, */ opts.Writer, opts.File, opts.Out)
	if err != nil {
		return err
	}
	opts.Writer = writer
	defer writerCleanUp()

	buf := new(bytes.Buffer)

	for i, n := range opts.Nodes {
		subOpts := writeNodeToOutputOpts{
			Node:           n,
			Parser:         opts.Parser,
			Writer:         buf,
			FormatTemplate: opts.FormatTemplate,
		}
		if err := writeNodeToOutput(subOpts, options...); err != nil {
			return fmt.Errorf("could not write node %d to output: %w", i, err)
		}
	}

	if _, err := io.Copy(opts.Writer, buf); err != nil {
		return fmt.Errorf("could not copy buffer to real output: %w", err)
	}

	return nil
}

func runSelectMultiCommand(rootNode *dasel.Node, opts selectOptions, writeParser WriteParser, writeOptions []ReadWriteOption) error {
	var results []*dasel.Node
	var err error
	if opts.Selector == "." {
		results = []*dasel.Node{rootNode}
	} else {
		results, err = rootNode.QueryMultiple(opts.Selector)
	}

	written, err := customErrorHandling(customErrorHandlingOpts{
		File:     opts.File,
		Writer:   opts.Writer,
		Err:      err,
		NullFlag: opts.NullValueNotFound,
	})
	if err != nil {
		return fmt.Errorf("could not query multiple node: %w", err)
	}
	if written {
		return nil
	}

	if opts.DisplayLength {
		return outputNodeLength(opts.Writer, results...)
	}

	if err := writeNodesToOutput(writeNodesToOutputOpts{
		Nodes:          results,
		Parser:         writeParser,
		Writer:         opts.Writer,
		FormatTemplate: opts.FormatTemplate,
	}, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}
	return nil
}

func runSelectCommand(opts selectOptions) error {
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

	if !rootNode.Value.IsValid() {
		rootNode = dasel.New(&BasicSingleDocument{
			Value: map[string]interface{}{},
		})
	}

	writeParser, err := getWriteParser(readParser, opts.WriteParser, opts.Parser, "-", opts.File, opts.FormatTemplate)
	if err != nil {
		return err
	}

	writeOptions := []ReadWriteOption{
		EscapeHTMLOption(opts.EscapeHTML),
	}

	if opts.Compact {
		writeOptions = append(writeOptions, PrettyPrintOption(false))
	}
	if opts.Colourise {
		writeOptions = append(writeOptions, ColouriseOption(true))
	}
	if opts.Multi {
		return runSelectMultiCommand(rootNode, opts, writeParser, writeOptions)
	}

	var res *dasel.Node

	if opts.Selector == "." {
		res = rootNode
	} else {
		res, err = rootNode.Query(opts.Selector)
		if err != nil {
			err = fmt.Errorf("could not query node: %w", err)
		}
	}

	written, err := customErrorHandling(customErrorHandlingOpts{
		File:     opts.File,
		Writer:   opts.Writer,
		Err:      err,
		NullFlag: opts.NullValueNotFound,
	})
	if err != nil {
		return err
	}
	if written {
		return nil
	}

	if opts.DisplayLength {
		return outputNodeLength(opts.Writer, res)
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:           res,
		Parser:         writeParser,
		Writer:         opts.Writer,
		FormatTemplate: opts.FormatTemplate,
	}, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}

// Find 查询字段值，select表示选择字段，inFormat输入数据解析格式yaml、json 、toml，outFormat输出数据解析格式yaml、json 、toml，输出格式可以和输入格式不一致
func Find(in []byte, selector string, inFormat string, outFormat string) ([]byte, error) {
	outputBuffer := bytes.NewBuffer([]byte{})

	err := runSelectCommand(selectOptions{
		Parser:      inFormat,
		WriteParser: outFormat,
		Selector:    selector,
		Reader:      bytes.NewReader(in),
		Writer:      outputBuffer,
		Multi:       true,
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

// FindYaml 查询yaml文档字段值，select表示选择字段，输入输出都是yaml格式
func FindYaml(in []byte, selector string) ([]byte, error) {
	return Find(in, selector, YamlFormat, YamlFormat)
}

// FindJSON 查询json文档字段值，select表示选择字段，输入输出都是json格式
func FindJSON(in []byte, selector string) ([]byte, error) {
	return Find(in, selector, JSONFormat, JSONFormat)
}

// Count 计算查询字段值数量，inFormat输入数据解析格式yaml、json 、toml，outFormat输出数据解析格式yaml、json 、toml
func Count(in []byte, selector string, inFormat string) (int, error) {
	outputBuffer := bytes.NewBuffer([]byte{})

	err := runSelectCommand(selectOptions{
		Parser:        inFormat,
		Selector:      selector,
		Reader:        bytes.NewReader(in),
		Writer:        outputBuffer,
		Multi:         true,
		DisplayLength: true,
	})
	if err != nil {
		return 0, err
	}

	output, err := io.ReadAll(outputBuffer)
	if err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(strings.Replace(string(output), "\n", "", -1))
	if err != nil {
		return 0, err
	}

	return n, nil
}

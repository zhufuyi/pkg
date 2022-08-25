package mconf

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tomwright/dasel"
)

type deleteOptions struct {
	File                string
	Parser              string
	ReadParser          string
	WriteParser         string
	Selector            string
	Reader              io.Reader
	Writer              io.Writer
	Multi               bool
	Compact             bool
	MergeInputDocuments bool
	Out                 string
	EscapeHTML          bool
}

func runDeleteMultiCommand(rootNode *dasel.Node, opts deleteOptions, writeParser WriteParser, writeOptions []ReadWriteOption) error {
	err := rootNode.DeleteMultiple(opts.Selector)

	written, err := customErrorHandling(customErrorHandlingOpts{
		File:     opts.File,
		Writer:   opts.Writer,
		Err:      err,
		NullFlag: false,
		Out:      opts.Out,
	})
	if err != nil {
		return fmt.Errorf("could not delete multiple node: %w", err)
	}
	if written {
		return nil
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   rootNode,
		Parser: writeParser,
		Writer: opts.Writer,
		File:   opts.File,
		Out:    opts.Out,
	}, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}
	return nil
}

func runDeleteCommand(opts deleteOptions) error {
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

	if opts.Multi {
		return runDeleteMultiCommand(rootNode, opts, writeParser, writeOptions)
	}

	if deleteErr := rootNode.Delete(opts.Selector); deleteErr != nil {
		err = fmt.Errorf("could not delete node: %w", deleteErr)
	}
	written, err := customErrorHandling(customErrorHandlingOpts{
		File:     opts.File,
		Writer:   opts.Writer,
		Err:      err,
		NullFlag: false,
		Out:      opts.Out,
	})
	if err != nil {
		return err
	}
	if written {
		return nil
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   rootNode,
		Parser: writeParser,
		Writer: opts.Writer,
		File:   opts.File,
		Out:    opts.Out,
	}, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}

// Delete 删除字段，参数inFormat输是入数据解析格式yaml、json 、toml，outFormat是输出数据解析格式yaml、json 、toml，输出格式可以和输入格式不一致
func Delete(in []byte, selector string, inFormat string, outFormat string) ([]byte, error) {
	outputBuffer := bytes.NewBuffer([]byte{})

	opts := deleteOptions{
		Parser:      inFormat,
		WriteParser: outFormat,
		Selector:    selector,
		Reader:      bytes.NewReader(in),
		Writer:      outputBuffer,
		Multi:       true,
	}

	err := runDeleteCommand(opts)
	if err != nil {
		return nil, err
	}

	output, err := io.ReadAll(outputBuffer)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DeleteYaml 删除字段，针对yaml文档
func DeleteYaml(in []byte, selector string) ([]byte, error) {
	return Delete(in, selector, YamlFormat, YamlFormat)
}

// DeleteJSON 删除字段，针对json文档
func DeleteJSON(in []byte, selector string) ([]byte, error) {
	return Delete(in, selector, JSONFormat, JSONFormat)
}

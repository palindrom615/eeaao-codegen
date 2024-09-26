package plugin

import (
	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
	"io"
	"log"
	"os"
)

type ProtobufPlugin struct {
	handler *reporter.Handler
}

func NewProtobufPlugin() *ProtobufPlugin {
	errReporter := func(err reporter.ErrorWithPos) error {
		log.Println(err)
		return nil
	}
	warnReporter := func(err reporter.ErrorWithPos) {
		log.Println(err)
	}
	handler := reporter.NewHandler(reporter.NewReporter(errReporter, warnReporter))
	return &ProtobufPlugin{handler: handler}
}

func (p *ProtobufPlugin) LoadSpecFile(path string) (SpecData, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileNode, err := parser.Parse(path, reader, p.handler)
	if err != nil {
		return nil, err
	}
	res, err := parser.ResultFromAST(fileNode, false, p.handler)
	if err != nil {
		return nil, err
	}
	return res.FileDescriptorProto(), nil
}

func (p *ProtobufPlugin) LoadSpec(reader io.Reader) (SpecData, error) {
	fileNode, err := parser.Parse("", reader, p.handler)
	if err != nil {
		return nil, err
	}
	res, err := parser.ResultFromAST(fileNode, false, p.handler)
	if err != nil {
		return nil, err
	}
	return res.FileDescriptorProto(), nil
}

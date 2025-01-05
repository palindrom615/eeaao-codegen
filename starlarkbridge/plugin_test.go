package starlarkbridge_test

import (
	"github.com/palindrom615/eeaao-codegen/plugin"
	"github.com/palindrom615/eeaao-codegen/starlarkbridge"
	"go.starlark.net/starlark"
	"io"
	"testing"
)

type dummyPlugin struct {
}

func (d *dummyPlugin) Name() string {
	return "dummy"
}

func (d *dummyPlugin) LoadSpecFile(path string) (plugin.SpecData, error) {
	return "dummy" + path, nil
}

func (d *dummyPlugin) LoadSpec(reader io.Reader) (plugin.SpecData, error) {
	return nil, nil
}

func (d *dummyPlugin) LoadSpecUrl(url string) (plugin.SpecData, error) {
	return "dummy" + url, nil
}

func TestNewPluginStarlark(t *testing.T) {
	d := &dummyPlugin{}
	s := starlarkbridge.NewPluginStarlark(d)
	loadSpecUrlCallable, err := s.Attr("loadSpecUrl")
	if err != nil {
		t.Errorf("loadSpecFile is nil: %v", err)
	}
	loadSpecUrl, ok := loadSpecUrlCallable.(starlark.Callable)
	if !ok {
		t.Errorf("loadSpecUrl is not a starlark.Callable: %v", loadSpecUrl)
	}
}

func TestPluginMethodStarlark_CallInternal(t *testing.T) {
	d := &dummyPlugin{}
	s := starlarkbridge.NewPluginStarlark(d)
	src := `
res = dummy.loadSpecUrl("aaa")
`
	res, err := starlark.ExecFile(&starlark.Thread{Name: "main"}, "test.star", src, starlark.StringDict{
		"dummy": s,
	})
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if res["res"].(starlark.String) != "dummyaaa" {
		t.Errorf("res should be 'dummyaaa': %v", res)
	}
}

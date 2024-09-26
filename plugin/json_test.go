package plugin_test

import (
	"bytes"
	"github.com/palindrom615/eeaao-codegen/plugin"
	"testing"
)

func TestLoadSpec(t *testing.T) {
	jsonPlugin := plugin.JsonPlugin{}
	b := bytes.Buffer{}
	b.Write([]byte("{}"))
	spec, err := jsonPlugin.LoadSpec(&b)
	if err != nil {
		t.Error(err)
	}
	m := spec.(map[string]any)
	if len(m) != 0 {
		t.Errorf("specData should be empty: %v", m)
	}
}

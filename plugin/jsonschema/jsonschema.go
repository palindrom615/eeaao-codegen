package jsonschema

import (
	"git.sr.ht/~emersion/go-jsonschema"
	"github.com/palindrom615/eeaao-codegen/plugin"
	"os"
)

type JSONSchemaPlugin struct {
}

func NewJSONSchemaPlugin() *JSONSchemaPlugin {
	return &JSONSchemaPlugin{}
}

func (p *JSONSchemaPlugin) LoadSpecFile(path string) (plugin.SpecData, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := jsonschema.Schema{}
	err = data.UnmarshalJSON(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

package plugin

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type YamlPlugin struct {
}

func (y *YamlPlugin) LoadSpecFile(path string) (SpecData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return y.LoadSpec(f)
}

func (y *YamlPlugin) LoadSpec(reader io.Reader) (SpecData, error) {
	specData := make(map[string]any)
	err := yaml.NewDecoder(reader).Decode(&specData)
	if err != nil {
		return nil, err
	}
	return specData, nil
}

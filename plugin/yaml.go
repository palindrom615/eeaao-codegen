package plugin

import (
	"gopkg.in/yaml.v3"
	"os"
)

type YamlPlugin struct {
}

func (y YamlPlugin) LoadSpecFile(path string) (SpecData, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = yaml.Unmarshal(f, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

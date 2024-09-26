package plugin

import (
	"encoding/json"
	"io"
	"os"
)

type JsonPlugin struct {
}

func (j *JsonPlugin) LoadSpecFile(path string) (SpecData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return j.LoadSpec(f)
}

func (j *JsonPlugin) LoadSpec(reader io.Reader) (SpecData, error) {
	specData := make(map[string]any)
	err := json.NewDecoder(reader).Decode(&specData)
	if err != nil {
		return nil, err
	}
	return specData, nil
}

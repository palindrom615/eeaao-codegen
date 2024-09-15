package plugin

import (
	"encoding/json"
	"os"
)

type JsonPlugin struct {
}

func (j JsonPlugin) LoadSpecFile(path string) (SpecData, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(f, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

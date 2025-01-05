package plugin

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type JsonPlugin struct {
	client *http.Client
}

func NewJsonPlugin() *JsonPlugin {
	return &JsonPlugin{client: http.DefaultClient}
}

func (j *JsonPlugin) Name() string {
	return "json"
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

func (j *JsonPlugin) LoadSpecUrl(url string) (SpecData, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	return j.LoadSpec(res.Body)
}

package plugin

import (
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
)

type YamlPlugin struct {
	client *http.Client
}

func (y *YamlPlugin) Name() string {
	return "yaml"
}

func NewYamlPlugin() *YamlPlugin {
	return &YamlPlugin{client: http.DefaultClient}
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

func (y *YamlPlugin) LoadSpecUrl(url string) (SpecData, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/yaml,application/x-yaml,text/yaml")
	res, err := y.client.Do(req)
	if err != nil {
		return nil, err
	}
	return y.LoadSpec(res.Body)
}

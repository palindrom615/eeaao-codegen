package plugin

import (
	"io"
)

type OpenApiPlugin struct {
	jsonPlugin *JsonPlugin
	yamlPlugin *YamlPlugin
}

func (o *OpenApiPlugin) LoadSpecFile(path string) (SpecData, error) {
	file, err := o.jsonPlugin.LoadSpecFile(path)
	if err != nil {
		return o.yamlPlugin.LoadSpecFile(path)
	}
	return file, nil
}

func (o *OpenApiPlugin) LoadSpec(reader io.Reader) (SpecData, error) {
	specdata, err := o.jsonPlugin.LoadSpec(reader)
	if err != nil {
		return o.yamlPlugin.LoadSpec(reader)
	}
	return specdata, nil
}

func (o *OpenApiPlugin) LoadSpecUrl(url string) (SpecData, error) {
	specdata, err := o.jsonPlugin.LoadSpecUrl(url)
	if err != nil {
		return o.yamlPlugin.LoadSpecUrl(url)
	}
	return specdata, nil
}

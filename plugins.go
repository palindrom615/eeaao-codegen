package eeaao_codegen

import (
	"github.com/palindrom615/eeaao-codegen/plugin"
	"github.com/palindrom615/eeaao-codegen/plugin/jsonschema"
)

type PluginMap map[string]plugin.Plugin

var plugins = PluginMap{}

func init() {
	plugins["openapi"] = &plugin.OpenApiPlugin{}
	plugins["json"] = &plugin.JsonPlugin{}
	plugins["yaml"] = &plugin.YamlPlugin{}
	plugins["proto"] = plugin.NewProtobufPlugin()
	plugins["json-schema"] = jsonschema.NewJSONSchemaPlugin()
}

func GetPlugin(pluginName string) plugin.Plugin {
	return plugins[pluginName]
}

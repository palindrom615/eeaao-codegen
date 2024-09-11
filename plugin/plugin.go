package plugin

type SpecData any

type Plugin interface {
	LoadSpecFile(path string) SpecData
}

type PluginMap map[string]Plugin

var plugins = PluginMap{}

func init() {
	plugins["openapi"] = OpenApiPlugin{}
}

func GetPlugin(pluginName string) Plugin {
	return plugins[pluginName]
}

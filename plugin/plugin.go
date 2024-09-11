package plugin

type SpecData any

type Plugin interface {
	LoadSpecFile(path string) SpecData
}

type PluginMap map[string]Plugin

var Plugins = PluginMap{}

func init() {
	Plugins["openapi"] = OpenApiPlugin{}
}

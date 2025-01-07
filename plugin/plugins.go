package plugin

type Plugins interface {
	// GetPlugin provide Plugin by its name
	GetPlugin(pluginName string) Plugin
}

type pluginsImpl struct {
	m map[string]Plugin
}

func NewPlugins() Plugins {
	m := make(map[string]Plugin)
	m["openapi"] = NewOpenApiPlugin()
	m["json"] = NewJsonPlugin()
	m["yaml"] = NewYamlPlugin()
	m["proto"] = NewProtobufPlugin()
	return &pluginsImpl{
		m: m,
	}
}

func (p *pluginsImpl) GetPlugin(pluginName string) Plugin {
	return p.m[pluginName]
}

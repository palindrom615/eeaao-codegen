// Package plugin provides the interface for plugins to load spec files.
// The plugins are registered in the init function.
//   - openapi: OpenAPI plugin
//   - json: JSON plugin
//   - yaml: YAML plugin

package plugin

import "io"

// SpecData represent the data of a specification. it can be any type, depending on the plugin. For example, OpenAPI plugin uses go-openapi/loads.Spec() to load a spec file.
//
// The SpecData loaded by plugin is flattened to a map[string]interface{} via encoding/json/json.Marshal and json.Unmarshal for security.
// Thus, the field of SpecData that is not exported will not be exposed on rendering.
type SpecData any

type Plugin interface {
	LoadSpecFile(path string) (SpecData, error)
	LoadSpec(reader io.Reader) (SpecData, error)
}

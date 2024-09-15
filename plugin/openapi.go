package plugin

import (
	"github.com/go-openapi/loads"
)

type OpenApiPlugin struct {
}

func (o OpenApiPlugin) LoadSpecFile(path string) (SpecData, error) {
	return loads.Spec(path)
}

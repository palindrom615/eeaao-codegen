package plugin

import (
	"github.com/go-openapi/loads"
)

type OpenApiPlugin struct {
}

func (o OpenApiPlugin) LoadSpecFile(path string) (SpecData, error) {
	doc, err := loads.Spec(path)
	if err != nil {
		return nil, err
	}
	return doc.Spec(), nil
}

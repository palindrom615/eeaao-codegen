package plugin

import (
	"github.com/go-openapi/loads"
	"log"
)

type OpenApiPlugin struct {
}

func (o OpenApiPlugin) LoadSpecFile(path string) SpecData {
	doc, err := loads.JSONSpec(path)
	if err != nil {
		log.Printf("Error reading doc file: %v\n", err)
		return nil
	}
	return doc.Spec()
}

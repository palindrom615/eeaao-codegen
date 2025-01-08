package eeaao_codegen

import (
	"testing"
)

func TestRender(t *testing.T) {
	a := NewApp(
		"./example/openapi-v3/kotlin-spring/build",
		"./example/openapi-v3/kotlin-spring/codelet",
		"",
	)
	a.Render()
}

func TestTemplate_Include(t *testing.T) {
	app := NewApp("test/out", "test/addTemplateFunc", "")
	rendered, err := app.tmpl.Include("add.tmpl", "test/addTemplateFunc/add.star")
	if err != nil {
		t.Errorf("Error rendering template: %v", err)
	}
	if rendered != "1 + 1 = 2" {
		t.Errorf("Expected `1 + 1 = 2`, got %s", rendered)
	}
}

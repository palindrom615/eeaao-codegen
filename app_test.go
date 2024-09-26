package eeaao_codegen_test

import (
	"github.com/palindrom615/eeaao-codegen"
	"testing"
)

func TestRender(t *testing.T) {
	a := eeaao_codegen.NewApp(
		"./example/openapi-v3/spec",
		"./example/openapi-v3/kotlin-spring/build",
		"./example/openapi-v3/kotlin-spring/codelet",
		"./example/openapi-v3/kotlin-spring/values.yaml",
	)
	a.Render()
}

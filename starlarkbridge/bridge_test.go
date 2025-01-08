package starlarkbridge_test

import (
	"github.com/palindrom615/eeaao-codegen/starlarkbridge"
	"go.starlark.net/starlark"
	"testing"
)

func TestConvertFromStarlarkValue(t *testing.T) {
	thread := &starlark.Thread{
		Name: "test",
	}
	v, err := starlarkbridge.ConvertFromStarlarkValue(thread, starlark.String("hello"))
	if err != nil {
		t.Errorf("Error converting from starlark value: %v", err)
	}
	if v != "hello" {
		t.Errorf("Expected `hello`, got %v", v)
	}

	v, err = starlarkbridge.ConvertFromStarlarkValue(thread, starlark.MakeInt(100))
	if err != nil {
		t.Errorf("Error converting from starlark value: %v", err)
	}
	if v != 100 {
		t.Errorf("Expected 100, got %v", v)
	}

	v, err = starlarkbridge.ConvertFromStarlarkValue(thread, starlark.MakeInt(10_000_000_000_000))
	if err != nil {
		t.Errorf("Error converting from starlark value: %v", err)
	}
	if v != int64(10_000_000_000_000) {
		t.Errorf("Expected 10_000_000_000_000, got %v", v)
	}

	v, err = starlarkbridge.ConvertFromStarlarkValue(thread, starlark.Bool(false))
	if err != nil {
		t.Errorf("Error converting from starlark value: %v", err)
	}
	if v != false {
		t.Errorf("Expected false, got %v", v)
	}

	v, err = starlarkbridge.ConvertFromStarlarkValue(thread, starlark.None)
	if err != nil {
		t.Errorf("Error converting from starlark value: %v", err)
	}
	if v != nil {
		t.Errorf("Expected nil, got %v", v)
	}
}

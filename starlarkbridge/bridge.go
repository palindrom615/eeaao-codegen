package starlarkbridge

import (
	"encoding/json"
	json2 "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
	"log"
)

// ConvertToStarlarkValue converts any value to starlark.Value
//
// Internally the value is serialized in go and then deserialized in starlark.
// JSON is used for serialization. so any properties that are not serializable to JSON will be lost.
func ConvertToStarlarkValue(thread *starlark.Thread, value any) (starlark.Value, error) {
	specStr, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return decodeWithStarlarkJson(thread, starlark.String(specStr))
}

// ConvertFromStarlarkValue converts starlark.Value to any value
//
// As ConvertToStarlarkValue, the value is serialized in go and then deserialized in starlark.
// JSON is used for serialization. so any properties that are not serializable to JSON will be lost.
func ConvertFromStarlarkValue(thread *starlark.Thread, value starlark.Value) (any, error) {
	switch t := value.(type) {
	case starlark.String:
		return t.GoString(), nil
	case starlark.Bool:
		return bool(t), nil
	case starlark.Int:
		i32, err := starlark.AsInt32(t)
		if err == nil {
			return i32, nil
		}
		i64, ok := t.Int64()
		if ok {
			return i64, nil
		}
		return t.BigInt(), nil
	case starlark.Float:
		return float64(t), nil
	case starlark.NoneType:
		return nil, nil
	}
	encoded, err := encodeWithStarlarkJson(thread, value)
	if err != nil {
		return nil, err
	}
	d := make(map[string]any)
	err = json.Unmarshal([]byte(encoded.(starlark.String)), &d)
	if err != nil {
		log.Printf("Error decoding starlark injected data: %v\n%v\n", encoded, err)
		return nil, err
	}
	return d, nil
}

func decodeWithStarlarkJson(thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	decode := json2.Module.Members["decode"].(*starlark.Builtin)
	return starlark.Call(thread, decode, starlark.Tuple{value}, nil)
}

func encodeWithStarlarkJson(thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	encode := json2.Module.Members["encode"].(*starlark.Builtin)
	return starlark.Call(thread, encode, starlark.Tuple{value}, nil)
}

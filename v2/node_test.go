package jd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJsonNodeTypes(t *testing.T) {
	// map[interface{}]interface{} with string keys (YAML-style)
	yamlMap := map[interface{}]interface{}{"key": "value"}
	n, err := NewJsonNode(yamlMap)
	if err != nil {
		t.Fatal(err)
	}
	if n.Json() != `{"key":"value"}` {
		t.Errorf("got %v", n.Json())
	}

	// map[interface{}]interface{} with non-string key
	badMap := map[interface{}]interface{}{42: "value"}
	_, err = NewJsonNode(badMap)
	if err == nil {
		t.Fatal("expected error for non-string key")
	}

	// Direct jsonObject passthrough
	obj := jsonObject(map[string]JsonNode{"a": jsonString("b")})
	n, err = NewJsonNode(obj)
	if err != nil || n.Json() != `{"a":"b"}` {
		t.Errorf("jsonObject passthrough failed")
	}

	// Direct jsonArray passthrough
	arr := jsonArray{jsonNumber(1)}
	n, err = NewJsonNode(arr)
	if err != nil || n.Json() != `[1]` {
		t.Errorf("jsonArray passthrough failed")
	}

	// int type
	n, err = NewJsonNode(42)
	if err != nil || n.Json() != `42` {
		t.Errorf("int conversion failed")
	}

	// Direct jsonNumber passthrough
	n, err = NewJsonNode(jsonNumber(3.14))
	if err != nil || n.Json() != `3.14` {
		t.Errorf("jsonNumber passthrough failed")
	}

	// Direct jsonString passthrough
	n, err = NewJsonNode(jsonString("hello"))
	if err != nil || n.Json() != `"hello"` {
		t.Errorf("jsonString passthrough failed")
	}

	// Direct jsonBool passthrough
	n, err = NewJsonNode(jsonBool(true))
	if err != nil || n.Json() != `true` {
		t.Errorf("jsonBool passthrough failed")
	}

	// Direct jsonNull passthrough
	n, err = NewJsonNode(jsonNull(nil))
	if err != nil || n.Json() != `null` {
		t.Errorf("jsonNull passthrough failed")
	}

	// Unsupported type
	_, err = NewJsonNode(complex(1, 2))
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}

	// map[interface{}]interface{} with JsonNode value
	yamlMapNode := map[interface{}]interface{}{"k": jsonString("v")}
	n, err = NewJsonNode(yamlMapNode)
	if err != nil {
		t.Fatal(err)
	}
	if n.Json() != `{"k":"v"}` {
		t.Errorf("got %v", n.Json())
	}
}

func TestNodeListEmpty(t *testing.T) {
	l := nodeList()
	if len(l) != 0 {
		t.Errorf("expected empty list, got %v", l)
	}
}

func TestNewJsonNode(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    JsonNode
		wantErr bool
	}{
		{
			name:  "float64 to jsonNumber",
			input: float64(1),
			want:  jsonNumber(1),
		},
		{
			name:  "jsonNumber passthrough",
			input: jsonNumber(1),
			want:  jsonNumber(1),
		},
		{
			name:  "int to jsonNumber",
			input: int(1),
			want:  jsonNumber(1),
		},
		{
			name:  "bool true to jsonBool",
			input: true,
			want:  jsonBool(true),
		},
		{
			name:  "bool false to jsonBool",
			input: false,
			want:  jsonBool(false),
		},
		{
			name:  "jsonBool passthrough",
			input: jsonBool(true),
			want:  jsonBool(true),
		},
		{
			name:  "string to jsonString",
			input: "foo",
			want:  jsonString("foo"),
		},
		{
			name:  "empty string to jsonString",
			input: "",
			want:  jsonString(""),
		},
		{
			name:  "jsonString passthrough",
			input: jsonString("foo"),
			want:  jsonString("foo"),
		},
		{
			name:  "nil to jsonNull",
			input: nil,
			want:  jsonNull(nil),
		},
		{
			name:  "jsonNull passthrough",
			input: jsonNull(nil),
			want:  jsonNull(nil),
		},
		{
			name:  "slice to jsonArray with numbers",
			input: []any{1, 2, 3},
			want:  jsonArray{jsonNumber(1), jsonNumber(2), jsonNumber(3)},
		},
		{
			name:  "slice with jsonNumber to jsonArray",
			input: []any{jsonNumber(1)},
			want:  jsonArray{jsonNumber(1)},
		},
		{
			name:  "empty slice to jsonArray",
			input: []any{},
			want:  jsonArray{},
		},
		{
			name:  "jsonArray passthrough",
			input: jsonArray{},
			want:  jsonArray{},
		},
		{
			name:  "map to jsonObject",
			input: map[string]any{"foo": 1},
			want:  jsonObject{"foo": jsonNumber(1)},
		},
		{
			name:  "map with jsonNumber to jsonObject",
			input: map[string]any{"foo": jsonNumber(1)},
			want:  jsonObject{"foo": jsonNumber(1)},
		},
		{
			name:  "empty map to jsonObject",
			input: map[string]any{},
			want:  jsonObject{},
		},
		{
			name:  "jsonObject passthrough",
			input: jsonObject{},
			want:  jsonObject{},
		},
		{
			name:  "mixed types in slice",
			input: []any{1, "hello", true, nil},
			want:  jsonArray{jsonNumber(1), jsonString("hello"), jsonBool(true), jsonNull(nil)},
		},
		{
			name:  "nested map in slice",
			input: []any{map[string]any{"key": "value"}},
			want:  jsonArray{jsonObject{"key": jsonString("value")}},
		},
		{
			name:  "nested slice in map",
			input: map[string]any{"array": []any{1, 2}},
			want:  jsonObject{"array": jsonArray{jsonNumber(1), jsonNumber(2)}},
		},
		{
			name:  "complex nested structure",
			input: map[string]any{"users": []any{map[string]any{"id": 1, "name": "John"}}},
			want:  jsonObject{"users": jsonArray{jsonObject{"id": jsonNumber(1), "name": jsonString("John")}}},
		},
		{
			name:  "negative number",
			input: -42,
			want:  jsonNumber(-42),
		},
		{
			name:  "decimal number",
			input: 3.14,
			want:  jsonNumber(3.14),
		},
		{
			name:  "zero value",
			input: 0,
			want:  jsonNumber(0),
		},
		{
			name:  "large number",
			input: 1000000,
			want:  jsonNumber(1000000),
		},
		{
			name:  "string with special characters",
			input: "hello\nworld",
			want:  jsonString("hello\nworld"),
		},
		{
			name:  "unicode string",
			input: "héllo wørld",
			want:  jsonString("héllo wørld"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJsonNode(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.True(t, tt.want.Equals(got))
			}
		})
	}
}

func TestNewJsonNodeNestedErrors(t *testing.T) {
	// map[string]interface{} with unsupported nested value
	_, err := NewJsonNode(map[string]interface{}{"k": complex(1, 2)})
	if err == nil {
		t.Fatal("expected error for unsupported nested type in map[string]interface{}")
	}
	// map[interface{}]interface{} with unsupported nested value
	_, err = NewJsonNode(map[interface{}]interface{}{"k": complex(1, 2)})
	if err == nil {
		t.Fatal("expected error for unsupported nested type in map[interface{}]interface{}")
	}
	// []interface{} with unsupported element
	_, err = NewJsonNode([]interface{}{complex(1, 2)})
	if err == nil {
		t.Fatal("expected error for unsupported element type in []interface{}")
	}
}

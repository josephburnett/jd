package jd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

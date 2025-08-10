package jd

import (
	"github.com/josephburnett/jd/v2/internal/types"
)

// JsonNode is a JSON value, collection of values, or a void representing
// the absense of a value. JSON values can be a Number, String, Boolean
// or Null. Collections can be an Object, native JSON array, ordered
// List, unordered Set or Multiset. JsonNodes are created with the
// NewJsonNode function or ReadJson* and ReadYaml* functions.
type JsonNode = types.JsonNode

// Note: All node implementation is now in internal/types

// NewJsonNode constructs a JsonNode from native Golang objects. See the
// function source for supported types and conversions. Slices are always
// placed into native JSON Arrays and interpretated as Lists, Sets or
// Multisets based on Metadata provided during Equals and Diff
// operations.
func NewJsonNode(n interface{}) (JsonNode, error) {
	return types.NewJsonNode(n)
}

func nodeList(n ...JsonNode) []JsonNode {
	return types.NodeList(n...)
}

func ReadJsonString(s string) (JsonNode, error) {
	return types.ReadJsonString(s)
}

func ReadYamlString(s string) (JsonNode, error) {
	return types.ReadYamlString(s)
}

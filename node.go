package jd

import (
	"fmt"
)

type JsonNode interface {
	diff(b JsonNode, path Path) Diff
	equals(b JsonNode) bool
}

func newJsonNode(n interface{}) JsonNode {
	switch t := n.(type) {
	case map[string]interface{}:
		return JsonStruct(t)
	case []interface{}:
		return JsonList(t)
	case float64:
		return JsonNumber(t)
	case string:
		return JsonString(t)
	default:
		panic(fmt.Sprintf("Unexpected type %v", t))
	}
}

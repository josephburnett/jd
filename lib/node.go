package jd

import (
	"encoding/json"
	"errors"
	"fmt"
)

type JsonNode interface {
	Json(metadata ...Metadata) string
	Equals(n JsonNode, metadata ...Metadata) bool
	hashCode(metadata []Metadata) [8]byte
	Diff(n JsonNode, metadata ...Metadata) Diff
	diff(n JsonNode, p path, metadata []Metadata) Diff
	Patch(d Diff) (JsonNode, error)
	patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode) (JsonNode, error)
}

func NewJsonNode(n interface{}) (JsonNode, error) {
	switch t := n.(type) {
	case map[string]interface{}:
		m := jsonObject{
			properties: make(map[string]JsonNode),
			idKeys:     make(map[string]bool),
		}
		for k, v := range t {
			if _, ok := v.(JsonNode); !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				m.properties[k] = e
			}
		}
		return m, nil
	case []interface{}:
		l := make(jsonArray, len(t))
		for i, v := range t {
			if _, ok := v.(JsonNode); !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				l[i] = e
			}
		}
		return l, nil
	case float64:
		return jsonNumber(t), nil
	case string:
		return jsonString(t), nil
	case bool:
		return jsonBool(t), nil
	case nil:
		return jsonNull(nil), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported type %v", t))
	}
}

func nodeList(n ...JsonNode) []JsonNode {
	l := []JsonNode{}
	if len(n) == 0 {
		return l
	}
	if n[0].Equals(voidNode{}) {
		return l
	}
	return append(l, n...)
}

func renderJson(n JsonNode) string {
	s, _ := json.Marshal(n)
	// Errors are ignored because JsonNode types are
	// private and known to marshal without error.
	return string(s)
}

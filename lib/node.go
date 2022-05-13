package jd

import (
	"errors"
	"fmt"
)

type JsonNode interface {
	Json(metadata ...Metadata) string
	Yaml(metadata ...Metadata) string
	raw(metadata []Metadata) interface{}
	Equals(n JsonNode, metadata ...Metadata) bool
	hashCode(metadata []Metadata) [8]byte
	Diff(n JsonNode, metadata ...Metadata) Diff
	diff(n JsonNode, p path, metadata []Metadata, strategy patchStrategy) Diff
	Patch(d Diff) (JsonNode, error)
	patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode, strategy patchStrategy) (JsonNode, error)
}

func NewJsonNode(n interface{}) (JsonNode, error) {
	switch t := n.(type) {
	case map[string]interface{}:
		m := newJsonObject()
		for k, v := range t {
			n, ok := v.(JsonNode)
			if !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				n = e
			}
			m.properties[k] = n
		}
		return m, nil
	case map[interface{}]interface{}:
		m := newJsonObject()
		for k, v := range t {
			s, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("Unsupported key type %T", k)
			}
			if _, ok := v.(JsonNode); !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				m.properties[s] = e
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
	case int:
		return jsonNumber(t), nil
	case string:
		return jsonString(t), nil
	case bool:
		return jsonBool(t), nil
	case nil:
		return jsonNull(nil), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported type %T", t))
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

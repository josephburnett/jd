package jd

import (
	"errors"
	"fmt"
)

type JsonNode interface {
	Equals(n JsonNode) bool
	Diff(n JsonNode) Diff
	diff(n JsonNode, p Path) Diff
}

func NewJsonNode(n interface{}) (JsonNode, error) {
	switch t := n.(type) {
	case map[string]interface{}:
		m := make(jsonStruct)
		for k, v := range t {
			if _, ok := v.(JsonNode); !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				m[k] = e
			}
		}
		return m, nil
	case []interface{}:
		l := make(jsonList, len(t))
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
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported type %v", t))
	}
}

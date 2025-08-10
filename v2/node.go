package jd

import (
	"fmt"
)

// JsonNode is a JSON value, collection of values, or a void representing
// the absense of a value. JSON values can be a Number, String, Boolean
// or Null. Collections can be an Object, native JSON array, ordered
// List, unordered Set or Multiset. JsonNodes are created with the
// NewJsonNode function or ReadJson* and ReadYaml* functions.
type JsonNode interface {

	// Json renders a JsonNode as a JSON string.
	Json(renderOptions ...Option) string

	// Yaml renders a JsonNode as a YAML string in block format.
	Yaml(renderOptions ...Option) string

	// Equals returns true if the JsonNodes are equal according to
	// the provided Metadata. The default behavior (no Metadata) is
	// to compare the entire structure down to scalar values treating
	// Arrays as orders Lists. The SET and MULTISET Metadata will
	// treat Arrays as sets or multisets (bags) respectively. To deep
	// compare objects in an array irrespective of order, the SetKeys
	// function will construct Metadata to compare objects by a set
	// of keys. If two JsonNodes are equal, then Diff with the same
	// Metadata will produce an empty Diff. And vice versa.
	Equals(n JsonNode, options ...Option) bool

	// Diff produces a list of differences (Diff) between two
	// JsonNodes such that if the output Diff were applied to the
	// first JsonNode (Patch) then the two JsonNodes would be
	// Equal. The necessary Metadata is embeded in the Diff itself so
	// only the Diff is required to Patch a JsonNode.
	Diff(n JsonNode, options ...Option) Diff

	// Patch applies a Diff to a JsonNode. No Metadata is provided
	// because the original interpretation of the structure is
	// embedded in the Diff itself.
	Patch(d Diff) (JsonNode, error)

	jsonNodeInternals
}

type jsonNodeInternals interface {
	raw() interface{}
	hashCode(opts *options) [8]byte
	equals(n JsonNode, o *options) bool
	diff(n JsonNode, p Path, opts *options, strategy patchStrategy) Diff
	patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error)
}

// NewJsonNode constructs a JsonNode from native Golang objects. See the
// function source for supported types and conversions. Slices are always
// placed into native JSON Arrays and interpretated as Lists, Sets or
// Multisets based on Metadata provided during Equals and Diff
// operations.
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
			m[k] = n
		}
		return m, nil
	case map[interface{}]interface{}:
		m := newJsonObject()
		for k, v := range t {
			s, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("unsupported key type %T", k)
			}
			n, ok := v.(JsonNode)
			if !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				n = e
			}
			m[s] = n
		}
		return m, nil
	case jsonObject:
		return t, nil
	case []interface{}:
		l := make(jsonArray, len(t))
		for i, v := range t {
			n, ok := v.(JsonNode)
			if !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				n = e
			}
			l[i] = n
		}
		return l, nil
	case jsonArray:
		return t, nil
	case float64:
		return jsonNumber(t), nil
	case int:
		return jsonNumber(t), nil
	case jsonNumber:
		return t, nil
	case string:
		return jsonString(t), nil
	case jsonString:
		return t, nil
	case bool:
		return jsonBool(t), nil
	case jsonBool:
		return t, nil
	case nil:
		return jsonNull(nil), nil
	case jsonNull:
		return t, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", t)
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

package jd

import (
	"fmt"
	"sort"
)

type jsonObject map[string]JsonNode

var _ JsonNode = jsonObject{}

func newJsonObject() jsonObject {
	return jsonObject{}
}

func (o jsonObject) Json(_ ...Option) string {
	return renderJson(o.raw())
}

func (o jsonObject) MarshalJSON() ([]byte, error) {
	return []byte(o.Json()), nil
}

func (o jsonObject) Yaml(_ ...Option) string {
	return renderYaml(o.raw())
}

func (o jsonObject) raw() interface{} {
	j := make(map[string]interface{})
	for k, v := range o {
		j[k] = v.raw()
	}
	return j
}

func (o1 jsonObject) Equals(n JsonNode, opts ...Option) bool {
	o := refine(&options{retain: opts}, nil)
	return o1.equals(n, o)
}

func (o1 jsonObject) equals(n JsonNode, o *options) bool {
	o2, ok := n.(jsonObject)
	if !ok {
		return false
	}
	if len(o1) != len(o2) {
		return false
	}

	for key1, val1 := range o1 {
		val2, ok := o2[key1]
		if !ok {
			return false
		}
		ret := val1.equals(val2, o)
		if !ret {
			return false
		}
	}
	return true
}

func (o jsonObject) hashCode(opts *options) [8]byte {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	a := []byte{0x00, 0x5D, 0x39, 0xA4, 0x18, 0x10, 0xEA, 0xD5} // random bytes
	for _, k := range keys {
		keyHash := hash([]byte(k))
		a = append(a, keyHash[:]...)
		valueHash := o[k].hashCode(opts)
		a = append(a, valueHash[:]...)
	}
	return hash(a)
}

// ident is the identity of the json object based on either the hash of a
// given set of keys or the full object if no keys are present.
func (o jsonObject) ident(opts *options) [8]byte {
	keys, ok := getOption[setKeysOption](opts)
	if !ok {
		return o.hashCode(opts)
	}
	hashes := hashCodes{
		// We start with a constant hash to distinguish between
		// an empty object and an empty array.
		[8]byte{0x4B, 0x08, 0xD2, 0x0F, 0xBD, 0xC8, 0xDE, 0x9A}, // randomly chosen bytes
	}
	for _, key := range []string(*keys) {
		v, ok := o[key]
		if ok {
			hashes = append(hashes, v.hashCode(opts))
		}
	}
	if len(hashes) == 0 {
		return o.hashCode(opts)
	}
	return hashes.combine()
}

func (o jsonObject) pathIdent(pathObject jsonObject, opts *options) [8]byte {
	keys := []string{}
	for k := range pathObject {
		keys = append(keys, k)
	}
	id := make(map[string]interface{})
	for _, key := range keys {
		if value, ok := o[key]; ok {
			id[key] = value
		}
	}
	e, _ := NewJsonNode(id)
	return e.hashCode(&options{})
}

func (o jsonObject) Diff(n JsonNode, opts ...Option) Diff {
	op := &options{retain: opts}
	return o.diff(n, make(Path, 0), op, getPatchStrategy(op))
}

func (o1 jsonObject) diff(
	n JsonNode,
	path Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	d := make(Diff, 0)
	o2, ok := n.(jsonObject)
	if !ok {
		// Different types
		var e DiffElement
		switch strategy {
		case mergePatchStrategy:
			e = DiffElement{
				Metadata: Metadata{
					Merge: true,
				},
				Path: path.clone(),
				Add:  []JsonNode{n},
			}
		default:
			e = DiffElement{
				Path:   path.clone(),
				Remove: []JsonNode{o1},
				Add:    []JsonNode{n},
			}
		}
		return append(d, e)
	}
	o1Keys := make([]string, 0, len(o1))
	for k := range o1 {
		o1Keys = append(o1Keys, k)
	}
	sort.Strings(o1Keys)
	o2Keys := make([]string, 0, len(o2))
	for k := range o2 {
		o2Keys = append(o2Keys, k)
	}
	sort.Strings(o2Keys)
	for _, k1 := range o1Keys {
		v1 := o1[k1]
		if v2, ok := o2[k1]; ok {
			// Both keys are present
			o := refine(opts, PathKey(k1))
			subDiff := v1.diff(v2, append(path, PathKey(k1)), o, strategy)
			d = append(d, subDiff...)
		} else {
			// O2 missing key
			var e DiffElement
			switch strategy {
			case mergePatchStrategy:
				e = DiffElement{
					Metadata: Metadata{
						Merge: true,
					},
					Path: append(path, PathKey(k1)).clone(),
					Add:  []JsonNode{voidNode{}},
				}
			default:
				e = DiffElement{
					Path:   append(path, PathKey(k1)).clone(),
					Remove: nodeList(v1),
					Add:    nodeList(),
				}
			}
			d = append(d, e)
		}
	}
	for _, k2 := range o2Keys {
		v2 := o2[k2]
		if _, ok := o1[k2]; !ok {
			// O1 missing key
			var e DiffElement
			switch strategy {
			case mergePatchStrategy:
				e = DiffElement{
					Metadata: Metadata{
						Merge: true,
					},
					Path:   append(path, PathKey(k2)).clone(),
					Remove: nodeList(),
					Add:    nodeList(v2),
				}
			default:
				e = DiffElement{
					Path:   append(path, PathKey(k2)).clone(),
					Remove: nodeList(),
					Add:    nodeList(v2),
				}
			}
			d = append(d, e)
		}
	}
	return d
}

func (o jsonObject) Patch(d Diff) (JsonNode, error) {
	return patchAll(o, d)
}

func (o jsonObject) patch(
	pathBehind, pathAhead Path,
	before, oldValues, newValues, after []JsonNode,
	strategy patchStrategy,
) (JsonNode, error) {
	if (len(pathAhead) == 0) && (len(oldValues) > 1 || len(newValues) > 1) {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	// Base case
	if len(pathAhead) == 0 {
		newValue := singleValue(newValues)
		if strategy == mergePatchStrategy {
			return newValue, nil
		}
		oldValue := singleValue(oldValues)
		if !o.Equals(oldValue) {
			return patchErrExpectValue(oldValue, o, pathBehind)
		}
		return newValue, nil
	}
	// Recursive case
	n, _, rest := pathAhead.next()
	pe, ok := n.(PathKey)
	if !ok {
		return nil, fmt.Errorf(
			"found %v at %v: expected JSON object",
			o.Json(), pathBehind)
	}
	nextNode, ok := o[string(pe)]
	if !ok {
		switch strategy {
		case mergePatchStrategy:
			// Create objects
			if len(rest) == 0 {
				nextNode = voidNode{}
			} else {
				nextNode = newJsonObject()
			}
		case strictPatchStrategy:
			nextNode = voidNode{}
		default:
			return patchErrUnsupportedPatchStrategy(pathBehind, strategy)
		}
	}
	patchedNode, err := nextNode.patch(append(pathBehind, pe), rest, before, oldValues, newValues, after, strategy)
	if err != nil {
		return nil, err
	}
	if isVoid(patchedNode) {
		// Delete a pair
		delete(o, string(pe))
	} else {
		// Add or replace a pair
		o[string(pe)] = patchedNode
	}
	return o, nil
}

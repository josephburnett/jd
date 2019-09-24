package jd

import (
	"encoding/json"
	"fmt"
	"sort"
)

type jsonObject struct {
	properties map[string]JsonNode
	// TODO: drop idKeys.
	idKeys map[string]bool
}

var _ JsonNode = jsonObject{}

func (o jsonObject) Json(metadata ...Metadata) string {
	j := make(map[string]interface{})
	for k, v := range o.properties {
		j[k] = v
	}
	s, _ := json.Marshal(j)
	return string(s)
}

func (o jsonObject) MarshalJSON() ([]byte, error) {
	return []byte(o.Json()), nil
}

func (o1 jsonObject) Equals(n JsonNode, metadata ...Metadata) bool {
	o2, ok := n.(jsonObject)
	if !ok {
		return false
	}
	if len(o1.properties) != len(o2.properties) {
		return false
	}

	for key1, val1 := range o1.properties {
		val2, ok := o2.properties[key1]
		if !ok {
			return false
		}
		ret := val1.Equals(val2)
		if !ret {
			return false
		}
	}
	return true
}

func (o jsonObject) hashCode(metadata []Metadata) [8]byte {
	keys := make([]string, 0, len(o.properties))
	for k := range o.properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	a := make([]byte, 0, len(o.properties)*16)
	for _, k := range keys {
		keyHash := hash([]byte(k))
		a = append(a, keyHash[:]...)
		valueHash := o.properties[k].hashCode(metadata)
		a = append(a, valueHash[:]...)
	}
	return hash(a)
}

// ident is the identity of the json object based on either the hash of a
// given set of keys or the full object if no keys are present.
func (o jsonObject) ident(metadata []Metadata) [8]byte {
	keys := getSetkeysMetadata(metadata).mergeKeys(o.idKeys)
	if len(keys) == 0 {
		return o.hashCode(metadata)
	}
	hashes := make(hashCodes, 0)
	for key := range keys {
		v, ok := o.properties[key]
		if ok {
			hashes = append(hashes, v.hashCode(metadata))
		}
	}
	if len(hashes) == 0 {
		return o.hashCode(metadata)
	}
	return hashes.combine()
}

func (o jsonObject) pathIdent(metadata []Metadata) JsonNode {
	keys := getSetkeysMetadata(metadata).mergeKeys(o.idKeys)
	id := make(map[string]interface{})
	for key := range keys {
		if value, ok := o.properties[key]; ok {
			id[key] = value
		}
	}
	e, _ := NewJsonNode(id)
	return e
}

func (k1 *setkeysMetadata) mergeKeys(k2 map[string]bool) map[string]bool {
	if k1 == nil {
		// Nothing to merge
		return k2
	}
	k3 := make(map[string]bool)
	for k := range k1.keys {
		k3[k] = true
	}
	for k := range k2 {
		k3[k] = true
	}
	return k3
}

func (o jsonObject) Diff(n JsonNode, metadata ...Metadata) Diff {
	return o.diff(n, make(path, 0), metadata)
}

func (o1 jsonObject) diff(n JsonNode, path path, metadata []Metadata) Diff {
	d := make(Diff, 0)
	o2, ok := n.(jsonObject)
	if !ok {
		// Different types
		e := DiffElement{
			Path:      path.clone(),
			OldValues: []JsonNode{o1},
			NewValues: []JsonNode{n},
		}
		return append(d, e)
	}
	o1Keys := make([]string, 0, len(o1.properties))
	for k := range o1.properties {
		o1Keys = append(o1Keys, k)
	}
	sort.Strings(o1Keys)
	o2Keys := make([]string, 0, len(o2.properties))
	for k := range o2.properties {
		o2Keys = append(o2Keys, k)
	}
	sort.Strings(o2Keys)
	for _, k1 := range o1Keys {
		v1 := o1.properties[k1]
		if v2, ok := o2.properties[k1]; ok {
			// Both keys are present
			subDiff := v1.diff(v2, append(path, jsonString(k1)), metadata)
			d = append(d, subDiff...)
		} else {
			// O2 missing key
			e := DiffElement{
				Path:      append(path, jsonString(k1)).clone(),
				OldValues: nodeList(v1),
				NewValues: nodeList(),
			}
			d = append(d, e)
		}
	}
	for _, k2 := range o2Keys {
		v2 := o2.properties[k2]
		if _, ok := o1.properties[k2]; !ok {
			// O1 missing key
			e := DiffElement{
				Path:      append(path, jsonString(k2)).clone(),
				OldValues: nodeList(),
				NewValues: nodeList(v2),
			}
			d = append(d, e)
		}
	}
	return d
}

func (o jsonObject) Patch(d Diff) (JsonNode, error) {
	return patchAll(o, d)
}

func (o jsonObject) patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode) (JsonNode, error) {
	if (len(pathAhead) == 0) && (len(oldValues) > 1 || len(newValues) > 1) {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	// Base case
	if len(pathAhead) == 0 {
		oldValue := singleValue(oldValues)
		newValue := singleValue(newValues)
		if !o.Equals(oldValue) {
			return patchErrExpectValue(oldValue, o, pathBehind)
		}
		return newValue, nil
	}
	// Recursive case
	n, _, rest := pathAhead.next()
	pe, ok := n.(jsonString)
	if !ok {
		return nil, fmt.Errorf(
			"Found %v at %v. Expected JSON object.",
			o.Json(), pathBehind)
	}
	nextNode, ok := o.properties[string(pe)]
	if !ok {
		nextNode = voidNode{}
	}
	patchedNode, err := nextNode.patch(append(pathBehind, pe), rest, oldValues, newValues)
	if err != nil {
		return nil, err
	}
	if isVoid(patchedNode) {
		// Delete a pair
		delete(o.properties, string(pe))
	} else {
		// Add or replace a pair
		o.properties[string(pe)] = patchedNode
	}
	return o, nil
}

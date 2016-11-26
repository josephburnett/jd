package jd

import (
	"fmt"
)

type jsonNull struct{}

var _ JsonNode = jsonNull{}

func (n jsonNull) Json() string {
	return renderJson(nil)
}

func (n jsonNull) Equals(o JsonNode) bool {
	switch o.(type) {
	case jsonNull:
		return true
	default:
		return false
	}
}

func (n jsonNull) Diff(o JsonNode) Diff {
	return n.diff(o, Path{})
}

func (n jsonNull) diff(o JsonNode, path Path) Diff {
	d := make(Diff, 0)
	if n.Equals(o) {
		return d
	}
	e := DiffElement{
		Path:     path.clone(),
		OldValue: n,
		NewValue: o,
	}
	return append(d, e)
}

func (n jsonNull) Patch(d Diff) (JsonNode, error) {
	return patchAll(n, d)
}

func (n jsonNull) patch(pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(n, pathAhead[0])
	}
	if !n.Equals(oldValue) {
		return nil, fmt.Errorf(
			"Found %v at %v. Expected %v.",
			n.Json(), pathBehind, oldValue.Json())
	}
	return newValue, nil
}

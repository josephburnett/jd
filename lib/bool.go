package jd

import (
	"fmt"
)

type jsonBool bool

var _ JsonNode = jsonBool(true)

func (b jsonBool) Json() string {
	return renderJson(b)
}

func (b1 jsonBool) Equals(n JsonNode) bool {
	b2, ok := n.(jsonBool)
	if !ok {
		return false
	}
	return b1 == b2
}

func (b jsonBool) Diff(n JsonNode) Diff {
	return b.diff(n, Path{})
}

func (b jsonBool) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	if b.Equals(n) {
		return d
	}
	e := DiffElement{
		Path:     path.clone(),
		OldValue: b,
		NewValue: n,
	}
	return append(d, e)
}

func (b jsonBool) Patch(d Diff) (JsonNode, error) {
	return patchAll(b, d)
}

func (b jsonBool) patch(pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(b, pathAhead[0])
	}
	if !b.Equals(oldValue) {
		return nil, fmt.Errorf(
			"Found %v at %v. Expected %v.",
			b.Json(), pathBehind, oldValue.Json())
	}
	return newValue, nil
}

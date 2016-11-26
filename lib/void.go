package jd

import (
	"fmt"
)

type voidNode struct{}

var _ JsonNode = voidNode{}

func (e voidNode) Json() string {
	return ""
}

func (e voidNode) Equals(n JsonNode) bool {
	switch n.(type) {
	case voidNode:
		return true
	default:
		return false
	}
}

func isVoid(n JsonNode) bool {
	if n == nil {
		return false
	}
	if _, ok := n.(voidNode); ok {
		return true
	}
	return false
}

func (e voidNode) Diff(n JsonNode) Diff {
	return e.diff(n, Path{})
}

func (e voidNode) diff(n JsonNode, p Path) Diff {
	d := make(Diff, 0)
	if n.Equals(e) {
		return d
	}
	de := DiffElement{
		Path:     p,
		OldValue: e,
		NewValue: n,
	}
	return append(d, de)
}

func (e voidNode) Patch(d Diff) (JsonNode, error) {
	return patchAll(e, d)
}

func (v voidNode) patch(pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(v, pathBehind[0])
	}
	if !v.Equals(oldValue) {
		return nil, fmt.Errorf(
			"Found %v at %v. Expected %v.",
			v.Json(), pathBehind, oldValue.Json())
	}
	return newValue, nil
}

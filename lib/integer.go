package jd

import (
	"bytes"
	"encoding/binary"
)

type jsonInt int64

var _ JsonNode = jsonInt(0)

func (n jsonInt) Json() string {
	return renderJson(n)
}

func (n1 jsonInt) Equals(node JsonNode) bool {
	n2, ok := node.(jsonInt)
	if !ok {
		return false
	}
	if n1 != n2 {
		return false
	}
	return true
}

func (n jsonInt) hashCode() [8]byte {
	a := make([]byte, 0, 8)
	b := bytes.NewBuffer(a)
	binary.Write(b, binary.LittleEndian, n)
	return hash(b.Bytes())
}

func (n jsonInt) Diff(node JsonNode) Diff {
	return n.diff(node, Path{})
}

func (n jsonInt) diff(node JsonNode, path Path) Diff {
	d := make(Diff, 0)
	if n.Equals(node) {
		return d
	}
	e := DiffElement{
		Path:      path.clone(),
		OldValues: nodeList(n),
		NewValues: nodeList(node),
	}
	return append(d, e)
}

func (n jsonInt) Patch(d Diff) (JsonNode, error) {
	return patchAll(n, d)
}

func (n jsonInt) patch(pathBehind, pathAhead Path, oldValues, newValues []JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(n, pathAhead[0])
	}
	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)
	if !n.Equals(oldValue) {
		return patchErrExpectValue(oldValue, n, pathBehind)
	}
	return newValue, nil
}

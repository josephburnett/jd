package jd

type voidNode struct{}

var _ JsonNode = voidNode{}

func isVoid(n JsonNode) bool {
	if n == nil {
		return false
	}
	if _, ok := n.(voidNode); ok {
		return true
	}
	return false
}

func (v voidNode) Json() string {
	return ""
}

func (v voidNode) Equals(n JsonNode) bool {
	switch n.(type) {
	case voidNode:
		return true
	default:
		return false
	}
}

func (v voidNode) hashCode() [8]byte {
	return hash([]byte{0xF3, 0x97, 0x6B, 0x21, 0x91, 0x26, 0x8D, 0x96}) // Random bytes
}

func (v voidNode) Diff(n JsonNode) Diff {
	return v.diff(n, Path{})
}

func (v voidNode) diff(n JsonNode, p Path) Diff {
	d := make(Diff, 0)
	if v.Equals(n) {
		return d
	}
	de := DiffElement{
		Path:      p,
		OldValues: nodeList(v),
		NewValues: nodeList(n),
	}
	return append(d, de)
}

func (v voidNode) Patch(d Diff) (JsonNode, error) {
	return patchAll(v, d)
}

func (v voidNode) patch(pathBehind, pathAhead Path, oldValues, newValues []JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(v, pathBehind[0])
	}
	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)
	if !v.Equals(oldValue) {
		return patchErrExpectValue(oldValue, v, pathBehind)
	}
	return newValue, nil
}

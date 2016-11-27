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

func (v voidNode) Diff(n JsonNode) Diff {
	return v.diff(n, Path{})
}

func (v voidNode) diff(n JsonNode, p Path) Diff {
	d := make(Diff, 0)
	if v.Equals(n) {
		return d
	}
	de := DiffElement{
		Path:     p,
		OldValue: v,
		NewValue: n,
	}
	return append(d, de)
}

func (v voidNode) Patch(d Diff) (JsonNode, error) {
	return patchAll(v, d)
}

func (v voidNode) patch(pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(v, pathBehind[0])
	}
	if !v.Equals(oldValue) {
		return patchErrExpectValue(oldValue, v, pathBehind)
	}
	return newValue, nil
}

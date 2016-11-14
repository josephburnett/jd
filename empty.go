package jd

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

func (e voidNode) Diff(n JsonNode) Diff {
	return e.diff(n, Path{})
}

func (e voidNode) diff(n JsonNode, p Path) Diff {
	de := DiffElement{
		Path:     p,
		OldValue: e,
		NewValue: n,
	}
	return Diff{de}
}

func (e voidNode) Patch(d Diff) (JsonNode, error) {
	return patch(e, d)
}

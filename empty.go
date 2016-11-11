package jd

type emptyNode struct{}

var _ JsonNode = emptyNode{}

func (e emptyNode) Json() string {
	return ""
}

func (e emptyNode) Equals(n JsonNode) bool {
	switch n.(type) {
	case emptyNode:
		return true
	default:
		return false
	}
}

func (e emptyNode) Diff(n JsonNode) Diff {
	return e.diff(n, Path{})
}

func (e emptyNode) diff(n JsonNode, p Path) Diff {
	de := DiffElement{
		Path:     p,
		OldValue: e,
		NewValue: n,
	}
	return Diff{de}
}

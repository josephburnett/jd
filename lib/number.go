package jd

type jsonNumber float64

var _ JsonNode = jsonNumber(0)

func (n jsonNumber) Json() string {
	return renderJson(n)
}

func (n1 jsonNumber) Equals(node JsonNode) bool {
	n2, ok := node.(jsonNumber)
	if !ok {
		return false
	}
	if n1 != n2 {
		return false
	}
	return true
}

func (n jsonNumber) Diff(node JsonNode) Diff {
	return n.diff(node, Path{})
}

func (n jsonNumber) diff(node JsonNode, path Path) Diff {
	d := make(Diff, 0)
	if n.Equals(node) {
		return d
	}
	e := DiffElement{
		Path:     path.clone(),
		OldValue: n,
		NewValue: node,
	}
	return append(d, e)
}

func (n jsonNumber) Patch(d Diff) (JsonNode, error) {
	return patchAll(n, d)
}

func (n jsonNumber) patch(pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(n, pathAhead[0])
	}
	if !n.Equals(oldValue) {
		return patchErrExpectValue(oldValue, n, pathBehind)
	}
	return newValue, nil
}

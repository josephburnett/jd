package jd

type jsonNumber float64

var _ JsonNode = jsonNumber(0)

func (n jsonNumber) Json() string {
	return renderJson(n)
}

func (n1 jsonNumber) Equals(n JsonNode) bool {
	n2, ok := n.(jsonNumber)
	if !ok {
		return false
	}
	if n1 != n2 {
		return false
	}
	return true
}

func (n1 jsonNumber) Diff(n JsonNode) Diff {
	return n1.diff(n, Path{})
}

func (n1 jsonNumber) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	if n1.Equals(n) {
		return d
	}
	e := DiffElement{
		Path:     path.clone(),
		OldValue: n1,
		NewValue: n,
	}
	return append(d, e)
}

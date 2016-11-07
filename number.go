package jd

type JsonNumber float64

func (n1 JsonNumber) equals(n JsonNode) bool {
	n2, ok := n.(JsonNumber)
	if !ok {
		return false
	}
	if n1 != n2 {
		return false
	}
	return true
}

func (n1 JsonNumber) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	if n1.equals(n) {
		return d
	}
	e := DiffElement{
		Path:     path.clone(),
		OldValue: n1,
		NewValue: n,
	}
	return append(d, e)
}

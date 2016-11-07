package jd

type JsonString string

func (s1 JsonString) equals(n JsonNode) bool {
	s2, ok := n.(JsonString)
	if !ok {
		return false
	}
	return s1 == s2
}

func (s1 JsonString) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	if s1.equals(n) {
		return d
	}
	e := DiffElement{
		Path:     path.clone(),
		OldValue: s1,
		NewValue: n,
	}
	return append(d, e)
}

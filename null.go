package jd

type jsonNull struct{}

var _ JsonNode = jsonNull{}

func (n jsonNull) Json() string {
	return renderJson(nil)
}

func (n jsonNull) Equals(o JsonNode) bool {
	switch o.(type) {
	case jsonNull:
		return true
	default:
		return false
	}
}

func (n jsonNull) Diff(o JsonNode) Diff {
	return n.diff(o, Path{})
}

func (n jsonNull) diff(o JsonNode, path Path) Diff {
	d := make(Diff, 0)
	if n.Equals(o) {
		return d
	}
	e := DiffElement{
		Path:     path.clone(),
		OldValue: n,
		NewValue: o,
	}
	return append(d, e)
}

func (n jsonNull) Patch(d Diff) (JsonNode, error) {
	return patch(n, d)
}

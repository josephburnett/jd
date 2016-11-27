package jd

type jsonString string

var _ JsonNode = jsonString("")

func (s jsonString) Json() string {
	return renderJson(s)
}

func (s1 jsonString) Equals(n JsonNode) bool {
	s2, ok := n.(jsonString)
	if !ok {
		return false
	}
	return s1 == s2
}

func (s jsonString) hashCode() [8]byte {
	return hash([]byte(s))
}

func (s jsonString) Diff(n JsonNode) Diff {
	return s.diff(n, Path{})
}

func (s1 jsonString) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	if s1.Equals(n) {
		return d
	}
	e := DiffElement{
		Path:     path.clone(),
		OldValue: s1,
		NewValue: n,
	}
	return append(d, e)
}

func (s jsonString) Patch(d Diff) (JsonNode, error) {
	return patchAll(s, d)
}

func (s jsonString) patch(pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(s, pathBehind[0])
	}
	if !s.Equals(oldValue) {
		return patchErrExpectValue(oldValue, s, pathBehind)
	}
	return newValue, nil
}

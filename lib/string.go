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
		Path:      path.clone(),
		OldValues: nodeList(s1),
		NewValues: nodeList(n),
	}
	return append(d, e)
}

func (s jsonString) Patch(d Diff) (JsonNode, error) {
	return patchAll(s, d)
}

func (s jsonString) patch(pathBehind, pathAhead Path, oldValues, newValues []JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(s, pathBehind[0])
	}
	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)
	if !s.Equals(oldValue) {
		return patchErrExpectValue(oldValue, s, pathBehind)
	}
	return newValue, nil
}

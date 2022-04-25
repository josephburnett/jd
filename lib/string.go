package jd

type jsonString string

var _ JsonNode = jsonString("")

func (s jsonString) Json(metadata ...Metadata) string {
	return renderJson(s.raw(metadata))
}

func (s jsonString) Yaml(metadata ...Metadata) string {
	return renderYaml(s.raw(metadata))
}

func (s jsonString) raw(_ []Metadata) interface{} {
	return string(s)
}

func (s1 jsonString) Equals(n JsonNode, metadata ...Metadata) bool {
	s2, ok := n.(jsonString)
	if !ok {
		return false
	}
	return s1 == s2
}

func (s jsonString) hashCode(metadata []Metadata) [8]byte {
	return hash([]byte(s))
}

func (s jsonString) Diff(n JsonNode, metadata ...Metadata) Diff {
	return s.diff(n, make(path, 0), metadata)
}

func (s1 jsonString) diff(n JsonNode, path path, metadata []Metadata) Diff {
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

func (s jsonString) patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode, strategy patchStrategy) (JsonNode, error) {
	if !pathAhead.isLeaf() {
		return patchErrExpectColl(s, pathBehind[0])
	}
	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)
	switch strategy {
	case mergePatchStrategy:
		if !isVoid(oldValue) {
			return patchErrMergeWithOldValue(pathBehind, oldValue)
		}
		if isNull(newValue) {
			// Null deletes a node
			return voidNode{}, nil
		}
	case strictPatchStrategy:
		if !s.Equals(oldValue) {
			return patchErrExpectValue(oldValue, s, pathBehind)
		}
	default:
		return patchErrUnsupportedPatchStrategy(pathBehind, strategy)
	}
	return newValue, nil
}

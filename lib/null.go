package jd

type jsonNull []byte

var _ JsonNode = jsonNull{}

func (n jsonNull) Json(metadata ...Metadata) string {
	return renderJson(n)
}

func (n jsonNull) Equals(node JsonNode, metadata ...Metadata) bool {
	switch node.(type) {
	case jsonNull:
		return true
	default:
		return false
	}
}

func (n jsonNull) hashCode(metadata []Metadata) [8]byte {
	return hash([]byte{0xFE, 0x73, 0xAB, 0xCC, 0xE6, 0x32, 0xE0, 0x88}) // random bytes
}

func (n jsonNull) Diff(node JsonNode, metadata ...Metadata) Diff {
	return n.diff(node, make(path, 0), metadata)
}

func (n jsonNull) diff(node JsonNode, path path, metadata []Metadata) Diff {
	d := make(Diff, 0)
	if n.Equals(node) {
		return d
	}
	e := DiffElement{
		Path:      path.clone(),
		OldValues: nodeList(n),
		NewValues: nodeList(node),
	}
	return append(d, e)
}

func (n jsonNull) Patch(d Diff) (JsonNode, error) {
	return patchAll(n, d)
}

func (n jsonNull) patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(n, pathAhead[0])
	}
	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)
	if !n.Equals(oldValue) {
		return patchErrExpectValue(oldValue, n, pathBehind)
	}
	return newValue, nil
}

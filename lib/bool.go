package jd

type jsonBool bool

var _ JsonNode = jsonBool(true)

func (b jsonBool) Json(metadata ...Metadata) string {
	return renderJson(b)
}

func (b1 jsonBool) Equals(n JsonNode, metadata ...Metadata) bool {
	b2, ok := n.(jsonBool)
	if !ok {
		return false
	}
	return b1 == b2
}

func (b jsonBool) hashCode(metadata []Metadata) [8]byte {
	if b {
		return [8]byte{0x24, 0x6B, 0xE3, 0xE4, 0xAF, 0x59, 0xDC, 0x1C} // Random bytes
	} else {
		return [8]byte{0xC6, 0x38, 0x77, 0xD1, 0x0A, 0x7E, 0x1F, 0xBF} // Random bytes
	}
}

func (b jsonBool) Diff(n JsonNode, metadata ...Metadata) Diff {
	return b.diff(n, make(path, 0), metadata)
}

func (b jsonBool) diff(n JsonNode, path path, metadata []Metadata) Diff {
	d := make(Diff, 0)
	if b.Equals(n) {
		return d
	}
	e := DiffElement{
		Path:      path.clone(),
		OldValues: nodeList(b),
		NewValues: nodeList(n),
	}
	return append(d, e)
}

func (b jsonBool) Patch(d Diff) (JsonNode, error) {
	return patchAll(b, d)
}

func (b jsonBool) patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode) (JsonNode, error) {
	if len(pathAhead) != 0 {
		return patchErrExpectColl(b, pathAhead[0])
	}
	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)
	if !b.Equals(oldValue) {
		return patchErrExpectValue(oldValue, b, pathBehind)
	}
	return newValue, nil
}

package jd

type voidNode struct{}

var _ JsonNode = voidNode{}

func isVoid(n JsonNode) bool {
	if n == nil {
		return false
	}
	if _, ok := n.(voidNode); ok {
		return true
	}
	return false
}

func isNull(n JsonNode) bool {
	if n == nil {
		return false
	}
	if _, ok := n.(jsonNull); ok {
		return true
	}
	return false
}

func (v voidNode) Json(metadata ...Metadata) string {
	return ""
}

func (v voidNode) Yaml(metadata ...Metadata) string {
	return ""
}

func (v voidNode) raw(metadata []Metadata) interface{} {
	return ""
}

func (v voidNode) Equals(n JsonNode, metadata ...Metadata) bool {
	switch n.(type) {
	case voidNode:
		return true
	default:
		return false
	}
}

func (v voidNode) hashCode(metadata []Metadata) [8]byte {
	return hash([]byte{0xF3, 0x97, 0x6B, 0x21, 0x91, 0x26, 0x8D, 0x96}) // Random bytes
}

func (v voidNode) Diff(n JsonNode, metadata ...Metadata) Diff {
	return v.diff(n, make(path, 0), metadata, getPatchStrategy(metadata))
}

func (v voidNode) diff(n JsonNode, p path, metadata []Metadata, strategy patchStrategy) Diff {
	d := make(Diff, 0)
	if v.Equals(n) {
		return d
	}
	var de DiffElement
	switch strategy {
	case mergePatchStrategy:
		de = DiffElement{
			Path:      p.prependMetadataMerge(),
			NewValues: nodeList(n),
		}
	default:
		de = DiffElement{
			Path:      p.clone(),
			OldValues: nodeList(v),
			NewValues: nodeList(n),
		}
	}
	return append(d, de)
}

func (v voidNode) Patch(d Diff) (JsonNode, error) {
	return patchAll(v, d)
}

func (v voidNode) patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode, strategy patchStrategy) (JsonNode, error) {
	if !pathAhead.isLeaf() {
		return patchErrExpectColl(v, pathBehind[len(pathBehind)-1])
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
		if !v.Equals(oldValue) {
			return patchErrExpectValue(oldValue, v, pathBehind)
		}
	default:
		return patchErrUnsupportedPatchStrategy(pathBehind, strategy)
	}
	return newValue, nil
}

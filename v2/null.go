package jd

type jsonNull []byte

var _ JsonNode = jsonNull{}

func (n jsonNull) Json(_ ...Option) string {
	return renderJson(n.raw())
}

func (n jsonNull) Yaml(_ ...Option) string {
	return renderJson(n.raw())
}

func (n jsonNull) raw() interface{} {
	return nil
}

func (n jsonNull) Equals(node JsonNode, options ...Option) bool {
	switch node.(type) {
	case jsonNull:
		return true
	default:
		return false
	}
}

func (n jsonNull) hashCode(_ []Option) [8]byte {
	return hash([]byte{0xFE, 0x73, 0xAB, 0xCC, 0xE6, 0x32, 0xE0, 0x88}) // random bytes
}

func (n jsonNull) Diff(node JsonNode, options ...Option) Diff {
	return n.diff(node, make(path, 0), metadata, getPatchStrategy(metadata))
}

func (n jsonNull) diff(
	node JsonNode,
	path Path,
	options []Option,
	strategy patchStrategy,
) Diff {
	return diff(n, node, path, metadata, strategy)
}

func (n jsonNull) Patch(d Diff) (JsonNode, error) {
	return patchAll(n, d)
}

func (n jsonNull) patch(
	pathBehind, pathAhead Path,
	oldValues, newValues []JsonNode,
	strategy patchStrategy,
) (JsonNode, error) {
	return patch(n, pathBehind, pathAhead, oldValues, newValues, strategy)
}

package jd

type jsonBool bool

var _ JsonNode = jsonBool(true)

func (b jsonBool) Json(_ ...Option) string {
	return renderJson(b.raw())
}

func (b jsonBool) Yaml(_ ...Option) string {
	return renderYaml(b.raw())
}

func (b jsonBool) raw() interface{} {
	return bool(b)
}

func (b1 jsonBool) Equals(n JsonNode, options ...Option) bool {
	b2, ok := n.(jsonBool)
	if !ok {
		return false
	}
	return b1 == b2
}

func (b jsonBool) hashCode(_ []Option) [8]byte {
	if b {
		return [8]byte{0x24, 0x6B, 0xE3, 0xE4, 0xAF, 0x59, 0xDC, 0x1C} // Randomly chosen bytes
	} else {
		return [8]byte{0xC6, 0x38, 0x77, 0xD1, 0x0A, 0x7E, 0x1F, 0xBF} // Randomly chosen  bytes
	}
}

func (b jsonBool) Diff(n JsonNode, options ...Option) Diff {
	strategy := getPatchStrategy(options)
	return b.diff(n, make(Path, 0), options, strategy)
}

func (b jsonBool) diff(
	n JsonNode,
	path Path,
	options []Option,
	strategy patchStrategy,
) Diff {
	return diff(b, n, path, options, strategy)
}

func (b jsonBool) Patch(d Diff) (JsonNode, error) {
	return patchAll(b, d)
}

func (b jsonBool) patch(
	pathBehind, pathAhead Path,
	before, oldValues, newValues, after []JsonNode,
	strategy patchStrategy,
) (JsonNode, error) {
	return patch(b, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

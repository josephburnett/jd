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

func (v voidNode) Json(_ ...Option) string {
	return ""
}

func (v voidNode) Yaml(_ ...Option) string {
	return ""
}

func (v voidNode) raw() interface{} {
	return ""
}

func (v voidNode) Equals(n JsonNode, opts ...Option) bool {
	o := refine(&options{retain: opts}, nil)
	return v.equals(n, o)
}

func (v voidNode) equals(n JsonNode, o *options) bool {
	switch n.(type) {
	case voidNode:
		return true
	default:
		return false
	}
}

func (v voidNode) hashCode(_ *options) [8]byte {
	return hash([]byte{0xF3, 0x97, 0x6B, 0x21, 0x91, 0x26, 0x8D, 0x96}) // Randomly chosen bytes
}

func (v voidNode) Diff(n JsonNode, opts ...Option) Diff {
	o := refine(&options{retain: opts}, nil)
	return v.diff(n, make(Path, 0), o, getPatchStrategy(o))
}

func (v voidNode) diff(
	n JsonNode,
	p Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	return diff(v, n, p, opts, strategy)
}

func (v voidNode) Patch(d Diff) (JsonNode, error) {
	return patchAll(v, d)
}

func (v voidNode) patch(
	pathBehind, pathAhead Path,
	before, oldValues, newValues, after []JsonNode,
	strategy patchStrategy,
) (JsonNode, error) {
	return patch(v, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

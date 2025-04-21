package jd

type jsonString string

var _ JsonNode = jsonString("")

func (s jsonString) Json(_ ...Option) string {
	return renderJson(s.raw())
}

func (s jsonString) Yaml(_ ...Option) string {
	return renderYaml(s.raw())
}

func (s jsonString) raw() interface{} {
	return string(s)
}

func (s1 jsonString) Equals(n JsonNode, opts ...Option) bool {
	o := refine(&options{retain: opts}, nil)
	return s1.equals(n, o)
}

func (s1 jsonString) equals(n JsonNode, o *options) bool {
	s2, ok := n.(jsonString)
	if !ok {
		return false
	}
	return s1 == s2
}

func (s jsonString) hashCode(_ *options) [8]byte {
	return hash([]byte(s))
}

func (s jsonString) Diff(n JsonNode, opts ...Option) Diff {
	o := refine(&options{retain: opts}, nil)
	return s.diff(n, make(Path, 0), o, getPatchStrategy(o))
}

func (s1 jsonString) diff(
	n JsonNode,
	path Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	return diff(s1, n, path, opts, strategy)
}

func (s jsonString) Patch(d Diff) (JsonNode, error) {
	return patchAll(s, d)
}

func (s jsonString) patch(
	pathBehind, pathAhead Path,
	before, oldValues, newValues, after []JsonNode,
	strategy patchStrategy,
) (JsonNode, error) {
	return patch(s, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

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

func (s1 jsonString) Equals(n JsonNode, options ...Option) bool {
	s2, ok := n.(jsonString)
	if !ok {
		return false
	}
	return s1 == s2
}

func (s jsonString) hashCode(_ []Option) [8]byte {
	return hash([]byte(s))
}

func (s jsonString) Diff(n JsonNode, options ...Option) Diff {
	return s.diff(n, make(Path, 0), options, getPatchStrategy(options))
}

func (s1 jsonString) diff(
	n JsonNode,
	path Path,
	options []Option,
	strategy patchStrategy,
) Diff {
	return diff(s1, n, path, options, strategy)
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

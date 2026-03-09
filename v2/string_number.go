package jd

type jsonStringNumber string

var _ JsonNode = jsonStringNumber("")

func (s jsonStringNumber) Json(_ ...Option) string {
	return string(s)
}

func (s jsonStringNumber) Yaml(_ ...Option) string {
	return string(s)
}

func (s jsonStringNumber) raw() interface{} {
	return string(s)
}

func (s1 jsonStringNumber) Equals(n JsonNode, opts ...Option) bool {
	o := refine(&options{retain: opts}, nil)
	return s1.equals(n, o)
}

func (s1 jsonStringNumber) equals(n JsonNode, o *options) bool {
	s2, ok := n.(jsonStringNumber)
	if !ok {
		return false
	}
	return s1 == s2
}

func (s jsonStringNumber) hashCode(_ *options) [8]byte {
	return hash([]byte(s))
}

func (s jsonStringNumber) Diff(n JsonNode, opts ...Option) Diff {
	o := refine(newOptions(opts), nil)
	return s.diff(n, make(Path, 0), o, getPatchStrategy(o))
}

func (s1 jsonStringNumber) diff(
	n JsonNode,
	path Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	// Use event-driven diff architecture
	events := generateSimpleEvents(s1, n, opts)
	processor := newSimpleDiffProcessor(path, opts, strategy)
	return processor.ProcessEvents(events)
}

func (s jsonStringNumber) Patch(d Diff) (JsonNode, error) {
	return patchAll(s, d)
}

func (s jsonStringNumber) patch(
	pathBehind, pathAhead Path,
	before, oldValues, newValues, after []JsonNode,
	strategy patchStrategy,
) (JsonNode, error) {
	return patch(s, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

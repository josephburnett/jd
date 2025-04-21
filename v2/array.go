package jd

// jsonArray is a polymorphic type representing a concrete JSON array. It
// dispatches to list, set or multiset semantics.
type jsonArray []JsonNode

var _ JsonNode = jsonArray(nil)

func (a jsonArray) Json(renderOptions ...Option) string {
	o := refine(&options{retain: renderOptions}, nil)
	n := dispatch(a, o)
	return n.Json(renderOptions...)
}

func (a jsonArray) Yaml(renderOptions ...Option) string {
	o := refine(&options{retain: renderOptions}, nil)
	n := dispatch(a, o)
	return n.Yaml(renderOptions...)
}

func (a jsonArray) raw() interface{} {
	r := make([]interface{}, len(a))
	for i, n := range a {
		r[i] = n.raw()
	}
	return r
}

func (a1 jsonArray) Equals(n JsonNode, opts ...Option) bool {
	o := refine(&options{retain: opts}, nil)
	return a1.equals(n, o)
}

func (a1 jsonArray) equals(n JsonNode, o *options) bool {
	n1 := dispatch(a1, o)
	n2 := dispatch(n, o)
	return n1.equals(n2, o)
}

func (a jsonArray) hashCode(opts *options) [8]byte {
	n := dispatch(a, opts)
	return n.hashCode(opts)
}

func (a jsonArray) Diff(n JsonNode, opts ...Option) Diff {
	o := refine(&options{retain: opts}, nil)
	n1 := dispatch(a, o)
	n2 := dispatch(n, o)
	strategy := getPatchStrategy(o)
	return n1.diff(n2, make(Path, 0), o, strategy)
}

func (a jsonArray) diff(
	n JsonNode,
	path Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	n1 := dispatch(a, opts)
	n2 := dispatch(n, opts)
	return n1.diff(n2, path, opts, strategy)
}

func (a jsonArray) Patch(d Diff) (JsonNode, error) {
	return patchAll(a, d)
}

func (a jsonArray) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {
	_, metadata, _ := pathAhead.next()
	o := refine(&options{retain: metadata}, nil)
	n := dispatch(a, o)
	return n.patch(pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

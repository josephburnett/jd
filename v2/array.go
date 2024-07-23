package jd

// jsonArray is a polymorphic type representing a concrete JSON array. It
// dispatches to list, set or multiset semantics.
type jsonArray []JsonNode

var _ JsonNode = jsonArray(nil)

func (a jsonArray) Json(renderOptions ...Option) string {
	n := dispatch(a, renderOptions)
	return n.Json(renderOptions...)
}

func (a jsonArray) Yaml(renderOptions ...Option) string {
	n := dispatch(a, renderOptions)
	return n.Yaml(renderOptions...)
}

func (a jsonArray) raw() interface{} {
	r := make([]interface{}, len(a))
	for i, n := range a {
		r[i] = n.raw()
	}
	return r
}

func (a1 jsonArray) Equals(n JsonNode, options ...Option) bool {
	n1 := dispatch(a1, options)
	n2 := dispatch(n, options)
	return n1.Equals(n2, options...)
}

func (a jsonArray) hashCode(options []Option) [8]byte {
	n := dispatch(a, options)
	return n.hashCode(options)
}

func (a jsonArray) Diff(n JsonNode, options ...Option) Diff {
	n1 := dispatch(a, options)
	n2 := dispatch(n, options)
	strategy := getPatchStrategy(options)
	return n1.diff(n2, make(Path, 0), options, strategy)
}

func (a jsonArray) diff(
	n JsonNode,
	path Path,
	options []Option,
	strategy patchStrategy,
) Diff {
	n1 := dispatch(a, options)
	n2 := dispatch(n, options)
	return n1.diff(n2, path, options, strategy)
}

func (a jsonArray) Patch(d Diff) (JsonNode, error) {
	return patchAll(a, d)
}

func (a jsonArray) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {
	_, metadata, _ := pathAhead.next()
	n := dispatch(a, metadata)
	return n.patch(pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

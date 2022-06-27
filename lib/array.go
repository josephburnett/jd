package jd

// jsonArray is a polymorphic type representing a concrete JSON array. It
// dispatches to list, set or multiset semantics.
type jsonArray []JsonNode

var _ JsonNode = jsonArray(nil)

func (a jsonArray) Json(opts ...RenderOption) string {
	n := dispatchRenderOptions(a, opts)
	return n.Json(opts...)
}

func (a jsonArray) Yaml(opts ...RenderOption) string {
	n := dispatchRenderOptions(a, opts)
	return n.Yaml(opts...)
}

func (a jsonArray) raw() interface{} {
	r := make([]interface{}, len(a))
	for i, n := range a {
		r[i] = n.raw()
	}
	return r
}

func (a1 jsonArray) Equals(n JsonNode, metadata ...Metadata) bool {
	n1 := dispatch(a1, metadata)
	n2 := dispatch(n, metadata)
	return n1.Equals(n2, metadata...)
}

func (a jsonArray) hashCode(metadata []Metadata) [8]byte {
	n := dispatch(a, metadata)
	return n.hashCode(metadata)
}

func (a jsonArray) Diff(n JsonNode, metadata ...Metadata) Diff {
	n1 := dispatch(a, metadata)
	n2 := dispatch(n, metadata)
	strategy := getPatchStrategy(metadata)
	return n1.diff(n2, make(path, 0), metadata, strategy)
}

func (a jsonArray) diff(n JsonNode, path path, metadata []Metadata, strategy patchStrategy) Diff {
	n1 := dispatch(a, metadata)
	n2 := dispatch(n, metadata)
	return n1.diff(n2, path, metadata, strategy)
}

func (a jsonArray) Patch(d Diff) (JsonNode, error) {
	return patchAll(a, d)
}

func (a jsonArray) patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode, strategy patchStrategy) (JsonNode, error) {
	_, metadata, _ := pathAhead.next()
	n := dispatch(a, metadata)
	return n.patch(pathBehind, pathAhead, oldValues, newValues, strategy)
}

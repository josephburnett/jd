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
	// We need to refine to extract global options (SET, MULTISET, etc.) for dispatch,
	// but we want to preserve PathOptions. So we do a selective refine.
	o := a.refineForArrayDispatch(&options{retain: opts})
	n1 := dispatch(a, o)
	n2 := dispatch(n, o)
	strategy := getPatchStrategy(o)
	return n1.diff(n2, make(Path, 0), o, strategy)
}

// refineForArrayDispatch extracts global options for dispatch while preserving PathOptions
func (a jsonArray) refineForArrayDispatch(opts *options) *options {
	var apply, retain []Option

	for _, opt := range opts.retain {
		switch o := opt.(type) {
		// Global options - extract to apply for dispatch to work
		case mergeOption, setOption, multisetOption, colorOption, precisionOption, setKeysOption:
			apply = append(apply, o)
			retain = append(retain, o)
		case pathOption:
			// Special case: PathOption with empty path should have its options applied for dispatch
			if len(o.At) == 0 {
				// Extract dispatch-relevant options from the PathOption
				for _, thenOpt := range o.Then {
					switch thenOpt.(type) {
					case setOption, multisetOption:
						apply = append(apply, thenOpt)
					}
				}
			}
			// Always keep PathOptions in retain to preserve them for inner diffing
			retain = append(retain, o)
		}
	}

	return &options{
		apply:  apply,
		retain: retain,
	}
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

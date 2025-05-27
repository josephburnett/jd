package array

import (
	"github.com/josephburnett/jd/v2/internal/types"
)

func dispatch(n types.JsonNode, opts *types.Options) types.JsonNode {
	switch n := n.(type) {
	case JsonArray:
		for _, o := range opts.Apply {
			switch o.(type) {
			case types.SetOption, types.SetKeysOption:
				return JsonSet(n)
			case types.MultisetOption:
				return JsonMultiset(n)
			}
		}
		return JsonList(n)
	}
	return n
}

// JsonArray is a polymorphic type representing a concrete JSON array. It
// dispatches to list, set or multiset semantics.
type JsonArray []types.JsonNode

var _ types.JsonNode = JsonArray(nil)

func (a JsonArray) Json(renderOptions ...types.Option) string {
	o := types.Refine(&types.Options{Retain: renderOptions}, nil)
	n := dispatch(a, o)
	return n.Json(renderOptions...)
}

func (a JsonArray) Yaml(renderOptions ...types.Option) string {
	o := types.Refine(&types.Options{Retain: renderOptions}, nil)
	n := dispatch(a, o)
	return n.Yaml(renderOptions...)
}

func (a JsonArray) raw() interface{} {
	r := make([]interface{}, len(a))
	for i, n := range a {
		r[i] = n.raw()
	}
	return r
}

func (a1 JsonArray) Equals(n types.JsonNode, opts ...types.Option) bool {
	o := types.Refine(&types.Options{Retain: opts}, nil)
	return a1.equals(n, o)
}

func (a1 JsonArray) equals(n types.JsonNode, o *types.Options) bool {
	n1 := dispatch(a1, o)
	n2 := dispatch(n, o)
	return n1.equals(n2, o)
}

func (a JsonArray) hashCode(opts *types.Options) [8]byte {
	n := dispatch(a, opts)
	return n.hashCode(opts)
}

func (a JsonArray) Diff(n types.JsonNode, opts ...types.Option) Diff {
	o := types.Refine(&types.Options{Retain: opts}, nil)
	n1 := dispatch(a, o)
	n2 := dispatch(n, o)
	strategy := getPatchStrategy(o)
	return n1.diff(n2, make(Path, 0), o, strategy)
}

func (a JsonArray) diff(
	n types.JsonNode,
	path Path,
	opts *types.Options,
	strategy patchStrategy,
) Diff {
	n1 := dispatch(a, opts)
	n2 := dispatch(n, opts)
	return n1.diff(n2, path, opts, strategy)
}

func (a JsonArray) Patch(d Diff) (types.JsonNode, error) {
	return patchAll(a, d)
}

func (a JsonArray) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []types.JsonNode, strategy patchStrategy) (types.JsonNode, error) {
	_, metadata, _ := pathAhead.next()
	o := types.Refine(&types.Options{Retain: metadata}, nil)
	n := dispatch(a, o)
	return n.patch(pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

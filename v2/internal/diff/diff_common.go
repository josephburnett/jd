package jd

import "github.com/josephburnett/jd/v2/internal/types"

type JsonNode = types.JsonNode
type Path = types.Path 
type options = types.Options
type patchStrategy = types.PatchStrategy
type Diff = types.Diff
type DiffElement = types.DiffElement
type Metadata = types.Metadata

const (
	mergePatchStrategy = types.MergePatchStrategy
	strictPatchStrategy = types.StrictPatchStrategy
)

func diff(
	a, b JsonNode,
	p Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	d := make(Diff, 0)
	if a.equals(b, opts) {
		return d
	}
	var de DiffElement
	switch strategy {
	case mergePatchStrategy:
		de = DiffElement{
			Metadata: Metadata{
				Merge: true,
			},
			Path: p.clone(),
			Add:  jsonArray{b},
		}
	default:
		de = DiffElement{
			Path:   p.clone(),
			Remove: nodeList(a),
			Add:    nodeList(b),
		}
	}
	return append(d, de)
}

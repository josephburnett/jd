package jd

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

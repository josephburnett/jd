package jd

type path []JsonNode

func (p path) appendIndex(o jsonObject, metadata []Metadata) path {
	// Append metadata.
	meta := make(jsonArray, 0)
	if checkMetadata(SET, metadata) {
		meta = append(meta, jsonString(SET.string()))
	}
	if checkMetadata(MULTISET, metadata) {
		meta = append(meta, jsonString(MULTISET.string()))
	}
	sk := getSetkeysMetadata(metadata)
	if sk != nil {
		meta = append(meta, jsonString(sk.string()))
	}
	p = append(p, meta)
	// Append index.
	return append(p, o)
}

func (p path) withMetadata(metadata []Metadata) path {
	if checkMetadata(MERGE, metadata) {
		c := make(path, len(p)+1)
		c[0] = jsonArray{jsonString(MERGE.string())}
		for i, e := range p {
			c[i+1] = e
		}
		return c
	}
	c := make(path, len(p))
	copy(c, p)
	return c
}

func (p path) clone() path {
	c := make(path, len(p))
	copy(c, p)
	return c
}

func (p path) next() (JsonNode, []Metadata, path) {
	var metadata []Metadata
	for i, n := range p {
		switch n := n.(type) {
		case jsonArray:
			for _, meta := range n {
				// TODO: parse metadata cleanly.
				if s, ok := meta.(jsonString); ok {
					if string(s) == SET.string() {
						metadata = append(metadata, SET)
					}
					if string(s) == MULTISET.string() {
						metadata = append(metadata, MULTISET)
					}
					if string(s) == MERGE.string() {
						metadata = append(metadata, MERGE)
					}
				}
				// Ignore unrecognized metadata.
			}
		case jsonObject:
			// JSON object implies a set.
			if !checkMetadata(SET, metadata) && !checkMetadata(MULTISET, metadata) {
				metadata = append(metadata, SET)
			}
			return n, metadata, p[i+1:]
		default:
			return n, metadata, p[i+1:]
		}
	}
	return voidNode{}, metadata, nil
}

func (p path) getPatchStrategy() patchStrategy {
	var defaultStrategy = strictPatchStrategy
	if len(p) == 0 {
		return defaultStrategy
	}
	a, ok := p[0].(jsonArray)
	if !ok {
		return defaultStrategy
	}
	for _, n := range a {
		s, ok := n.(jsonString)
		if !ok {
			continue
		}
		if string(s) == MERGE.string() {
			return mergePatchStrategy
		}
	}
	return defaultStrategy
}

func (p path) isLeaf() bool {
	if len(p) == 0 {
		return true
	}
	if len(p) == 1 {
		if _, ok := p[0].(jsonArray); ok {
			// The only path element is metadata.
			return true
		}
	}
	return false
}

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

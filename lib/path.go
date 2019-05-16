package jd

type path []JsonNode

func (p path) appendSetIndex(o jsonObject, metadata []Metadata) path {
	// Append metadata.
	meta := make(jsonArray)
	meta = append(meta, jsonString(SET.string()))
	sk := getSetkeysMetadata(metadata)
	if sk != nil {
		meta = append(meta, jsonString(sk.string()))
	}
	p = append(p, meta)
	// Append index.
	return append(p, o)
}

func (p path) next() (JsonNode, []Metadata, path) {
	var metadata []Metadata
	for i, n := range path {
		switch n := n.(type) {
		case jsonArray:
			for _, meta := range n {
				// TODO: parse metadata cleanly.
				if s, ok := meta.(jsonString); ok {
					if s == SET.string() {
						metadata = append(metadata, SET)
					}
					if s == MULTISET.string() {
						metadata = append(metadata, MULTISET)
					}
				}
				// Ignore unrecognized metadata.
			}
		default:
			return n, metadata, path[i:]
		}
	}
	return jsonVoid(nil), metadata, nil
}

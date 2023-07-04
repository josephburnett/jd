package jd

type Metadata struct {
	Version int
	Merge   bool
}

func dispatch(n JsonNode, options []Option) JsonNode {
	switch n := n.(type) {
	case jsonArray:
		if metadata.Set {
			return jsonSet(n)
		}
		if metadata.Multiset {
			return jsonMultiset(n)
		}
		return jsonList(n)
	}
	return n
}

func checkMetadata(want Metadata, metadata []Metadata) bool {
	for _, o := range metadata {
		if o == want {
			return true
		}
	}
	return false
}

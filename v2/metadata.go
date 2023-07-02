package jd

// Metadata is a closed set of values which modify Diff and Equals
// semantics.
type PatchMetadata struct {
	Version int
	Merge   bool
}

type DiffMetadata struct {
	Version  int
	Merge    bool
	Set      bool
	Multiset bool
}

type patchStrategy string

const (
	mergePatchStrategy  patchStrategy = "merge"
	strictPatchStrategy patchStrategy = "strict"
)

func getPatchStrategy(metadata Metadata) patchStrategy {
	if metadata.Merge {
		return mergePatchStrategy
	}
	return strictPatchStrategy
}

func dispatch(n JsonNode, metadata Metadata) JsonNode {
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

func dispatchRenderOptions(n JsonNode, opts []RenderOption) JsonNode {
	metadata := []Metadata{}
	for _, o := range opts {
		if m, ok := o.(Metadata); ok {
			metadata = append(metadata, m)
		}
	}
	return dispatch(n, metadata)
}

func checkMetadata(want Metadata, metadata []Metadata) bool {
	for _, o := range metadata {
		if o == want {
			return true
		}
	}
	return false
}

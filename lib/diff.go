package jd

type DiffElement struct {
	Path      []JsonNode
	OldValues []JsonNode
	NewValues []JsonNode
}

type Diff []DiffElement

// JSON Patch (RFC 6902)
type patchElement struct {
	Op    string      `json:"op"`   // "add", "test" or "remove"
	Path  string      `json:"path"` // JSON Pointer (RFC 6901)
	Value interface{} `json:"value"`
}

package jd

type DiffElement struct {
	Path      []JsonNode
	OldValues []JsonNode
	NewValues []JsonNode
}

type Diff []DiffElement

package jd

type DiffElement struct {
	Path     Path
	OldValue JsonNode
	NewValue JsonNode
}
type Diff []DiffElement

package jd

import "github.com/josephburnett/jd/v2/internal/types"

// Type aliases for main package compatibility
type Metadata = types.Metadata
type DiffElement = types.DiffElement  
type Diff = types.Diff
type Path = types.Path
type PathElement = types.PathElement
type PathIndex = types.PathIndex
type PathKey = types.PathKey
type PathSet = types.PathSet
type PathMultiset = types.PathMultiset
type PathSetKeys = types.PathSetKeys
type PathMultisetKeys = types.PathMultisetKeys
type PathAllKeys = types.PathAllKeys
type PathAllValues = types.PathAllValues

// Note: DiffElement, Diff and related types are now aliases to internal/types

// JSON Patch (RFC 6902)
type patchElement struct {
	Op    string      `json:"op"`   // "add", "test" or "remove"
	Path  string      `json:"path"` // JSON Pointer (RFC 6901)
	Value interface{} `json:"value"`
}

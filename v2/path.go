package jd

import "github.com/josephburnett/jd/v2/internal/types"

// Note: Path and PathElement types are now aliases to internal/types (defined in diff.go)

func NewPath(n JsonNode) (Path, error) {
	return types.NewPath(n)
}
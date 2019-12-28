package jd

import (
	"fmt"
)

func patchAll(n JsonNode, d Diff) (JsonNode, error) {
	var err error
	for _, de := range d {
		n, err = n.patch(make(path, 0), de.Path, de.OldValues, de.NewValues)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}

func singleValue(nodes []JsonNode) JsonNode {
	if len(nodes) == 0 {
		return voidNode{}
	}
	return nodes[0]
}

func patchErrExpectColl(n JsonNode, pe interface{}) (JsonNode, error) {
	switch pe := pe.(type) {
	case string:
		return nil, fmt.Errorf(
			"Found %v at %v. Expected JSON object.",
			// TODO: plumb through metadata.
			n.Json(), pe)
	case float64:
		return nil, fmt.Errorf(
			"Found %v at %v. Expected JSON array.",
			n.Json(), pe)
	default:
		return nil, fmt.Errorf("Invalid path element %v.", pe)
	}

}

func patchErrNonSetDiff(oldValues, newValues []JsonNode, path path) (JsonNode, error) {
	if len(oldValues) > 1 {
		return nil, fmt.Errorf(
			"Invalid diff: Multiple removals from non-set at %v.",
			path)
	} else {
		return nil, fmt.Errorf(
			"Invalid diff: Multiple additions to a non-set at %v.",
			path)
	}
}

func patchErrExpectValue(want, found JsonNode, path path) (JsonNode, error) {
	return nil, fmt.Errorf(
		"Found %v at %v. Expected %v.",
		found.Json(), path, want.Json())
}

package jd

import (
	"fmt"
)

func patchAll(n JsonNode, d Diff) (JsonNode, error) {
	var err error
	for _, de := range d {
		n, err = n.patch(Path{}, de.Path, de.OldValue, de.NewValue)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}

func patchErrExpectColl(n JsonNode, pe interface{}) (JsonNode, error) {
	switch pe := pe.(type) {
	case string:
		return nil, fmt.Errorf(
			"Found %v at %v. Expected JSON object.",
			n.Json(), pe)
	case float64:
		return nil, fmt.Errorf(
			"Found %v at %v. Expected JSON array.",
			n.Json(), pe)
	default:
		panic(fmt.Sprintf("Invalid path element %v.", pe))
	}

}

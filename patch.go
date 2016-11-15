package jd

import (
	"fmt"
)

func patch(n JsonNode, d Diff) (JsonNode, error) {
	var err error
	for _, de := range d {
		n, err = patchInternal(n, Path{}, de.Path, de.OldValue, de.NewValue)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}

func patchInternal(n JsonNode, pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error) {
	if oldValue == nil && newValue == nil {
		return nil, fmt.Errorf(
			"Invalid diff element. No old or new value provided.")
	}
	// Base case
	if len(pathAhead) == 0 {
		if oldValue != nil && !n.Equals(oldValue) {
			return nil, fmt.Errorf(
				"Found %v at %v. Expected %v.",
				n.Json(), pathBehind, oldValue.Json())
		}
		return newValue, nil
	}
	// Recursive case
	switch pe := pathAhead[0].(type) {
	case string:
		// JSON object
		s, ok := n.(jsonObject)
		if !ok {
			return nil, fmt.Errorf(
				"Found %v at %v. Expected JSON object.",
				n.Json(), pathBehind)
		}
		nextNode, ok := s[pe]
		if !ok {
			nextNode = voidNode{}
		}
		patchedNode, err := patchInternal(nextNode, append(pathBehind, pe), pathAhead[1:], oldValue, newValue)
		if err != nil {
			return nil, err
		}
		if patchedNode == nil {
			// Deletion of a pair.
			delete(s, pe)
		} else {
			// Addition or replacement of a pair.
			s[pe] = patchedNode
		}
		return s, nil
	case float64:
		// JSON array
		s, ok := n.(jsonArray)
		if !ok {
			return nil, fmt.Errorf(
				"Found %v at %v. Expected JSON array.",
				n.Json(), pathBehind)
		}
		i := int(pe)
		var nextNode JsonNode = voidNode{}
		if len(s) > i {
			nextNode = s[i]
		}
		patchedNode, err := patchInternal(nextNode, append(pathBehind, pe), pathAhead[1:], oldValue, newValue)
		if err != nil {
			return nil, err
		}
		if patchedNode == nil {
			if i != len(s)-1 {
				return nil, fmt.Errorf(
					"Removal of non-terminal element of array.")
			}
			// Deletion of an element.
			s = s[:len(s)-1]
		} else {
			if i > len(s) {
				return nil, fmt.Errorf(
					"Addition of non-terminal element of array.")
			}
			if i == len(s) {
				// Addition of an element.
				s = append(s, patchedNode)
			} else {
				// Replacement of an element.
				if oldValue == nil {
					return nil, fmt.Errorf(
						"Overwrite of an unknown value.")
				}
				s[int(pe)] = patchedNode
			}
		}
		return s, nil
	default:
		panic(fmt.Sprintf("Invalid path element %v", pe))
	}
}

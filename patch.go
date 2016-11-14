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
			delete(s, pe)
		} else {
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
		s[int(pe)] = patchedNode
		return s, nil
	default:
		panic(fmt.Sprintf("Invalid path element %v", pe))
	}
}

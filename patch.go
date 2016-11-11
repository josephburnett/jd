package jd

import (
	"fmt"
)

func patch(n JsonNode, pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error) {
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
		s, ok := n.(jsonStruct)
		if !ok {
			return nil, fmt.Errorf(
				"Found %v at %v. Expected JSON struct.",
				n.Json(), pathBehind)
		}
		nextNode, ok := s[pe]
		if !ok {
			nextNode = emptyNode{}
		}
		patchedNode, err := patch(nextNode, append(pathBehind, pe), pathAhead[1:], oldValue, newValue)
		if err != nil {
			return nil, err
		}
		s[pe] = patchedNode
		return s, nil
	case float64:
		s, ok := n.(jsonList)
		if !ok {
			return nil, fmt.Errorf(
				"Found %v at %v. Expected JSON list.",
				n.Json(), pathBehind)
		}
		i := int(pe)
		var nextNode JsonNode = emptyNode{}
		if len(s) > i {
			nextNode = s[i]
		}
		patchedNode, err := patch(nextNode, append(pathBehind, pe), pathAhead[1:], oldValue, newValue)
		if err != nil {
			return nil, err
		}
		s[int(pe)] = patchedNode
		return s, nil
	default:
		panic(fmt.Sprintf("Invalid path element %v", pe))
	}
}

package jd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/jsonpointer"
)

func readPointer(s string) ([]JsonNode, error) {
	pointer, err := jsonpointer.New(s)
	if err != nil {
		return nil, err
	}
	tokens := pointer.DecodedTokens()
	path := make([]JsonNode, len(tokens))
	for i, t := range tokens {
		var element JsonNode
		var err error
		number, err := strconv.Atoi(t)
		if err == nil {
			element, err = NewJsonNode(number)
		} else {
			element, err = NewJsonNode(t)
		}
		if err != nil {
			return nil, err
		}
		path[i] = element
	}
	return path, nil
}

func writePointer(path []JsonNode) (string, error) {
	var b strings.Builder
	for _, element := range path {
		b.WriteString("/")
		switch e := element.(type) {
		case jsonNumber:
			b.WriteString(jsonpointer.Escape(strconv.Itoa(int(e))))
		case jsonString:
			b.WriteString(string(e))
		case jsonArray:
			return "", fmt.Errorf("JSON Pointer does not support jd metadata.")
		default:
			return "", fmt.Errorf("Unsupported type: %T", e)
		}
	}
	return b.String(), nil
}

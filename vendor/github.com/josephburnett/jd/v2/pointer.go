package jd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/jsonpointer"
)

func readPointer(s string) (Path, error) {
	pointer, err := jsonpointer.New(s)
	if err != nil {
		return nil, err
	}
	tokens := pointer.DecodedTokens()
	path := make(jsonArray, len(tokens))
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
		if s, ok := element.(jsonString); ok && s == "-" {
			element, _ = NewJsonNode(-1)
		}
		path[i] = element
	}
	return NewPath(path)
}

func writePointer(path []JsonNode) (string, error) {
	var b strings.Builder
	for _, element := range path {
		b.WriteString("/")
		switch e := element.(type) {
		case jsonNumber:
			if int(e) == -1 {
				b.WriteString("-")
			} else {
				b.WriteString(jsonpointer.Escape(strconv.Itoa(int(e))))
			}
		case jsonString:
			if _, err := strconv.Atoi(string(e)); err == nil {
				return "", fmt.Errorf("JSON Pointer does not support object keys that look like numbers: %v", e)
			}
			if string(e) == "-" {
				return "", fmt.Errorf("JSON Pointer does not support object key '-'")
			}
			s := jsonpointer.Escape(string(e))
			b.WriteString(s)
		case jsonArray:
			return "", fmt.Errorf("JSON Pointer does not support jd metadata")
		default:
			return "", fmt.Errorf("unsupported type: %T", e)
		}
	}
	return b.String(), nil
}

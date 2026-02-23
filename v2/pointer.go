package jd

import (
	"fmt"
	"strconv"
	"strings"
)

func readPointer(s string) (Path, error) {
	tokens, err := jsonPointerParse(s)
	if err != nil {
		return nil, err
	}
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
				b.WriteString(jsonPointerEscape(strconv.Itoa(int(e))))
			}
		case jsonString:
			if _, err := strconv.Atoi(string(e)); err == nil {
				return "", fmt.Errorf("JSON Pointer does not support object keys that look like numbers: %v", e)
			}
			if string(e) == "-" {
				return "", fmt.Errorf("JSON Pointer does not support object key '-'")
			}
			s := jsonPointerEscape(string(e))
			b.WriteString(s)
		case jsonObject:
			return "", fmt.Errorf("JSON Pointer does not support set-based paths. Use jd format instead of patch")
		case jsonArray:
			return "", fmt.Errorf("JSON Pointer does not support jd metadata")
		default:
			return "", fmt.Errorf("unsupported type: %T", e)
		}
	}
	return b.String(), nil
}

var jsonPointerUnescaper = strings.NewReplacer("~1", "/", "~0", "~")
var jsonPointerEscaper = strings.NewReplacer("~", "~0", "/", "~1")

// jsonPointerParse validates and parses a JSON Pointer (RFC 6901) string
// into decoded reference tokens.
func jsonPointerParse(s string) ([]string, error) {
	if s == "" {
		return nil, nil
	}
	if s[0] != '/' {
		return nil, fmt.Errorf("non-empty JSON pointer must start with '/'")
	}
	tokens := strings.Split(s[1:], "/")
	for i, t := range tokens {
		tokens[i] = jsonPointerUnescaper.Replace(t)
	}
	return tokens, nil
}

// jsonPointerEscape escapes a reference token per RFC 6901:
// '~' → '~0', '/' → '~1'.
func jsonPointerEscape(token string) string {
	return jsonPointerEscaper.Replace(token)
}

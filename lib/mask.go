package jd

import (
	"bytes"
	"fmt"
	"strings"
)

type MaskElement struct {
	Include bool
	Path    []JsonNode
}

type Mask []MaskElement

var _ Metadata = Mask{}

func (m Mask) is_metadata() {}

func (m Mask) string() string {
	return m.Render()
}

func (d Mask) Render() string {
	b := bytes.NewBuffer(nil)
	for _, e := range d {
		b.WriteString(e.Render())
		b.WriteString("\n")
	}
	return b.String()
}

func (e MaskElement) Render() string {
	b := bytes.NewBuffer(nil)
	if e.Include {
		b.WriteString("+ ")
	} else {
		b.WriteString("- ")
	}
	b.WriteString(jsonArray(e.Path).Json())
	return b.String()
}

func ReadMaskString(s string) (Mask, error) {
	mask := Mask{}
	maskLines := strings.Split(s, "\n")
	for _, line := range maskLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		element := MaskElement{}
		switch line[0] {
		case '+':
			element.Include = true
		case '-':
			element.Include = false
		default:
			return nil, fmt.Errorf("Unexpected %c. Expecting + or -.", line[0])
		}
		line = strings.TrimSpace(line[1:])
		if len(line) == 0 {
			return nil, fmt.Errorf("Unexpected end of line. Expecting JSON path")
		}
		path, err := ReadJsonString(line)
		if err != nil {
			return nil, err
		}
		array, ok := path.(jsonArray)
		if !ok {
			return nil, fmt.Errorf("Unexpected path. Expecting JSON array.")
		}
		element.Path = array
		mask = append(mask, element)
	}
	return mask, nil
}

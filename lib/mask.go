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

func getMask(metadata []Metadata) Mask {
	for _, m := range metadata {
		if n, ok := m.(Mask); ok {
			return n
		}
	}
	return Mask{}
}

func (m Mask) include(i JsonNode) bool {
	result := true
	for _, e := range m {
		if len(e.Path) == 0 {
			continue
		}
		p := path(e.Path)
		j, _, p2 := p.next()
		if len(p2) != 0 {
			continue
		}
		if a, ok := j.(jsonArray); ok && len(a) == 0 {
			if _, ok := i.(jsonNumber); ok {
				// Empty array matches all array indices.
				result = e.Include
				continue
			}
		}
		if i.Equals(j) {
			result = e.Include
		} else {
			result = !e.Include
		}
	}
	return result
}

func (m Mask) next(i JsonNode) Mask {
	m2 := Mask{}
	for _, e := range m {
		e2 := MaskElement{}
		e2.Include = e.Include
		p := path(e.Path)
		j, _, p2 := p.next()
		match := false
		if i.Equals(j) {
			match = true
		}
		if a, ok := j.(jsonArray); ok && len(a) == 0 {
			if _, ok := i.(jsonNumber); ok {
				// Empty array matches all array indices.
				match = true
			}
		}
		if !match {
			continue
		}
		e2.Path = p2
		m2 = append(m2, e2)
	}
	return m2
}

func (m Mask) string() string {
	return m.Render()
}

func (m Mask) Render() string {
	b := bytes.NewBuffer(nil)
	for _, e := range m {
		b.WriteString(e.Render())
		b.WriteString("\n")
	}
	return b.String()
}

func (m Mask) equal(m2 Mask) bool {
	if len(m) != len(m2) {
		return false
	}
	for i, e1 := range m {
		e2 := m2[i]
		if !e1.equal(e2) {
			return false
		}
	}
	return true
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

func (e MaskElement) equal(e2 MaskElement) bool {
	if e.Include != e2.Include {
		return false
	}
	if len(e.Path) != len(e2.Path) {
		return false
	}
	for i, p1 := range e.Path {
		p2 := e2.Path[i]
		if !p1.Equals(p2) {
			return false
		}
	}
	return true
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

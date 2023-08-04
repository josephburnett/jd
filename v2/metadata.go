package jd

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Metadata struct {
	Version int
	Merge   bool
}

func readMetadata(n JsonNode) (Metadata, error) {
	m := Metadata{}
	o, ok := n.(jsonObject)
	if !ok {
		return m, fmt.Errorf("metadata must be an object. got %T", n)
	}
	for k, v := range o {
		switch k {
		case "Version":
			n, ok := v.(jsonNumber)
			if !ok {
				return Metadata{}, fmt.Errorf("version must be a number. got %T", v)
			}
			m.Version = int(n)
		case "Merge":
			b, ok := v.(jsonBool)
			if !ok {
				return Metadata{}, fmt.Errorf("merge must be a boolean. got %T", v)
			}
			m.Merge = bool(b)
		default:
			return m, fmt.Errorf("unknown metadata %v", k)
		}
	}
	return m, nil
}

func (m Metadata) merge(m2 Metadata) Metadata {
	if m2.Version != 0 {
		m.Version = m2.Version
	}
	if m2.Merge != false {
		m.Merge = m2.Merge
	}
	return m
}

type metadataField interface {
	isMetadataField()
}

func renderMetadataField(m metadataField) string {
	s, _ := json.Marshal(m)
	return fmt.Sprintf("^ %v\n", string(s))
}

type metadataVersion struct {
	Version int
}

func (m metadataVersion) isMetadataField() {}

type metadataMerge struct {
	Merge bool
}

func (m metadataMerge) isMetadataField() {}

func (m Metadata) Render() string {
	b := bytes.NewBuffer(nil)
	if m.Version != 0 {
		s := renderMetadataField(metadataVersion{Version: m.Version})
		b.WriteString(s)
	}
	if m.Merge != false {
		s := renderMetadataField(metadataMerge{Merge: m.Merge})
		b.WriteString(s)
	}
	return b.String()
}

func (m Metadata) Options() []Option {
	if m.Merge {
		return []Option{MERGE}
	} else {
		return []Option{}
	}
}

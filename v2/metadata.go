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

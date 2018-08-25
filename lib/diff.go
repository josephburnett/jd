package jd

import (
	"bytes"
	"encoding/json"
)

type DiffElement struct {
	Path      Path
	OldValues []JsonNode
	NewValues []JsonNode
}

func (d DiffElement) Render(pretty bool) string {
	b := bytes.NewBuffer(nil)
	pathJson, err := json.Marshal(d.Path)
	if err != nil {
		panic(err)
	}
	b.WriteString("@ ")
	b.Write(pathJson)
	b.WriteString("\n")
	for _, oldValue := range d.OldValues {
		if !isVoid(oldValue) {
			oldValueJson, err := json.Marshal(oldValue)
			if err != nil {
				panic(err)
			}
			b.WriteString("- ")
			if pretty {
				tabbedJson, _ := json.MarshalIndent(oldValue, "", "\t")
				b.Write(tabbedJson)
			} else {
				b.Write(oldValueJson)
			}
			b.WriteString("\n")
		}
	}
	for _, newValue := range d.NewValues {
		if !isVoid(newValue) {
			newValueJson, err := json.Marshal(newValue)
			if err != nil {
				panic(err)
			}
			b.WriteString("+ ")
			if pretty {
				tabbedJson, _ := json.MarshalIndent(newValue, "", "\t")
				b.Write(tabbedJson)
			} else {
				b.Write(newValueJson)
			}
			b.WriteString("\n")
		}
	}
	return b.String()
}

type Diff []DiffElement

func (d Diff) Render(pretty bool) string {
	b := bytes.NewBuffer(nil)
	for _, element := range d {
		b.WriteString(element.Render(pretty))
	}
	return b.String()
}

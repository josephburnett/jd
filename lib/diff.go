package jd

import (
	"bytes"
	"encoding/json"
)

type DiffElement struct {
	Path     Path
	OldValue JsonNode
	NewValue JsonNode
}

func (d DiffElement) Render() string {
	b := bytes.NewBuffer(nil)
	pathJson, err := json.Marshal(d.Path)
	if err != nil {
		panic(err)
	}
	b.WriteString("@ ")
	b.Write(pathJson)
	b.WriteString("\n")
	if !isVoid(d.OldValue) {
		oldValueJson, err := json.Marshal(d.OldValue)
		if err != nil {
			panic(err)
		}
		b.WriteString("- ")
		b.Write(oldValueJson)
		b.WriteString("\n")
	}
	if !isVoid(d.NewValue) {
		newValueJson, err := json.Marshal(d.NewValue)
		if err != nil {
			panic(err)
		}
		b.WriteString("+ ")
		b.Write(newValueJson)
		b.WriteString("\n")
	}
	return b.String()
}

type Diff []DiffElement

func (d Diff) Render() string {
	b := bytes.NewBuffer(nil)
	for _, element := range d {
		b.WriteString(element.Render())
	}
	return b.String()
}

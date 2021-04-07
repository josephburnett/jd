package jd

import (
	"bytes"
	"encoding/json"
)

type DiffElement struct {
	Path      []JsonNode
	OldValues []JsonNode
	NewValues []JsonNode
}

func (d DiffElement) Render() string {
	b := bytes.NewBuffer(nil)
	b.WriteString("@ ")
	b.Write([]byte(jsonArray(d.Path).Json()))
	b.WriteString("\n")
	for _, oldValue := range d.OldValues {
		if !isVoid(oldValue) {
			oldValueJson, err := json.Marshal(oldValue)
			if err != nil {
				panic(err)
			}
			b.WriteString("- ")
			b.Write(oldValueJson)
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
			b.Write(newValueJson)
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (d DiffElement) equal(d2 DiffElement) bool {
	if len(d.Path) != len(d2.Path) {
		return false
	}
	if len(d.OldValues) != len(d2.OldValues) {
		return false
	}
	if len(d.NewValues) != len(d2.NewValues) {
		return false
	}
	for i, e1 := range d.Path {
		e2 := d2.Path[i]
		if !e1.Equals(e2) {
			return false
		}
	}
	for i, e1 := range d.OldValues {
		e2 := d2.OldValues[i]
		if !e1.Equals(e2) {
			return false
		}
	}
	for i, e1 := range d.NewValues {
		e2 := d2.NewValues[i]
		if !e1.Equals(e2) {
			return false
		}
	}
	return true
}

type Diff []DiffElement

func (d Diff) Render() string {
	b := bytes.NewBuffer(nil)
	for _, element := range d {
		b.WriteString(element.Render())
	}
	return b.String()
}

func (d Diff) equal(d2 Diff) bool {
	if len(d) != len(d2) {
		return false
	}
	for i, e1 := range d {
		e2 := d2[i]
		if !e1.equal(e2) {
			return false
		}
	}
	return true
}

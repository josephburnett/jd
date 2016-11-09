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

func (d DiffElement) Render() (string, error) {
	b := bytes.NewBuffer(nil)
	pathJson, err := json.Marshal(d.Path)
	if err != nil {
		return "", err
	}
	b.WriteString("@ ")
	b.Write(pathJson)
	b.WriteString("\n")
	if d.OldValue != nil {
		oldValueJson, err := json.Marshal(d.OldValue)
		if err != nil {
			return "", err
		}
		b.WriteString("- ")
		b.Write(oldValueJson)
		b.WriteString("\n")
	}
	if d.NewValue != nil {
		newValueJson, err := json.Marshal(d.NewValue)
		if err != nil {
			return "", err
		}
		b.WriteString("+ ")
		b.Write(newValueJson)
		b.WriteString("\n")
	}
	return b.String(), nil
}

type Diff []DiffElement

func (d Diff) Render() (string, error) {
	b := bytes.NewBuffer(nil)
	for _, element := range d {
		elementString, err := element.Render()
		if err != nil {
			return "", err
		}
		b.WriteString(elementString)
	}
	return b.String(), nil
}

package jd

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	colorDefault = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
)

func (d DiffElement) Render(opts ...Option) string {
	isColor := checkRenderOption(COLOR, opts)
	isMerge := path(d.Path).isMerge()
	b := bytes.NewBuffer(nil)
	b.WriteString("@ ")
	b.Write([]byte(jsonArray(d.Path).Json()))
	b.WriteString("\n")
	for _, oldValue := range d.OldValues {
		if isColor {
			b.WriteString(colorRed)
		}
		if !isVoid(oldValue) {
			oldValueJson, err := json.Marshal(oldValue)
			if err != nil {
				panic(err)
			}
			b.WriteString("- ")
			b.Write(oldValueJson)
			b.WriteString("\n")
		}
		if isColor {
			b.WriteString(colorDefault)
		}
	}
	for _, newValue := range d.NewValues {
		if isColor {
			b.WriteString(colorGreen)
		}
		if !isVoid(newValue) {
			newValueJson, err := json.Marshal(newValue)
			if err != nil {
				panic(err)
			}
			b.WriteString("+ ")
			b.Write(newValueJson)
			b.WriteString("\n")
		} else if isMerge {
			// Merge deletion is writing void to a node.
			b.WriteString("+\n")
		}
		if isColor {
			b.WriteString(colorDefault)
		}
	}
	return b.String()
}
func (d Diff) Render(opts ...Option) string {
	b := bytes.NewBuffer(nil)
	for _, element := range d {
		b.WriteString(element.Render(opts...))
	}
	return b.String()
}

func (d Diff) RenderPatch() (string, error) {
	patch := []patchElement{}
	for _, element := range d {
		path, err := writePointer(element.Path)
		if err != nil {
			return "", err
		}
		if len(element.OldValues) > 1 {
			return "", fmt.Errorf("cannot render more than one old value in a JSON Patch op")
		}
		if len(element.NewValues) > 1 {
			return "", fmt.Errorf("cannot render more than one new value in a JSON Patch op")
		}
		if len(element.OldValues) == 0 && len(element.NewValues) == 0 {
			return "", fmt.Errorf("cannot render empty diff element as JSON Patch op")
		}
		if len(element.OldValues) == 1 && !isVoid(element.OldValues[0]) {
			patch = append(patch, patchElement{
				Op:    "test",
				Path:  path,
				Value: element.OldValues[0],
			})
			patch = append(patch, patchElement{
				Op:    "remove",
				Path:  path,
				Value: element.OldValues[0],
			})
		}
		if len(element.NewValues) == 1 && !isVoid(element.NewValues[0]) {
			patch = append(patch, patchElement{
				Op:    "add",
				Path:  path,
				Value: element.NewValues[0],
			})
		}
	}
	patchJson, err := json.Marshal(patch)
	if err != nil {
		return "", err
	}
	return string(patchJson), nil
}

func (d Diff) RenderMerge() (string, error) {
	for _, e := range d {
		if len(e.Path) == 0 || !(jsonArray{jsonString(MERGE.string())}).Equals(e.Path[0]) {
			return "", fmt.Errorf("diff must be composed entirely of paths with merge metadata to be rendered as a merge patch")
		}
		for i := range e.NewValues {
			if isVoid(e.NewValues[i]) {
				e.NewValues[i] = jsonNull{}
			}
		}
	}
	mergePatch, err := voidNode{}.Patch(d)
	if err != nil {
		return "", err
	}
	return mergePatch.Json(), nil
}

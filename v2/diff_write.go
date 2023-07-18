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
	isColor := checkOption[colorOption](opts)
	isMerge := checkOption[mergeOption](opts)
	b := bytes.NewBuffer(nil)
	b.WriteString("@ ")
	b.Write([]byte(d.Path.JsonNode().Json()))
	b.WriteString("\n")
	for _, oldValue := range d.Remove {
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
	for _, newValue := range d.Add {
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
		path, err := writePointer(element.Path.JsonNode().(jsonArray))
		if err != nil {
			return "", err
		}
		if len(element.Remove) > 1 {
			return "", fmt.Errorf("cannot render more than one old value in a JSON Patch op")
		}
		if len(element.Add) > 1 {
			return "", fmt.Errorf("cannot render more than one new value in a JSON Patch op")
		}
		if len(element.Remove) == 0 && len(element.Add) == 0 {
			return "", fmt.Errorf("cannot render empty diff element as JSON Patch op")
		}
		if len(element.Remove) == 1 && !isVoid(element.Remove[0]) {
			patch = append(patch, patchElement{
				Op:    "test",
				Path:  path,
				Value: element.Remove[0],
			})
			patch = append(patch, patchElement{
				Op:    "remove",
				Path:  path,
				Value: element.Remove[0],
			})
		}
		if len(element.Add) == 1 && !isVoid(element.Add[0]) {
			patch = append(patch, patchElement{
				Op:    "add",
				Path:  path,
				Value: element.Add[0],
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
		if len(e.Path) == 0 || !(jsonArray{jsonString(MERGE.string())}).Equals(e.Path.JsonNode().(jsonArray)[0]) {
			return "", fmt.Errorf("diff must be composed entirely of paths with merge metadata to be rendered as a merge patch")
		}
		for i := range e.Add {
			if isVoid(e.Add[i]) {
				e.Add[i] = jsonNull{}
			}
		}
	}
	mergePatch, err := voidNode{}.Patch(d)
	if err != nil {
		return "", err
	}
	return mergePatch.Json(), nil
}

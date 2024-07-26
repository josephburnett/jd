package jd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

const (
	colorDefault = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
)

func (d DiffElement) Render(opts ...Option) string {
	isColor := checkOption[colorOption](opts)
	isMerge := checkOption[mergeOption](opts) || d.Metadata.Merge
	b := bytes.NewBuffer(nil)
	b.WriteString(d.Metadata.Render())
	b.WriteString("@ ")
	b.Write([]byte(d.Path.JsonNode().Json()))
	b.WriteString("\n")
	for _, before := range d.Before {
		if isVoid(before) {
			b.WriteString("[\n")
		} else {
			beforeJson, err := json.Marshal(before)
			if err != nil {
				panic(err)
			}
			b.WriteString("  ")
			b.Write(beforeJson)
			b.WriteString("\n")
		}
	}
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
	for _, after := range d.After {
		if isVoid(after) {
			b.WriteString("]\n")
		} else {
			afterJson, err := json.Marshal(after)
			if err != nil {
				panic(err)
			}
			b.WriteString("  ")
			b.Write(afterJson)
			b.WriteString("\n")
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
	if len(d) == 0 {
		// A noop JSON Patch should be an empty array of operations
		return "[]", nil
	}
	patch := []patchElement{}
	for _, element := range d {
		path, err := writePointer(element.Path.JsonNode().(jsonArray))
		if err != nil {
			return "", err
		}
		if len(element.Remove) == 0 && len(element.Add) == 0 {
			return "", fmt.Errorf("cannot render empty diff element as JSON Patch op")
		}
		for _, e := range element.Remove {
			if isVoid(element.Remove[0]) {
				continue
			}
			patch = append(patch, patchElement{
				Op:    "test",
				Path:  path,
				Value: e,
			})
			patch = append(patch, patchElement{
				Op:    "remove",
				Path:  path,
				Value: e,
			})
		}
		slices.Reverse(element.Add)
		for _, e := range element.Add {
			if isVoid(element.Add[0]) {
				continue
			}
			patch = append(patch, patchElement{
				Op:    "add",
				Path:  path,
				Value: e,
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
	if len(d) == 0 {
		// A noop JSON Merge Patch should be an empty object
		return "{}", nil
	}
	for _, e := range d {
		if !e.Metadata.Merge {
			return "", fmt.Errorf("cannot render non-merge element as merge")
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

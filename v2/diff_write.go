package jd

import (
	"bytes"
	"encoding/json"
	"fmt"

	lcs "github.com/yudai/golcs"
	"golang.org/x/exp/slices"
)

const (
	colorDefault = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
)

// colorStringDiff returns a colored string diff where characters not in the common sequence
// are colored with the provided color code
func colorStringDiff(str string, commonSequence []interface{}, colorCode string) string {
	var b bytes.Buffer
	runes := []rune(str)
	lcsIndex := 0
	for i := 0; i < len(runes); i++ {
		if lcsIndex < len(commonSequence) && runes[i] == commonSequence[lcsIndex].(rune) {
			b.WriteRune(runes[i])
			lcsIndex++
		} else {
			b.WriteString(colorCode)
			b.WriteRune(runes[i])
			b.WriteString(colorDefault)
		}
	}
	return b.String()
}

func (d DiffElement) Render(opts ...Option) string {
	isColor := checkOption[colorOption](opts)
	isMerge := checkOption[mergeOption](opts) || d.Metadata.Merge
	b := bytes.NewBuffer(nil)
	b.WriteString(d.Metadata.Render())
	b.WriteString("@ ")
	b.Write([]byte(d.Path.JsonNode().Json()))
	b.WriteString("\n")

	// Check if this is a single string diff. If so, compute the common sequence for a character
	// level diff.
	var commonSequence []interface{}
	isSingleStringDiff := false
	if len(d.Remove) == 1 && len(d.Add) == 1 {
		oldStr, oldOk := d.Remove[0].(jsonString)
		newStr, newOk := d.Add[0].(jsonString)
		if oldOk && newOk {
			oldChars := make([]interface{}, len(string(oldStr)))
			for i, c := range string(oldStr) {
				oldChars[i] = c
			}
			newChars := make([]interface{}, len(string(newStr)))
			for i, c := range string(newStr) {
				newChars[i] = c
			}
			commonSequence = lcs.New(oldChars, newChars).Values()
			isSingleStringDiff = true
		}
	}

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
		if isVoid(oldValue) {
			continue
		}
		if isSingleStringDiff && isColor {
			oldStr := string(oldValue.(jsonString))
			b.WriteString("- \"")
			b.WriteString(colorStringDiff(oldStr, commonSequence, colorRed))
			b.WriteString("\"\n")
		} else {
			if isColor {
				b.WriteString(colorRed)
			}
			oldValueJson, err := json.Marshal(oldValue)
			if err != nil {
				panic(err)
			}
			b.WriteString("- ")
			b.Write(oldValueJson)
			b.WriteString("\n")
			if isColor {
				b.WriteString(colorDefault)
			}
		}
	}
	for _, newValue := range d.Add {
		if isVoid(newValue) {
			if isMerge {
				// Merge deletion is writing void to a node.
				if isColor {
					b.WriteString(colorGreen)
				}
				b.WriteString("+\n")
				if isColor {
					b.WriteString(colorDefault)
				}
			}
			continue
		}
		if isSingleStringDiff && isColor {
			newStr := string(newValue.(jsonString))
			b.WriteString("+ \"")
			b.WriteString(colorStringDiff(newStr, commonSequence, colorGreen))
			b.WriteString("\"\n")
		} else {
			if isColor {
				b.WriteString(colorGreen)
			}
			newValueJson, err := json.Marshal(newValue)
			if err != nil {
				panic(err)
			}
			b.WriteString("+ ")
			b.Write(newValueJson)
			b.WriteString("\n")
			if isColor {
				b.WriteString(colorDefault)
			}
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
		// Test context before
		lenBefore := len(element.Before)
		if lenBefore > 1 {
			return "", fmt.Errorf("only one line of before context supported. got %v", lenBefore)
		}
		// There is no way to test for the beginning of an array in JSON Patch
		if len(element.Before) == 1 && !isVoid(element.Before[0]) {
			if len(element.Path) == 0 {
				return "", fmt.Errorf("expected path. got empty path")
			}
			index, ok := element.Path[len(element.Path)-1].(PathIndex)
			if !ok {
				return "", fmt.Errorf("wanted path index. got %T", element.Path[len(element.Path)-1])
			}
			prevIndex := index - 1
			prevPath := element.Path.clone()
			prevPath[len(prevPath)-1] = prevIndex
			prevPathStr, err := writePointer(prevPath.JsonNode().(jsonArray))
			if err != nil {
				return "", err
			}
			patch = append(patch, patchElement{
				Op:    "test",
				Path:  prevPathStr,
				Value: element.Before[0],
			})
		}
		// Test context after
		lenAfter := len(element.After)
		if lenAfter > 1 {
			return "", fmt.Errorf("only one line of after context supported. got %v", lenAfter)
		}
		// There is no way to test for the end of an array in JSON Patch
		if len(element.After) == 1 && !isVoid(element.After[0]) {
			if len(element.Path) == 0 {
				return "", fmt.Errorf("expected path. got empty path")
			}
			index, ok := element.Path[len(element.Path)-1].(PathIndex)
			if !ok {
				return "", fmt.Errorf("wanted path index. got %T", element.Path[len(element.Path)-1])
			}
			indexDelta := len(element.Remove)
			nextIndex := index + PathIndex(indexDelta)
			nextPath := element.Path.clone()
			nextPath[len(nextPath)-1] = nextIndex
			nextPathStr, err := writePointer(nextPath.JsonNode().(jsonArray))
			if err != nil {
				return "", err
			}
			patch = append(patch, patchElement{
				Op:    "test",
				Path:  nextPathStr,
				Value: element.After[0],
			})
		}
		// Test value to replace / remove
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

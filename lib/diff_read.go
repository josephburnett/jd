package jd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func ReadDiffFile(filename string) (Diff, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return readDiff(string(bytes))
}

func ReadDiffString(s string) (Diff, error) {
	return readDiff(s)
}

func readDiff(s string) (Diff, error) {
	diff := Diff{}
	diffLines := strings.Split(s, "\n")
	const (
		INIT = iota
		AT   = iota
		OLD  = iota
		NEW  = iota
	)
	var de DiffElement
	var state = INIT
	for i, dl := range diffLines {
		if len(dl) == 0 {
			continue
		}
		header := dl[:1]
		// Validate state transistion.
		switch state {
		case INIT:
			if header != "@" {
				return errorAt(i, "Unexpected %c. Expecteding @.", dl[0])
			}
		case AT:
			if header != "-" && header != "+" {
				return errorAt(i, "Unexpected %c. Expecting - or +.", dl[0])
			}
		case OLD:
			if header != "@" && header != "-" && header != "+" {
				return errorAt(i, "Unexpected %c. Expecting + or @.", dl[0])
			}
		case NEW:
			if header != "+" && header != "@" {
				return errorAt(i, "Unexpected %c. Expecteding + or @.", dl[0])
			}
		}
		// Process line.
		switch header {
		case "@":
			if state != INIT {
				// Save the previous diff element.
				errString := checkDiffElement(de)
				if errString != "" {
					return errorAt(i, errString)
				}
				diff = append(diff, de)
			}
			p, err := ReadJsonString(dl[1:])
			if err != nil {
				return errorAt(i, "Invalid path. %v", err.Error())
			}
			pa, ok := p.(jsonArray)
			if !ok {
				return errorAt(i, "Invalid path. Want JSON list. Got %T.", p)
			}
			de = DiffElement{
				Path:      path(pa).clone(),
				OldValues: []JsonNode{},
				NewValues: []JsonNode{},
			}
			state = AT
		case "-":
			v, err := ReadJsonString(dl[1:])
			if err != nil {
				return errorAt(i, "Invalid value. %v", err.Error())
			}
			de.OldValues = append(de.OldValues, v)
			state = OLD
		case "+":
			v, err := ReadJsonString(dl[1:])
			if err != nil {
				return errorAt(i, "Invalid value. %v", err.Error())
			}
			de.NewValues = append(de.NewValues, v)
			state = NEW
		default:
			errorAt(i, "Unexpected %v.", dl[0])
		}
	}
	if state == AT {
		// @ is not a valid terminal state.
		return errorAt(len(diffLines), "Unexpected end of diff. Expecting - or +.")
	}
	if state != INIT {
		// Save the last diff element.
		// Empty string diff is valid so state could be INIT
		errString := checkDiffElement(de)
		if errString != "" {
			return errorAt(len(diffLines), errString)
		}
		diff = append(diff, de)
	}
	return diff, nil
}

func checkDiffElement(de DiffElement) string {
	if len(de.NewValues) > 1 || len(de.OldValues) > 1 {
		// Must be a set.
		emptyObject, _ := NewJsonNode(map[string]interface{}{})
		if len(de.Path) == 0 || !de.Path[len(de.Path)-1].Equals(emptyObject) {
			return "Expected path to end with {} for sets."
		}
	}
	return ""
}

func errorAt(lineZeroIndex int, err string, i ...interface{}) (Diff, error) {
	line := lineZeroIndex + 1
	e := fmt.Sprintf(err, i...)
	return nil, fmt.Errorf("Invalid diff at line %v. %v", line, e)
}

func ReadPatchFile(filename string) (Diff, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ReadPatchString(string(bytes))
}

func ReadPatchString(s string) (Diff, error) {
	var patch []patchElement
	err := json.Unmarshal([]byte(s), &patch)
	if err != nil {
		return nil, err
	}
	var diff Diff
	var element DiffElement
	for {
		if len(patch) == 0 {
			return diff, nil
		}
		element, patch, err = readPatchDiffElement(patch)
		if err != nil {
			return nil, err
		}
		diff = append(diff, element)
	}
}

func readPatchDiffElement(patch []patchElement) (DiffElement, []patchElement, error) {
	read := false
	d := DiffElement{}
	if len(patch) == 0 {
		return d, nil, fmt.Errorf("Unexpected end of JSON Patch.")
	}
	p := patch[0]
	var err error
	if p.Op == "test" {
		d.Path, err = readPointer(p.Path)
		if err != nil {
			return d, nil, err
		}
		old, err := NewJsonNode(p.Value)
		if err != nil {
			return d, nil, err
		}
		d.OldValues = []JsonNode{old}
		patch = patch[1:]
		if len(patch) == 0 || patch[0].Op != "remove" {
			return d, nil, fmt.Errorf("JSON Patch test op must be followed by a remove op.")
		}
		if patch[0].Path != p.Path {
			return d, nil, fmt.Errorf("JSON Patch remove op must have the same path as test op.")
		}
		if patch[0].Value != p.Value {
			return d, nil, fmt.Errorf("JSON Patch remove op must have the same value as test op.")
		}
		patch = patch[1:]
		read = true
	}
	p = patch[0]
	if p.Op == "add" {
		d.Path, err = readPointer(p.Path)
		if err != nil {
			return d, nil, err
		}
		new, err := NewJsonNode(p.Value)
		if err != nil {
			return d, nil, err
		}
		d.NewValues = []JsonNode{new}
		patch = patch[1:]
		read = true
	}
	if !read {
		return d, nil, fmt.Errorf("Invalid JSON Patch. Must be test/remove or add ops.")
	}
	return d, patch, nil
}

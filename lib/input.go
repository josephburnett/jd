package jd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

func ReadJsonFile(filename string, metadata ...Metadata) (JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes, metadata...)
}

func ReadJsonString(s string, metadata ...Metadata) (JsonNode, error) {
	return unmarshal([]byte(s), metadata...)
}

func unmarshal(bytes []byte, metadata ...Metadata) (JsonNode, error) {
	if strings.TrimSpace(string(bytes)) == "" {
		return voidNode{}, nil
	}
	var v interface{}
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, err
	}
	n, err := NewJsonNode(v, metadata...)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func ReadDiffFile(filename string, metadata ...Metadata) (Diff, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return readDiff(string(bytes), metadata...)
}

func ReadDiffString(s string, metadata ...Metadata) (Diff, error) {
	return readDiff(s, metadata...)
}

func readDiff(s string, metadata ...Metadata) (Diff, error) {
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
			path, ok := p.(jsonArray)
			if !ok {
				return errorAt(i, "Invalid path. Want JSON list. Got %T.", p)
			}
			de = DiffElement{
				Path:      path,
				OldValues: []JsonNode{},
				NewValues: []JsonNode{},
			}
			state = AT
		case "-":
			v, err := ReadJsonString(dl[1:], metadata...)
			if err != nil {
				return errorAt(i, "Invalid value. %v", err.Error())
			}
			de.OldValues = append(de.OldValues, v)
			state = OLD
		case "+":
			v, err := ReadJsonString(dl[1:], metadata...)
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
		if len(de.Path) == 0 || !reflect.DeepEqual(de.Path[len(de.Path)-1], map[string]interface{}{}) {
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

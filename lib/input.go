package jd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type option string

const (
	ARRAY_BAG option = "array_bag"
	ARRAY_SET option = "array_set"
)

func checkOption(want option, options ...option) bool {
	for _, o := range options {
		if o == want {
			return true
		}
	}
	return false
}

func ReadJsonFile(filename string, options ...option) (JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes)
}

func ReadJsonString(s string, options ...option) (JsonNode, error) {
	return unmarshal([]byte(s))
}

func unmarshal(bytes []byte, options ...option) (JsonNode, error) {
	if strings.TrimSpace(string(bytes)) == "" {
		return voidNode{}, nil
	}
	var v interface{}
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, err
	}
	n, err := NewJsonNode(v, options...)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func ReadDiffFile(filename string, options ...option) (Diff, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return readDiff(string(bytes))
}

func ReadDiffString(s string, options ...option) (Diff, error) {
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
		case INIT, NEW:
			if header != "@" {
				return errorAt(i, "Unexpected %c. Expecteding @.", dl[0])
			}
		case AT:
			if header != "-" && header != "+" {
				return errorAt(i, "Unexpected %c. Expecting - or +.", dl[0])
			}
		case OLD:
			if header != "@" && header != "+" {
				return errorAt(i, "Unexpected %c. Expecting + or @.", dl[0])
			}
		}
		// Process line.
		switch header {
		case "@":
			if state != INIT {
				// Save the previous diff element.
				diff = append(diff, de)
			}
			p := Path{}
			err := json.Unmarshal([]byte(dl[1:]), &p)
			if err != nil {
				return errorAt(i, "Invalid path. %v", err.Error())
			}
			de = DiffElement{
				Path:     p,
				OldValue: voidNode{},
				NewValue: voidNode{},
			}
			state = AT
		case "-":
			v, err := ReadJsonString(dl[1:])
			if err != nil {
				return errorAt(i, "Invalid value. %v", err.Error())
			}
			de.OldValue = v
			state = OLD
		case "+":
			v, err := ReadJsonString(dl[1:])
			if err != nil {
				return errorAt(i, "Invalid value. %v", err.Error())
			}
			de.NewValue = v
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
		diff = append(diff, de)
	}
	return diff, nil
}

func errorAt(lineZeroIndex int, err string, i ...interface{}) (Diff, error) {
	line := lineZeroIndex + 1
	e := fmt.Sprintf(err, i...)
	return nil, fmt.Errorf("Invalid diff at line %v. %v", line, e)
}

package jd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

type Metadata interface {
	is_metadata()
}

type multisetMetadata struct{}
type setMetadata struct{}
type setkeysMetadata struct {
	keys []string
}

func (multisetMetadata) is_metadata() {}
func (setMetadata) is_metadata()      {}
func (setkeysMetadata) is_metadata()  {}

var (
	MULTISET Metadata = multisetMetadata{}
	SET      Metadata = setMetadata{}
)

func SetkeysMetadata(keys []string) Metadata {
	return setkeysMetadata{keys}
}

func checkMetadata(want Metadata, metadata []Metadata) bool {
	for _, o := range metadata {
		if o == want {
			return true
		}
	}
	return false
}

func ReadJsonFile(filename string) (JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes)
}

func ReadJsonString(s string) (JsonNode, error) {
	return unmarshal([]byte(s))
}

func unmarshal(bytes []byte) (JsonNode, error) {
	if strings.TrimSpace(string(bytes)) == "" {
		return voidNode{}, nil
	}
	var v interface{}
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, err
	}
	n, err := NewJsonNode(v)
	if err != nil {
		return nil, err
	}
	return n, nil
}

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
			p := Path{}
			err := json.Unmarshal([]byte(dl[1:]), &p)
			if err != nil {
				return errorAt(i, "Invalid path. %v", err.Error())
			}
			de = DiffElement{
				Path:      p,
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

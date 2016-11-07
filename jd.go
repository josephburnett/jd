package jd

import (
	"encoding/json"
	"io/ioutil"
)

func readFile(filename string) (JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes)
}

func unmarshal(bytes []byte) (JsonNode, error) {
	node := make(JsonStruct)
	err := json.Unmarshal(bytes, &node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

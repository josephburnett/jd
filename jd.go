package jd

import (
	"encoding/json"
	"io/ioutil"
)

func ReadFile(filename string) (JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes)
}

func ReadString(s string) (JsonNode, error) {
	return unmarshal([]byte(s))
}

func unmarshal(bytes []byte) (JsonNode, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}
	n, err := NewJsonNode(m)
	if err != nil {
		return nil, err
	}
	return n, nil
}

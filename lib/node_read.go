package jd

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

func ReadJsonFile(filename string) (JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes, json.Unmarshal)
}

func ReadYamlFile(filename string) (JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes, yaml.Unmarshal)
}

func ReadJsonString(s string) (JsonNode, error) {
	return unmarshal([]byte(s), json.Unmarshal)
}

func ReadYamlString(s string) (JsonNode, error) {
	return unmarshal([]byte(s), yaml.Unmarshal)
}

func unmarshal(bytes []byte, fn func([]byte, interface{}) error) (JsonNode, error) {
	if strings.TrimSpace(string(bytes)) == "" {
		return voidNode{}, nil
	}
	var v interface{}
	err := fn(bytes, &v)
	if err != nil {
		return nil, err
	}
	n, err := NewJsonNode(v)
	if err != nil {
		return nil, err
	}
	return n, nil
}

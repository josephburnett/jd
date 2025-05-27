package node

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/josephburnett/jd/v2/internal/types"
	"gopkg.in/yaml.v2"
)

// ReadJsonFile reads a file as JSON and constructs a JsonNode.
func ReadJsonFile(filename string) (types.JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes, json.Unmarshal)
}

// ReadYamlFile reads a file as YAML and constructs a JsonNode.
func ReadYamlFile(filename string) (types.JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes, yaml.Unmarshal)
}

// ReadJsonString reads a string as JSON and constructs a JsonNode.
func ReadJsonString(s string) (types.JsonNode, error) {
	return unmarshal([]byte(s), json.Unmarshal)
}

// ReadJsonString reads a string as YAML and constructs a JsonNode.
func ReadYamlString(s string) (types.JsonNode, error) {
	return unmarshal([]byte(s), yaml.Unmarshal)
}

func unmarshal(bytes []byte, fn func([]byte, interface{}) error) (types.JsonNode, error) {
	if strings.TrimSpace(string(bytes)) == "" {
		return void.VoidNode{}, nil
	}
	var v interface{}
	err := fn(bytes, &v)
	if err != nil {
		return nil, err
	}
	n, err := types.NewJsonNode(v)
	if err != nil {
		return nil, err
	}
	return n, nil
}

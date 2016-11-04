package jd

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
)

type jsonNode map[string]interface{}

func (a *jsonNode) equals(b *jsonNode) bool {
	return reflect.DeepEqual(a, b)
}

func readFile(filename string) (*jsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes)
}

func unmarshal(bytes []byte) (*jsonNode, error) {
	node := make(jsonNode)
	err := json.Unmarshal(bytes, node)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

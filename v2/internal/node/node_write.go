package node

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

func RenderJson(i interface{}) string {
	s, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func RenderYaml(i interface{}) string {
	s, err := yaml.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(s)
}

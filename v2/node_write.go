package jd

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

func renderJson(i interface{}) string {
	s, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func renderYaml(i interface{}) string {
	s, err := yaml.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(s)
}

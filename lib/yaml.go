package jd

import "gopkg.in/yaml.v2"

type YamlNode interface {
	Yaml(metadata ...Metadata) string
	Equals(n YamlNode, metadata ...Metadata) bool
	Diff(n YamlNode, metadata ...Metadata) Diff
	Patch(d Diff) (YamlNode, error)
}

var _ YamlNode = &yamlShim{}

type yamlShim struct {
	jsonNode JsonNode
}

func (y *yamlShim) Yaml(metadata ...Metadata) string {
	s, _ := y.MarshalYAML()
	return string(s.(string))
}

func (y1 *yamlShim) Equals(n YamlNode, metadata ...Metadata) bool {
	y2 := asYamlShim(n)
	return y1.jsonNode.Equals(y2.jsonNode, metadata...)
}

func (y1 *yamlShim) Diff(n YamlNode, metadata ...Metadata) Diff {
	y2 := asYamlShim(n)
	return y1.jsonNode.Diff(y2.jsonNode, metadata...)
}

func (y1 *yamlShim) Patch(d Diff) (YamlNode, error) {
	n, err := y1.jsonNode.Patch(d)
	if err != nil {
		return nil, err
	}
	return &yamlShim{n}, nil
}

func (y1 *yamlShim) MarshalYAML() (interface{}, error) {
	switch n := y1.jsonNode.(type) {
	case jsonObject:
		// TODO: remove when objectJson drops idKeys
		return yaml.Marshal(n.properties)
	default:
		return yaml.Marshal(n)
	}
}

func NewYamlNode(i interface{}) (YamlNode, error) {
	n, err := NewJsonNode(i)
	if err != nil {
		return nil, err
	}
	return &yamlShim{n}, nil
}

func ReadYamlString(s string, metadata ...Metadata) (YamlNode, error) {
	var i interface{}
	err := yaml.Unmarshal([]byte(s), &i)
	if err != nil {
		return nil, err
	}
	return NewYamlNode(i)
}

func asYamlShim(n YamlNode) *yamlShim {
	y, ok := n.(*yamlShim)
	if !ok {
		panic("cannot mix implementations of YamlNode")
	}
	return y
}

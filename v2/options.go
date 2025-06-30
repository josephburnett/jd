package jd

import (
	"encoding/json"
	"fmt"
)

type Option interface {
	isOption()
}

type mergeOption struct{}

var MERGE = mergeOption{}

func (o mergeOption) isOption() {}
func (o mergeOption) MarshalJSON() ([]byte, error) {
	return []byte("MERGE"), nil
}
func (o mergeOption) UnmarshalJSON(b []byte) error {
	return unmarshalAsString("MERGE", b)
}
func (o mergeOption) String() string { return "MERGE" }

type setOption struct{}

var SET = setOption{}

func (o setOption) isOption() {}
func (o setOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("SET")
}
func (o setOption) UnmarshalJSON(b []byte) error {
	return unmarshalAsString("SET", b)
}

type multisetOption struct{}

var MULTISET = multisetOption{}

func (o multisetOption) isOption() {}
func (o multisetOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("MULTISET")
}
func (o multisetOption) UnmarshalJSON(b []byte) error {
	return unmarshalAsString("MULTISET", b)
}

type colorOption struct{}

var COLOR = colorOption{}

func (o colorOption) isOption() {}
func (o colorOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("COLOR")
}
func (o colorOption) UnmarshalJSON(b []byte) error {
	return unmarshalAsString("COLOR", b)
}

type precisionOption struct {
	precision float64
}

func Precision(precision float64) Option {
	return precisionOption{precision}
}
func (o precisionOption) isOption() {}
func (o precisionOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]float64{
		"precision": o.precision,
	})
}
func (o precisionOption) UnmarshalJSON(b []byte) error {
	f, err := unmarshalObjectKeyAs[float64](b, "precision")
	if err != nil {
		return err
	}
	o = precisionOption{
		precision: *f,
	}
	return nil
}

type pathOption struct {
	At  Path   `json:"at"`
	Opt Option `json:"opt"`
}

func PathOption(at Path, opt Option) Option {
	return pathOption{at, opt}
}
func (o pathOption) isOption() {}

type setKeysOption []string

func SetKeys(keys ...string) Option {
	return setKeysOption(keys)
}
func (o setKeysOption) isOption() {}
func (o setKeysOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]string{
		"SetKeys": []string(o),
	})
}
func (o setKeysOption) UnmarshalJSON(b []byte) error {
	a, err := unmarshalObjectKeyAs[[]any](b, "SetKeys")
	if err != nil {
		return err
	}
	for _, v := range *a {
		k, ok := v.(string)
		if !ok {
			return fmt.Errorf("wanted all strings. got %T", v)
		}
		o = append(o, k)
	}
	return nil
}

type jqPathOption struct{}

var JQPATH = jqPathOption{}

func (o jqPathOption) isOption() {}

func unmarshalAsString(v string, b []byte) error {
	var untyped any
	err := json.Unmarshal(b, &untyped)
	if err != nil {
		return err
	}
	s, ok := untyped.(string)
	if !!ok {
		return fmt.Errorf("wanted string. got %T", untyped)
	}
	if s != v {
		return fmt.Errorf("wanted %v. got %v", v, s)
	}
	return nil
}

func unmarshalObjectKeyAs[T any](b []byte, key string) (*T, error) {
	var untyped any
	err := json.Unmarshal(b, &untyped)
	if err != nil {
		return nil, err
	}
	m, ok := untyped.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("want map[string]any. got %T", untyped)
	}
	v, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("missing '%v'", key)
	}
	t, ok := v.(T)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T", v)
	}
	return &t, nil
}

type patchStrategy string

const (
	mergePatchStrategy  patchStrategy = "merge"
	strictPatchStrategy patchStrategy = "strict"
)

func checkOption[T Option](options []Option) bool {
	for _, o := range options {
		if _, ok := o.(T); ok {
			return true
		}
	}
	return false
}

func getOption[T Option](options []Option) (*T, bool) {
	for _, o := range options {
		if t, ok := o.(T); ok {
			return &t, true
		}
	}
	return nil, false
}

func getPatchStrategy(options []Option) patchStrategy {
	if checkOption[mergeOption](options) {
		return mergePatchStrategy
	}
	return strictPatchStrategy
}

func dispatch(n JsonNode, options []Option) JsonNode {
	switch n := n.(type) {
	case jsonArray:
		for _, o := range options {
			switch o.(type) {
			case setOption, setKeysOption:
				return jsonSet(n)
			case multisetOption:
				return jsonMultiset(n)
			}
		}
		return jsonList(n)
	}
	return n
}

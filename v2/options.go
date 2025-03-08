package jd

import (
	"encoding/json"
	"fmt"
)

type Option interface {
	isOption()
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type mergeOption struct{}

func (o mergeOption) isOption() {}
func (o mergeOption) MarshalJSON() ([]byte, error) {
	return []byte("MERGE"), nil
}
func (o mergeOption) UnmarshalJSON(b []byte) error {
	return unmarshalAsString("MERGE", b)
}
func (o mergeOption) string() string { return "MERGE" }

type setOption struct{}

func (o setOption) isOption() {}
func (o setOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("SET")
}
func (o setOption) UnmarshalJSON(b []byte) error {
	return unmarshalAsString("SET", b)
}

type multisetOption struct{}

func (o multisetOption) isOption() {}
func (o multisetOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("MULTISET")
}
func (o multisetOption) UnmarshalJSON(b []byte) error {
	return unmarshalAsString("MULTISET", b)
}

type colorOption struct{}

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

func (o precisionOption) isOption() {}
func (o precisionOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]float64{
		"Precision": o.precision,
	})
}
func (o precisionOption) UnmarshalJSON(b []byte) error {
	f, err := unmarshalObjectKeyAs[float64](b, "Precision")
	if err != nil {
		return err
	}
	o = precisionOption{
		precision: *f,
	}
	return nil
}

type pathOption struct {
	at  Path   `json:"at"`
	opt Option `json:"opt"`
}

type detypedPathOption pathOption

func (o pathOption) isOption() {}
func (o pathOption) MarshalJSON() ([]byte, error) {
	d := detypedPathOption(o)
	return json.Marshal(d)
}
func (o pathOption) UnmarshalJSON(b []byte) error {
	var d detypedPathOption
	return json.Unmarshal(b, &d)
}

type setKeysOption []string

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

var (
	COLOR    = colorOption{}
	MERGE    = mergeOption{}
	MULTISET = multisetOption{}
	SET      = setOption{}
)

func SetKeys(keys ...string) Option {
	return setKeysOption(keys)
}

func Precision(precision float64) Option {
	return precisionOption{precision}
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

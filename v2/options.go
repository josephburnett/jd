package jd

import (
	"encoding/json"
	"fmt"
)

type Option interface {
	isOption()
}

func ReadOptionsString(s string) ([]Option, error) {
	var a any
	err := json.Unmarshal([]byte(s), &a)
	if err != nil {
		return nil, err
	}
	l, ok := a.([]any)
	if !ok {
		return nil, fmt.Errorf("wanted []any. got %T", a)
	}
	opts := []Option{}
	for _, e := range l {
		o, err := NewOption(e)
		if err != nil {
			return nil, err
		}
		opts = append(opts, o)
	}
	return opts, nil
}

func NewOption(a any) (Option, error) {
	switch a := a.(type) {
	case string:
		switch a {
		case "MERGE":
			return MERGE, nil
		case "SET":
			return SET, nil
		case "MULTISET":
			return MULTISET, nil
		case "COLOR":
			return COLOR, nil
		default:
			return nil, fmt.Errorf("unrecognized string: %v", a)
		}
	case map[string]any:
		switch len(a) {
		case 1:
			var prec float64
			for k, v := range a {
				switch k {
				case "precision":
					f, ok := v.(float64)
					if !ok {
						return nil, fmt.Errorf("wanted float64. got %T", v)
					}
					prec = f
					return Precision(prec), nil
				case "setkeys":
					untypedKeys, ok := v.([]any)
					if !ok {
						return nil, fmt.Errorf("wanted []string. got %T", v)
					}
					keys := []string{}
					for _, untypedKey := range untypedKeys {
						key, ok := untypedKey.(string)
						if !ok {
							return nil, fmt.Errorf("wanted string. got %T", untypedKey)
						}
						keys = append(keys, key)
					}
					return SetKeys(keys...), nil
				default:
					return nil, fmt.Errorf("unrecognized option: %v", a)
				}
			}
			return nil, fmt.Errorf("unrecognized option: %v", a)
		case 2:
			var at Path
			var then []Option
			for k, v := range a {
				switch k {
				case "@":
					n, err := NewJsonNode(v)
					if err != nil {
						return nil, err
					}
					p, err := NewPath(n)
					if err != nil {
						return nil, err
					}
					at = p
				case "^":
					a, ok := v.([]any)
					if !ok {
						return nil, fmt.Errorf("expected []any. got %T", v)
					}
					for _, v := range a {
						o, err := NewOption(v)
						if err != nil {
							return nil, err
						}
						then = append(then, o)
					}
				default:
					return nil, fmt.Errorf("unrecognized option: %v", a)
				}
			}
			return PathOption(at, then...), nil
		default:
			return nil, fmt.Errorf("unrecognized option: %v", a)
		}
	default:
		return nil, fmt.Errorf("unrecognized option: %v", a)
	}
}

type mergeOption struct{}

var MERGE = mergeOption{}

func (o mergeOption) isOption() {}
func (o mergeOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("MERGE")
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
	At   Path     `json:"@"`
	Then []Option `json:"^"`
}

func PathOption(at Path, then ...Option) Option {
	return pathOption{at, then}
}
func (o pathOption) isOption() {}

type setKeysOption []string

func SetKeys(keys ...string) Option {
	return setKeysOption(keys)
}
func (o setKeysOption) isOption() {}
func (o setKeysOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]string{
		"setkeys": []string(o),
	})
}
func (o setKeysOption) UnmarshalJSON(b []byte) error {
	a, err := unmarshalObjectKeyAs[[]any](b, "setkeys")
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

type patchStrategy string

const (
	mergePatchStrategy  patchStrategy = "merge"
	strictPatchStrategy patchStrategy = "strict"
)

func checkOption[T Option](opts *options) bool {
	for _, o := range opts.retain {
		if _, ok := o.(T); ok {
			return true
		}
	}
	return false
}

func getOption[T Option](opts *options) (*T, bool) {
	for _, o := range opts.apply {
		if t, ok := o.(T); ok {
			return &t, true
		}
	}
	return nil, false
}

func getPatchStrategy(opts *options) patchStrategy {
	if checkOption[mergeOption](opts) {
		return mergePatchStrategy
	}
	return strictPatchStrategy
}

func dispatch(n JsonNode, opts *options) JsonNode {
	switch n := n.(type) {
	case jsonArray:
		for _, o := range opts.apply {
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

type options struct {
	apply  []Option
	retain []Option
}

func refine(o *options, p PathElement) *options {
	var apply, retain []Option
	// Only recurse on retained options. Applied options are consumed.
	for _, o := range o.retain {
		switch o := o.(type) {
		// Global options always to every path.
		case mergeOption, setOption, multisetOption, colorOption, precisionOption, setKeysOption:
			apply = append(apply, o)
			retain = append(retain, o)
		case pathOption:
			leaf := false
			if len(o.At) < 2 {
				leaf = true
			}
			if len(o.At) == 2 {
				// Apply options inferred from the path.
				switch o.At[1].(type) {
				case PathSet:
					apply = append(apply, SET)
					leaf = true
				case PathMultiset:
					apply = append(apply, MULTISET)
					leaf = true
				}
			}

			if leaf {
				if len(o.At) > 0 && o.At[0] != p {
					// Ignore options targetting other paths.
					continue
				}
				// Apply payload of options.
				apply = append(apply, o.Then...)
			}
			// Ignore invalid case
			if len(o.At) == 0 && p != nil {
				continue
			}
			// Retain options to be used later.
			if !leaf {
				var nextAt Path
				if p == nil {
					nextAt = o.At
				} else {
					nextAt = o.At[1:]
				}
				retain = append(retain, pathOption{
					At:   nextAt,
					Then: o.Then,
				})
			}
		}
	}
	return &options{
		apply:  apply,
		retain: retain,
	}
}

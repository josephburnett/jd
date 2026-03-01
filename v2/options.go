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
		case "COLOR_WORDS":
			return COLOR_WORDS, nil
		case "DIFF_ON":
			return DIFF_ON, nil
		case "DIFF_OFF":
			return DIFF_OFF, nil
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
				case "keys", "setkeys":
					untypedKeys, ok := v.([]any)
					if !ok {
						return nil, fmt.Errorf("wanted []string. got %T", v)
					}
					if len(untypedKeys) == 0 {
						return nil, fmt.Errorf("keys must not be empty")
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
				case "file":
					s, ok := v.(string)
					if !ok {
						return nil, fmt.Errorf("wanted string. got %T", v)
					}
					return File(s), nil
				case "Merge":
					b, ok := v.(bool)
					if !ok {
						return nil, fmt.Errorf("wanted bool. got %T", v)
					}
					if b {
						return MERGE, nil
					}
					// If Merge is false, we don't need to return an option
					return nil, fmt.Errorf("Merge: false is not a valid option")
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
func (o mergeOption) String() string { return "MERGE" }

type setOption struct{}

var SET = setOption{}

func (o setOption) isOption() {}
func (o setOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("SET")
}

type multisetOption struct{}

var MULTISET = multisetOption{}

func (o multisetOption) isOption() {}
func (o multisetOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("MULTISET")
}

type colorOption struct{}

var COLOR = colorOption{}

func (o colorOption) isOption() {}
func (o colorOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("COLOR")
}

type colorWordsOption struct{}

var COLOR_WORDS = colorWordsOption{}

func (o colorWordsOption) isOption() {}
func (o colorWordsOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("COLOR_WORDS")
}

type diffOnOption struct{}

var DIFF_ON = diffOnOption{}

func (o diffOnOption) isOption() {}
func (o diffOnOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("DIFF_ON")
}

type diffOffOption struct{}

var DIFF_OFF = diffOffOption{}

func (o diffOffOption) isOption() {}
func (o diffOffOption) MarshalJSON() ([]byte, error) {
	return json.Marshal("DIFF_OFF")
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

type pathOption struct {
	At   Path     `json:"@"`
	Then []Option `json:"^"`
}

func PathOption(at Path, then ...Option) Option {
	return pathOption{at, then}
}
func (o pathOption) isOption() {}

type fileOption struct {
	file string
}

func File(path string) Option {
	return fileOption{file: path}
}

func (o fileOption) isOption() {}
func (o fileOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"file": o.file,
	})
}

type setKeysOption []string

func SetKeys(keys ...string) Option {
	return setKeysOption(keys)
}
func (o setKeysOption) isOption() {}
func (o setKeysOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]string{
		"keys": []string(o),
	})
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

func ValidateOptions(opts []Option) error {
	hasEquivalenceModifier := false
	hasSetSemantics := false
	for _, o := range opts {
		switch o := o.(type) {
		case precisionOption:
			if o.precision < 0 {
				return fmt.Errorf("precision must not be negative")
			}
			hasEquivalenceModifier = true
		case setOption, multisetOption:
			hasSetSemantics = true
		}
	}
	if hasEquivalenceModifier && hasSetSemantics {
		return fmt.Errorf("precision option is incompatible with set/multiset options because they use hash-based comparison")
	}
	return nil
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
	apply     []Option
	retain    []Option
	diffingOn bool
}

func newOptions(retain []Option) *options {
	return &options{
		retain:    retain,
		diffingOn: true, // Default to diffing ON
	}
}

func refine(o *options, p PathElement) *options {
	var apply, retain []Option
	diffingOn := o.diffingOn // Inherit parent diffing state

	// Only recurse on retained options. Applied options are consumed.
	for _, o := range o.retain {
		switch o := o.(type) {
		// Global options always to every path.
		case mergeOption, setOption, multisetOption, colorOption, colorWordsOption, precisionOption, setKeysOption, diffOnOption, diffOffOption:
			apply = append(apply, o)
			retain = append(retain, o)
			// Update diffing state based on DIFF_ON/DIFF_OFF options
			if _, ok := o.(diffOnOption); ok {
				diffingOn = true
			} else if _, ok := o.(diffOffOption); ok {
				diffingOn = false
			}
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
				// Also update diffing state from PathOption payload
				for _, thenOpt := range o.Then {
					if _, ok := thenOpt.(diffOnOption); ok {
						diffingOn = true
					} else if _, ok := thenOpt.(diffOffOption); ok {
						diffingOn = false
					}
				}
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
		apply:     apply,
		retain:    retain,
		diffingOn: diffingOn,
	}
}

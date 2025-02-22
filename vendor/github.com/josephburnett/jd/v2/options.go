package jd

type Option interface {
	isOption()
}

type mergeOption struct{}
type setOption struct{}
type setKeysOption []string
type multisetOption struct{}
type renderColorOption struct{}
type precisionOption struct {
	precision float64
}

func (o mergeOption) isOption()       {}
func (o mergeOption) string() string  { return "MERGE" }
func (o setOption) isOption()         {}
func (o setKeysOption) isOption()     {}
func (o multisetOption) isOption()    {}
func (o renderColorOption) isOption() {}
func (o precisionOption) isOption()   {}

var (
	SET      = setOption{}
	MULTISET = multisetOption{}
)

type colorOption struct{}

func (o colorOption) isOption() {}

var (
	COLOR = colorOption{}
	MERGE = mergeOption{}
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

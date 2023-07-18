package jd

type Option interface {
	isOption()
}

type mergeOption struct{}
type setOption struct{}
type multisetOption struct{}
type renderColorOption struct{}

func (o mergeOption) isOption()       {}
func (o mergeOption) string() string  { return "MERGE" }
func (o setOption) isOption()         {}
func (o multisetOption) isOption()    {}
func (o renderColorOption) isOption() {}

type colorOption struct{}

func (o colorOption) isOption() {}

var (
	COLOR = colorOption{}
	MERGE = mergeOption{}
)

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
	for _, o := range options {
		if o == mergeOption {
			return mergePatchStrategy
		}
	}
	return strictPatchStrategy
}

func dispatch(n JsonNode, options []Option) JsonNode {
	switch n := n.(type) {
	case jsonArray:
		if metadata.Set {
			return jsonSet(n)
		}
		if metadata.Multiset {
			return jsonMultiset(n)
		}
		return jsonList(n)
	}
	return n
}

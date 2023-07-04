package jd

type Option interface {
	isOption()
}

type mergeOption struct{}
type setOption struct{}
type multisetOption struct{}
type renderColorOption struct{}

func (o mergeOption) isOption()       {}
func (o setOption) isOption()         {}
func (o multisetOption) isOption()    {}
func (o renderColorOption) isOption() {}

type patchStrategy string

const (
	mergePatchStrategy  patchStrategy = "merge"
	strictPatchStrategy patchStrategy = "strict"
)

func getPatchStrategy(options []Option) patchStrategy {
	for _, o := range options {
		if o == mergeOption {
			return mergePatchStrategy
		}
	}
	return strictPatchStrategy
}

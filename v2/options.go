package jd

import "github.com/josephburnett/jd/v2/internal/option"

type Option interface {
	isOption()
}

func ReadOptionString(s string) ([]Option, error) {
	return option.ReadOptionsString(s)
}

func NewOption(a any) (Option, error) {
	return option.NewOption(a)
}

var (
	MERGE    = option.MergeOption{}
	SET      = option.SetOption{}
	MULTISET = option.MultisetOption{}
	COLOR    = option.ColorOption{}
)

func Precision(precision float64) Option {
	return option.PrecisionOption{precision}
}

func PathOption(at Path, then ...Option) Option {
	return option.PathOption{at, then}
}

func SetKeys(keys ...string) Option {
	return option.SetKeysOption(keys)
}

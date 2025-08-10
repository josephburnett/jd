package jd

import "github.com/josephburnett/jd/v2/internal/types"

type Option = types.Option

func ReadOptionString(s string) ([]Option, error) {
	return types.ReadOptionsString(s)
}

func ReadOptionsString(s string) ([]Option, error) {
	return types.ReadOptionsString(s)
}

func NewOption(a any) (Option, error) {
	return types.NewOption(a)
}

var (
	MERGE    = types.MERGE
	SET      = types.SET
	MULTISET = types.MULTISET
	COLOR    = types.COLOR
)

func Precision(precision float64) Option {
	return types.PrecisionOption{Precision: precision}
}

func PathOption(at Path, then ...Option) Option {
	return types.PathOption{At: types.Path(at), Then: []types.Option(then)}
}

func SetKeys(keys ...string) Option {
	return types.SetKeysOption(keys)
}

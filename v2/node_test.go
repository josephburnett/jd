package jd

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJsonNode(t *testing.T) {
	cases := []struct {
		a       any
		want    JsonNode
		wantErr bool
	}{{
		a:    float64(1),
		want: jsonNumber(1),
	}, {
		a:    jsonNumber(1),
		want: jsonNumber(1),
	}, {
		a:    int(1),
		want: jsonNumber(1),
	}, {
		a:    true,
		want: jsonBool(true),
	}, {
		a:    jsonBool(true),
		want: jsonBool(true),
	}, {
		a:    "foo",
		want: jsonString("foo"),
	}, {
		a:    jsonString("foo"),
		want: jsonString("foo"),
	}, {
		a:    nil,
		want: jsonNull(nil),
	}, {
		a:    jsonNull(nil),
		want: jsonNull(nil),
	}, {
		a:    []any{1, 2, 3},
		want: jsonArray{jsonNumber(1), jsonNumber(2), jsonNumber(3)},
	}, {
		a:    []any{jsonNumber(1)},
		want: jsonArray{jsonNumber(1)},
	}, {
		a:    jsonArray{},
		want: jsonArray{},
	}, {
		a:    map[string]any{"foo": 1},
		want: jsonObject{"foo": jsonNumber(1)},
	}, {
		a:    map[string]any{"foo": jsonNumber(1)},
		want: jsonObject{"foo": jsonNumber(1)},
	}, {
		a:    jsonObject{},
		want: jsonObject{},
	}}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, err := NewJsonNode(c.a)
			if c.wantErr {
				require.Error(t, err)
				require.Nil(t, b)
			} else {
				require.NoError(t, err)
				require.True(t, c.want.Equals(b))
			}
		})
	}
}

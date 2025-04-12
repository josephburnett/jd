package jd

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptionJSON(t *testing.T) {
	cases := []struct {
		json   string
		option Option
	}{{
		json:   `"MERGE"`,
		option: MERGE,
	}, {
		json:   `"SET"`,
		option: SET,
	}, {
		json:   `"MULTISET"`,
		option: MULTISET,
	}, {
		json:   `"COLOR"`,
		option: COLOR,
	}, {
		json:   `{"precision":1.01}`,
		option: Precision(1.01),
	}, {
		json:   `{"setkeys":["foo","bar"]}`,
		option: SetKeys("foo", "bar"),
	}, {
		json:   `{"at":["foo"],"then":["SET"]}`,
		option: PathOption(Path{PathKey("foo")}, SET),
	}, {
		json:   `{"at":["foo"],"then":[{"at":["bar"],"then":["SET"]}]}`,
		option: PathOption(Path{PathKey("foo")}, PathOption(Path{PathKey("bar")}, SET)),
	}}
	for _, c := range cases {
		t.Run(c.json, func(t *testing.T) {
			o, err := unmarshalOption([]byte(c.json))
			require.NoError(t, err)
			s, err := json.Marshal(o)
			require.NoError(t, err)
			require.Equal(t, c.json, string(s))
			s, err = json.Marshal(c.option)
			require.NoError(t, err)
			require.Equal(t, c.json, string(s))
		})
	}
}

func TestRecurse(t *testing.T) {
	cases := []struct {
		name      string
		opts      []Option
		element   PathElement
		wantApply []Option
		wantRest  []Option
	}{{
		name:      "recurse applies and keeps gloval options",
		opts:      []Option{SET, COLOR},
		element:   PathKey("foo"),
		wantApply: []Option{SET, COLOR},
		wantRest:  []Option{SET, COLOR},
	}, {
		name:      "recurse consumes on element of path",
		opts:      []Option{PathOption(Path{PathKey("foo"), PathKey("bar")}, SET)},
		element:   PathKey("foo"),
		wantApply: nil,
		wantRest:  []Option{PathOption(Path{PathKey("bar")}, SET)},
	}, {
		name:      "path option ending in set",
		opts:      []Option{PathOption(Path{PathKey("foo"), PathSet{}})},
		element:   PathKey("foo"),
		wantApply: []Option{SET},
		wantRest:  nil,
	}, {
		name:      "path option delivering set in payload",
		opts:      []Option{PathOption(Path{PathKey("foo")}, SET)},
		element:   PathKey("foo"),
		wantApply: []Option{SET},
		wantRest:  nil,
	}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			apply, rest := recurse(c.element, c.opts...)
			require.Equal(t, c.wantApply, apply)
			require.Equal(t, c.wantRest, rest)
		})
	}
}

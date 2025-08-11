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
		json:   `["MERGE"]`,
		option: MERGE,
	}, {
		json:   `["SET"]`,
		option: SET,
	}, {
		json:   `["MULTISET"]`,
		option: MULTISET,
	}, {
		json:   `["COLOR"]`,
		option: COLOR,
	}, {
		json:   `[{"precision":1.01}]`,
		option: Precision(1.01),
	}, {
		json:   `[{"setkeys":["foo","bar"]}]`,
		option: SetKeys("foo", "bar"),
	}, {
		json:   `[{"@":["foo"],"^":["SET"]}]`,
		option: PathOption(Path{PathKey("foo")}, SET),
	}, {
		json:   `[{"@":["foo"],"^":[{"@":["bar"],"^":["SET"]}]}]`,
		option: PathOption(Path{PathKey("foo")}, PathOption(Path{PathKey("bar")}, SET)),
	}}
	for _, c := range cases {
		t.Run(c.json, func(t *testing.T) {
			opts, err := ReadOptionsString(c.json)
			require.NoError(t, err)
			s, err := json.Marshal(opts)
			require.NoError(t, err)
			require.Equal(t, c.json, string(s))
			gotOpts := []any{c.option}
			s, err = json.Marshal(gotOpts)
			require.NoError(t, err)
			require.Equal(t, c.json, string(s))
		})
	}
}

func TestRefine(t *testing.T) {
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
			o := refine(&options{retain: c.opts}, c.element)
			require.Equal(t, c.wantApply, o.apply)
			require.Equal(t, c.wantRest, o.retain)
		})
	}
}

func TestPathOption(t *testing.T) {
	cases := []struct {
		name string
		opts string
		a, b string
	}{{
		name: "Precision on a number in an object",
		opts: `[{"@":["foo"],"^":[{"precision":0.1}]}]`,
		a:    `{"foo":1.0}`,
		b:    `{"foo":1.001}`,
	}, {
		name: "Precision on a number alone",
		opts: `[{"@":[],"^":[{"precision":0.1}]}]`,
		a:    `1.0`,
		b:    `1.001`,
		// }, {
		// 	name: "Precision in a list index",
		// 	opts: `[{"@":[0],"^":[{"precision":0.1}]}]`,
		// 	a:    `[1.0]`,
		// 	b:    `[1.001]`,
	}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			controlA, err := ReadJsonString(c.a)
			require.NoError(t, err)
			controlB, err := ReadJsonString(c.b)
			require.NoError(t, err)
			controlDiff := controlA.Diff(controlB)
			require.NotEmpty(t, controlDiff)

			a, err := ReadJsonString(c.a)
			require.NoError(t, err)
			b, err := ReadJsonString(c.b)
			require.NoError(t, err)
			o, err := ReadOptionsString(c.opts)
			require.NoError(t, err)
			diff := a.Diff(b, o...)
			require.Empty(t, diff)
		})
	}
}

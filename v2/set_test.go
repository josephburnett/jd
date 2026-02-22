package jd

import (
	"testing"
)

func TestSetJson(t *testing.T) {
	cases := []struct {
		name  string
		given string
		want  string
	}{{
		name:  "array with no space",
		given: `[]`,
		want:  `[]`,
	}, {
		name:  "array with space",
		given: ` [ ] `,
		want:  `[]`,
	}, {
		name:  "array with numbers out of order",
		given: `[2,1,3]`,
		want:  `[3,2,1]`,
	}, {
		name:  "array with numbers in order",
		given: `[3,2,1]`,
		want:  `[3,2,1]`,
	}, {
		name:  "array with spaced numbers",
		given: ` [1, 2, 3] `,
		want:  `[3,2,1]`,
	}, {
		name:  "duplicate entries",
		given: `[1,1,1]`,
		want:  `[1]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withOptions(setOption{})
			checkJson(ctx, c.given, c.want)
		})
	}
}

func TestSetEquals(t *testing.T) {
	cases := []struct {
		name string
		a    string
		b    string
	}{{
		name: "empty arrays",
		a:    `[]`,
		b:    `[]`,
	}, {
		name: "array with numbers unordered 1",
		a:    `[1,2,3]`,
		b:    `[3,2,1]`,
	}, {
		name: "array with numbers unordered 2",
		a:    `[1,2,3]`,
		b:    `[2,3,1]`,
	}, {
		name: "array with numbers unordered 3",
		a:    `[1,2,3]`,
		b:    `[1,3,2]`,
	}, {
		name: "array with empty objects",
		a:    `[{},{}]`,
		b:    `[{},{}]`,
	}, {
		name: "nested unordered arrays",
		a:    `[[1,2],[3,4]]`,
		b:    `[[2,1],[4,3]]`,
	}, {
		name: "repeated numbers",
		a:    `[1,1,1]`,
		b:    `[1]`,
	}, {
		name: "array with numbers repeated and unordered",
		a:    `[1,2,1]`,
		b:    `[2,1,2]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withOptions(setOption{})
			checkEqual(ctx, c.a, c.b)
		})
	}
}

func TestSetNotEquals(t *testing.T) {
	cases := []struct {
		name string
		a    string
		b    string
	}{{
		name: "empty and non-empty sets",
		a:    `[]`,
		b:    `[1]`,
	}, {
		name: "sets with unique and repeated elements",
		a:    `[1,2,3]`,
		b:    `[1,2,2]`,
	}, {
		name: "sets of different sizes",
		a:    `[1,2,3]`,
		b:    `[1,2]`,
	}, {
		name: "nested sets with different elements",
		a:    `[[],[1]]`,
		b:    `[[],[2]]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withOptions(setOption{})
			checkNotEqual(ctx, c.a, c.b)
		})
	}
}

func TestSetDiff(t *testing.T) {
	cases := []struct {
		name    string
		options []Option
		a       string
		b       string
		want    []string
	}{{
		name:    "empty sets no diff",
		options: m(setOption{}),
		a:       `[]`,
		b:       `[]`,
		want:    ss(),
	}, {
		name:    "add a number",
		options: m(setOption{}),
		a:       `[1]`,
		b:       `[1,2]`,
		want: ss(
			`@ [{}]`,
			`+ 2`,
		),
	}, {
		name:    "sets with same numbers",
		options: m(setOption{}),
		a:       `[1,2]`,
		b:       `[1,2]`,
		want:    ss(),
	}, {
		name:    "add a number multiple times",
		options: m(setOption{}),
		a:       `[1]`,
		b:       `[1,2,2]`,
		want: ss(
			`@ [{}]`,
			`+ 2`,
		),
	}, {
		name:    "remove a number",
		options: m(setOption{}),
		a:       `[1,2,3]`,
		b:       `[1,3]`,
		want: ss(
			`@ [{}]`,
			`- 2`,
		),
	}, {
		name:    "replace one object with another",
		options: m(setOption{}),
		a:       `[{"a":1}]`,
		b:       `[{"a":2}]`,
		want: ss(
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
	}, {
		name:    "replace one repeated object with another",
		options: m(setOption{}),
		a:       `[{"a":1},{"a":1}]`,
		b:       `[{"a":2}]`,
		want: ss(
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
	}, {
		name:    "remove two strings and add one",
		options: m(setOption{}),
		a:       `["foo","foo","bar"]`,
		b:       `["baz"]`,
		want: ss(
			`@ [{}]`,
			`- "bar"`,
			`- "foo"`,
			`+ "baz"`,
		),
	}, {
		name:    "remove one string and add two repeated",
		options: m(setOption{}),
		a:       `["foo"]`,
		b:       `["bar","baz","bar"]`,
		want: ss(
			`@ [{}]`,
			`- "foo"`,
			`+ "bar"`,
			`+ "baz"`,
		),
	}, {
		name:    "remove object and add array",
		options: m(setOption{}),
		a:       `{}`,
		b:       `[]`,
		want: ss(
			`@ []`,
			`- {}`,
			`+ []`,
		),
	}, {
		name: "add property to object in set",
		options: m(
			setOption{},
			SetKeys("id"),
		),
		a: `[{"id":"foo"}]`,
		b: `[{"id":"foo","bar":"baz"}]`,
		want: ss(
			`@ [{"id":"foo"},"bar"]`,
			`+ "baz"`,
		),
	}, {
		name: "find object by id among empty objects",
		options: m(
			setOption{},
			SetKeys("id"),
		),
		a: `[{},{},{"id":"foo"},{}]`,
		b: `[{},{"id":"foo","bar":"baz"},{},{}]`,
		want: ss(
			`@ [{"id":"foo"},"bar"]`,
			`+ "baz"`,
		),
	}, {
		name: "find object by multiple ids",
		options: m(
			setOption{},
			SetKeys("id1", "id2"),
		),
		a: `[{},{"id1":"foo","id2":"zap"},{}]`,
		b: `[{},{"id1":"foo","id2":"zap","bar":"baz"},{}]`,
		want: ss(
			`@ [{"id1":"foo","id2":"zap"},"bar"]`,
			`+ "baz"`,
		),
	}, {
		name: "find object by id among others",
		options: m(
			setOption{},
			SetKeys("id"),
		),
		a: `[{"id":"foo"},{"id":"bar"}]`,
		b: `[{"id":"foo","baz":"zap"},{"id":"bar"}]`,
		want: ss(
			`@ [{"id":"foo"},"baz"]`,
			`+ "zap"`,
		),
	}, {
		name: "two objects with different ids being exchanged",
		options: m(
			setOption{},
			SetKeys("id"),
		),
		a: `[{"id":"foo"}]`,
		b: `[{"id":"bar"}]`,
		want: ss(
			`@ [{}]`,
			`- {"id":"foo"}`,
			`+ {"id":"bar"}`,
		),
	}, {
		name:    "set options applies to array in object",
		options: m(setOption{}),
		a:       `{"a":[1,2]}`,
		b:       `{"a":[2,1]}`,
		want:    ss(),
	}, {
		name:    "merge different types produces only new values",
		options: m(MERGE, setOption{}),
		a:       `[1,2,3]`,
		b:       `{}`,
		want: ss(
			`^ {"Merge":true}`,
			`@ []`,
			`+ {}`,
		),
	}, {
		name:    "merge outputs no diff when equal",
		options: m(MERGE, setOption{}),
		a:       `[1,2,3]`,
		b:       `[2,1,3]`,
		want:    ss(),
	}, {
		name:    "merge replaces entire set when not equal",
		options: m(MERGE, setOption{}),
		a:       `[1,2,3]`,
		b:       `[2,1,4]`,
		want: ss(
			`^ {"Merge":true}`,
			`@ []`,
			`+ [2,1,4]`,
		),
	}, {
		name: "pod set diff regression", // github.com/josephburnett/jd/issues/82
		options: m(
			SET,
			SetKeys("name"),
		),
		a: `{"apiVersion":"v1","kind":"Pod","metadata":{"name":"nginx"},"spec":{"containers":[{"name":"nginx","image":"nginx:1.14.2","ports":[{"containerPort":80}]}]}}`,
		b: `{"apiVersion":"v1","kind":"Pod","metadata":{"name":"nginx"},"spec":{"containers":[{"name":"nginx","image":"nginx:1.14.2","ports":[{"containerPort":8080}]}]}}`,
		want: ss(
			`@ ["spec","containers",{"name":"nginx"},"ports",{}]`,
			`- {"containerPort":80}`,
			`+ {"containerPort":8080}`,
		),
	}, {
		name: "objects missing all setkeys fields are distinct",
		options: m(
			setOption{},
			SetKeys("id"),
		),
		a: `[]`,
		b: `[{},{"":"0"}]`,
		want: ss(
			`@ [{}]`,
			`+ {}`,
			`+ {"":"0"}`,
		),
	}, {
		name: "complex set key diff",
		options: m(
			SET,
			SetKeys("a", "b"),
		),
		a: `[{"a":1,"b":2,"c":3},{"a":1,"b":5,"c":3},{"a":1,"c":6}]`,
		b: `[{"a":1,"b":2,"c":4},{"a":1,"b":5,"c":4},{"a":1,"c":7}]`,
		want: ss(
			`@ [{"a":1,"b":5},"c"]`,
			`- 3`,
			`+ 4`,
			`@ [{"a":1,"b":2},"c"]`,
			`- 3`,
			`+ 4`,
			`@ [{"a":1,"b":null},"c"]`,
			`- 6`,
			`+ 7`,
		),
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withOptions(c.options...)
			checkDiff(ctx, c.a, c.b, c.want...)
		})
	}
}

func TestSetPatch(t *testing.T) {
	cases := []struct {
		name  string
		given string
		patch []string
		want  string
	}{{
		name:  "empty patch on empty set does nothing",
		given: `[]`,
		patch: ss(``),
		want:  `[]`,
	}, {
		name:  "add a number",
		given: `[1]`,
		patch: ss(
			`@ [{}]`,
			`+ 2`,
		),
		want: `[1,2]`,
	}, {
		name:  "empty patch on set with numbers does nothing",
		given: `[1,2]`,
		patch: ss(``),
		want:  `[1,2]`,
	}, {
		name:  "remove a number from a set",
		given: `[1,2,3]`,
		patch: ss(
			`@ [{}]`,
			`- 2`,
		),
		want: `[1,3]`,
	}, {
		name:  "replace one object with another",
		given: `[{"a":1}]`,
		patch: ss(
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
		want: `[{"a":2}]`,
	}, {
		name:  "replace one repeated object with another",
		given: `[{"a":1},{"a":1}]`,
		patch: ss(
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
		want: `[{"a":2}]`,
	}, {
		name:  "replace two strings with one string",
		given: `["foo","foo","bar"]`,
		patch: ss(
			`@ [{}]`,
			`- "bar"`,
			`- "foo"`,
			`+ "baz"`,
		),
		want: `["baz"]`,
	}, {
		name:  "replace one string with two strings",
		given: `["foo"]`,
		patch: ss(
			`@ [{}]`,
			`- "foo"`,
			`+ "bar"`,
			`+ "baz"`,
		),
		want: `["bar","baz","bar"]`,
	}, {
		name:  "replace object with array",
		given: `{}`,
		patch: ss(
			`@ []`,
			`- {}`,
			`+ []`,
		),
		want: `[]`,
	}, {
		name:  "patch property to object in set",
		given: `[{"id":"foo"}]`,
		patch: ss(
			`@ [{"id":"foo"},"bar"]`,
			`+ "baz"`,
		),
		want: `[{"id":"foo","bar":"baz"}]`,
	}, {
		name:  "patch object among empty objects",
		given: `[{},{},{"id":"foo"},{}]`,
		patch: ss(
			`@ [{"id":"foo"},"bar"]`,
			`+ "baz"`,
		),
		want: `[{},{"id":"foo","bar":"baz"},{},{}]`,
	}, {
		name:  "patch object by multiple ids",
		given: `[{},{"id1":"foo","id2":"zap"},{}]`,
		patch: ss(
			`@ [{"id1":"foo","id2":"zap"},"bar"]`,
			`+ "baz"`,
		),
		want: `[{},{"id1":"foo","id2":"zap","bar":"baz"},{}]`,
	}, {
		name:  "patch object by id among other",
		given: `[{"id":"foo"},{"id":"bar"}]`,
		patch: ss(
			`@ [{"id":"foo"},"baz"]`,
			`+ "zap"`,
		),
		want: `[{"id":"foo","baz":"zap"},{"id":"bar"}]`,
	}, {
		name:  "replace two objects with diffent ids",
		given: `[{"id":"foo"}]`,
		patch: ss(
			`@ [{}]`,
			`- {"id":"foo"}`,
			`+ {"id":"bar"}`,
		),
		want: `[{"id":"bar"}]`,
	}, {
		name:  "merge replaces entire set",
		given: `[1,2,3]`,
		patch: ss(
			`^ {"Merge":true}`,
			`@ [{}]`,
			`+ [4,5,6]`,
		),
		want: `[4,5,6]`,
	}, {
		name:  "void deletes a node",
		given: `[1,2,3]`,
		patch: ss(
			`^ {"Merge":true}`,
			`@ [{}]`,
			`+`,
		),
		want: ``,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withOptions(setOption{})
			checkPatch(ctx, c.given, c.want, c.patch...)
		})
	}
}

func TestSetPatchError(t *testing.T) {
	cases := []struct {
		name  string
		given string
		patch []string
	}{{
		name:  "removing number from empty set",
		given: `[]`,
		patch: ss(
			`@ [{}]`,
			`- 1`,
		),
	}, {
		name:  "removing number from set twice",
		given: `[1]`,
		patch: ss(
			`@ [{}]`,
			`- 1`,
			`- 1`,
		),
	}, {
		name:  "removing object from empty set",
		given: `[]`,
		patch: ss(
			`@ []`,
			`- {}`,
		),
	}, {
		name:  "removing number from set twice added twice",
		given: `[]`,
		patch: ss(
			`@ [{}]`,
			`+ 1`,
			`+ 1`,
			`@ [{}]`,
			`- 1`,
			`- 1`,
		),
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withOptions(setOption{})
			checkPatchError(ctx, c.given, c.patch...)
		})
	}
}

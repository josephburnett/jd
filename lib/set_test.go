package jd

import (
	"testing"
)

func TestSetJson(t *testing.T) {
	cases := []struct {
		name  string
		opt   RenderOption
		given string
		want  string
	}{{
		name:  "array with no space",
		opt:   SET,
		given: `[]`,
		want:  `[]`,
	}, {
		name:  "array with space",
		opt:   SET,
		given: ` [ ] `,
		want:  `[]`,
	}, {
		name:  "array with numbers out of order",
		opt:   SET,
		given: `[2,1,3]`,
		want:  `[3,2,1]`,
	}, {
		name:  "array with numbers in order",
		opt:   SET,
		given: `[3,2,1]`,
		want:  `[3,2,1]`,
	}, {
		name:  "array with spaced numbers",
		opt:   SET,
		given: ` [1, 2, 3] `,
		want:  `[3,2,1]`,
	}, {
		name:  "duplicate entries",
		opt:   SET,
		given: `[1,1,1]`,
		want:  `[1]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withRenderOption(c.opt)
			checkJson(ctx, c.given, c.want)
		})
	}
}

func TestSetEquals(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		a        string
		b        string
	}{{
		name:     "empty arrays",
		metadata: SET,
		a:        `[]`,
		b:        `[]`,
	}, {
		name:     "array with numbers unordered 1",
		metadata: SET,
		a:        `[1,2,3]`,
		b:        `[3,2,1]`,
	}, {
		name:     "array with numbers unordered 2",
		metadata: SET,
		a:        `[1,2,3]`,
		b:        `[2,3,1]`,
	}, {
		name:     "array with numbers unordered 3",
		metadata: SET,
		a:        `[1,2,3]`,
		b:        `[1,3,2]`,
	}, {
		name:     "array with empty objects",
		metadata: SET,
		a:        `[{},{}]`,
		b:        `[{},{}]`,
	}, {
		name:     "nested unordered arrays",
		metadata: SET,
		a:        `[[1,2],[3,4]]`,
		b:        `[[2,1],[4,3]]`,
	}, {
		name:     "repeated numbers",
		metadata: SET,
		a:        `[1,1,1]`,
		b:        `[1]`,
	}, {
		name:     "array with numbers repeated and unordered",
		metadata: SET,
		a:        `[1,2,1]`,
		b:        `[2,1,2]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withMetadata(c.metadata)
			checkEqual(ctx, c.a, c.b)
		})
	}
}

func TestSetNotEquals(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		a        string
		b        string
	}{{
		name:     "empty and non-empty sets",
		metadata: SET,
		a:        `[]`,
		b:        `[1]`,
	}, {
		name:     "sets with unique and repeated elements",
		metadata: SET,
		a:        `[1,2,3]`,
		b:        `[1,2,2]`,
	}, {
		name:     "sets of different sizes",
		metadata: SET,
		a:        `[1,2,3]`,
		b:        `[1,2]`,
	}, {
		name:     "nested sets with different elements",
		metadata: SET,
		a:        `[[],[1]]`,
		b:        `[[],[2]]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withMetadata(c.metadata)
			checkNotEqual(ctx, c.a, c.b)
		})
	}
}

func TestSetDiff(t *testing.T) {
	cases := []struct {
		name     string
		metadata []Metadata
		a        string
		b        string
		want     []string
	}{{
		name:     "empty sets no diff",
		metadata: m(SET),
		a:        `[]`,
		b:        `[]`,
		want:     ss(),
	}, {
		name:     "add a number",
		metadata: m(SET),
		a:        `[1]`,
		b:        `[1,2]`,
		want: ss(
			`@ [["set"],{}]`,
			`+ 2`,
		),
	}, {
		name:     "sets with same numbers",
		metadata: m(SET),
		a:        `[1,2]`,
		b:        `[1,2]`,
		want:     ss(),
	}, {
		name:     "add a number multiple times",
		metadata: m(SET),
		a:        `[1]`,
		b:        `[1,2,2]`,
		want: ss(
			`@ [["set"],{}]`,
			`+ 2`,
		),
	}, {
		name:     "remove a number",
		metadata: m(SET),
		a:        `[1,2,3]`,
		b:        `[1,3]`,
		want: ss(
			`@ [["set"],{}]`,
			`- 2`,
		),
	}, {
		name:     "replace one object with another",
		metadata: m(SET),
		a:        `[{"a":1}]`,
		b:        `[{"a":2}]`,
		want: ss(
			`@ [["set"],{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
	}, {
		name:     "replace one repeated object with another",
		metadata: m(SET),
		a:        `[{"a":1},{"a":1}]`,
		b:        `[{"a":2}]`,
		want: ss(
			`@ [["set"],{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
	}, {
		name:     "remove two strings and add one",
		metadata: m(SET),
		a:        `["foo","foo","bar"]`,
		b:        `["baz"]`,
		want: ss(
			`@ [["set"],{}]`,
			`- "bar"`,
			`- "foo"`,
			`+ "baz"`,
		),
	}, {
		name:     "remove one string and add two repeated",
		metadata: m(SET),
		a:        `["foo"]`,
		b:        `["bar","baz","bar"]`,
		want: ss(
			`@ [["set"],{}]`,
			`- "foo"`,
			`+ "bar"`,
			`+ "baz"`,
		),
	}, {
		name:     "remove object and add array",
		metadata: m(SET),
		a:        `{}`,
		b:        `[]`,
		want: ss(
			`@ []`,
			`- {}`,
			`+ []`,
		),
	}, {
		name: "add property to object in set",
		metadata: m(
			SET,
			Setkeys("id"),
		),
		a: `[{"id":"foo"}]`,
		b: `[{"id":"foo","bar":"baz"}]`,
		want: ss(
			`@ [["set","setkeys=id"],{"id":"foo"},"bar"]`,
			`+ "baz"`,
		),
	}, {
		name: "find object by id among empty objects",
		metadata: m(
			SET,
			Setkeys("id"),
		),
		a: `[{},{},{"id":"foo"},{}]`,
		b: `[{},{"id":"foo","bar":"baz"},{},{}]`,
		want: ss(
			`@ [["set","setkeys=id"],{"id":"foo"},"bar"]`,
			`+ "baz"`,
		),
	}, {
		name: "find object by multiple ids",
		metadata: m(
			SET,
			Setkeys("id1", "id2"),
		),
		a: `[{},{"id1":"foo","id2":"zap"},{}]`,
		b: `[{},{"id1":"foo","id2":"zap","bar":"baz"},{}]`,
		want: ss(
			`@ [["set","setkeys=id1,id2"],{"id1":"foo","id2":"zap"},"bar"]`,
			`+ "baz"`,
		),
	}, {
		name: "find object by id among others",
		metadata: m(
			SET,
			Setkeys("id"),
		),
		a: `[{"id":"foo"},{"id":"bar"}]`,
		b: `[{"id":"foo","baz":"zap"},{"id":"bar"}]`,
		want: ss(
			`@ [["set","setkeys=id"],{"id":"foo"},"baz"]`,
			`+ "zap"`,
		),
	}, {
		name: "two objects with different ids being exchanged",
		metadata: m(
			SET,
			Setkeys("id"),
		),
		a: `[{"id":"foo"}]`,
		b: `[{"id":"bar"}]`,
		want: ss(
			`@ [["set","setkeys=id"],{}]`,
			`- {"id":"foo"}`,
			`+ {"id":"bar"}`,
		),
	}, {
		name:     "set metadata applies to array in object",
		metadata: m(SET),
		a:        `{"a":[1,2]}`,
		b:        `{"a":[2,1]}`,
		want:     ss(),
	}, {
		name:     "merge different types produces only new values",
		metadata: m(MERGE, SET),
		a:        `[1,2,3]`,
		b:        `{}`,
		want: ss(
			`@ [["MERGE"]]`,
			`+ {}`,
		),
	}, {
		name:     "merge outputs no diff when equal",
		metadata: m(MERGE, SET),
		a:        `[1,2,3]`,
		b:        `[2,1,3]`,
		want:     ss(),
	}, {
		name:     "merge replaces entire set when not equal",
		metadata: m(MERGE, SET),
		a:        `[1,2,3]`,
		b:        `[2,1,4]`,
		want: ss(
			`@ [["MERGE"]]`,
			`+ [2,1,4]`,
		),
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withMetadata(c.metadata...)
			checkDiff(ctx, c.a, c.b, c.want...)
		})
	}
}

func TestSetPatch(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		given    string
		patch    []string
		want     string
	}{{
		name:     "empty patch on empty set does nothing",
		metadata: SET,
		given:    `[]`,
		patch:    ss(``),
		want:     `[]`,
	}, {
		name:     "add a number",
		metadata: SET,
		given:    `[1]`,
		patch: ss(
			`@ [{}]`,
			`+ 2`,
		),
		want: `[1,2]`,
	}, {
		name:     "empty patch on set with numbers does nothing",
		metadata: SET,
		given:    `[1,2]`,
		patch:    ss(``),
		want:     `[1,2]`,
	}, {
		name:     "remove a number from a set",
		metadata: SET,
		given:    `[1,2,3]`,
		patch: ss(
			`@ [{}]`,
			`- 2`,
		),
		want: `[1,3]`,
	}, {
		name:     "replace one object with another",
		metadata: SET,
		given:    `[{"a":1}]`,
		patch: ss(
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
		want: `[{"a":2}]`,
	}, {
		name:     "replace one repeated object with another",
		metadata: SET,
		given:    `[{"a":1},{"a":1}]`,
		patch: ss(
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
		want: `[{"a":2}]`,
	}, {
		name:     "replace two strings with one string",
		metadata: SET,
		given:    `["foo","foo","bar"]`,
		patch: ss(
			`@ [{}]`,
			`- "bar"`,
			`- "foo"`,
			`+ "baz"`,
		),
		want: `["baz"]`,
	}, {
		name:     "replace one string with two strings",
		metadata: SET,
		given:    `["foo"]`,
		patch: ss(
			`@ [{}]`,
			`- "foo"`,
			`+ "bar"`,
			`+ "baz"`,
		),
		want: `["bar","baz","bar"]`,
	}, {
		name:     "replace object with array",
		metadata: SET,
		given:    `{}`,
		patch: ss(
			`@ []`,
			`- {}`,
			`+ []`,
		),
		want: `[]`,
	}, {
		name:     "patch property to object in set",
		metadata: SET,
		given:    `[{"id":"foo"}]`,
		patch: ss(
			`@ [["set"],{"id":"foo"},"bar"]`,
			`+ "baz"`,
		),
		want: `[{"id":"foo","bar":"baz"}]`,
	}, {
		name:     "patch object among empty objects",
		metadata: SET,
		given:    `[{},{},{"id":"foo"},{}]`,
		patch: ss(
			`@ [["set","setkeys=id"],{"id":"foo"},"bar"]`,
			`+ "baz"`,
		),
		want: `[{},{"id":"foo","bar":"baz"},{},{}]`,
	}, {
		name:     "patch object by multiple ids",
		metadata: SET,
		given:    `[{},{"id1":"foo","id2":"zap"},{}]`,
		patch: ss(
			`@ [["set","setkeys=id1,id2"],{"id1":"foo","id2":"zap"},"bar"]`,
			`+ "baz"`,
		),
		want: `[{},{"id1":"foo","id2":"zap","bar":"baz"},{}]`,
	}, {
		name:     "patch object by id among other",
		metadata: SET,
		given:    `[{"id":"foo"},{"id":"bar"}]`,
		patch: ss(
			`@ [["set","setkeys=id"],{"id":"foo"},"baz"]`,
			`+ "zap"`,
		),
		want: `[{"id":"foo","baz":"zap"},{"id":"bar"}]`,
	}, {
		name:     "replace two objects with diffent ids",
		metadata: SET,
		given:    `[{"id":"foo"}]`,
		patch: ss(
			`@ [["set","setkeys=id"],{}]`,
			`- {"id":"foo"}`,
			`+ {"id":"bar"}`,
		),
		want: `[{"id":"bar"}]`,
	}, {
		name:  "merge replaces entire set",
		given: `[1,2,3]`,
		patch: ss(
			`@ [["MERGE","set"]]`,
			`+ [4,5,6]`,
		),
		want: `[4,5,6]`,
	}, {
		name:  "void deletes a node",
		given: `[1,2,3]`,
		patch: ss(
			`@ [["MERGE","set"]]`,
			`+`,
		),
		want: ``,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withMetadata(c.metadata)
			checkPatch(ctx, c.given, c.want, c.patch...)
		})
	}
}

func TestSetPatchError(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		given    string
		patch    []string
	}{{
		name:     "removing number from empty set",
		metadata: SET,
		given:    `[]`,
		patch: ss(
			`@ [{}]`,
			`- 1`,
		),
	}, {
		name:     "removing number from set twice",
		metadata: SET,
		given:    `[1]`,
		patch: ss(
			`@ [{}]`,
			`- 1`,
			`- 1`,
		),
	}, {
		name:     "removing object from empty set",
		metadata: SET,
		given:    `[]`,
		patch: ss(
			`@ []`,
			`- {}`,
		),
	}, {
		name:     "removing number from set twice added twice",
		metadata: SET,
		given:    `[]`,
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
				withMetadata(c.metadata)
			checkPatchError(ctx, c.given, c.patch...)
		})
	}
}

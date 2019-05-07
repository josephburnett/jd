package jd

import (
	"testing"
)

func TestSetJson(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		given    string
		want     string
	}{{
		name:     "array with no space",
		metadata: SET,
		given:    `[]`,
		want:     `[]`,
	}, {
		name:     "array with space",
		metadata: SET,
		given:    ` [ ] `,
		want:     `[]`,
	}, {
		name:     "array with numbers out of order",
		metadata: SET,
		given:    `[2,1,3]`,
		want:     `[3,2,1]`,
	}, {
		name:     "array with numbers in order",
		metadata: SET,
		given:    `[3,2,1]`,
		want:     `[3,2,1]`,
	}, {
		name:     "array with spaced numbers",
		metadata: SET,
		given:    ` [1, 2, 3] `,
		want:     `[3,2,1]`,
	}, {
		name:     "duplicate entries",
		metadata: SET,
		given:    `[1,1,1]`,
		want:     `[1]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withReadMetadata(c.metadata)
			checkJson(ctx, c.given, c.want)
			ctx = newTestContext(t).
				withApplyMetadata(c.metadata)
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
				withReadMetadata(c.metadata)
			checkEqual(ctx, c.a, c.b)
			ctx = newTestContext(t).
				withApplyMetadata(c.metadata)
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
				withReadMetadata(c.metadata)
			checkNotEqual(ctx, c.a, c.b)
			ctx = newTestContext(t).
				withApplyMetadata(c.metadata)
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
		want:     s(``),
	}, {
		name:     "add a number",
		metadata: m(SET),
		a:        `[1]`,
		b:        `[1,2]`,
		want: s(
			`@ [{}]`,
			`+ 2`,
		),
	}, {
		name:     "sets with same numbers",
		metadata: m(SET),
		a:        `[1,2]`,
		b:        `[1,2]`,
		want:     s(``),
	}, {
		name:     "add a number multiple times",
		metadata: m(SET),
		a:        `[1]`,
		b:        `[1,2,2]`,
		want: s(
			`@ [{}]`,
			`+ 2`,
		),
	}, {
		name:     "remove a number",
		metadata: m(SET),
		a:        `[1,2,3]`,
		b:        `[1,3]`,
		want: s(
			`@ [{}]`,
			`- 2`,
		),
	}, {
		name:     "replace one object with another",
		metadata: m(SET),
		a:        `[{"a":1}]`,
		b:        `[{"a":2}]`,
		want: s(
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
	}, {
		name:     "replace one repeated object with another",
		metadata: m(SET),
		a:        `[{"a":1},{"a":1}]`,
		b:        `[{"a":2}]`,
		want: s(
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
	}, {
		name:     "remove two strings and add one",
		metadata: m(SET),
		a:        `["foo","foo","bar"]`,
		b:        `["baz"]`,
		want: s(
			`@ [{}]`,
			`- "bar"`,
			`- "foo"`,
			`+ "baz"`,
		),
	}, {
		name:     "remove one string and add two repeated",
		metadata: m(SET),
		a:        `["foo"]`,
		b:        `["bar","baz","bar"]`,
		want: s(
			`@ [{}]`,
			`- "foo"`,
			`+ "bar"`,
			`+ "baz"`,
		),
	}, {
		name:     "remove object and add array",
		metadata: m(SET),
		a:        `{}`,
		b:        `[]`,
		want: s(
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
		want: s(
			`@ [{"id":"foo"}, "bar"]`,
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
		want: s(
			`@ [{"id":"foo"}, "bar"]`,
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
		want: s(
			`@ [{"id1":"foo","id2":"zap"}, "bar"]`,
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
		want: s(
			`@ [{"id":"foo"}, "baz"]`,
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
		want: s(
			// TODO: emit set keys as path metadata
			//       e.g. `@ [["setkeys=id"],{}]`
			//       so that diffs will be self-describing.
			`@ [{}]`,
			`- {"id":"foo"}`,
			`+ {"id":"bar"}`,
		),
	}, {
		name:     "set metadata applies to array in object",
		metadata: m(SET),
		a:        `{"a":[1,2]}`,
		b:        `{"a":[2,1]}`,
		want:     s(``),
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// ctx := newTestContext(t).
			// 	withReadMetadata(c.metadata...)
			// checkDiff(ctx, c.a, c.b, c.want...)
			ctx := newTestContext(t).
				withApplyMetadata(c.metadata...)
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
		patch:    s(``),
		want:     `[]`,
	}, {
		name:     "add a number",
		metadata: SET,
		given:    `[1]`,
		patch: s(
			`@ [{}]`,
			`+ 2`,
		),
		want: `[1,2]`,
	}, {
		name:     "empty patch on set with numbers does nothing",
		metadata: SET,
		given:    `[1,2]`,
		patch:    s(``),
		want:     `[1,2]`,
	}, {
		name:     "remove a number from a set",
		metadata: SET,
		given:    `[1,2,3]`,
		patch: s(
			`@ [{}]`,
			`- 2`,
		),
		want: `[1,3]`,
	}, {
		name:     "replace one object with another",
		metadata: SET,
		given:    `[{"a":1}]`,
		patch: s(
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
		want: `[{"a":2}]`,
	}, {
		name:     "replace one repeated object with another",
		metadata: SET,
		given:    `[{"a":1},{"a":1}]`,
		patch: s(
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
		want: `[{"a":2}]`,
	}, {
		name:     "replace two strings with one string",
		metadata: SET,
		given:    `["foo","foo","bar"]`,
		patch: s(
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
		patch: s(
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
		patch: s(
			`@ []`,
			`- {}`,
			`+ []`,
		),
		want: `[]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withReadMetadata(c.metadata)
			checkPatch(ctx, c.given, c.want, c.patch...)
			// TODO: implement set patch with metadata.
			// ctx = newTestContext(t).
			// 	withApplyMetadata(c.metadata)
			// checkPatch(ctx, c.given, c.want, c.patch...)
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
		patch: s(
			`@ [{}]`,
			`- 1`,
		),
	}, {
		name:     "removing number from set twice",
		metadata: SET,
		given:    `[1]`,
		patch: s(
			`@ [{}]`,
			`- 1`,
			`- 1`,
		),
	}, {
		name:     "removing object from empty set",
		metadata: SET,
		given:    `[]`,
		patch: s(
			`@ []`,
			`- {}`,
		),
	}, {
		name:     "removing number from set twice added twice",
		metadata: SET,
		given:    `[]`,
		patch: s(
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
				withReadMetadata(c.metadata)
			checkPatchError(ctx, c.given, c.patch...)
			ctx = newTestContext(t).
				withApplyMetadata(c.metadata)
			checkPatchError(ctx, c.given, c.patch...)
		})
	}
}

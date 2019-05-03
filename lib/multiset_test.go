package jd

import (
	"testing"
)

func TestMultisetJson(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		given    string
		want     string
	}{{
		name:     "empty mulitset",
		metadata: MULTISET,
		given:    `[]`,
		want:     `[]`,
	}, {
		name:     "empty multiset with space",
		metadata: MULTISET,
		given:    ` [ ] `,
		want:     `[]`,
	}, {
		name:     "ordered multiset",
		metadata: MULTISET,
		given:    `[1,2,3]`,
		want:     `[1,2,3]`,
	}, {
		name:     "ordered multiset with space",
		metadata: MULTISET,
		given:    ` [1, 2, 3] `,
		want:     `[1,2,3]`,
	}, {
		name:     "multset with multiple duplicates",
		metadata: MULTISET,
		given:    `[1,1,1]`,
		want:     `[1,1,1]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withReadMetadata(c.metadata)
			checkJson(ctx, c.given, c.want)
			ctx := newTestContext(t).
				withApplyMetadata(c.metadata)
			checkJson(ctx, c.given, c.want)
		})
	}
}

func TestMultisetEquals(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		a        string
		b        string
	}{{
		name:     "empty multisets",
		metadata: MULTISET,
		a:        `[]`,
		b:        `[]`,
	}, {
		name:     "different ordered multisets 1",
		metadata: MULTISET,
		a:        `[1,2,3]`,
		b:        `[3,2,1]`,
	}, {
		name:     "different ordered multisets 2",
		metadata: MULTISET,
		a:        `[1,2,3]`,
		b:        `[2,3,1]`,
	}, {
		name:     "different ordered multisets 2",
		metadata: MULTISET,
		a:        `[1,2,3]`,
		b:        `[1,3,2]`,
	}, {
		name:     "multsets with empty objects",
		metadata: MULTISET,
		a:        `[{},{}]`,
		b:        `[{},{}]`,
	}, {
		name:     "nested multisets",
		metadata: MULTISET,
		a:        `[[1,2],[3,4]]`,
		b:        `[[2,1],[4,3]]`,
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

func TestMultisetNotEquals(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		a        string
		b        string
	}{{
		name:     "empty multiset and multiset with number",
		metadata: MULTISET,
		a:        `[]`,
		b:        `[1]`,
	}, {
		name:     "multisets with different numbers",
		metadata: MULTISET,
		a:        `[1,2,3]`,
		b:        `[1,2,2]`,
	}, {
		name:     "multiset missing a number",
		metadata: MULTISET,
		a:        `[1,2,3]`,
		b:        `[1,2]`,
	}, {
		name:     "nested multisets with different numbers",
		a:        `[[],[1]]`,
		b:        `[[],[2]]`,
		metadata: MULTISET,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withReadMetadata(c.metadata)
			checkNotEqual(ctx, c.a, c.b)
			ctx := newTestContext(t).
				withApplyMetadata(c.metadata)
			checkNotEqual(ctx, c.a, c.b)
		})
	}
}

func TestMultisetDiff(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		a        string
		b        string
		want     []string
	}{{
		name:     "two empty multisets",
		metadata: MULTISET,
		a:        `[]`,
		b:        `[]`,
		want:     {``},
	}, {
		name:     "two multisets with different numbers",
		metadata: MULTISET,
		a:        `[1]`,
		b:        `[1,2]`,
		want: {
			`@ q[{}]`,
			`+ 2`,
		},
	}, {
		name:     "two multisets with the same number",
		metadata: MULTISET,
		a:        `[1,2]`,
		b:        `[1,2]`,
		want:     {``},
	}, {
		name:     "adding two numbers",
		metadata: MULTISET,
		a:        `[1]`,
		b:        `[1,2,2]`,
		want: {
			`@ [{}]`,
			`+ 2`,
			`+ 2`,
		},
	}, {
		name:     "removing a number",
		metadata: MULTISET,
		a:        `[1,2,3]`,
		b:        `[1,3]`,
		want: {
			`@ [{}]`,
			`- 2`,
		},
	}, {
		name:     "replacing one object with another",
		metadata: MULTISET,
		a:        `[{"a":1}]`,
		b:        `[{"a":2}]`,
		want: {
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		},
	}, {
		name:     "replacing two objects with one object",
		metadata: MULTISET,
		a:        `[{"a":1},{"a":1}]`,
		b:        `[{"a":2}]`,
		want: {
			`@ [{}]`,
			`- {"a":1}`,
			`- {"a":1}`,
			`+ {"a":2}`,
		},
	}, {
		name:     "replacing three strings repeated with one string",
		metadata: MULTISET,
		a:        `["foo","foo","bar"]`,
		b:        `["baz"]`,
		want: {
			`@ [{}]`,
			`- "bar"`,
			`- "foo"`,
			`- "foo"`,
			`+ "baz"`,
		},
	}, {
		name:     "replacing one string with three repeated",
		metadata: MULTISET,
		a:        `["foo"]`,
		b:        `["bar","baz","bar"]`,
		want: {
			`@ [{}]`,
			`- "foo"`,
			`+ "bar"`,
			`+ "bar"`,
			`+ "baz"`,
		},
	}, {
		name:     "replacing multiset with array",
		metadata: MULTISET,
		a:        `{}`,
		b:        `[]`,
		want: {
			`@ []`,
			`- {}`,
			`+ []`,
		},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withReadMetadata(c.metadata)
			checkDiff(ctx, c.a, c.b, c.want...)
			ctx = newTestContext(t).
				withApplyMetadata(c.metadata)
			checkDiff(ctx, c.a, c.b, c.want...)
		})
	}
}

func TestMultisetPatch(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		given    string
		patch    []string
		want     string
	}{{
		name:     "empty patch on empty multiset",
		metadata: MULTISET,
		given:    `[]`,
		patch:    {``},
		want:     `[]`,
	}, {
		name:     "add a number",
		metadata: MULTISET,
		given:    `[1]`,
		patch: {
			`@ [{}]`,
			`+ 2`,
		},
		want: `[1,2]`,
	}, {
		name:     "empty patch on multiset with numbers",
		metadata: MULTISET,
		given:    `[1,2]`,
		patch:    {``},
		want:     `[1,2]`,
	}, {
		name:     "add two numbers",
		metadata: MULTISET,
		given:    `[1]`,
		patch: {
			`@ [{}]`,
			`+ 2`,
			`+ 2`,
		},
		want: `[1,2,2]`,
	}, {
		name:     "remove a number",
		metadata: MULTISET,
		given:    `[1,2,3]`,
		patch: {
			`@ [{}]`,
			`- 2`,
		},
		want: `[1,3]`,
	}, {
		want:     "replace one object with another",
		metadata: MULTISET,
		given:    `[{"a":1}]`,
		patch: {
			`@ [{}]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		},
		want: `[{"a":2}]`,
	}, {
		name:     "remove two objects and add one",
		metadata: MULTISET,
		given:    `[{"a":1},{"a":1}]`,
		patch: {
			`@ [{}]`,
			`- {"a":1}`,
			`- {"a":1}`,
			`+ {"a":2}`,
		},
		want: `[{"a":2}]`,
	}, {
		name:     "remove three objects repeated and add one",
		metadata: MULTISET,
		given:    `["foo","foo","bar"]`,
		patch: {
			`@ [{}]`,
			`- "bar"`,
			`- "foo"`,
			`- "foo"`,
			`+ "baz"`,
		},
		want: `["baz"]`,
	}, {
		name:     "remove one object and add three repeated",
		metadata: MULTISET,
		given:    `["foo"]`,
		patch: {
			`@ [{}]`,
			`- "foo"`,
			`+ "bar"`,
			`+ "bar"`,
			`+ "baz"`,
		},
		want: `["bar","baz","bar"]`,
	}, {
		name:     "replace multiset with array",
		metadata: MULTISET,
		given:    `{}`,
		patch: {
			`@ []`,
			`- {}`,
			`+ []`,
		},
		want: `[]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withReadMetadata(c.metadata)
			checkPatch(ctx, c.given, c.want, c.patch...)
			ctx = newTestContext(t).
				withApplyMetadata(c.metadata)
			checkPatch(ctx, c.given, c.want, c.patch...)
		})
	}
}

func TestMultisetPatchError(t *testing.T) {
	cases := []struct {
		name     string
		metadata Metadata
		given    string
		patch    []string
	}{{
		name:     "remove number from empty multiset",
		metadata: MULTISET,
		given:    `[]`,
		patch: {
			`@ [{}]`,
			`- 1`,
		},
	}, {
		name:     "remove a single number twice",
		metadata: MULTISET,
		given:    `[1]`,
		patch: {
			`@ [{}]`,
			`- 1`,
			`- 1`,
		},
	}, {
		name:     "remove an object when there is a multiset",
		metadata: MULTISET,
		given:    `[]`,
		patch: {
			`@ []`,
			`- {}`,
		},
	}}
}

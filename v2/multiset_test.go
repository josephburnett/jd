package jd

import (
	"testing"
)

func TestMultisetJson(t *testing.T) {
	ctx := newTestContext(t).
		withOptions(multisetOption{})
	cases := []struct {
		name  string
		given string
		want  string
	}{{
		name:  "empty mulitset",
		given: `[]`,
		want:  `[]`,
	}, {
		name:  "empty multiset with space",
		given: ` [ ] `,
		want:  `[]`,
	}, {
		name:  "ordered multiset",
		given: `[1,2,3]`,
		want:  `[1,2,3]`,
	}, {
		name:  "ordered multiset with space",
		given: ` [1, 2, 3] `,
		want:  `[1,2,3]`,
	}, {
		name:  "multset with multiple duplicates",
		given: `[1,1,1]`,
		want:  `[1,1,1]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checkJson(ctx, c.given, c.want)
		})
	}
}

func TestMultisetEquals(t *testing.T) {
	ctx := newTestContext(t).
		withOptions(multisetOption{})
	cases := []struct {
		name string
		a    string
		b    string
	}{{
		name: "empty multisets",
		a:    `[]`,
		b:    `[]`,
	}, {
		name: "different ordered multisets 1",
		a:    `[1,2,3]`,
		b:    `[3,2,1]`,
	}, {
		name: "different ordered multisets 2",
		a:    `[1,2,3]`,
		b:    `[2,3,1]`,
	}, {
		name: "different ordered multisets 2",
		a:    `[1,2,3]`,
		b:    `[1,3,2]`,
	}, {
		name: "multsets with empty objects",
		a:    `[{},{}]`,
		b:    `[{},{}]`,
	}, {
		name: "nested multisets",
		a:    `[[1,2],[3,4]]`,
		b:    `[[2,1],[4,3]]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checkEqual(ctx, c.a, c.b)
		})
	}
}

func TestMultisetNotEquals(t *testing.T) {
	ctx := newTestContext(t).
		withOptions(multisetOption{})
	cases := []struct {
		name     string
		metadata Metadata
		a        string
		b        string
	}{{
		name: "empty multiset and multiset with number",
		a:    `[]`,
		b:    `[1]`,
	}, {
		name: "multisets with different numbers",
		a:    `[1,2,3]`,
		b:    `[1,2,2]`,
	}, {
		name: "multiset missing a number",
		a:    `[1,2,3]`,
		b:    `[1,2]`,
	}, {
		name: "nested multisets with different numbers",
		a:    `[[],[1]]`,
		b:    `[[],[2]]`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checkNotEqual(ctx, c.a, c.b)
		})
	}
}

func TestMultisetDiff(t *testing.T) {
	ctx := newTestContext(t).
		withOptions(multisetOption{})
	cases := []struct {
		name string
		a    string
		b    string
		want []string
		ctx  *testContext
	}{{
		name: "two empty multisets",
		a:    `[]`,
		b:    `[]`,
		want: ss(),
	}, {
		name: "two multisets with different numbers",
		a:    `[1]`,
		b:    `[1,2]`,
		want: ss(
			`@ [[]]`,
			`+ 2`,
		),
	}, {
		name: "two multisets with the same number",
		a:    `[1,2]`,
		b:    `[1,2]`,
		want: ss(),
	}, {
		name: "adding two numbers",
		a:    `[1]`,
		b:    `[1,2,2]`,
		want: ss(
			`@ [[]]`,
			`+ 2`,
			`+ 2`,
		),
	}, {
		name: "removing a number",
		a:    `[1,2,3]`,
		b:    `[1,3]`,
		want: ss(
			`@ [[]]`,
			`- 2`,
		),
	}, {
		name: "replacing one object with another",
		a:    `[{"a":1}]`,
		b:    `[{"a":2}]`,
		want: ss(
			`@ [[]]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
	}, {
		name: "replacing two objects with one object",
		a:    `[{"a":1},{"a":1}]`,
		b:    `[{"a":2}]`,
		want: ss(
			`@ [[]]`,
			`- {"a":1}`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
	}, {
		name: "replacing three strings repeated with one string",
		a:    `["foo","foo","bar"]`,
		b:    `["baz"]`,
		want: ss(
			`@ [[]]`,
			`- "bar"`,
			`- "foo"`,
			`- "foo"`,
			`+ "baz"`,
		),
	}, {
		name: "replacing one string with three repeated",
		a:    `["foo"]`,
		b:    `["bar","baz","bar"]`,
		want: ss(
			`@ [[]]`,
			`- "foo"`,
			`+ "bar"`,
			`+ "bar"`,
			`+ "baz"`,
		),
	}, {
		name: "replacing multiset with array",
		a:    `{}`,
		b:    `[]`,
		want: ss(
			`@ []`,
			`- {}`,
			`+ []`,
		),
	}, {
		name: "merge different types produces only new values",
		a:    `[1,2,2,3]`,
		b:    `{}`,
		want: ss(
			`^ {"Merge":true}`,
			`@ []`,
			`+ {}`,
		),
		ctx: newTestContext(t).withOptions(MERGE, multisetOption{}),
	}, {
		name: "merge outputs no diff when equal",
		a:    `[1,2,2,3]`,
		b:    `[2,1,3,2]`,
		want: ss(),
		ctx:  newTestContext(t).withOptions(MERGE, multisetOption{}),
	}, {
		name: "merge replaces entire multiset when not equal",
		a:    `[1,2,2,3]`,
		b:    `[2,1,3,3]`,
		want: ss(
			`^ {"Merge":true}`,
			`@ []`,
			`+ [2,1,3,3]`,
		),
		ctx: newTestContext(t).withOptions(MERGE, multisetOption{}),
	}}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.ctx
			if c == nil {
				c = ctx
			}
			checkDiff(c, tt.a, tt.b, tt.want...)
		})
	}
}

func TestMultisetPatch(t *testing.T) {
	cases := []struct {
		name  string
		given string
		patch []string
		want  string
	}{{
		name:  "empty patch on empty multiset",
		given: `[]`,
		patch: ss(``),
		want:  `[]`,
	}, {
		name:  "add a number",
		given: `[1]`,
		patch: ss(
			`@ [[]]`,
			`+ 2`,
		),
		want: `[1,2]`,
	}, {
		name:  "empty patch on multiset with numbers",
		given: `[1,2]`,
		patch: ss(``),
		want:  `[1,2]`,
	}, {
		name:  "add two numbers",
		given: `[1]`,
		patch: ss(
			`@ [[]]`,
			`+ 2`,
			`+ 2`,
		),
		want: `[1,2,2]`,
	}, {
		name:  "remove a number",
		given: `[1,2,3]`,
		patch: ss(
			`@ [[]]`,
			`- 2`,
		),
		want: `[1,3]`,
	}, {
		name:  "replace one object with another",
		given: `[{"a":1}]`,
		patch: ss(
			`@ [[]]`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
		want: `[{"a":2}]`,
	}, {
		name:  "remove two objects and add one",
		given: `[{"a":1},{"a":1}]`,
		patch: ss(
			`@ [[]]`,
			`- {"a":1}`,
			`- {"a":1}`,
			`+ {"a":2}`,
		),
		want: `[{"a":2}]`,
	}, {
		name:  "remove three objects repeated and add one",
		given: `["foo","foo","bar"]`,
		patch: ss(
			`@ [[]]`,
			`- "bar"`,
			`- "foo"`,
			`- "foo"`,
			`+ "baz"`,
		),
		want: `["baz"]`,
	}, {
		name:  "remove one object and add three repeated",
		given: `["foo"]`,
		patch: ss(
			`@ [[]]`,
			`- "foo"`,
			`+ "bar"`,
			`+ "bar"`,
			`+ "baz"`,
		),
		want: `["bar","baz","bar"]`,
	}, {
		name:  "replace multiset with array",
		given: `{}`,
		patch: ss(
			`@ []`,
			`- {}`,
			`+ []`,
		),
		want: `[]`,
	}, {
		name:  "merge replaces entire set",
		given: `[1,2,3]`,
		patch: ss(
			`^ {"Merge":true}`,
			`@ [[]]`,
			`+ [4,5,6]`,
		),
		want: `[4,5,6]`,
	}, {
		name:  "void deletes a node",
		given: `[1,2,3]`,
		patch: ss(
			`^ {"Merge":true}`,
			`@ [[]]`,
			`+`,
		),
		want: ``,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withOptions(multisetOption{})
			checkPatch(ctx, c.given, c.want, c.patch...)
		})
	}
}

func TestMultisetDiffOff(t *testing.T) {
	// DIFF_OFF disables multiset diffing
	a, _ := ReadJsonString(`{"items":[1,2,3]}`)
	b, _ := ReadJsonString(`{"items":[1,3,4]}`)
	d := a.Diff(b, MULTISET, PathOption(Path{PathKey("items")}, DIFF_OFF))
	if len(d) != 0 {
		t.Errorf("expected empty diff with DIFF_OFF, got: %v", d.Render())
	}
}

func TestMultisetPatchError(t *testing.T) {
	cases := []struct {
		name  string
		given string
		patch []string
	}{{
		name:  "remove number from empty multiset",
		given: `[]`,
		patch: ss(
			`@ [[]]`,
			`- 1`,
		),
	}, {
		name:  "remove a single number twice",
		given: `[1]`,
		patch: ss(
			`@ [[]]`,
			`- 1`,
			`- 1`,
		),
	}, {
		name:  "remove an object when there is a multiset",
		given: `[]`,
		patch: ss(
			`@ [[]]`,
			`- {}`,
		),
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := newTestContext(t).
				withOptions(multisetOption{})
			checkPatchError(ctx, c.given, c.patch...)
		})
	}
}

func TestMultisetPatchMultipleValues(t *testing.T) {
	// Multiple values at root triggers patchErrNonSetDiff
	node := jsonMultiset{jsonNumber(1)}
	d := Diff{{
		Path:   Path{},
		Remove: []JsonNode{jsonNumber(1), jsonNumber(2)},
	}}
	_, err := node.Patch(d)
	if err == nil {
		t.Fatal("expected error for multiple removals")
	}
}

func TestMultisetPatchBaseCases(t *testing.T) {
	// Base case: single value replace at root
	ms := jsonMultiset{jsonNumber(1), jsonNumber(2)}
	result, err := ms.patch(nil, Path{}, nil,
		[]JsonNode{ms},
		[]JsonNode{jsonMultiset{jsonNumber(3)}},
		nil, strictPatchStrategy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := jsonMultiset{jsonNumber(3)}
	if !result.Equals(expected) {
		t.Errorf("got %v, want %v", result.Json(), expected.Json())
	}
	// Base case: old value mismatch
	ms2 := jsonMultiset{jsonNumber(1)}
	_, err = ms2.patch(nil, Path{}, nil,
		[]JsonNode{jsonMultiset{jsonNumber(99)}},
		[]JsonNode{jsonMultiset{jsonNumber(2)}},
		nil, strictPatchStrategy)
	if err == nil {
		t.Fatal("expected error for old value mismatch")
	}
	// Non-PathMultiset path element error
	ms3 := jsonMultiset{jsonNumber(1)}
	_, err = ms3.patch(nil, Path{PathKey("x")}, nil,
		[]JsonNode{jsonNumber(1)},
		[]JsonNode{jsonNumber(2)},
		nil, strictPatchStrategy)
	if err == nil {
		t.Fatal("expected error for non-PathMultiset path element")
	}
}

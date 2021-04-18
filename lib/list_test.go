package jd

import "testing"

func TestListDiffMask(t *testing.T) {
	cases := []struct {
		name string
		a    JsonNode
		b    JsonNode
		mask Mask
		want Diff
	}{{
		name: "empty mask",
		a:    mustParseJson(`[1,2,3]`),
		b:    mustParseJson(`[1,2,4]`),
		mask: mustParseMask(``),
		want: mustParseDiff(
			`@ [2]`,
			`- 3`,
			`+ 4`,
		),
	}, {
		name: "specific element mask",
		a:    mustParseJson(`[1,2,3]`),
		b:    mustParseJson(`[1,2,4]`),
		mask: mustParseMask(`- [2]`),
		want: mustParseDiff(``),
	}, {
		name: "mask elements within a list",
		a:    mustParseJson(`[[1,2],[3,4]]`),
		b:    mustParseJson(`[[5,6],[7,8]]`),
		mask: mustParseMask(`- [[],0]`),
		want: mustParseDiff(
			`@ [0,1]`,
			`- 2`,
			`+ 6`,
			`@ [1,1]`,
			`- 4`,
			`+ 8`,
			// ignore the change at index 0's
		),
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a.Diff(tc.b, tc.mask)
			if !got.equal(tc.want) {
				t.Errorf("Wanted %v. Got %v", tc.want, got)
			}
		})
	}
}

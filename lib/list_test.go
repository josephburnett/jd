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

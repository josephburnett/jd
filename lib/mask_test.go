package jd

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadMaskString(t *testing.T) {

	cases := []struct {
		name    string
		mask    []string
		wantErr bool
		want    Mask
	}{{
		name: "empty mask",
		mask: []string{},
		want: Mask{},
	}, {
		name: "include single path",
		mask: []string{
			`+["foo"]`,
		},
		want: Mask{
			MaskElement{
				Include: true,
				Path:    mustParseJsonArray(`["foo"]`).(jsonArray),
			},
		},
	}, {
		name: "exclude single path",
		mask: []string{
			`-["foo"]`,
		},
		want: Mask{
			MaskElement{
				Include: false,
				Path:    mustParseJsonArray(`["foo"]`).(jsonArray),
			},
		},
	}, {
		name: "ignore whitespace",
		mask: []string{
			`  +  ["foo"]  `,
		},
		want: Mask{
			MaskElement{
				Include: true,
				Path:    mustParseJsonArray(`["foo"]`).(jsonArray),
			},
		},
	}, {
		name: "multiple and longer paths",
		mask: []string{
			`+["foo","bar"]`,
			`-["baz","boo"]`,
		},
		want: Mask{
			MaskElement{
				Include: true,
				Path:    mustParseJsonArray(`["foo","bar"]`).(jsonArray),
			},
			MaskElement{
				Include: false,
				Path:    mustParseJsonArray(`["baz","boo"]`).(jsonArray),
			},
		},
	}, {
		name: "path without inclusion sign",
		mask: []string{
			`["foo"]`,
		},
		wantErr: true,
	}, {
		name: "inclusion sign without path",
		mask: []string{
			`+`,
		},
		wantErr: true,
	}, {
		name: "double inclusion sign",
		mask: []string{
			`++["foo"]`,
		},
		wantErr: true,
	}, {
		name: "extra json",
		mask: []string{
			`+["foo"]["bar"]`,
		},
		wantErr: true,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mask := strings.Join(tc.mask, "\n")
			got, err := ReadMaskString(mask)
			if tc.wantErr && err == nil {
				t.Errorf("Wanted err. Got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Wanted no err. Got %q", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Wanted %v. Got %v", tc.want, got)
			}
		})
	}
}

func TestMaskInclude(t *testing.T) {
	cases := []struct {
		name string
		mask Mask
		i    JsonNode
		want bool
	}{{
		name: "empty mask",
		mask: Mask{},
		i:    mustParseJson(`"foo"`),
		want: true,
	}, {
		name: "single negative mask",
		mask: mustParseMask(
			`- ["foo"]`,
		),
		i:    mustParseJson(`"foo"`),
		want: false,
	}, {
		name: "single positive mask",
		mask: mustParseMask(
			`+ ["foo"]`,
		),
		i:    mustParseJson(`"foo"`),
		want: true,
	}, {
		name: "positive mask with inapplicable negative mask",
		mask: mustParseMask(
			`+ ["foo"]`,
			`- ["foo","bar"]`,
		),
		i:    mustParseJson(`"foo"`),
		want: true,
	}, {
		name: "positive mask with overriding negative mask",
		mask: mustParseMask(
			`+ ["foo"]`,
			`- ["foo"]`,
		),
		i:    mustParseJson(`"foo"`),
		want: false,
	}, {
		name: "array index negative masked",
		mask: mustParseMask(
			`- [0]`,
		),
		i:    mustParseJson(`0`),
		want: false,
	}, {
		name: "array index not negative masked",
		mask: mustParseMask(
			`- [0]`,
		),
		i:    mustParseJson(`1`),
		want: true,
	}, {
		name: "array all indexes negative masked",
		mask: mustParseMask(
			`- [[]]`,
		),
		i:    mustParseJson(`0`),
		want: false,
	}, {
		name: "array all indexes positive masked",
		mask: mustParseMask(
			`+ [[]]`,
		),
		i:    mustParseJson(`0`),
		want: true,
	}, {
		name: "set with matching key",
		mask: mustParseMask(
			`- [{"key":"value"}]`,
		),
		i:    mustParseJson(`{"key":"value"}`),
		want: false,
	}, {
		name: "set with non-matching key",
		mask: mustParseMask(
			`- [{"key":"value"}]`,
		),
		i:    mustParseJson(`{"key":"not"}`),
		want: true,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.mask.include(tc.i)
			if tc.want != got {
				t.Errorf("Wanted %v. Got %v", tc.want, got)
			}
		})
	}
}

func TestDiffMask(t *testing.T) {
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
	}, {
		name: "mask property of all objects in set",
		a:    mustParseJson(`[{"foo":"bar"},{"baz":"bam"}]`),
		b:    mustParseJson(`[{"foo":"boo"},{"baz":"hoo"}]`),
		mask: mustParseMask(`- [["set"],"foo"]`),
		want: mustParseDiff(
			`@ [{}]`,
			`{"baz":"bam"}`,
			`{"baz":"hoo"}`,
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

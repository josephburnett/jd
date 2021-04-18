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

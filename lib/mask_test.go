package jd

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadMaskString(t *testing.T) {
	mustParse := func(s string) []JsonNode {
		n, err := ReadJsonString(s)
		if err != nil {
			panic(err)
		}
		return n.(jsonArray)
	}

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
				Path:    mustParse(`["foo"]`),
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
				Path:    mustParse(`["foo"]`),
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
				Path:    mustParse(`["foo"]`),
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
				Path:    mustParse(`["foo","bar"]`),
			},
			MaskElement{
				Include: false,
				Path:    mustParse(`["baz","boo"]`),
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

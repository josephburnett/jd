package jd

import "github.com/josephburnett/jd/v2/internal/types"

func ReadDiffString(s string) (Diff, error) {
	return types.ReadDiffString(s)
}

func ReadPatchString(s string) (Diff, error) {
	return types.ReadPatchString(s)
}

func ReadMergeString(s string) (Diff, error) {
	return types.ReadMergeString(s)
}
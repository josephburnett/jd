package jd

import (
	"testing"
)

func TestNumberEqual(t *testing.T) {
	checkEqual(t, `0`, `0`)
	checkEqual(t, `0`, `0.0`)
	checkEqual(t, `0.0001`, `0.0001`)
	checkEqual(t, `123`, `123`)
}

func TestNumberNotEqual(t *testing.T) {
	checkNotEqual(t, `0`, `1`)
	checkNotEqual(t, `0`, `0.0001`)
	checkNotEqual(t, `1234`, `1235`)
}

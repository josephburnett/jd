package jd

import (
	"testing"
)

func TestBoolEqual(t *testing.T) {
	checkEqual(t, `true`, `true`)
	checkEqual(t, `false`, `false`)
}

func TestBoolNotEqual(t *testing.T) {
	checkNotEqual(t, `true`, `false`)
	checkNotEqual(t, `false`, `true`)
	checkNotEqual(t, `false`, `[]`)
	checkNotEqual(t, `true`, `"true"`)
}

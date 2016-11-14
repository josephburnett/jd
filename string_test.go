package jd

import (
	"testing"
)

func TestStringEqual(t *testing.T) {
	checkEqual(t, `""`, `""`)
	checkEqual(t, `"a"`, `"a"`)
	checkEqual(t, `"123"`, `"123"`)
}

func TestStringNotEqual(t *testing.T) {
	checkNotEqual(t, `""`, `"a"`)
	checkNotEqual(t, `""`, `[]`)
	checkNotEqual(t, `""`, `{}`)
	checkNotEqual(t, `""`, `0`)
}

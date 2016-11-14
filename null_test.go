package jd

import (
	"testing"
)

func TestNullEqual(t *testing.T) {
	checkEqual(t, `null`, `null`)
}

func TestNullNotEqual(t *testing.T) {
	checkNotEqual(t, `null`, `0`)
	checkNotEqual(t, `null`, `[]`)
	checkNotEqual(t, `null`, `{}`)
}

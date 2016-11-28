package jd

import (
	"testing"
)

func TestArrayBagJson(t *testing.T) {
	checkJson(t, `[]`, `[]`)
	checkJson(t, ` [ ] `, `[]`)
	checkJson(t, `[1,2,3]`, `[1,2,3]`)
	checkJson(t, ` [1, 2, 3] `, `[1,2,3]`)
}

func TestArrayBagEquals(t *testing.T) {
	checkEqual(t, `[]`, `[]`, ARRAY_BAG)
	checkEqual(t, `[1,2,3]`, `[3,2,1]`, ARRAY_BAG)
	checkEqual(t, `[1,2,3]`, `[2,3,1]`, ARRAY_BAG)
	checkEqual(t, `[1,2,3]`, `[1,3,2]`, ARRAY_BAG)
	checkEqual(t, `[{},{}]`, `[{},{}]`, ARRAY_BAG)
	checkEqual(t, `[[1,2],[3,4]]`, `[[2,1],[4,3]]`, ARRAY_BAG)
}

func TestArrayBagNotEquals(t *testing.T) {
	checkNotEqual(t, `[]`, `[1]`, ARRAY_BAG)
	checkNotEqual(t, `[1,2,3]`, `[1,2,2]`, ARRAY_BAG)
	checkNotEqual(t, `[1,2,3]`, `[1,2]`, ARRAY_BAG)
	checkNotEqual(t, `[[],[1]]`, `[[],[2]]`, ARRAY_BAG)
}

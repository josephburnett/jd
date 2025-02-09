package jd

import (
	"fmt"
	"testing"
)

func TestExampleJsonNode_Diff(t *testing.T) {
	a, _ := ReadJsonString(`{"foo":["bar"]}`)
	b, _ := ReadJsonString(`{"foo":["baz"]}`)
	fmt.Print(a.Diff(b).Render())
	// Output:
	// @ ["foo",0]
	// [
	// - "bar"
	// + "baz"
	// ]
}

func TestExampleJsonNode_Patch(t *testing.T) {
	a, _ := ReadJsonString(`["foo"]`)
	diff, _ := ReadDiffString(`
@ [1]
  "foo"
+ "bar"
]
`)
	b, _ := a.Patch(diff)
	fmt.Print(b.Json())
	// Output:
	// ["foo","bar"]
}

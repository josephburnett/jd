package jd

import (
	"fmt"
)

func ExampleJsonNode_Diff() {
	a, _ := ReadJsonString(`{"foo":"bar"}`)
	b, _ := ReadJsonString(`{"foo":"baz"}`)
	fmt.Print(a.Diff(b).Render())
	// Output:
	// @ ["foo"]
	// - "bar"
	// + "baz"
}

func ExampleJsonNode_Patch() {
	a, _ := ReadJsonString(`["foo"]`)
	diff, _ := ReadDiffString(`` +
		`@ [1]` + "\n" +
		`+ "bar"` + "\n")
	b, _ := a.Patch(diff)
	fmt.Print(b.Json())
	// Output:
	// ["foo","bar"]
}

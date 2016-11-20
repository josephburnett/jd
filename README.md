# JSON diff and patch

`jd` is a commandline utility and Go library for diffing and patching JSON values.

## Command line usage

Download [latest release](https://github.com/josephburnett/jd/releases/latest) or `go get github.com/josephburnett/jd`

```
Usage: jd [OPTION]... FILE1 [FILE2]
Diff and patch JSON files.

Prints the diff of FILE1 and FILE2 to STDOUT.
When FILE2 is omitted the second input is read from STDIN.
When patching (-p) FILE1 is a diff.

Options:
  -p        Apply patch FILE1 to FILE2 or STDIN.
  -o=FILE3  Write to FILE3 instead of STDOUT.

Examples:
  jd a.json b.json
  cat b.json | jd a.json
  jd -o patch a.json b.json; jd patch a.json
```

## Library usage

`go get github.com/josephburnett/jd`

```Go
import (
	"fmt"
	jd "github.com/josephburnett/jd/lib"
)

func ExampleJsonNode_Diff() {
	a, _ := jd.ReadJsonString(`{"foo":"bar"}`)
	b, _ := jd.ReadJsonString(`{"foo":"baz"}`)
	fmt.Print(a.Diff(b).Render())
	// Output:
	// @ ["foo"]
	// - "bar"
	// + "baz"
}

func ExampleJsonNode_Patch() {
	a, _ := jd.ReadJsonString(`["foo"]`)
	diff, _ := jd.ReadDiffString(`` +
		`@ [1]` + "\n" +
		`+ "bar"` + "\n")
	b, _ := a.Patch(diff)
	fmt.Print(b.Json())
	// Output:
	// ["foo","bar"]
}
```

## Diff language

![Railroad diagram of EBNF](/ebnf.png)

### EBNF

```EBNF
Diff ::= ( '@' '[' ( 'JSON String' | 'JSON Number' )* ']' '\n' ( '+' | '-' ( 'JSON Value' '\n' '+' )? ) 'JSON Value' '\n' )*
```

### Examples

```JSON
@ ["a"]
- 1
+ 2
```

```JSON
@ [2]
+ {"foo":"bar"}
```

```JSON
@ ["Movies",67,"Title"]
- "Dr. Strangelove"
+ "Dr. Evil Love"
@ ["Movies",67,"Actors","Dr. Strangelove"]
- "Peter Sellers"
+ "Mike Myers"
@ ["Movies",102]
+ {"Title":"Austin Powers","Actors":{"Austin Powers":"Mike Myers"}}
```

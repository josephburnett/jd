# JSON diff and patch

`jd` is a commandline utility and Go library for diffing and patching JSON values.

## Try it out

https://jd-tool.appspot.com

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
  -set      Treat arrays as sets.
  -mset     Treat arrays as multisets (bags).

Examples:
  jd a.json b.json
  cat b.json | jd a.json
  jd -o patch a.json b.json; jd patch a.json
  jd -set a.json b.json
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

- A diff is zero or more sections
- Sections start with a `@` header and the path to a node
- A path is a JSON list of zero or more elements accessing collections
- A JSON number element (e.g. `0`) accesses an array
- A JSON string element (e.g. `"foo"`) accesses an object
- An empty JSON object element (`{}`) accesses an array as a set or multiset
- After the path is one or more removals or additions, removals first
- Removals start with `-` and then the JSON value to be removed
- Additions start with `+` and then the JSON value to added

### EBNF

```EBNF
Diff ::= ( '@' '[' ( 'JSON String' | 'JSON Number' | 'Empty JSON Object' )* ']' '\n' ( ( '-' 'JSON Value' '\n' )+ | '+' 'JSON Value' '\n' ) ( '+' 'JSON Value' '\n' )* )*
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

```JSON
@ ["Movies",67,"Tags",{}]
- "Romance"
+ "Action"
+ "Comedy"
```

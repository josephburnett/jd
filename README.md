# jd -- JSON diff and patch

`jd` is a commandline utility and Go library for diffing and patching JSON values.

## Command line usage

Download [latest release](https://github.com/josephburnett/jd/releases) or `go get github.com/josephburnett/jd`

```
Usage: jd [OPTION]... FILE1 [FILE2]
Diff and patch JSON files.

When FILE2 is omitted the second input is read from STDIN.

Options:
  -p  Apply patch FILE1 to FILE2 or STDIN.

Examples:
  jd a.json b.json
  cat b.json | jd a.json
  jd a.json b.json > patch; jd patch a.json
```

## Libarary usage

`go get github.com/josephburnett/jd`

```Go
func ExampleJsonNode_Diff() {
	a, _ := ReadJsonString(`{"foo":"bar"}`)
	b, _ := ReadJsonString(`{"foo":"baz"}`)
	fmt.Print(a.Diff(b).Render())
	// Output:
	// @ ["foo"]
	// - "bar"
	// + "baz"
}
```

```Go
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
```

## Diff language

![Railroad diagram of EBNF](/ebnf.png)

### EBNF

```
Diff ::= ( '@' '[' ( 'JSON String' | 'JSON Number' )* ']' '\n' ( '+' | '-' ( 'JSON Value' '\n' '+' )? ) 'JSON Value' '\n' )*
```

### Examples

```
@ ["a"]
- 1
+ 2
```

```
@ [2]
+ {"foo":"bar"}
```

```
@ ["Movies",67,"Title"]
- "Dr. Strangelove"
+ "Dr. Evil Love"
@ ["Movies",67,"Actors","Dr. Strangelove"]
- "Peter Sellers"
+ "Mike Myers"
@ ["Movies",102]
+ {"Title":"Austin Powers","Actors":{"Austin Powers":"Mike Myers"}}
```

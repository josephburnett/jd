# jd -- JSON diff and patch

## Command line usage

`go get github.com/josephburnett/jd`

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

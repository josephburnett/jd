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

### EBNF

```
diff         = { diff element } ;
diff element = path line, add | replace | remove ;
path line    = path header, path, "\n" ;
path header  = "@" ;
path         = "[", { path element }, "]" ;
path element = json string | json number ;
add          = add line ;
replace      = remove line, add line ;
remove       = remove line ;
add line     = "+", json value, "\n" ;
remove line  = "-", json value, "\n" ;
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

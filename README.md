# jd -- JSON diff and patch

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

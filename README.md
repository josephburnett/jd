[![Go Report Card](https://goreportcard.com/badge/josephburnett/jd)](https://goreportcard.com/report/josephburnett/jd)

# JSON diff and patch

`jd` is a commandline utility and Go library for diffing and patching
JSON and YAML values. It supports a native `jd` format (similar to
unified format) as well as JSON Merge Patch ([RFC
7386](https://datatracker.ietf.org/doc/html/rfc7386)) and a subset of
JSON Patch ([RFC
6902](https://datatracker.ietf.org/doc/html/rfc6902)). Try it out at
http://play.jd-tool.io/.

## Example

Diff `jd a.json b.json`:

```JSON
{"foo":["bar","baz"]}
```

```JSON
{"foo":["bar","bam","boom"]}
```

Output:

```DIFF
@ ["foo",1]
  "bar"
- "baz"
+ "bam"
+ "boom"
]
```

## Features

1. Human-friendly format, similar to Unified Diff.
2. Produces a minimal diff between array elements using LCS algorithm.
3. Adds context before and after when modifying an array to prevent bad patches.
4. Create and apply structural patches in jd, patch (RFC 6902) and merge (RFC 7386) patch formats.
5. Translates between patch formats.
6. Includes Web Assembly-based UI (no network calls).

## Installation

GitHub Action:

```yaml
    - name: Diff A and B
      id: diff
      uses: josephburnett/jd@v2.1.2
      with:
        args: a.json b.json
    - name: Print the diff
      run: echo '${{ steps.diff.outputs.output }}'
    - name: Check the exit code
      run: if [ "${{ steps.diff.outputs.exit_code }}" != "1" ]; then exit 1; fi
```

To get the `jd` commandline utility:
* run `brew install jd`, or
* run `go install github.com/josephburnett/jd/v2/jd@latest`, or
* visit https://github.com/josephburnett/jd/releases/latest and download the pre-built binary for your architecture/os, or
* run in a Docker image `jd(){ docker run --rm -i -v $PWD:$PWD -w $PWD josephburnett/jd "$@"; }`.

To use the `jd` web UI:
* visit http://play.jd-tool.io/, or
* run `jd -port 8080` and visit http://localhost:8080.

Note: to include the UI when building from source, use the Makefile.

## Command line usage

```
Usage: jd [OPTION]... FILE1 [FILE2]
Diff and patch JSON files.

Prints the diff of FILE1 and FILE2 to STDOUT.
When FILE2 is omitted the second input is read from STDIN.
When patching (-p) FILE1 is a diff.

Options:
  -color       Print color diff.
  -p           Apply patch FILE1 to FILE2 or STDIN.
  -o=FILE3     Write to FILE3 instead of STDOUT.
  -opts='[]'   JSON array of options. Supports global options and PathOptions.
               Global: ["SET"], ["MULTISET"], [{"precision":0.1}], [{"setkeys":["id"]}], ["DIFF_ON"], ["DIFF_OFF"]
               PathOptions target specific paths: [{"@":["path"],"^":["SET"]}]
               Example: [{"@":["users"],"^":["SET"]},{"@":["scores",0],"^":[{"precision":0.1}]}]
  -set         Treat arrays as sets. Same as -opts='["SET"]'.
  -mset        Treat arrays as multisets (bags). Same as -opts='["MULTISET"]'.
  -setkeys     Keys to identify set objects. Same as -opts='[{"setkeys":["key1","key2"]}]'.
  -yaml        Read and write YAML instead of JSON.
  -port=N      Serve web UI on port N
  -precision=N Maximum absolute difference for numbers to be equal.
               Same as -opts='[{"precision":N}]'. Example: -precision=0.00001
  -f=FORMAT    Read and write diff in FORMAT "jd" (default), "patch" (RFC 6902) or
               "merge" (RFC 7386)
  -t=FORMATS   Translate FILE1 between FORMATS. Supported formats are "jd",
               "patch" (RFC 6902), "merge" (RFC 7386), "json" and "yaml".
               FORMATS are provided as a pair separated by "2". E.g.
               "yaml2json" or "jd2patch".

Examples:
  jd a.json b.json
  cat b.json | jd a.json
  jd -o patch a.json b.json; jd patch a.json
  jd -set a.json b.json
  jd -f patch a.json b.json
  jd -f merge a.json b.json
  jd -opts='[{"@":["items"],"^":["SET"]}]' a.json b.json
  jd -opts='[{"@":["temperature"],"^":[{"precision":0.1}]}]' a.json b.json
```

#### Command Line Option Details

`setkeys` This option determines what keys are used to decide if two
objects 'match'. Then the matched objects are compared, which will
return a diff if there are differences in the objects themselves,
their keys and/or values. You shouldn't expect this option to mask or
ignore non-specified keys, it is not intended as a way to 'ignore'
some differences between objects.

#### PathOptions: Targeted Comparison Options

PathOptions allow you to apply different comparison semantics to specific paths in your JSON/YAML data. This enables precise control over how different parts of your data are compared.

**PathOption Syntax:**
```json
{"@": ["path", "to", "target"], "^": [options]}
```

- `@` (At): JSON path array specifying where to apply the option
- `^` (Then): Array of options to apply at that path

**Supported Options:**
- `"SET"`: Treat array as a set (ignore order and duplicates)
- `"MULTISET"`: Treat array as a multiset (ignore order, count duplicates)  
- `{"precision": N}`: Numbers within N are considered equal
- `{"setkeys": ["key1", "key2"]}`: Match objects by specified keys
- `"DIFF_ON"`: Enable diffing at this path (default behavior)
- `"DIFF_OFF"`: Disable diffing at this path, ignore all changes

**Examples:**

Treat specific array as a set while others remain as lists:
```bash
jd -opts='[{"@":["tags"],"^":["SET"]}]' a.json b.json
```

Apply precision to specific temperature field:
```bash
jd -opts='[{"@":["sensor","temperature"],"^":[{"precision":0.1}]}]' a.json b.json
```

Multiple PathOptions - SET on one path, precision on another:
```bash
jd -opts='[{"@":["items"],"^":["SET"]}, {"@":["price"],"^":[{"precision":0.01}]}]' a.json b.json
```

Target specific array index:
```bash
jd -opts='[{"@":["measurements", 0],"^":[{"precision":0.05}]}]' a.json b.json
```

Apply to root level:
```bash
jd -opts='[{"@":[],"^":["SET"]}]' a.json b.json
```

Ignore specific fields (deny-list approach):
```bash
jd -opts='[{"@":["timestamp"],"^":["DIFF_OFF"]}, {"@":["metadata","generated"],"^":["DIFF_OFF"]}]' a.json b.json
```

Allow-list approach - ignore everything except specific fields:
```bash
jd -opts='[{"@":[],"^":["DIFF_OFF"]}, {"@":["userdata"],"^":["DIFF_ON"]}]' a.json b.json
```

Nested override - ignore parent but include specific child:
```bash
jd -opts='[{"@":["config"],"^":["DIFF_OFF"]}, {"@":["config","user_settings"],"^":["DIFF_ON"]}]' a.json b.json
```

## Library usage

Note: import only release commits (`v2.Y.Z`) because `master` can be unstable.

Note: the `v2` library replaces the v1 (`lib`) library. V2 adds diff
context, minimal array diffs and hunk-level metadata. However the
format is not backward compatable. You should use `v2`.

```GO
import (
	"fmt"
	jd "github.com/josephburnett/jd/v2"
)

func ExampleJsonNode_Diff() {
	a, _ := jd.ReadJsonString(`{"foo":["bar"]}`)
	b, _ := jd.ReadJsonString(`{"foo":["baz"]}`)
	fmt.Print(a.Diff(b).Render())
	// Output:
	// @ ["foo",0]
	// [
	// - "bar"
	// + "baz"
	// ]
}

func ExampleJsonNode_Patch() {
	a, _ := jd.ReadJsonString(`["foo"]`)
	diff, _ := jd.ReadDiffString(`
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

func ExamplePathOptions() {
	// Apply SET semantics to specific array path
	a, _ := jd.ReadJsonString(`{"tags":["red","blue","green"], "items":[1,2,3]}`)
	b, _ := jd.ReadJsonString(`{"tags":["green","red","blue"], "items":[3,2,1]}`)
	
	// Only treat "tags" as a set, "items" remain as list
	opts, _ := jd.ReadOptionsString(`[{"@":["tags"],"^":["SET"]}]`)
	diff := a.Diff(b, opts...)
	fmt.Print(diff.Render())
	// Output:
	// @ ["items",0]
	// [
	// + 3
	// + 2
	//   1
	// @ ["items",3]
	//   1
	// - 2
	// - 3
	// ]
}

func ExampleMultiplePathOptions() {
	a, _ := jd.ReadJsonString(`{"temp":20.12, "pressure":1013.25, "tags":["A","B","C"]}`)
	b, _ := jd.ReadJsonString(`{"temp":20.15, "pressure":1013.30, "tags":["C","A","B"]}`)
	
	// Apply precision to temp, exact match to pressure, SET semantics to tags
	opts, _ := jd.ReadOptionsString(`[
		{"@":["temp"],"^":[{"precision":0.1}]},
		{"@":["tags"],"^":["SET"]}
	]`)
	diff := a.Diff(b, opts...)
	fmt.Print(diff.Render())
	// Output:
	// @ ["pressure"]
	// - 1013.25
	// + 1013.3
}

func ExampleSelectiveDiffing() {
	a, _ := jd.ReadJsonString(`{"userdata":"important","system":"ignore1","timestamp":"2023-01-01"}`)
	b, _ := jd.ReadJsonString(`{"userdata":"changed","system":"ignore2","timestamp":"2023-01-02"}`)
	
	// Allow-list approach: ignore everything except userdata
	opts, _ := jd.ReadOptionsString(`[
		{"@":[],"^":["DIFF_OFF"]},
		{"@":["userdata"],"^":["DIFF_ON"]}
	]`)
	diff := a.Diff(b, opts...)
	fmt.Print(diff.Render())
	// Output:
	// @ ["userdata"]
	// - "important"
	// + "changed"
}

func ExampleNestedOverride() {
	a, _ := jd.ReadJsonString(`{"config":{"system":"val1","user_settings":"setting1"}}`)
	b, _ := jd.ReadJsonString(`{"config":{"system":"val2","user_settings":"setting2"}}`)
	
	// Ignore config changes except for user_settings
	opts, _ := jd.ReadOptionsString(`[
		{"@":["config"],"^":["DIFF_OFF"]},
		{"@":["config","user_settings"],"^":["DIFF_ON"]}
	]`)
	diff := a.Diff(b, opts...)
	fmt.Print(diff.Render())
	// Output:
	// @ ["config","user_settings"]
	// - "setting1"
	// + "setting2"
}
```

## Diff Language (v2)

The jd v2 diff format is a human-readable structural diff format with context and metadata support.

### Format Overview

A diff consists of:
- **Options header** (optional): Shows the options used to create the diff
- **Metadata lines** (optional): Start with `^` and specify hunk-level metadata  
- **Diff hunks**: Start with `@` and specify the path, followed by changes and context

### Options Header

When options are provided to `jd`, they are displayed at the beginning of the diff to show how it was produced. Each option appears on its own line starting with `^ `:

```diff
^ "SET"
^ {"precision":0.001}
@ ["items",{}]
- "old-item"
+ "new-item"
```

This feature helps understand:
- Whether arrays were treated as sets (`"SET"`) or multisets (`"MULTISET"`) 
- What precision was used for number comparisons (`{"precision":N}`)
- Which keys identify set objects (`{"setkeys":["key1","key2"]}`)
- Path-specific options (`{"@":["path"],"^":["OPTION"]}`)
- Whether merge semantics were applied (`"MERGE"`)
- If color output was requested (`"COLOR"`)

The options header is informational and helps with debugging diff behavior. Note that diffs with options headers can still be parsed and applied as patches.

### EBNF Grammar

```EBNF
Diff ::= OptionsHeader* (MetadataLine | DiffHunk)*

OptionsHeader ::= '^' SP JsonValue NEWLINE

MetadataLine ::= '^' SP JsonObject NEWLINE

DiffHunk ::= '@' SP JsonArray NEWLINE
             ContextLine*
             (RemoveLine | AddLine)*
             ContextLine*

ContextLine ::= SP SP JsonValue NEWLINE

RemoveLine ::= '-' SP JsonValue NEWLINE

AddLine ::= '+' SP JsonValue NEWLINE

JsonArray ::= '[' (PathElement (',' PathElement)*)? ']'

PathElement ::= JsonString        // Object key: "foo"
              | JsonNumber        // Array index: 0 
              | EmptyObject       // Set marker: {}
              | EmptyArray        // List marker: [] 
              | ObjectWithKeys    // Set keys: {"id":"value"}
              | ArrayWithObject   // Multiset: [{}] or [{"id":"value"}]
```

*Note: Railroad diagram at /ebnf.png needs updating for v2 format.*

### Path Elements Reference

| Element | Description | Example Path |
|---------|-------------|--------------|
| `"key"` | Object field access | `["user","name"]` |
| `0`, `1`, etc. | Array index access | `["items",0]` |
| `{}` | Treat array as set (ignore order/duplicates) | `["tags",{}]` |
| `[]` | Explicit list marker | `["values",[]]` |
| `{"id":"val"}` | Match objects by specific key | `["users",{"id":"123"}]` |
| `[{}]` | Treat as multiset (ignore order, count duplicates) | `["counts",[{}]]` |
| `[{"key":"val"}]` | Match multiset objects by key | `["items",[{"id":"456"}]]` |

### Line Types

- **`@ [path]`**: Diff hunk header specifying the location
- **`^ {metadata}`**: Metadata for the following hunks (inherits downward)  
- **`  value`**: Context lines (spaces) - elements that provide context
- **`- value`**: Remove lines - values being removed
- **`+ value`**: Add lines - values being added

### Core Examples

#### Simple Object Change
```diff
@ ["name"]
- "Alice"
+ "Bob"
```

#### Array Element with Context
```diff
@ ["items",1]
  "apple"
+ "banana" 
  "cherry"
```

#### Set Operations (Ignore Order)
```diff
@ ["tags",{}]
- "urgent"
+ "completed"
+ "reviewed"
```

#### Object Identification by Key
```diff
@ ["users",{"id":"123"},"status"]
- "pending"
+ "active"
```

#### Multiset Operations
```diff
@ ["scores",[{}]]
- 85
- 92
+ 88
+ 95
+ 95
```

### Advanced Examples

#### Merge Patch Metadata
```diff
^ {"Merge":true}
@ ["config"]
- {"timeout":30,"retries":3}
+ {"timeout":60,"retries":5,"debug":true}
```

#### Complex List Context
```diff
@ ["matrix",1,2]
  [[1,2,3],[4,5,6]]
- 6
+ 9
  [7,8,9]
]
```

#### Nested Set with PathOptions
```diff
@ ["department","employees",{"employeeId":"E123"},"projects",{}]
- "ProjectA"
+ "ProjectB" 
+ "ProjectC"
```

#### Multiple Hunks with Inheritance
```diff
^ {"Merge":true}
@ ["user","preferences"] 
+ {"theme":"dark","notifications":true}
@ ["user","lastLogin"]
+ "2023-12-01T10:30:00Z"
```

### Integration with PathOptions

The path syntax directly corresponds to PathOption targeting:
- Diff path `["users",{}]` ↔ PathOption `{"@":["users"],"^":["SET"]}`
- Diff path `["items",{"id":"123"}]` ↔ PathOption with SetKeys targeting
- Diff path `["scores",[{}]]` ↔ PathOption `{"@":["scores"],"^":["MULTISET"]}`

This allows fine-grained control over how different parts of your data structures are compared and diffed.

## Cookbook

### Use git diff to produce a structural diff:

#### Option 1: Direct external tool (simple, one-off usage)
```bash
git difftool -yx jd @ -- foo.json
```
Use this when you want a quick structural diff without any setup. Works immediately if `jd` is in your PATH.

#### Option 2: Configure as git diff driver (integrated, regular usage)
```bash
# One-time setup
git config diff.jd.command 'jd --git-diff-driver'
echo "*.json diff=jd" >> .gitattributes

# Then use with any git diff command:
git diff foo.json
git difftool -t jd foo.json
```
Use this approach if you regularly work with JSON files. It integrates `jd` into git's diff system, automatically using structural diffs for JSON files across all git commands. The `.gitattributes` file can be committed to share this behavior with your team.

Example output from either approach:
```diff
@ ["foo"]
- "bar"
+ "baz"
```

### See what changes in a Kubernetes Deployment:
```bash
kubectl get deployment example -oyaml > a.yaml
kubectl edit deployment example
# change cpu resource from 100m to 200m
kubectl get deployment example -oyaml | jd -yaml a.yaml
```
output:
```diff
@ ["metadata","annotations","deployment.kubernetes.io/revision"]
- "2"
+ "3"
@ ["metadata","generation"]
- 2
+ 3
@ ["metadata","resourceVersion"]
- "4661"
+ "5179"
@ ["spec","template","spec","containers",0,"resources","requests","cpu"]
- "100m"
+ "200m"
@ ["status","conditions",1,"lastUpdateTime"]
- "2021-12-23T09:40:39Z"
+ "2021-12-23T09:41:49Z"
@ ["status","conditions",1,"message"]
- "ReplicaSet \"nginx-deployment-787d795676\" has successfully progressed."
+ "ReplicaSet \"nginx-deployment-795c7f5bb\" has successfully progressed."
@ ["status","observedGeneration"]
- 2
+ 3
```
apply these change to another deployment:
```bash
# edit file "patch" to contain only the hunk updating cpu request
kubectl patch deployment example2 --type json --patch "$(jd -t jd2patch ~/patch)"
```


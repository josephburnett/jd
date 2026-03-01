# Test Cases

Implementation-agnostic test data for validating structural diff format implementations. These files define what the format produces and accepts, without reference to any specific CLI, library, or programming language.

## Operations

Test cases exercise two operations:

- **diff**: Given two JSON documents, produce a structural diff. This is the default operation when `operation` is omitted.
- **patch**: Given a structural diff and a JSON document, produce a patched document.

## Test Case Schema

Each JSON file contains an array of test case objects:

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique identifier for the test case |
| `description` | string | Human-readable description |
| `category` | string | Grouping: `diff`, `options`, `path_options`, or `errors` |
| `operation` | string | `"diff"` (default if omitted) or `"patch"` |
| `content_a` | string | First input document (JSON text) |
| `content_b` | string | Second input document (JSON text) |
| `options` | array | Format options, e.g. `["SET"]`, `[{"precision": 0.1}]` |
| `expected_output` | string | Exact expected output text |
| `accepted_lines` | array | Groups of lines for order-insensitive matching (alternative to `expected_output`) |
| `expect_error` | bool | `true` if the operation should fail |

## Categories

### `diff` (diff.json)

Diff generation across all JSON types: objects, arrays, strings, numbers, booleans, and null. Covers additions, removals, replacements, nesting, unicode, and edge cases like empty structures.

### `options` and `path_options` (options.json)

Options that modify comparison behavior: `SET`, `MULTISET`, `precision`, `keys`, `DIFF_OFF`/`DIFF_ON`. Includes both global options and path-scoped `PathOptions`.

### `errors` (errors.json)

Error conditions: malformed JSON, invalid options, conflicting options, and patch application failures.

## Output Matching

For diff tests:
- If `expected_output` is set, the output must match exactly (after normalizing line endings and trimming trailing whitespace). Metadata lines (`^ ...`) in the output that are not present in `expected_output` are ignored.
- If `accepted_lines` is set, the output is matched group-by-group. Each group is a set of lines that may appear in any order within that group, but groups are consumed sequentially.
- If neither is set and `expect_error` is false, the output should be empty (documents are equal).

## Round-Trip Property

For any diff test that produces non-empty output, applying that output as a patch to `content_a` should yield `content_b`. This property provides implicit patch coverage without needing separate patch test cases.

## Files

- **diff.json** — Diff generation tests
- **options.json** — Options and PathOptions tests
- **errors.json** — Error condition tests

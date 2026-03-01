# Test Cases

Implementation-agnostic test data for validating structural diff format implementations. These files define what the format produces and accepts, without reference to any specific CLI, library, or programming language.

## Operations

The format defines two operations:

- **diff**: Given two JSON documents and options, produce a structural diff.
- **patch**: Given a structural diff and a JSON document, produce a patched document.

Both operations are fundamental to the format. Diff produces the format's output; patch consumes it.

## Verification Procedure

### Non-error test cases

For each test case where `expect_error` is false (or absent):

1. **Diff**: Compute `diff(content_a, content_b, options)`. The result must match `expected_output` or `accepted_lines` (see Output Matching below). If neither is set, the result must be empty.

2. **Patch**: If the diff produced non-empty output, apply it as a patch: compute `patch(diff_output, content_a)`. The result must equal `content_b`.

3. **No-diff**: Compute `diff(patched_result, content_b, options)`. The output must be empty, confirming the patched result is identical to the target document.

### Error test cases

For each test case where `expect_error` is true, the operation must produce an error. How errors surface (exceptions, error codes, result types) is implementation-dependent.

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
| `expect_error` | bool | `true` if the operation must fail |

## Output Matching

- If `expected_output` is set, the output must match exactly after normalizing line endings and trimming trailing whitespace. Metadata lines (`^ ...`) in the output that are not present in `expected_output` are ignored.
- If `accepted_lines` is set, the output is matched group-by-group. Each group is a set of lines that may appear in any order within that group, but groups are consumed sequentially.
- If neither is set and `expect_error` is false, the output must be empty (documents are equal).

## Files

- **core.json** — Core format tests: diff generation across all JSON types with round-trip patch verification
- **options.json** — Options and PathOptions tests
- **errors.json** — Error condition tests

# jd Spec Test Suite

This directory contains a test suite for validating jd format implementations. The tests are implementation-agnostic: test cases express format-level concepts (operations, options, expected output) and the test runner maps them to CLI flags for a given binary.

## Quick Start

```bash
# Build the test runner
go build -o test-runner .

# Test the reference implementation
./test-runner ../../v2/jd/jd

# Test with verbose output
./test-runner -verbose ../../v2/jd/jd

# Test only patch tests
./test-runner -category=patch ../../v2/jd/jd
```

## Test Runner Options

```
Usage: ./test-runner [options] <binary-path>

Options:
  -binary string
        Path to jd binary to test (required)
  -testdata string
        Path to test data directory (default "./testdata")
  -verbose
        Verbose output with detailed results
  -category string
        Filter by test category
  -timeout duration
        Timeout per test case (default 30s)
  -fail-fast
        Stop on first failure
  -opts-flag string
        Flag for passing options to binary (default "-opts")
  -patch-flag string
        Flag for invoking patch mode (default "-p")
  -error-exit int
        Exit code meaning error (default 2)
  -diff-exit int
        Exit code meaning differences found (default 1)
```

The `-opts-flag`, `-patch-flag`, `-error-exit`, and `-diff-exit` flags allow adapting the runner to different CLI implementations that may use different flag names or exit codes.

## Test Categories

### `diff` (diff.json)
Basic diff generation covering all JSON types and structures:
- Object property changes (simple, nested, addition, removal)
- Array operations (append, prepend, remove, context preservation)
- Type changes (object/array/null/boolean/number/string conversions)
- Unicode and special character handling
- Edge cases (empty structures, deeply nested, large numbers, etc.)

### `options` and `path_options` (options.json)
Options system tests:
- SET, MULTISET, precision, setkeys global options
- PathOptions targeting specific paths
- Multiple PathOptions combinations
- DIFF_OFF/DIFF_ON filtering
- Edge cases (boundary precision, null values in sets, duplicate keys)

### `patch` (patch.json)
Patch application tests:
- Object key addition, removal, and value changes
- Array element changes with context validation
- Array append, prepend, and remove operations
- Patching with SET, MULTISET, precision, and SetKeys options
- Deeply nested changes and multiple hunks
- No-op patches, type changes, null handling
- Error cases (context mismatch, invalid syntax)

### `errors` (errors.json)
Format-level error condition tests:
- Malformed JSON input
- Option conflicts (precision with set/multiset)
- Invalid option values
- Patch application failures

## Test Data Structure

```
test/
├── main.go              # Test runner implementation
├── go.mod               # Go module definition
├── README.md            # This file
├── cases/               # Test case definitions (JSON)
│   ├── diff.json        # Diff generation tests
│   ├── options.json     # Options and PathOptions tests
│   ├── patch.json       # Patch application tests
│   └── errors.json      # Error condition tests
└── testdata/            # Test data files
    ├── simple_a.json    # Simple test input A
    ├── simple_b.json    # Simple test input B
    ├── array_a.json     # Array test input A
    ├── array_b.json     # Array test input B
    ├── complex_a.json   # Complex test input A
    └── complex_b.json   # Complex test input B
```

## Test Case Format

Test cases are defined in JSON files with this structure:

```json
{
  "name": "test_case_name",
  "description": "Human readable description",
  "category": "test_category",
  "operation": "diff",
  "content_a": "JSON content for input A",
  "content_b": "JSON content for input B",
  "options": ["SET", {"precision": 0.1}],
  "expected_output": "expected diff output",
  "expect_error": false
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique test case name |
| `description` | string | Human readable description |
| `category` | string | Test category (diff, options, path_options, patch, errors) |
| `operation` | string | `"diff"` (default if omitted) or `"patch"` |
| `content_a` | string | First input (JSON document for diff, diff text for patch) |
| `content_b` | string | Second input (JSON document) |
| `options` | array | Options using metadata line format (e.g. `["SET"]`, `[{"precision": 0.1}]`) |
| `expected_output` | string | Expected output from the operation |
| `accepted_lines` | array | Alternative to `expected_output` for order-insensitive matching |
| `expect_error` | bool | Whether the operation should produce an error |

Either `content_a`/`content_b` or `file_a`/`file_b` should be specified, not both.

### Exit Code Derivation

Exit codes are not specified in test cases. The runner derives them:
- `expect_error: true` → expects the configured error exit code (default 2)
- Diff with non-empty output → expects the configured diff exit code (default 1)
- Diff with empty output → expects 0
- Patch (non-error) → expects 0

## Exit Codes

The test runner uses these exit codes:
- **0**: All tests passed
- **1**: Some tests failed
- **64**: Usage error (invalid arguments)
- **65**: Input data error (malformed test cases)
- **66**: Cannot access binary or test files

## Writing New Tests

To add new test cases:

1. **Identify the category** (diff, options, patch, errors)
2. **Add to appropriate JSON file** in `cases/`
3. **Set `operation`** to `"patch"` for patch tests (diff is default)
4. **Use `options`** for format-level options (not CLI flags)
5. **Provide `expected_output`** for positive test cases
6. **Set `expect_error: true`** for negative test cases
7. **Test your test** by running it against the reference implementation

### Example Test Cases

Diff test:
```json
{
  "name": "my_diff_test",
  "description": "Tests a specific edge case",
  "category": "diff",
  "content_a": "{\"test\": \"value\"}",
  "content_b": "{\"test\": \"modified\"}",
  "expected_output": "@ [\"test\"]\n- \"value\"\n+ \"modified\"\n"
}
```

Options test:
```json
{
  "name": "my_set_test",
  "description": "Tests set comparison",
  "category": "options",
  "content_a": "[1, 2, 3]",
  "content_b": "[3, 1, 2]",
  "options": ["SET"],
  "expected_output": ""
}
```

Patch test:
```json
{
  "name": "my_patch_test",
  "description": "Tests patch application",
  "category": "patch",
  "operation": "patch",
  "content_a": "@ [\"a\"]\n- 1\n+ 2\n",
  "content_b": "{\"a\":1}",
  "expected_output": "{\"a\":2}\n"
}
```

## Common Implementation Issues

Based on testing the reference implementation, watch out for:

1. **Line ending consistency** - Always use LF (`\n`), not CRLF
2. **JSON number formatting** - Normalize `1.0` to `1`
3. **Array context boundaries** - Properly handle `[` and `]` markers
4. **Unicode handling** - Ensure proper UTF-8 encoding/decoding
5. **Option conflicts** - Detect precision with set/multiset conflicts
6. **Path resolution** - Handle edge cases in path element parsing

## Performance Testing

For performance validation, test with:
- Large documents (>1MB JSON files)
- Deep nesting (>100 levels)
- Large arrays (>10000 elements)
- Complex PathOptions configurations

The test runner includes timeout protection (default 30s per test) to catch performance regressions.

## License

This test suite is provided under the same license as the project and can be used freely to validate any structural format implementation.

# jd Blackbox Test Suite

This directory contains a test suite for validating jd format implementations. The tests are designed to be run against any jd-compatible binary to ensure conformance with the specification.

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
```

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
- YAML input mode
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
Error condition tests:
- Malformed JSON input
- Option conflicts (precision with set/multiset)
- Invalid option values
- Patch application failures
- Resource and argument errors

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
  "content_a": "JSON content for input A",
  "content_b": "JSON content for input B",
  "args": ["-set", "-precision=0.1"],
  "expected_diff": "expected diff output",
  "expected_exit": 1,
  "should_error": false
}
```

Either `content_a`/`content_b` OR `file_a`/`file_b` should be specified, not both.

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
3. **Provide expected output** for positive test cases
4. **Set should_error: true** for negative test cases
5. **Test your test** by running it against the reference implementation

### Example Test Case

```json
{
  "name": "my_new_test",
  "description": "Tests a specific edge case I discovered",
  "category": "diff",
  "content_a": "{\"test\": \"value\"}",
  "content_b": "{\"test\": \"modified\"}",
  "expected_diff": "@ [\"test\"]\n- \"value\"\n+ \"modified\"\n",
  "expected_exit": 1
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
7. **Error exit codes** - Use appropriate exit codes for different error types

## Performance Testing

For performance validation, test with:
- Large documents (>1MB JSON files)
- Deep nesting (>100 levels)
- Large arrays (>10000 elements)
- Complex PathOptions configurations

The test runner includes timeout protection (default 30s per test) to catch performance regressions.

## License

This test suite is provided under the same license as the project and can be used freely to validate any structural format implementation.

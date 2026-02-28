# jd Blackbox Test Suite

This directory contains a complete test suite for validating jd format implementations. The tests are designed to be run against any jd-compatible binary to ensure compliance with the specification.

## Quick Start

```bash
# Build the test runner
go build -o test-runner .

# Test the reference implementation
./test-runner ../../v2/jd/jd

# Test with verbose output
./test-runner -verbose ../../v2/jd/jd

# Test only core compliance
./test-runner -core-only ../../v2/jd/jd
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
  -core-only
        Test core compliance only
  -extended-only  
        Test extended compliance only
  -format-only
        Test format compliance only
  -category string
        Filter by test category
  -timeout duration
        Timeout per test case (default 30s)
  -fail-fast
        Stop on first failure
```

## Compliance Levels

### Core Compliance
Basic jd functionality required for all implementations:
- Simple diff generation (object/array changes)
- Path navigation (object keys, array indices)
- Array context preservation
- Basic error handling
- JSON input/output

### Extended Compliance
Advanced features for full-featured implementations:
- Options support (SET, MULTISET, precision, setkeys)
- PathOptions with inheritance rules
- DIFF_ON/DIFF_OFF functionality
- Options header rendering
- YAML input/output

### Format Compliance
Interoperability with other patch formats:
- jd ↔ JSON Patch (RFC 6902) translation
- jd ↔ JSON Merge Patch (RFC 7386) translation
- jd ↔ JSON/YAML translation
- Patch application mode (-p flag)

## Test Categories

### `core/`
- `simple_object_value_change` - Basic property changes
- `nested_object_value_change` - Deep object modifications
- `array_element_change` - Array operations with context
- `type_changes` - Converting between JSON types
- `unicode_support` - Unicode string handling
- `empty_void_states` - Empty documents and void operations

### `options/`
- `set_option_*` - SET option behavior
- `multiset_option_*` - MULTISET option behavior  
- `precision_option_*` - Numeric precision tolerance
- `setkeys_option_*` - Object matching by keys

### `path_options/`
- `path_options_set` - Targeted SET operations
- `path_options_precision` - Path-specific precision
- `multiple_path_options` - Multiple PathOptions
- `diff_off_option` - Ignoring specific paths
- `diff_on_allowlist` - Allow-list approaches

### `format/`
- `patch_format_*` - JSON Patch output generation
- `merge_format_*` - JSON Merge Patch output
- `translate_*` - Format translation operations
- `patch_mode_*` - Patch application testing

### `errors/`
- `invalid_json_*` - Malformed input handling
- `option_conflicts` - Incompatible option combinations
- `patch_errors` - Patch application failures
- `resource_limits` - Resource exhaustion scenarios

### `edge_cases/`
- `unicode_*` - Unicode edge cases
- `numeric_precision` - Floating-point edge cases
- `empty_structures` - Empty objects/arrays
- `special_characters` - Keys with special chars
- `deeply_nested` - Deep structure nesting

## Test Data Structure

```
test/
├── main.go              # Test runner implementation
├── go.mod               # Go module definition
├── README.md            # This file
├── cases/               # Test case definitions (JSON)
│   ├── core.json        # Core functionality tests
│   ├── extended.json    # Extended feature tests  
│   ├── format.json      # Format translation tests
│   ├── errors.json      # Error condition tests
│   └── edge_cases.json  # Edge case tests
├── testdata/            # Test data files
│   ├── simple_a.json    # Simple test input A
│   ├── simple_b.json    # Simple test input B
│   ├── array_a.json     # Array test input A
│   ├── array_b.json     # Array test input B
│   ├── complex_a.json   # Complex test input A
│   └── complex_b.json   # Complex test input B
└── golden/              # Expected output files (future use)
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
  "file_a": "path/to/file_a.json",
  "file_b": "path/to/file_b.json", 
  "args": ["-set", "-precision=0.1"],
  "expected_diff": "expected diff output",
  "expected_exit": 1,
  "should_error": false,
  "compliance_level": "core"
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

1. **Identify the category** (core, extended, format, errors, edge_cases)
2. **Add to appropriate JSON file** in `cases/`
3. **Set compliance level** appropriately
4. **Provide expected output** for positive test cases
5. **Set should_error: true** for negative test cases
6. **Test your test** by running it against the reference implementation

### Example Test Case

```json
{
  "name": "my_new_test",
  "description": "Tests a specific edge case I discovered",
  "category": "edge_cases", 
  "content_a": "{\"test\": \"value\"}",
  "content_b": "{\"test\": \"modified\"}",
  "expected_diff": "@ [\"test\"]\n- \"value\"\n+ \"modified\"\n",
  "expected_exit": 1,
  "compliance_level": "core"
}
```

## Implementation Validation

To validate your jd implementation:

1. **Start with core compliance**: `./test-runner -core-only your-binary`
2. **Add extended features**: `./test-runner -extended-only your-binary`
3. **Implement format translation**: `./test-runner -format-only your-binary`
4. **Run full suite**: `./test-runner your-binary`
5. **Fix any failures** and retest

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

## Extending the Test Suite

The test suite is designed to grow. Consider adding:
- Additional edge cases as they're discovered
- Performance regression tests  
- Memory usage validation
- Compatibility tests with other JSON diff tools
- Fuzzing-based test case generation

## License

This test suite is provided under the same license as the project and can be used freely to validate any structural format implementation.
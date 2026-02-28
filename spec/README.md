# Structural Format Specification

This directory contains the formal specification for the structural JSON diff format, providing the complete syntax, semantics, and behavior definitions needed for implementation across programming languages.

## Documents

### Core Specification
- **[jd-format.md](jd-format.md)** - Complete formal specification of the structural diff format
- **[grammar.md](grammar.md)** - ABNF grammar definition for parsing and generation
- **[semantics.md](semantics.md)** - Semantic definitions for diff operations and data structures
- **[errors.md](errors.md)** - Error handling procedures and taxonomy

### Reference Materials
- **[examples.md](examples.md)** - Complete examples covering all features and edge cases

### Test Suite
- **[test/](test/)** - Blackbox test suite for implementation validation

## Implementation Guide

To implement the structural format:

1. **Start with [jd-format.md](jd-format.md)** - Read the complete specification
2. **Study [grammar.md](grammar.md)** - Implement the parser using the ABNF grammar
3. **Review [semantics.md](semantics.md)** - Understand operational semantics
4. **Handle [errors.md](errors.md)** - Implement proper error handling
5. **Test with [test/](test/)** - Validate your implementation

## Testing Your Implementation

The test suite in `test/` provides validation across all features:

```bash
cd test
go build -o test-runner .
./test-runner /path/to/your/structural/binary
```

Test cases are organized by feature:
- **diff.json** - Diff generation tests
- **options.json** - Options and PathOptions tests
- **patch.json** - Patch application tests
- **errors.json** - Error handling tests

Exit code 0 indicates all tests passed. Non-zero indicates failures with detailed reporting.

## About the Structural Format

The structural format is a human-readable diff format for JSON and YAML data with these features:

- **Human-readable**: Unified diff-style output
- **Context-aware**: Shows surrounding elements for change location clarity
- **Set semantics**: Treats arrays as sets or multisets when order doesn't matter
- **Configurable**: Supports numeric precision tolerance for floating-point comparisons
- **Flexible**: PathOptions enable fine-grained comparison control

## Version

This specification documents the structural format, as implemented in the reference Go library at https://github.com/josephburnett/jd/v2.

## License

This specification is provided under the same license as the reference implementation.

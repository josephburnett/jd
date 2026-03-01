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

### Test Data and Runner
- **[cases/](cases/)** - Implementation-agnostic test data (JSON files describing inputs and expected outputs)
- **[test/](test/)** - Reference test runner that executes test data against a CLI binary

## Structure

```
spec/
├── README.md
├── jd-format.md, grammar.md, semantics.md, errors.md, examples.md
├── cases/               # Test data (part of the spec)
│   ├── README.md
│   ├── core.json
│   ├── options.json
│   └── errors.json
└── test/                # Reference test runner (Go)
    ├── main.go
    ├── go.mod
    └── README.md
```

## Implementation Guide

To implement the structural format:

1. **Start with [jd-format.md](jd-format.md)** - Read the complete specification
2. **Study [grammar.md](grammar.md)** - Implement the parser using the ABNF grammar
3. **Review [semantics.md](semantics.md)** - Understand operational semantics
4. **Handle [errors.md](errors.md)** - Implement proper error handling
5. **Validate with [cases/](cases/)** - Test your implementation against the test data

## Testing Your Implementation

The test data in `cases/` defines expected behavior independent of any CLI, library, or language. See [cases/README.md](cases/README.md) for the verification procedure — it describes how to use the test cases to validate both diff and patch operations.

For CLI-based implementations, the reference test runner in `test/` automates this process:

```bash
cd test
go build -o test-runner .
./test-runner /path/to/your/binary
```

See [test/README.md](test/README.md) for runner flags and how to adapt it to different CLI interfaces.

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

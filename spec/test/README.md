# Reference Test Runner

A Go tool that runs the spec's test cases against a CLI binary. It maps format-level test concepts (operations, options, expected output) to CLI invocations with flags, temp files, and exit code checks.

## Building

```bash
go build -o test-runner .
```

## Usage

```bash
# Test the reference implementation
./test-runner /path/to/jd

# Verbose output
./test-runner -verbose /path/to/jd

# Filter by category
./test-runner -category=diff /path/to/jd

# Custom cases directory
./test-runner -cases /path/to/cases /path/to/jd
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-binary` | (positional) | Path to binary under test |
| `-cases` | `../cases` | Path to test cases directory |
| `-testdata` | `./testdata` | Path to file-based test data |
| `-verbose` | `false` | Print per-test results and generate report |
| `-category` | (all) | Filter by test category |
| `-timeout` | `30s` | Timeout per test case |
| `-fail-fast` | `false` | Stop on first failure |

## CLI Mapping Flags

These flags adapt the runner to binaries with different flag names or exit codes:

| Flag | Default | Description |
|------|---------|-------------|
| `-opts-flag` | `-opts` | Flag name for passing options |
| `-patch-flag` | `-p` | Flag name for patch mode |
| `-error-exit` | `2` | Exit code indicating error |
| `-diff-exit` | `1` | Exit code indicating differences found |

## How It Works

For each test case, the runner:

1. Writes `content_a` and `content_b` to temp files
2. Builds a command: `binary [opts-flag=options] [patch-flag] <file_a> <file_b>`
3. Checks the exit code matches expectations (derived from the test case, not specified directly)
4. Compares output against `expected_output` or `accepted_lines`
5. For diff tests that produce output, verifies the round-trip property: patching `content_a` with the diff output yields `content_b`

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All tests passed |
| 1 | Test failures |
| 64 | Usage error |
| 65 | Malformed test data |
| 66 | Cannot access binary |

## Using with Other Implementations

To test a non-Go implementation, build this runner and point it at your binary. Use the CLI mapping flags if your binary uses different flag names or exit codes:

```bash
./test-runner -opts-flag=--options -patch-flag=--patch -error-exit=1 /path/to/your/binary
```

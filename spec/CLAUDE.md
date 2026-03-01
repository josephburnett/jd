# Spec

The spec defines the JD diff format for any implementation — Go, Rust, JavaScript,
SQL functions, or anything else. It is not concerned with the reference implementation's
backward compatibility, CLI flags, or file-based tooling.

## What the Spec Covers

The core value of the JD format: human-readable structural diffs with set semantics
and structured patching. Both diffing and patching are equally important.

- Diff and patch semantics for JSON values
- Set and multiset collection types
- The `keys` option for object identity in arrays
- Precision option for numeric comparison (with restrictions — see below)
- PathOptions for targeting specific paths with different semantics
- Error conditions that implementations MUST reject
- Undefined behaviors that implementations MAY handle however they choose

## What the Spec Does NOT Cover

- **Translation** between formats (RFC 6902, RFC 7386) — out of scope
- **MERGE** — the spec doesn't define how merge diffs are created, so it shouldn't
  define how they are applied
- **COLOR** or any presentation concern — per-implementation metadata
- **File metadata** (`^ {"file":"..."}`) — not all implementations are file-based
- **CLI flags, exit codes** — these are reference implementation concerns. The spec
  test data is expressed in format-level concepts (operations, options, expected output)

Any of these can exist as per-implementation `^` metadata, which the spec's
forward-compatibility mechanism handles by silently skipping unrecognized lines.

## Mathematical Constraints

Set semantics require a proper equivalence relation: reflexive, symmetric, and
**transitive**. Approximate equality (e.g. precision) breaks transitivity —
a ≈ b and b ≈ c does not imply a ≈ c. Therefore equivalence modifiers and set
semantics are fundamentally incompatible. This is not negotiable or implementation-
dependent; it follows from the math.

## Spec Strictness Philosophy

- **Static configuration errors** (empty keys, negative precision, incompatible option
  combinations): MUST reject. These are detectable before processing data.
- **Data-dependent runtime issues** (duplicate identity keys): undefined behavior.
  Document it explicitly, but don't mandate a specific response since not all APIs
  can return errors mid-operation.
- **Unknown metadata**: silently skip. This enables forward compatibility.

## Naming

The spec uses the canonical names going forward. For example, the option is `"keys"`,
not `"setkeys"`. The reference implementation may accept legacy aliases for backward
compatibility, but the spec text and test data should only use the canonical form.

## Test Data

Test data in `/spec/cases/` is implementation-agnostic. It expresses format-level
concepts: `operation`, `options`, `expected_output`, `expect_error`. No CLI flags,
no exit codes, no file paths.

The Go test binary in `/spec/test/` is a reference harness that maps test data to a
specific CLI binary. Other implementations should be able to consume the same test
data with their own harness.

Round-trip testing (diff, patch, diff-produces-no-output) is the preferred way to
verify patch semantics without writing separate patch test cases for every diff test.

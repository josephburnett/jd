# v2 — Reference Implementation

This is the Go reference implementation of the JD diff format. It serves as both
a usable library/CLI and the reference for spec compliance.

## Backward Compatibility

- Public API methods must not change without instruction. The API surface should be
  deliberate and minimal.
- The CLI accepts legacy option names (e.g. `-setkeys` alongside `-keys`) but the
  library always renders the canonical form (e.g. `"keys"` not `"setkeys"`).
- The `^` metadata parser tolerates unrecognized lines (skips them via `continue` in
  diff_read.go). This means older versions can read diffs produced by newer versions
  that include new metadata like `^ {"file":"..."}`.

## Testing

- `make test` — unit tests
- `make vet` — static analysis (excludes WASM package automatically)
- `make cover` — enforces 100% coverage on non-trivial code
- `make fuzz` — fuzz testing

Coverage enforcement uses `v2/cover_exclude.txt` to exclude trivial and unreachable
code (I/O wrappers, YAML support, defensive guards in closed type switches). Everything
else must be at 100%. When bugs are found, add tests that cover the gap.

## Key Internal Details

- **Myers diff** kicks in for arrays where both sides have >10 elements. Below that
  threshold, LCS is used directly.
- **Options plumbing**: `checkOption` reads `retain`, `getOption` reads `apply`. These
  are different fields in the options struct — easy to mix up.
- **`refine()`** propagates options during recursive descent. It silently drops unknown
  option types, which is correct — informational metadata like `fileOption` shouldn't
  affect diff computation.
- **`Diff()` returns `Diff`, not `(Diff, error)`.** Errors can't be returned during
  diffing, which is why data-dependent issues like duplicate identity keys are undefined
  behavior rather than errors.
- **WASM**: `go vet ./...` fails on `internal/web/ui` because it needs `GOOS=js
  GOARCH=wasm`. The Makefile targets handle this exclusion.

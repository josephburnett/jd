# jd

A commandline tool, Go library, and website for diffing and patching JSON and YAML.
The JD diff format is human-readable and supports set semantics — things no other
diff/patch format provides.

## Design Principles

1. **Spec and implementation are separate concerns.** The spec (`/spec`) defines
   forward-looking behavior for any implementation. The reference implementation (`/v2`)
   has its own backward compatibility and engineering concerns. Don't conflate the two.
2. **Minimal surface area.** Applies to the API, dependencies, and what the spec covers.
   Focus on what's uniquely valuable: human-readable diffs, set semantics, structured patching.
3. **Correctness is mathematical.** Set semantics require equivalence relations (reflexive,
   symmetric, transitive). Approximate equality (e.g. precision) breaks transitivity, so
   it is fundamentally incompatible with set operations. Constraints like these aren't
   arbitrary — they follow from the math.
4. **Permissive input, strict output.** Unknown `^` metadata is silently skipped so
   implementations can extend without breaking each other. But the spec is strict about
   what it does define: static configuration errors (e.g. empty keys, negative precision)
   MUST be rejected. Data-dependent runtime issues (e.g. duplicate identity keys) are
   documented as undefined behavior.

## Key Commands

Always use Makefile targets, not raw `go test`, `go vet`, etc. The Makefile handles
WASM exclusions and other project-specific build concerns.

- Test: `make test`
- Vet: `make vet`
- Format: `make go-fmt`
- Coverage: `make cover` (enforces 100% on non-trivial code)
- Fuzz: `make fuzz`
- Spec tests: `make spec-test`

## Code Standards

- MUST pass unit tests
- MUST be go formatted
- Add tests to existing table tests when possible
- When bugs are found, always add unit tests covering the gap
- Delete dead code rather than leaving it uncovered
- API surface should be deliberate and minimal — do not add or change public methods
  without instruction

## File Structure

- `README.md` - Installation, usage, and diff format documentation
- `Makefile` - All build, test, and validation targets
- `/v2` - Reference implementation (Go library + CLI). See `/v2/CLAUDE.md`
- `/spec` - Format specification and test data. See `/spec/CLAUDE.md`
- `/action.yml` - GitHub Action
- `/doc` - Plans and documents
- `/lib` - Deprecated v1 library (read-only)
- `main.go` - Deprecated v1 commandline (read-only)

## Supply Chain Security

- Dependencies MUST be minimal. Do not add new dependencies without justification.
- GitHub Actions MUST be pinned to full commit SHA with version comment.
- When bumping Go toolchain: update `GOTOOLCHAIN` in Makefile, `toolchain` in both go.mod
  files, and `FROM golang:` in Dockerfile. `make validate-toolchain` checks all of these.
- `wasm_exec.js` is copied from GOROOT at build time — no manual update needed.
- Run `govulncheck ./...` periodically to check for stdlib and dependency vulns.

## Advice

- Try and avoid creating temporary files and instead rely on modifying existing unit tests
  for debugging. This prevents getting stuck on tool approval. Use tools that don't need
  approval to unblock yourself.
- Don't use words like "comprehensive" or "robust" — they don't add specificity.
- Read `~/joe-voice.md` before writing GitHub issues, PR comments, or anything on the
  owner's behalf.
- Documentation is in `README.md` and `v2/jd/main.go` (see usage).

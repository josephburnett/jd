# Project Overview

This project provides both a commandline, library and website for diffing and patching JSON and YAML values.
It is focused on a human-readable diff format but has the ability to read and produce other diff formats.

## Key Commands

- Test: `make test`
- Format: `make go-fmt`
- Fuzz: `make fuzz`

## Code Standards

- MUST pass unit tests
- MUST be go formatted
- Add tests to existing table tests when possible
- Public methods MUST not change (maintain API backward compatability) *
- New public methods MUST NOT be created (do not expand the API) *

  * Unless otherwise instructed

## File Structure

- `README.md` - Information about installation, usage and diff format
- `Makefile` - Project commands
- `/v2` - V2 library
- `/v2/jd/main.go` - V2 commandline
- `/v2/web` - Website and associated tools
- `/action.yml` - GitHub Action
- `/doc` - Plans and documents
- `/lib` - Deprecated v1 library (read-only)
- `main.go` - Deprecated v1 commandline (read-only)

## Advice

- Try and avoid creating temporary files and instead rely on modifying existing unit tests for debugging.
  This will prevent you from getting stuck on asking for my approval. Try and use tools you don't need
  to ask approval for so you can unblock yourself.
- Don't use words like "comprehensive" or "robust" because they don't add anything to the specificity of
  a sentence. Each layer of testing adds a layer of probability to catch a bug, but few things are truely
  comprehensive. And robustness again is very relative and requires context, such as SLAs.
- Documentation is in `README.md` and `v2/jd/main.go` (see usage).

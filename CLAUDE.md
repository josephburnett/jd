# Project Overview

This project provides both a commandline, library and website for diffing and patching JSON and YAML values.
It is focused on a human-readable diff format but has the ability to read and produce other diff formats.

## Key Commands

- Test: `make test`
- Format: `go fmt ./...`
- Fuzz: `make fuzz`

## Code Standards

- MUST pass unit tests
- MUST be go formatted
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

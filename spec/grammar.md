# Structural Format Grammar

This document provides the formal ABNF grammar specification for the structural JSON diff format.

## Overview

The structural diff format is a text-based representation of changes between JSON documents. It consists of optional metadata headers followed by diff elements that specify locations and changes.

## Core Grammar

```abnf
; Top-level document structure
StructuralDiff = *MetadataLine *DiffElement

; Metadata/options headers
MetadataLine = "^" SP JsonValue CRLF

; Main diff elements
DiffElement = PathLine [ArrayOpen] *ContextLine *ChangeLine [*ContextLine] [ArrayClose]

; Path specification
PathLine = "@" SP JsonArray CRLF

; Array context markers
ArrayOpen = "[" CRLF
ArrayClose = "]" CRLF

; Content lines
ContextLine = SP SP JsonValue CRLF
ChangeLine = (AddLine / RemoveLine)
AddLine = "+" SP [JsonValue] CRLF
RemoveLine = "-" SP JsonValue CRLF

; JSON array for paths (restricted form)
JsonArray = "[" [PathElement *(", " PathElement)] "]"

; Path elements
PathElement = JsonString         ; Object key
            / JsonNumber         ; Array index  
            / EmptyObject        ; Set marker
            / EmptyArray         ; List marker
            / SetKeysObject      ; Set with matching keys
            / MultisetContainer  ; Multiset marker

; Special path element types
EmptyObject = "{}"
EmptyArray = "[]"
SetKeysObject = "{" KeyValuePair *(", " KeyValuePair) "}"
MultisetContainer = "[" [EmptyObject / SetKeysObject] "]"

; Key-value pairs for set keys
KeyValuePair = JsonString ":" JsonValue

; Standard JSON components (simplified for path elements)
JsonValue = JsonString / JsonNumber / JsonObject / JsonArray / JsonBool / JsonNull / JsonVoid
JsonString = DQUOTE *JsonChar DQUOTE
JsonNumber = ["-"] (("0" / (DIGIT1-9 *DIGIT)) ["." 1*DIGIT] [("e" / "E") ["+" / "-"] 1*DIGIT])
JsonObject = "{" [JsonString ":" JsonValue *(", " JsonString ":" JsonValue)] "}"
JsonArray = "[" [JsonValue *(", " JsonValue)] "]"
JsonBool = "true" / "false"
JsonNull = "null"
JsonVoid = ""  ; Empty value for merge operations

; Character definitions
JsonChar = %x20-21 / %x23-5B / %x5D-10FFFF / EscapeSequence
EscapeSequence = "\" ("\"" / "/" / "b" / "f" / "n" / "r" / "t" / UnicodeEscape)
UnicodeEscape = "u" 4HEXDIG

; Core ABNF rules  
CRLF = LF          ; Line ending (LF only, not CRLF)
LF = %x0A          ; Line feed
SP = %x20          ; Space
DQUOTE = %x22      ; Double quote
DIGIT = %x30-39    ; 0-9
DIGIT1-9 = %x31-39 ; 1-9
HEXDIG = DIGIT / "A" / "B" / "C" / "D" / "E" / "F" / "a" / "b" / "c" / "d" / "e" / "f"
```

## Metadata Line Grammar

Metadata lines provide options and configuration for the diff:

```abnf
; Option types
MetadataOption = SimpleOption / ObjectOption / PathOption

; Simple string options
SimpleOption = %s"SET" / %s"MULTISET" / %s"MERGE" / %s"COLOR" / %s"DIFF_ON" / %s"DIFF_OFF"

; Complex object options  
ObjectOption = PrecisionOption / SetKeysOption / LegacyMergeOption

PrecisionOption = "{" %s"\"precision\"" ":" JsonNumber "}"
SetKeysOption = "{" %s"\"setkeys\"" ":" JsonArray "}"
LegacyMergeOption = "{" %s"\"Merge\"" ":" JsonBool "}"

; Path-specific options
PathOption = "{" %s"\"@\"" ":" JsonArray ", " %s"\"^\"" ":" "[" MetadataOption *(", " MetadataOption) "]" "}"
```

## Path Element Specifications

### Object Keys
```abnf
ObjectKey = JsonString
; Examples: "name", "user_id", "nested.key"
```

### Array Indices
```abnf
ArrayIndex = JsonNumber
; Must be non-negative integer or -1 for append
; Examples: 0, 5, 42, -1
```

### Set Operations
```abnf
SetMarker = EmptyObject    ; {} - operate on any set element
SetWithKeys = SetKeysObject ; {"id":"value"} - match by specific keys
```

### Multiset Operations  
```abnf
MultisetMarker = "[" "]"                    ; [] - list marker
MultisetWithObject = "[" EmptyObject "]"    ; [{}] - multiset of any
MultisetWithKeys = "[" SetKeysObject "]"    ; [{"key":"val"}] - multiset with keys
```

## Line Type Specifications

### Path Lines
```abnf
PathLine = "@" SP "[" PathSequence "]" CRLF
PathSequence = [PathElement *(", " PathElement)]
```

### Change Lines
```abnf
; Addition (may have empty value for void)
AddLine = "+" SP [JsonValue] CRLF

; Removal (must have value)
RemoveLine = "-" SP JsonValue CRLF

; Context (two spaces, then value)
ContextLine = SP SP JsonValue CRLF
```

### Array Context
```abnf
; Array context markers (contextual usage)
ArrayOpen = "[" CRLF     ; Only when showing array beginning as context
ArrayClose = "]" CRLF    ; Only when showing array end as context
```

**Important**: Array boundary markers are contextual:
- `[` appears only when the diff shows changes at or near the beginning of an array
- `]` appears only when the diff shows changes at or near the end of an array  
- Middle array changes don't require boundary markers since array indices provide sufficient context

## Whitespace and Formatting Rules

1. **Line Endings**: LF only (`\n`), not CRLF
2. **Indentation**: 
   - Metadata lines: No indentation before `^`
   - Path lines: No indentation before `@`
   - Context lines: Exactly two spaces before content
   - Change lines: One space between `+`/`-` and content
3. **JSON Formatting**: Standard JSON syntax within JsonValue productions
4. **Unicode**: Full Unicode support with proper JSON string escaping

## Grammar Extensions

### Color Support
When `COLOR` option is present, implementations MAY add ANSI color codes to change lines while preserving the grammar structure.

### Legacy Compatibility  
- `{"Merge":true}` metadata is equivalent to `"MERGE"` option
- Implementations SHOULD normalize to modern format when rendering

## Validation Rules

1. **Path Validity**: Path elements must form valid JSON property/index chains
2. **Value Consistency**: JSON values must be syntactically valid
3. **Context Preservation**: Array contexts must maintain proper opening/closing
4. **Option Conflicts**: Implementations SHOULD detect conflicting options (e.g., precision with set operations)

## Implementation Notes

1. **Parser Requirements**: Must handle UTF-8 encoded text
2. **Memory Limits**: Implementations MAY impose reasonable limits on nesting depth and value sizes
3. **Error Handling**: Syntax errors SHOULD provide line and column information
4. **Extensibility**: Unknown options SHOULD be preserved but ignored during processing

This grammar provides the complete syntactic specification for parsing and generating structural diff format. For semantic interpretations of these constructs, see [semantics.md](semantics.md).
# Structural Format Error Handling Specification

This document defines error conditions, error messages, and error handling procedures for structural format implementations.

## Error Categories

### 1. Parsing Errors

Errors that occur during diff format parsing or JSON document parsing.

#### Diff Format Parsing Errors

| Error Code | Description | Example |
|------------|-------------|---------|
| `DIFF_SYNTAX_ERROR` | Invalid diff format syntax | Missing `@` in path line |
| `INVALID_PATH_SYNTAX` | Malformed path array | `@ [unclosed` |
| `INVALID_METADATA` | Malformed metadata line | `^ invalid-json` |
| `MISSING_PATH` | Diff element without path | Change lines without `@` line |
| `INVALID_LINE_PREFIX` | Unrecognized line prefix | `? unknown prefix` |

**Example Syntax Error:**
```
Input:  @ ["path"
Output: DIFF_SYNTAX_ERROR: Unterminated path array at line 1, column 10
```

#### JSON Parsing Errors

| Error Code | Description | Example |
|------------|-------------|---------|
| `JSON_SYNTAX_ERROR` | Invalid JSON in document | `{"key": invalid}` |
| `JSON_ENCODING_ERROR` | Invalid UTF-8 encoding | Binary data in JSON string |
| `JSON_NESTING_ERROR` | Excessive nesting depth | 1000+ levels deep |
| `JSON_SIZE_ERROR` | Document exceeds size limits | File > implementation limit |

### 2. Path Resolution Errors

Errors that occur when navigating document paths.

| Error Code | Description | Recovery Action |
|------------|-------------|----------------|
| `PATH_NOT_FOUND` | Path does not exist in document | Return null or error |
| `PATH_TYPE_MISMATCH` | Path expects different type | Convert or error |
| `ARRAY_INDEX_OUT_OF_BOUNDS` | Array index exceeds length | Extend array or error |
| `OBJECT_KEY_NOT_FOUND` | Object key does not exist | Create key or error |
| `PATH_ELEMENT_INVALID` | Invalid path element type | Reject path |

**Example Path Error:**
```
Document: {"users": [{"name": "Alice"}]}
Path: ["users", "invalid_index", "name"]  
Error: PATH_TYPE_MISMATCH: Expected array index, got string "invalid_index" at path ["users"]
```

### 3. Patch Application Errors

Errors during patch application to documents.

| Error Code | Description | When It Occurs |
|------------|-------------|----------------|
| `PATCH_CONTEXT_MISMATCH` | Context doesn't match document | Array context validation fails |
| `PATCH_PRECONDITION_FAILED` | Test operation failed | Value doesn't match expected |
| `PATCH_CONFLICT` | Conflicting changes | Multiple operations on same path |
| `PATCH_TYPE_ERROR` | Type mismatch during application | Adding string to number |
| `PATCH_INVALID_OPERATION` | Unsupported operation | Unknown operation type |

**Example Context Mismatch:**
```
Document: ["a", "b", "c"]
Diff:
@ [1]
[
  "x"    <- Expected context "a", found "a" ✓
- "y"    <- Expected "b", attempting to remove ✗
+ "z"
]

Error: PATCH_CONTEXT_MISMATCH: Expected context "y" at path [0], found "a"
```

### 4. Option Processing Errors

Errors related to option parsing and conflicts.

| Error Code | Description | Example |
|------------|-------------|---------|
| `OPTION_PARSE_ERROR` | Cannot parse option syntax | `{"precision": "invalid"}` |
| `OPTION_CONFLICT` | Incompatible options | Precision with SET |
| `OPTION_VALUE_ERROR` | Invalid option value | Negative precision |
| `PATH_OPTION_INVALID` | Invalid PathOption syntax | Missing `@` or `^` fields |
| `OPTION_NOT_SUPPORTED` | Unrecognized option | Unknown option name |

**Example Option Conflict:**
```
Options: ["SET", {"precision": 0.01}]
Error: OPTION_CONFLICT: Precision option incompatible with SET (uses hash comparison)
```

### 5. Resource Limit Errors

Errors due to resource exhaustion or limits.

| Error Code | Description | Typical Limit |
|------------|-------------|---------------|
| `MEMORY_LIMIT_EXCEEDED` | Out of memory | Implementation-dependent |
| `RECURSION_LIMIT_EXCEEDED` | Too much nesting | 1000 levels |
| `SIZE_LIMIT_EXCEEDED` | Document too large | 100MB |
| `TIME_LIMIT_EXCEEDED` | Operation timeout | 30 seconds |
| `COMPLEXITY_LIMIT_EXCEEDED` | Algorithm complexity too high | LCS on huge arrays |

## Error Message Format

### Standard Error Message Structure

```
ERROR_CODE: Brief description
  Location: file:line:column or path information
  Expected: what was expected (if applicable)  
  Found: what was actually encountered
  Context: surrounding context (first 50 chars)
```

### Examples

#### Parsing Error
```
DIFF_SYNTAX_ERROR: Invalid path array syntax
  Location: line 3, column 15
  Expected: ']' to close path array
  Found: end of line
  Context: @ ["users", 0, "name"
```

#### Path Resolution Error
```
PATH_TYPE_MISMATCH: Cannot index non-array value
  Location: path ["config", "settings", 0]
  Expected: array value for index access
  Found: string "development"
  Context: {"config": {"settings": "development"}}
```

#### Patch Application Error  
```
PATCH_CONTEXT_MISMATCH: Context validation failed
  Location: path ["items", 2]
  Expected: context value "apple"
  Found: "banana"
  Context: ["orange", "banana", "cherry"]
```

## Exit Codes

### Command-line Tool Exit Codes

| Exit Code | Meaning | Description |
|-----------|---------|-------------|
| 0 | Success | Operation completed successfully (no differences or successful operation) |
| 1 | Differences found | Normal diff operation found differences between inputs |
| 2 | Error occurred | Any error condition (parsing, file access, option conflicts, etc.) |

**Rationale**: The structural format uses a simplified exit code scheme rather than Unix-style codes. This provides sufficient information for automation while remaining simple and consistent.

### Library Error Codes

Libraries should use appropriate error handling mechanisms for their language:
- **Go**: Return error values with typed error structs
- **Python**: Raise specific exception classes
- **JavaScript**: Throw Error objects with `code` property
- **Java**: Throw checked exceptions with error codes
- **Rust**: Return `Result<T, JdError>` types

## Error Recovery Strategies

### 1. Graceful Degradation

When possible, implementations should continue operation with limited functionality:

- **Unknown options**: Ignore and warn, continue processing
- **Unsupported path elements**: Fall back to string comparison
- **Minor syntax errors**: Attempt to recover and continue

### 2. Fail-Fast Behavior

For critical errors, fail immediately:

- **Invalid JSON**: Cannot proceed with malformed input
- **Memory exhaustion**: Prevent system instability
- **Security violations**: Prevent potential exploits

### 3. Context Preservation

Provide sufficient context for error diagnosis:

- **Line/column numbers** for parsing errors
- **Path information** for resolution errors  
- **Surrounding data** for context understanding
- **Full error chain** for nested errors

## Debugging Support

### Debug Output Format

Implementations SHOULD support verbose error reporting:

```
--debug flag enables:
- Step-by-step operation logging
- Internal state dumps at error points  
- Full path resolution traces
- Option inheritance chain
- Algorithm decision points
```

### Error Context Information

Include in error reports:
- **Input documents** (or snippets)  
- **Applied options** and their sources
- **Path resolution history**
- **Algorithm choices** made
- **Resource usage** at error time

## Implementation Guidelines

### Error Handling Best Practices

1. **Validate early**: Check inputs before processing
2. **Provide context**: Include location and surrounding data
3. **Use specific codes**: Enable programmatic error handling
4. **Log selectively**: Balance verbosity with usefulness
5. **Test error paths**: Ensure error handling works correctly

### Error Message Localization

- Use error codes for programmatic handling
- Separate user-facing messages from codes
- Support internationalization where appropriate
- Maintain English as fallback language

### Performance Considerations

- **Lazy error collection**: Don't compute expensive error details unless needed
- **Error caching**: Cache repeated error condition checks
- **Resource cleanup**: Ensure errors don't leak resources
- **Controlled shutdown**: Handle interruption during error states

### Security Considerations

- **Information disclosure**: Don't expose sensitive data in error messages
- **Resource limits**: Prevent resource exhaustion attacks through error paths
- **Input validation**: Thoroughly validate all inputs to prevent injection attacks
- **Error amplification**: Prevent error conditions from cascading

## Testing Error Conditions

### Required Error Test Coverage

Implementations MUST test:
1. All documented error conditions
2. Error message format consistency
3. Proper exit code usage
4. Resource cleanup after errors
5. Error recovery mechanisms

### Error Test Categories

- **Boundary conditions**: Edge cases that trigger errors
- **Malformed inputs**: Invalid JSON, diff syntax
- **Resource exhaustion**: Large inputs, deep nesting
- **Option conflicts**: Incompatible option combinations
- **Context validation**: Mismatched patch contexts

This error specification ensures consistent, debuggable error handling across all structural format implementations while providing specific guidance for both implementors and users.
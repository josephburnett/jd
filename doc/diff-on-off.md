# Diffing On/Off Feature Implementation

## Overview
This document describes the implementation of a feature that allows turning diffing on and off using PathOptions. When diffing is turned off, differences are ignored and no diff events are generated. When diffing is turned on, differences are produced as they normally would be. This feature enables selective ignoring of parts of a JSON/YAML document (e.g., system-generated values) or conversely allows selective inclusion of parts by turning off diffing at the root and then selectively turning it back on.

## Key Requirements
1. **PathOptions Control**: Use PathOptions to turn diffing on/off at specific paths
2. **Nesting Support**: Diffing on/off can be nested - lower-level PathOptions can override higher-level settings
3. **Order Preservation**: PathOptions are processed preserving order, so the last directive takes precedence
4. **Default Behavior**: By default, diffing is ON (current behavior preserved)
5. **Selective Ignoring**: Can ignore system-generated values by turning diffing off for specific paths
6. **Allow-list Mode**: Can turn off diffing at root, then selectively turn it on for specific paths

## Architecture Overview

### Current Event-Driven System
The jd library uses an event-driven diffing system:
- Each node type generates diff events when comparing two values
- Events are processed by specialized processors that generate DiffElements
- PathOptions are refined as the diff traverses the structure

### New Components Added

#### 1. New Option Types
- `DIFF_ON`: Explicitly enables diffing for a path
- `DIFF_OFF`: Explicitly disables diffing for a path

#### 2. Options State Tracking
- Add `diffingOn bool` field to `options` struct
- Track current diffing state as PathOptions are processed
- Default to `true` (current behavior)

#### 3. Event Generation Control
- Check diffing state before generating events
- Return empty event lists when diffing is off
- Respect nested PathOption overrides

## Implementation Details

### 1. Option Types (v2/options.go)
```go
type diffOnOption struct{}
var DIFF_ON = diffOnOption{}

type diffOffOption struct{}  
var DIFF_OFF = diffOffOption{}
```

Added to `NewOption()` function:
```go
case "DIFF_ON":
    return DIFF_ON, nil
case "DIFF_OFF":
    return DIFF_OFF, nil
```

### 2. Options State Management
Updated `options` struct:
```go
type options struct {
    apply     []Option
    retain    []Option
    diffingOn bool  // New: tracks current diffing state
}
```

Updated `refine()` function to:
- Process DIFF_ON/DIFF_OFF options in order
- Set `diffingOn` field based on final state
- Handle these as global options like SET/MULTISET

### 3. Event Generation Updates
All event generation functions check diffing state:

**generateSimpleEvents()**:
```go
func generateSimpleEvents(a, b JsonNode, opts *options) []diffEvent {
    if !opts.diffingOn {
        return []diffEvent{} // No events when diffing is off
    }
    // Existing logic...
}
```

**generateObjectdiffEvents()**:
```go
func generateObjectdiffEvents(o1, o2 jsonObject, opts *options) []diffEvent {
    if !opts.diffingOn {
        return []diffEvent{} // No events when diffing is off
    }
    // Existing logic...
}
```

### 4. Core Diffing Logic
Updated `diff()` function in diff_common.go:
```go
func diff(a, b JsonNode, p Path, opts *options, strategy patchStrategy) Diff {
    if !opts.diffingOn {
        return Diff{} // Return empty diff when diffing is off
    }
    // Existing logic...
}
```

### 5. Recursive Diffing
- Use `refine(opts, pathElement)` to get refined options for children
- Children inherit parent's diffing state unless overridden by PathOption
- PathOptions processed in order, later ones take precedence

### 6. Array Handling
Special consideration for arrays since they dispatch to different types:
- Check diffing state in `refineForArrayDispatch()`
- Ensure SET/MULTISET dispatch still works when combined with DIFF_OFF
- Handle nested PathOptions correctly in array context

## Usage Examples

### Basic Usage
```json
// Turn off diffing for specific field
[{"@":["metadata","timestamp"],"^":["DIFF_OFF"]}]

// Turn off diffing at root (ignore all changes)
[{"@":[],"^":["DIFF_OFF"]}]

// Turn on diffing explicitly (redundant but valid)
[{"@":["data"],"^":["DIFF_ON"]}]
```

### Allow-list Approach
```json
// Ignore everything except user data
[
  {"@":[],"^":["DIFF_OFF"]}, 
  {"@":["userdata"],"^":["DIFF_ON"]}
]
```

### Deny-list Approach  
```json
// Allow all changes except system fields
[
  {"@":["metadata","generated"],"^":["DIFF_OFF"]},
  {"@":["metadata","timestamp"],"^":["DIFF_OFF"]}
]
```

### Complex Nesting
```json
// Turn off diffing for config, but allow changes to user-modifiable parts
[
  {"@":["config"],"^":["DIFF_OFF"]},
  {"@":["config","user_settings"],"^":["DIFF_ON"]}
]
```

### Combination with Other Options
```json
// Turn off diffing for tags array, but when diffing is on elsewhere, treat arrays as sets
[
  {"@":["tags"],"^":["DIFF_OFF"]},
  {"@":["categories"],"^":["SET"]}
]
```

## Testing Strategy

### Unit Tests
1. **Basic functionality**: DIFF_ON/DIFF_OFF at various paths
2. **Nesting**: Override parent diffing state with child PathOptions  
3. **Order precedence**: Last PathOption directive wins
4. **Empty diffs**: Verify empty Diff returned when diffing off
5. **Default behavior**: Ensure current behavior preserved without options

### Integration Tests
1. **Complex documents**: Real-world JSON with mixed selective diffing
2. **Performance**: Ensure no performance regression
3. **Combined options**: DIFF_OFF with SET, MULTISET, precision options
4. **Error handling**: Invalid PathOption combinations

### Edge Cases
1. **Empty paths**: DIFF_OFF at root path
2. **Nonexistent paths**: PathOptions targeting missing fields
3. **Type mismatches**: PathOptions on incompatible types
4. **Deeply nested**: Very deep nesting with multiple overrides

## Backward Compatibility
- No breaking changes to existing API
- Default behavior unchanged (diffingOn = true)
- Existing PathOptions work exactly as before
- New options only take effect when explicitly used
- All existing tests should pass without modification

## Performance Considerations
- Minimal overhead when not using DIFF_ON/DIFF_OFF options
- Early return in event generation saves computation when diffing off
- PathOption processing overhead unchanged
- No impact on memory usage for typical use cases

## Future Extensions
This implementation provides a foundation for:
- More granular diffing control (e.g., ignore only certain types of changes)
- Conditional diffing based on value patterns
- Integration with external ignore/allow rules
- Performance optimizations for large documents with many ignored sections
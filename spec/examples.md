# Structural Format Examples

This document provides complete examples of the structural diff format, demonstrating all features and edge cases for implementors.

## Basic Examples

### Simple Object Changes

**Input A:**
```json
{"name": "Alice", "age": 25}
```

**Input B:**
```json
{"name": "Bob", "age": 25}
```

**Diff Output:**
```diff
@ ["name"]
- "Alice"
+ "Bob"
```

### Array Element Changes with Context

**Input A:**
```json
{"items": ["apple", "banana", "cherry"]}
```

**Input B:**
```json
{"items": ["apple", "blueberry", "cherry"]}
```

**Diff Output:**
```diff
@ ["items",1]
  "apple"
- "banana"
+ "blueberry"
  "cherry"
]
```

### Nested Object Changes

**Input A:**
```json
{
  "user": {
    "profile": {
      "name": "Alice",
      "preferences": {
        "theme": "light"
      }
    }
  }
}
```

**Input B:**
```json
{
  "user": {
    "profile": {
      "name": "Alice",
      "preferences": {
        "theme": "dark"
      }
    }
  }
}
```

**Diff Output:**
```diff
@ ["user","profile","preferences","theme"]
- "light"
+ "dark"
```

## Array Diffing Examples

### LCS Algorithm Demonstration

**Input A:**
```json
[1, 2, 3, 4, 5]
```

**Input B:**
```json
[1, 3, 4, 6, 7]
```

**Diff Output:**
```diff
@ [1]
  1
- 2
  3
  4
@ [4]
  3
  4
- 5
+ 6
+ 7
]
```

### Complex Array Context

**Input A:**
```json
["red", "green", "blue", "yellow", "purple"]
```

**Input B:**
```json
["red", "orange", "magenta", "yellow", "purple"]
```

**Diff Output:**
```diff
@ [1]
  "red"
- "green"
- "blue"
+ "orange"
+ "magenta"
  "yellow"
  "purple"
```

## Options Examples

### SET Option

**Options:** `["SET"]`

**Input A:**
```json
{"tags": ["urgent", "bug", "frontend"]}
```

**Input B:**
```json
{"tags": ["frontend", "urgent", "enhancement"]}
```

**Diff Output:**
```diff
^ "SET"
@ ["tags",{}]
- "bug"
+ "enhancement"
```

### MULTISET Option

**Options:** `["MULTISET"]`

**Input A:**
```json
{"counts": ["a", "a", "b", "c", "c", "c"]}
```

**Input B:**
```json
{"counts": ["a", "b", "b", "c", "c", "d"]}
```

**Diff Output:**
```diff
^ "MULTISET"
@ ["counts",[]]
- "a"
- "c"
+ "b"
+ "d"
```

### Precision Option

**Options:** `[{"precision": 0.01}]`

**Input A:**
```json
{"temperature": 20.123, "humidity": 65.456}
```

**Input B:**
```json
{"temperature": 20.125, "humidity": 65.489}
```

**Diff Output:**
```diff
^ {"precision":0.01}
@ ["humidity"]
- 65.456
+ 65.489
```

### Keys Option

**Options:** `[{"keys": ["id"]}]`

**Input A:**
```json
{
  "users": [
    {"id": 1, "name": "Alice", "status": "active"},
    {"id": 2, "name": "Bob", "status": "inactive"}
  ]
}
```

**Input B:**
```json
{
  "users": [
    {"id": 1, "name": "Alice", "status": "inactive"},
    {"id": 2, "name": "Bob", "status": "active"}
  ]
}
```

**Diff Output:**
```diff
^ {"keys":["id"]}
@ ["users",{"id":1},"status"]
- "active"
+ "inactive"
@ ["users",{"id":2},"status"]
- "inactive"
+ "active"
```

## PathOptions Examples

### Targeted SET Operation

**Options:** `[{"@":["tags"],"^":["SET"]}]`

**Input A:**
```json
{
  "tags": ["red", "blue", "green"],
  "items": [1, 2, 3]
}
```

**Input B:**
```json
{
  "tags": ["green", "red", "blue"],
  "items": [3, 2, 1]
}
```

**Diff Output:**
```diff
^ {"@":["tags"],"^":["SET"]}
@ ["items",0]
[
- 1
- 2
- 3
+ 3
+ 2
+ 1
]
```

### Multiple PathOptions

**Options:** `[{"@":["coords"],"^":[{"precision":0.1}]},{"@":["labels"],"^":["SET"]}]`

**Input A:**
```json
{
  "coords": {"x": 10.15, "y": 20.25},
  "labels": ["A", "B", "C"],
  "values": [100, 200]
}
```

**Input B:**
```json
{
  "coords": {"x": 10.18, "y": 20.35},
  "labels": ["C", "A", "B"],
  "values": [150, 200]
}
```

**Diff Output:**
```diff
^ {"@":["coords"],"^":[{"precision":0.1}]}
^ {"@":["labels"],"^":["SET"]}
@ ["coords","y"]
- 20.25
+ 20.35
@ ["values",0]
[
- 100
+ 150
  200
]
```

### DIFF_OFF Example

**Options:** `[{"@":["metadata"],"^":["DIFF_OFF"]}]`

**Input A:**
```json
{
  "data": {"value": 42},
  "metadata": {
    "timestamp": "2023-01-01T10:00:00Z",
    "version": "1.0"
  }
}
```

**Input B:**
```json
{
  "data": {"value": 43},
  "metadata": {
    "timestamp": "2023-01-02T11:00:00Z",
    "version": "1.1"
  }
}
```

**Diff Output:**
```diff
^ {"@":["metadata"],"^":["DIFF_OFF"]}
@ ["data","value"]
- 42
+ 43
```

### Allow-list with DIFF_OFF/DIFF_ON

**Options:** `[{"@":[],"^":["DIFF_OFF"]},{"@":["userdata"],"^":["DIFF_ON"]}]`

**Input A:**
```json
{
  "userdata": {"name": "Alice"},
  "system": {"cpu": "high", "memory": "normal"},
  "timestamp": "2023-01-01"
}
```

**Input B:**
```json
{
  "userdata": {"name": "Bob"},
  "system": {"cpu": "low", "memory": "high"},
  "timestamp": "2023-01-02"
}
```

**Diff Output:**
```diff
^ {"@":[],"^":["DIFF_OFF"]}
^ {"@":["userdata"],"^":["DIFF_ON"]}
@ ["userdata","name"]
- "Alice"
+ "Bob"
```

## Advanced Path Elements

### Set with Matching Keys

**Input A:**
```json
{
  "items": [
    {"id": "apple", "color": "red", "size": "medium"},
    {"id": "banana", "color": "yellow", "size": "large"}
  ]
}
```

**Input B:**
```json
{
  "items": [
    {"id": "apple", "color": "green", "size": "medium"},
    {"id": "banana", "color": "yellow", "size": "small"}
  ]
}
```

**Options:** `[{"keys": ["id"]}]`

**Diff Output:**
```diff
^ {"keys":["id"]}
@ ["items",{"id":"apple"},"color"]
- "red"
+ "green"
@ ["items",{"id":"banana"},"size"]
- "large"
+ "small"
```

### Multiset with Objects

**Input A:**
```json
{
  "events": [
    {"type": "click", "count": 5},
    {"type": "hover", "count": 2},
    {"type": "click", "count": 3}
  ]
}
```

**Input B:**
```json
{
  "events": [
    {"type": "click", "count": 3},
    {"type": "scroll", "count": 1},
    {"type": "hover", "count": 2}
  ]
}
```

**Options:** `[{"@":["events"],"^":["MULTISET"]}]`

**Diff Output:**
```diff
^ {"@":["events"],"^":["MULTISET"]}
@ ["events",[]]
- {"count":5,"type":"click"}
+ {"count":1,"type":"scroll"}
```

## Type Conversion Examples

### Object to Array

**Input A:**
```json
{"data": {"x": 1, "y": 2}}
```

**Input B:**
```json
{"data": [1, 2]}
```

**Diff Output:**
```diff
@ ["data"]
- {"x":1,"y":2}
+ [1,2]
```

### Null Handling

**Input A:**
```json
{"optional": "value"}
```

**Input B:**
```json
{"optional": null}
```

**Diff Output:**
```diff
@ ["optional"]
- "value"
+ null
```

### Empty/Void States

**Input A (empty document):**
```
(empty)
```

**Input B:**
```json
{"new": "document"}
```

**Diff Output:**
```diff
@ []
+ {"new":"document"}
```

## Unicode and Special Characters

### Unicode Support

**Input A:**
```json
{"message": "こんにちは世界"}
```

**Input B:**
```json
{"message": "さようなら世界"}
```

**Diff Output:**
```diff
@ ["message"]
- "こんにちは世界"
+ "さようなら世界"
```

### Escaped Characters

**Input A:**
```json
{"path": "C:\\Users\\Alice\\Documents"}
```

**Input B:**
```json
{"path": "C:\\Users\\Bob\\Documents"}
```

**Diff Output:**
```diff
@ ["path"]
- "C:\\Users\\Alice\\Documents"
+ "C:\\Users\\Bob\\Documents"
```

## Complex Real-World Examples

### Configuration File Update

**Input A:**
```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "connections": {"min": 5, "max": 20}
  },
  "features": {
    "feature_flags": ["beta_ui", "new_auth"],
    "experimental": ["ml_suggestions"]
  },
  "logging": {
    "level": "info",
    "outputs": ["console", "file"]
  }
}
```

**Input B:**
```json
{
  "database": {
    "host": "db.example.com",
    "port": 5432,
    "connections": {"min": 10, "max": 50}
  },
  "features": {
    "feature_flags": ["new_auth", "beta_ui", "improved_search"],
    "experimental": ["ml_suggestions", "ai_chat"]
  },
  "logging": {
    "level": "debug",
    "outputs": ["console", "file", "remote"]
  }
}
```

**Options:** `[{"@":["features","feature_flags"],"^":["SET"]}]`

**Diff Output:**
```diff
^ {"@":["features","feature_flags"],"^":["SET"]}
@ ["database","connections","max"]
- 20
+ 50
@ ["database","connections","min"]
- 5
+ 10
@ ["database","host"]
- "localhost"
+ "db.example.com"
@ ["features","experimental",1]
[
  "ml_suggestions"
+ "ai_chat"
]
@ ["features","feature_flags",{}]
+ "improved_search"
@ ["logging","level"]
- "info"
+ "debug"
@ ["logging","outputs",2]
[
  "console"
  "file"
+ "remote"
]
```

### API Response Comparison

**Input A:**
```json
{
  "status": "success",
  "data": {
    "users": [
      {"id": 1, "name": "Alice", "last_login": "2023-01-01T10:00:00Z"},
      {"id": 2, "name": "Bob", "last_login": "2023-01-01T11:00:00Z"}
    ],
    "total": 2,
    "page": 1
  },
  "metadata": {
    "request_id": "req-123",
    "timestamp": "2023-01-01T12:00:00Z"
  }
}
```

**Input B:**
```json
{
  "status": "success",
  "data": {
    "users": [
      {"id": 1, "name": "Alice", "last_login": "2023-01-02T09:00:00Z"},
      {"id": 2, "name": "Bob", "last_login": "2023-01-01T11:00:00Z"},
      {"id": 3, "name": "Charlie", "last_login": "2023-01-02T10:00:00Z"}
    ],
    "total": 3,
    "page": 1
  },
  "metadata": {
    "request_id": "req-456",
    "timestamp": "2023-01-02T12:00:00Z"
  }
}
```

**Options:** `[{"@":["metadata"],"^":["DIFF_OFF"]},{"@":["data","users"],"^":[{"keys":["id"]}]}]`

**Diff Output:**
```diff
^ {"@":["metadata"],"^":["DIFF_OFF"]}
^ {"@":["data","users"],"^":[{"keys":["id"]}]}
@ ["data","total"]
- 2
+ 3
@ ["data","users",{"id":1},"last_login"]
- "2023-01-01T10:00:00Z"
+ "2023-01-02T09:00:00Z"
@ ["data","users",2]
+ {"id":3,"last_login":"2023-01-02T10:00:00Z","name":"Charlie"}
```

These examples demonstrate the full range of the structural format's capabilities and should provide implementors with complete understanding of expected behavior across all scenarios.

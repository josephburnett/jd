package jd

import (
	"strings"
	"testing"
)

func TestPatchErrExpectColl(t *testing.T) {
	// string path element
	_, err := patchErrExpectColl(jsonString("val"), "key")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "expected JSON object") {
		t.Errorf("expected 'expected JSON object', got: %v", err)
	}
	// float64 path element
	_, err = patchErrExpectColl(jsonString("val"), float64(0))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "expected JSON array") {
		t.Errorf("expected 'expected JSON array', got: %v", err)
	}
	// Unknown type path element
	_, err = patchErrExpectColl(jsonString("val"), true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPatchErrNonSetDiff(t *testing.T) {
	// Multiple removals
	_, err := patchErrNonSetDiff(
		[]JsonNode{jsonString("a"), jsonString("b")},
		[]JsonNode{jsonString("c")},
		Path{PathKey("foo")},
	)
	if err == nil || !strings.Contains(err.Error(), "multiple removals") {
		t.Errorf("expected 'multiple removals' error, got: %v", err)
	}
	// Multiple additions
	_, err = patchErrNonSetDiff(
		[]JsonNode{jsonString("a")},
		[]JsonNode{jsonString("b"), jsonString("c")},
		Path{PathKey("foo")},
	)
	if err == nil || !strings.Contains(err.Error(), "multiple additions") {
		t.Errorf("expected 'multiple additions' error, got: %v", err)
	}
}

func TestPatchCommonMergeEdgeCases(t *testing.T) {
	// Merge patch with non-PathKey path element should error
	_, err := patch(
		jsonNumber(1),
		nil,
		Path{PathIndex(0), PathKey("a")},
		nil, nil, []JsonNode{jsonNumber(2)}, nil,
		mergePatchStrategy,
	)
	if err == nil {
		t.Fatal("expected error for non-PathKey in merge patch path")
	}
	// Merge patch where value is void and rest is empty — key should be omitted
	result, err := patch(
		jsonNumber(1),
		nil,
		Path{PathKey("a")},
		nil, nil, []JsonNode{voidNode{}}, nil,
		mergePatchStrategy,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	o, ok := result.(jsonObject)
	if !ok {
		t.Fatalf("expected jsonObject, got %T", result)
	}
	if _, exists := o["a"]; exists {
		t.Error("expected key 'a' to be omitted for void value with empty rest")
	}
}

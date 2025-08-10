package types

import (
	"testing"
)

func TestPathMetadataBasics(t *testing.T) {
	// Test basic path option creation
	pathOpt := PathOption{
		At:   Path{PathKey("users")},
		Then: []Option{SET},
	}
	
	if len(pathOpt.At) != 1 {
		t.Errorf("Expected path length 1, got %d", len(pathOpt.At))
	}
	
	if len(pathOpt.Then) != 1 {
		t.Errorf("Expected 1 option in Then, got %d", len(pathOpt.Then))
	}
}

func TestOptionsRefine(t *testing.T) {
	// Test that options refine correctly at specific paths
	opts := &Options{
		Retain: []Option{
			PathOption{
				At:   Path{PathKey("users")},
				Then: []Option{SET},
			},
		},
	}
	
	// At root level, should have the path option in retain
	rootOpts := Refine(opts, nil)
	if len(rootOpts.Retain) != 1 {
		t.Errorf("Expected 1 retained option at root, got %d", len(rootOpts.Retain))
	}
	if len(rootOpts.Apply) != 0 {
		t.Errorf("Expected 0 applied options at root, got %d", len(rootOpts.Apply))
	}
	
	// At "users" path, should apply SET option
	usersOpts := Refine(opts, PathKey("users"))
	if len(usersOpts.Apply) != 1 {
		t.Errorf("Expected 1 applied option at users path, got %d", len(usersOpts.Apply))
	}
	
	// Check that the applied option is SET
	_, isSet := usersOpts.Apply[0].(SetOption)
	if !isSet {
		t.Errorf("Expected SET option to be applied at users path")
	}
}

func TestNestedPathOptions(t *testing.T) {
	// Test nested path options: users[0].tags should use SET semantics
	opts := &Options{
		Retain: []Option{
			PathOption{
				At:   Path{PathKey("users"), PathIndex(0), PathKey("tags")},
				Then: []Option{SET},
			},
		},
	}
	
	// Navigate down the path step by step
	step1 := Refine(opts, PathKey("users"))
	if len(step1.Apply) != 0 {
		t.Errorf("Expected no applied options at users level")
	}
	
	step2 := Refine(step1, PathIndex(0))
	if len(step2.Apply) != 0 {
		t.Errorf("Expected no applied options at users[0] level")
	}
	
	step3 := Refine(step2, PathKey("tags"))
	if len(step3.Apply) != 1 {
		t.Errorf("Expected 1 applied option at users[0].tags level, got %d", len(step3.Apply))
	}
	
	_, isSet := step3.Apply[0].(SetOption)
	if !isSet {
		t.Errorf("Expected SET option at users[0].tags level")
	}
}

// Note: LCS tests would go here but are in separate package to avoid import cycle

func TestBasicNodeEquality(t *testing.T) {
	cases := []struct {
		name     string
		a, b     JsonNode
		opts     *Options
		wantEq   bool
	}{{
		name:   "same strings equal",
		a:      jsonString("test"),
		b:      jsonString("test"), 
		opts:   &Options{},
		wantEq: true,
	}, {
		name:   "different strings not equal",
		a:      jsonString("test"),
		b:      jsonString("other"),
		opts:   &Options{},
		wantEq: false,
	}, {
		name:   "numbers without precision",
		a:      jsonNumber(1.0),
		b:      jsonNumber(1.00001),
		opts:   &Options{},
		wantEq: false,
	}, {
		name:   "numbers with precision option",
		a:      jsonNumber(1.0),
		b:      jsonNumber(1.00001),
		opts:   &Options{Apply: []Option{PrecisionOption{Precision: 0.001}}},
		wantEq: true,
	}, {
		name:   "numbers precision too small",
		a:      jsonNumber(1.0),
		b:      jsonNumber(1.1),
		opts:   &Options{Apply: []Option{PrecisionOption{Precision: 0.001}}},
		wantEq: false,
	}}
	
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.a.equals(c.b, c.opts)
			if got != c.wantEq {
				t.Errorf("equals() = %v, want %v", got, c.wantEq)
			}
		})
	}
}
package jd

import (
	"testing"
)

func TestNewPathEdgeCases(t *testing.T) {
	// nil input
	p, err := NewPath(nil)
	if err != nil || p != nil {
		t.Error("nil should return nil, nil")
	}
	// Non-array input
	_, err = NewPath(jsonString("not array"))
	if err == nil {
		t.Fatal("expected error for non-array")
	}
	// Empty object -> PathSet
	p, err = NewPath(jsonArray{jsonObject{}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := p[0].(PathSet); !ok {
		t.Errorf("expected PathSet, got %T", p[0])
	}
	// Empty nested array -> PathMultiset
	p, err = NewPath(jsonArray{jsonArray{}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := p[0].(PathMultiset); !ok {
		t.Errorf("expected PathMultiset, got %T", p[0])
	}
	// Nested array with object -> PathMultisetKeys
	p, err = NewPath(jsonArray{jsonArray{jsonObject{"id": jsonString("foo")}}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := p[0].(PathMultisetKeys); !ok {
		t.Errorf("expected PathMultisetKeys, got %T", p[0])
	}
	// Nested array with non-object -> error
	_, err = NewPath(jsonArray{jsonArray{jsonString("not obj")}})
	if err == nil {
		t.Fatal("expected error for non-object in multiset")
	}
	// Nested array with length > 1 -> error
	_, err = NewPath(jsonArray{jsonArray{jsonObject{}, jsonObject{}}})
	if err == nil {
		t.Fatal("expected error for multiset array length > 1")
	}
	// Unsupported element type
	_, err = NewPath(jsonArray{jsonBool(true)})
	if err == nil {
		t.Fatal("expected error for bool in path")
	}
}

func TestPathJsonNodeMultiset(t *testing.T) {
	// PathMultiset
	p := Path{PathMultiset{}}
	n := p.JsonNode()
	a, ok := n.(jsonArray)
	if !ok || len(a) != 1 {
		t.Fatal("expected array with 1 element")
	}
	inner, ok := a[0].(jsonArray)
	if !ok || len(inner) != 0 {
		t.Error("expected empty inner array for PathMultiset")
	}

	// PathMultisetKeys
	p = Path{PathMultisetKeys{"id": jsonString("foo")}}
	n = p.JsonNode()
	a, ok = n.(jsonArray)
	if !ok || len(a) != 1 {
		t.Fatal("expected array with 1 element")
	}
	inner, ok = a[0].(jsonArray)
	if !ok || len(inner) != 1 {
		t.Error("expected inner array with 1 object element")
	}
}

func TestPathNextMultiset(t *testing.T) {
	// PathMultiset
	p := Path{PathMultiset{}, PathKey("a")}
	elem, opts, rest := p.next()
	if _, ok := elem.(PathMultiset); !ok {
		t.Error("expected PathMultiset element")
	}
	if len(opts) != 1 {
		t.Error("expected 1 option")
	}
	if len(rest) != 1 {
		t.Error("expected 1 remaining path element")
	}

	// PathSetKeys
	p = Path{PathSetKeys{"id": jsonString("foo")}, PathKey("a")}
	elem, opts, rest = p.next()
	if _, ok := elem.(PathSetKeys); !ok {
		t.Error("expected PathSetKeys element")
	}
	if len(opts) != 1 {
		t.Error("expected 1 option")
	}

	// PathMultisetKeys
	p = Path{PathMultisetKeys{"id": jsonString("foo")}, PathKey("a")}
	elem, opts, rest = p.next()
	if _, ok := elem.(PathMultisetKeys); !ok {
		t.Error("expected PathMultisetKeys element")
	}
	if len(opts) != 1 {
		t.Error("expected 1 option")
	}
	_ = rest
}

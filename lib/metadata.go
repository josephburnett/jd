package jd

import (
	"sort"
	"strings"
)

// Metadata is a closed set of types which modify diff and patch semantics.
type Metadata interface {
	is_metadata()
	string() string
}

type setMetadata struct{}
type multisetMetadata struct{}
type setkeysMetadata struct {
	keys map[string]bool
}
type mergeMetadata struct{}

func (setMetadata) is_metadata()      {}
func (multisetMetadata) is_metadata() {}
func (setkeysMetadata) is_metadata()  {}
func (mergeMetadata) is_metadata()    {}

func (m setMetadata) string() string {
	return "set"
}

func (m multisetMetadata) string() string {
	return "multiset"
}

func (m setkeysMetadata) string() string {
	ks := make([]string, 0)
	for k := range m.keys {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	// TODO: escape commas.
	return "setkeys=" + strings.Join(ks, ",")
}

func (m mergeMetadata) string() string {
	return "merge"
}

var (
	MULTISET Metadata = multisetMetadata{}
	SET      Metadata = setMetadata{}
	MERGE    Metadata = mergeMetadata{}
)

func Setkeys(keys ...string) Metadata {
	m := setkeysMetadata{
		keys: make(map[string]bool),
	}
	for _, key := range keys {
		m.keys[key] = true
	}
	return m
}

type patchStrategy string

const (
	mergePatchStrategy  patchStrategy = "merge"
	strictPatchStrategy               = "strict"
)

func getPatchStrategy(metadata []Metadata) patchStrategy {
	if checkMetadata(MERGE, metadata) {
		return mergePatchStrategy
	}
	return strictPatchStrategy
}

func dispatch(n JsonNode, metadata []Metadata) JsonNode {
	switch n := n.(type) {
	case jsonArray:
		if checkMetadata(SET, metadata) {
			return jsonSet(n)
		}
		if checkMetadata(MULTISET, metadata) {
			return jsonMultiset(n)
		}
		return jsonList(n)
	}
	return n
}

func checkMetadata(want Metadata, metadata []Metadata) bool {
	for _, o := range metadata {
		if o == want {
			return true
		}
	}
	return false
}

func getSetkeysMetadata(metadata []Metadata) *setkeysMetadata {
	for _, o := range metadata {
		if s, ok := o.(setkeysMetadata); ok {
			return &s
		}
	}
	return nil
}

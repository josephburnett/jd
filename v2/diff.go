package jd

// DiffElement (hunk) is a way in which two JsonNodes differ at a given
// Path. OldValues can be removed and NewValues can be added. The exact
// Path and how to interpret the intervening structure is determined by a
// list of JsonNodes (path elements).
type DiffElement struct {

	// Metadata describes how this DiffElement should be
	// interpretted. It is also inherited by following
	// DiffElements until another Metadata is encountered.
	Metadata Metadata

	// Path elements can be strings to index Objects, numbers to
	// index Lists and objects to index Sets and Multisets.
	//
	// For example:
	//   ["foo","bar"]               // indexes to 1 in {"foo":{"bar":1}}
	//   ["foo",0]                   // indexes to 1 in {"foo":[1]}
	//   ["foo",{}]                  // indexes a set under "foo" in {"foo":[1]}
	//   ["foo",{"id":"bar"},"baz"]  // indexes to 1 in {"foo":[{"id":"bar","baz":1}]}
	Path Path

	// Before are the required context which should appear before
	// new and old values of a diff element. They are only used
	// for diffs in a list element.
	Before []JsonNode

	// Remove are removed from the JsonNode at the Path. Usually
	// only one old value is provided unless removing entries from
	// a Set or Multiset. When using merge semantics no old values
	// are provided (new values stomp old ones).
	Remove []JsonNode

	// Add are added to the JsonNode at the Path. Usually only one
	// new value is provided unless adding entries to a Set or
	// Multiset.
	Add []JsonNode

	// After are the required context which should appear after
	// new and old values of a diff element. They are only used
	// for diffs in a list element.
	After []JsonNode
}

// Diff describes how two JsonNodes differ from each other. A Diff is
// composed of DiffElements (hunks) which describe a difference at a
// given Path. Each hunk stands alone with all necessary Metadata
// embedded in the Path, so a Diff rendered in native jd format can
// easily be edited by hand. The elements of a Diff can be applied to
// a JsonNode by the Patch method.
type Diff []DiffElement

// JSON Patch (RFC 6902)
type patchElement struct {
	Op    string      `json:"op"`   // "add", "test" or "remove"
	Path  string      `json:"path"` // JSON Pointer (RFC 6901)
	Value interface{} `json:"value"`
}

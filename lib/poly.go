package jd

type Metadata interface {
	is_metadata()
}

type multisetMetadata struct{}
type setMetadata struct{}
type setkeysMetadata struct {
	keys map[string]bool
}

func (multisetMetadata) is_metadata() {}
func (setMetadata) is_metadata()      {}
func (setkeysMetadata) is_metadata()  {}

var (
	MULTISET Metadata = multisetMetadata{}
	SET      Metadata = setMetadata{}
)

func SetkeysMetadata(keys ...string) Metadata {
	m := setkeysMetadata{
		keys: make(map[string]bool),
	}
	for _, key := range keys {
		m.keys[key] = true
	}
	return m
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

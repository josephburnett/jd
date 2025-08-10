package jd

import "fmt"

type PathElement interface {
	isPathElement()
}

type PathIndex int
type PathKey string
type PathAllKeys struct{}
type PathSet struct{}
type PathMultiset struct{}
type PathSetKeys map[string]JsonNode
type PathMultisetKeys map[string]JsonNode
type PathAllValues struct{}

func (_ PathIndex) isPathElement()        {}
func (_ PathKey) isPathElement()          {}
func (_ PathAllKeys) isPathElement()      {}
func (_ PathSet) isPathElement()          {}
func (_ PathMultiset) isPathElement()     {}
func (_ PathSetKeys) isPathElement()      {}
func (_ PathMultisetKeys) isPathElement() {}
func (_ PathAllValues) isPathElement()    {}

func newPathSetKeys(o jsonObject, opts *options) PathSetKeys {
	setKeys, ok := getOption[setKeysOption](opts)
	if !ok || setKeys == nil {
		return PathSetKeys(o)
	}
	key := newJsonObject()
	for _, k := range *setKeys {
		v, ok := o[k]
		if ok {
			key[k] = v
		} else {
			key[k] = jsonNull{}
		}
	}
	return PathSetKeys(key)
}

type Path []PathElement

func NewPath(n JsonNode) (Path, error) {
	if n == nil {
		return nil, nil
	}
	a, ok := n.(jsonArray)
	if !ok {
		return nil, fmt.Errorf("path must be an array. got %T", n)
	}
	p := make(Path, len(a))
	for i, e := range a {
		switch e := e.(type) {
		case jsonString:
			p[i] = PathKey(e)
		case jsonNumber:
			p[i] = PathIndex(e)
		case jsonObject:
			if len(e) == 0 {
				p[i] = PathSet{}
			} else {
				p[i] = PathSetKeys(e)
			}
		case jsonArray:
			switch len(e) {
			case 0:
				p[i] = PathMultiset{}
			case 1:
				o, ok := e[0].(jsonObject)
				if !ok {
					return nil, fmt.Errorf("multiset keys must be an object. got %T", e[0])
				}
				p[i] = PathMultisetKeys(o)
			default:
				return nil, fmt.Errorf("multiset path element must have length 0 or 1. got %v", len(e))
			}
		default:
			return nil, fmt.Errorf("path element must be a number, object or array. got %T", e)
		}
	}
	return p, nil
}

func (p Path) JsonNode() JsonNode {
	a := make(jsonArray, len(p))
	for i, e := range p {
		switch e := e.(type) {
		case PathKey:
			a[i] = jsonString(e)
		case PathIndex:
			a[i] = jsonNumber(e)
		case PathSet:
			a[i] = jsonObject{}
		case PathMultiset:
			a[i] = jsonArray{}
		case PathSetKeys:
			a[i] = jsonObject(e)
		case PathMultisetKeys:
			a[i] = jsonArray{jsonObject(e)}
		default:
			panic(fmt.Sprintf("path element should be a closed set. got %T", e))
		}
	}
	return a
}

func (p Path) next() (PathElement, []Option, Path) {
	if len(p) == 0 {
		return nil, nil, nil
	}
	rest := p[1:]
	switch e := p[0].(type) {
	case PathKey:
		return p[0], []Option{}, rest
	case PathIndex:
		return p[0], nil, rest
	case PathSet:
		return p[0], []Option{setOption{}}, rest
	case PathMultiset:
		return p[0], []Option{multisetOption{}}, rest
	case PathSetKeys:
		return p[0], []Option{setOption{}}, rest
	case PathMultisetKeys:
		return p[0], []Option{multisetOption{}}, rest
	default:
		panic(fmt.Sprintf("path element should be a closed set. got %T", e))
	}
}

func (p Path) isLeaf() bool {
	if len(p) == 0 {
		return true
	}
	if len(p) > 1 {
		return false
	}
	switch p[0].(type) {
	case PathSet, PathSetKeys, PathMultiset, PathMultisetKeys:
		return true
	default:
		return false
	}
}

func (p Path) clone() Path {
	p2 := make(Path, len(p))
	for i, e := range p {
		p2[i] = e
	}
	return p2
}

func (p Path) drop() Path {
	if len(p) > 0 {
		return p[:len(p)-1]
	}
	return p
}

package jd

import "fmt"

type PathElement interface {
	isPathElement()
}

type pathIndex int
type pathKey string
type pathSet struct{}
type pathMultiset struct{}
type pathSetKeys map[string]JsonNode
type pathMultisetKeys map[string]JsonNode

func (_ pathIndex) isPathElement() {}
func (_ pathKey) isPathElement() {}
func (_ pathSet) isPathElement() {}
func (_ pathMultiset) isPathElement()  {}
func (_ pathSetKeys) isPathElement() {}
func (_ pathMultisetKeys) isPathElement()  {}

type Path []PathElement

func NewPath(n JsonNode) (Path, error) {
	if n == nil {
		return 
	}
	a, ok := n.(jsonArray)
	if !ok {
		return nil, fmt.Errorf("path must be an array. got %T", n)
	}
	p := make(Path, len(a))
	for i, e := range a {
		switch e := e.(type) {
		case jsonNumber:
			p[i] = pathIndex(jsonNumber)
		case jsonObject:
			if len(e) == 0 {
				p[i] = pathSet{}
			} else {
				p[i] = pathSetKeys(e)
			}
		case jsonArray:
			switch len(e) {
			case 0:
				p[i] = pathMultiset{}
			case 1:
				o, ok := e[0].(jsonObject)
				if !ok {
					return nil, fmt.Errorf("multiset keys must be an object. got %T", e[0])
				}
				p[i] = pathMultisetKeys(e[0])
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
		case pathIndex:
			a[i] = jsonNumber(e)
		case pathSet:
			a[i] = jsonObject{}
		case pathMultiset:
			a[i] = jsonMultiset{}
		case pathSetKeys:
			a[i] = jsonObject(e)
		case pathMultisetKeys:
			a[i] = jsonArray{jsonOject(e)}
		default:
			panic(fmt.Sprintf("path element should be a closed set. got %T", e))
		}
	}
	return a
}

func (p Path) next() (JsonNode, []Option, Path) {
	if len(p) == 0 {
		return jsonVoid{}, nil, nil
	}
	rest := p[1:]
	switch e := p[0].(type) {
	case jsonObject:
		return p[0], []Option{setOption{}}, rest
	case jsonArray:
		if len(e) == 0 {
			return p[0], []Option{multisetOption{}}, rest
		}
		if len(e) == 1 && _, ok := e[0].(jsonObject); ok {
			return p[0], []Option{multisetOption{}}, rest
		}
	case jsonNumber:
		return p[0], nil, rest
	default:
		panic(fmt.Sprintf("path element should be a closed set. got %T", e))
	}
}

func (p Path) clone() Path {
	p2 := make(Path, len(p))
	for i, e := range p {
		p2[i] = e
	}
	return p2
}

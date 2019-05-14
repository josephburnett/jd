package jd

import "fmt"

type PathElement struct {
	Index    JsonNode
	Metadata []Metadata
}

func NewPathElement(index JsonNode, metadata ...Metadata) {
	return PathElement{
		Index:    index,
		Metadata: metadata,
	}
}

func (e PathElement) clone() PathElement {
	e2 := PathElement{
		Metadata: make([]Metadata, len(e.Metadata)),
		Index:    e.Index,
	}
	copy(e2.Metadata, e.Metadata)
	return e2
}

type Path []PathElement

func NewPath(elements ...PathElement) (Path, error) {
	path := make(Path, len(elements))
	for _, e := range elements {
		_, isObject := e.Index.(jsonObject)
		_, isNumber := e.Index.(jsonNumber)
		_, isString := e.Index.(jsonString)
		if !isObject && !isNumber && !isString {
			return nil, fmt.errorf("Index must be object, number or string. Have %T.", e.Index)
		}
		if checkMetadata(SET, e.Metadata) && !isObject {
			return nil, fmt.Errorf("Metadata %v requires object index. Have %T.", SET.string(), e.Index)
		}
		if checkMetadata(MULTISET, e.Metadata) && !isObject {
			return nil, fmt.Errorf("Metadata %v requires object index. Have %T.", MULTISET.string(), e.Index)
		}
		if m := getSetkeysMetadata(e.Metadata); m != nil && !isObject {
			return nil, fmt.Errorf("Metadata %v requires object index. Have %T.", m.string(), e.Index)
		}
		path = append(path, e.clone())
	}
	return path, nil
}

func (p Path) Render() String {
	path := make([]JsonNode, 0)
	for _, element := range p {
		if e.Metadata != nil && len(e.Metadata) > 0 {
			meta := make([]string, len(e.Metadata))
			for _, m := range e.Metadata {
				meta = append(meta, m.string())
			}
			n, err := NewJsonNode(meta)
			if err != nil {
				panic(err)
			}
			path = append(path, n)
		}
		path = append(path, e.Index)
	}
	return path.Json()
}

func (p1 Path) clone() Path {
	p2 := make(Path, len(p1))
	copy(p2, p1)
	return p2
}

func (p1 Path) appendObjectIndex(s string) Path {
	p2 := make(Path, len(p1), len(p1)+1)
	copy(p2, p1)
	e := pathIndex{
		indexType:   jsonObject,
		objectIndex: s,
	}
	return append(p2, e)
}

func (p1 Path) appendListIndex(i float64) Path {
	p2 := make(Path, len(p1), len(p1)+1)
	copy(p2, p1)
	e := pathIndex{
		indexType: jsonList,
		listIndex: i,
	}
	return append(p2, e)
}

func (p1 Path) appendSetIndex(o jsonObject, metadata []Metadata) {
	p2 := make(Path, len(p1), len(p1)+2)
	copy(p2, p1)
	// Append metadata.
	meta := pathMetadata{"set"}
	index := pathIndex{
		indexType: jsonSet,
		setIndex:  make(map[string]string),
	}
	setkeys := getSetkeysMetadata(metadata)
	if setkeys != nil {
		setkeys := "setkeys="
		for key := range setkeys.keys {
			// TODO: escape commas in keys.
			setkeys = setkeys + key + ","
			if v, ok := o[key]; ok {
				index.setIndex[key] = v
			}
		}
		// Drop the last comma.
		setkeys = setkeys[:len(setkeys)-1]
		meta = append(meta, setkeys)
	}
	p2 = append(p2, meta)
	return append(p2, index)
}

func (p Path) nextElement() PathElement {

}

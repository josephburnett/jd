package jd

import "reflect"

type PathElement interface {
	isPathElement()
}

type pathMetadata []string
type pathIndex struct {
	indexType     reflect.Type
	objectIndex   string
	listIndex     float64
	setIndex      map[string]string
	multisetIndex map[string]string
}

func (p pathMetadata) isPathElement() {}
func (p pathIndex) isPathElement()    {}

type Path []PathElement

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

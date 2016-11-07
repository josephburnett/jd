package jd

type PathElement interface{}
type Path []PathElement

func (p1 Path) clone() Path {
	p2 := make(Path, len(p1), len(p1)+1)
	copy(p2, p1)
	return p2
}

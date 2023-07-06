package jd

type PathElement interface {
	isPathElement()
}

func (o jsonObject) isPathElement() {}
func (a jsonArray) isPathElement()  {}
func (n jsonNumber) isPathElement() {}

type Path []PathElement

func NewPath(n JsonNode) (Path, error) {
	if n == nil {
		return 
	}
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
	}
	panic("path element should be closed set")
}

func (p Path) clone() Path {
	p2 := make(Path, len(p))
	for i, e := range p {
		p2[i] = e
	}
	return p2
}

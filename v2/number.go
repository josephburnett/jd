package jd

import (
	"bytes"
	"encoding/binary"
	"math"
)

type jsonNumber float64

var _ JsonNode = jsonNumber(0)

func (n jsonNumber) Json(_ ...Option) string {
	return renderJson(n.raw())
}

func (n jsonNumber) Yaml(_ ...Option) string {
	return renderYaml(n.raw())
}

func (n jsonNumber) raw() interface{} {
	return float64(n)
}

func (n1 jsonNumber) Equals(node JsonNode, opts ...Option) bool {
	o := refine(&options{retain: opts}, nil)
	return n1.equals(node, o)
}

func (n1 jsonNumber) equals(node JsonNode, o *options) bool {
	precision := 0.0
	if p, ok := getOption[precisionOption](o); ok {
		precision = p.precision
	}

	n2, ok := node.(jsonNumber)
	if !ok {
		return false
	}

	return math.Abs(float64(n1)-float64(n2)) <= precision
}

func (n jsonNumber) hashCode(opts *options) [8]byte {
	a := make([]byte, 0, 8)
	b := bytes.NewBuffer(a)
	binary.Write(b, binary.LittleEndian, n)
	return hash(b.Bytes())
}

func (n jsonNumber) Diff(node JsonNode, opts ...Option) Diff {
	o := refine(&options{retain: opts}, nil)
	return n.diff(node, make(Path, 0), o, getPatchStrategy(o))
}

func (n jsonNumber) diff(
	node JsonNode,
	path Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	return diff(n, node, path, opts, strategy)
}

func (n jsonNumber) Patch(d Diff) (JsonNode, error) {
	return patchAll(n, d)
}

func (n jsonNumber) patch(
	pathBehind, pathAhead Path,
	before, oldValues, newValues, after []JsonNode,
	strategy patchStrategy,
) (JsonNode, error) {
	return patch(n, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

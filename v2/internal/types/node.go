package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// JsonNode is a JSON value, collection of values, or a void representing
// the absense of a value. JSON values can be a Number, String, Boolean
// or Null. Collections can be an Object, native JSON array, ordered
// List, unordered Set or Multiset. JsonNodes are created with the
// NewJsonNode function or ReadJson* and ReadYaml* functions.
type JsonNode interface {

	// Json renders a JsonNode as a JSON string.
	Json(renderOptions ...Option) string

	// Yaml renders a JsonNode as a YAML string in block format.
	Yaml(renderOptions ...Option) string

	// Equals returns true if the JsonNodes are equal according to
	// the provided Metadata. The default behavior (no Metadata) is
	// to compare the entire structure down to scalar values treating
	// Arrays as orders Lists. The SET and MULTISET Metadata will
	// treat Arrays as sets or multisets (bags) respectively. To deep
	// compare objects in an array irrespective of order, the SetKeys
	// function will construct Metadata to compare objects by a set
	// of keys. If two JsonNodes are equal, then Diff with the same
	// Metadata will produce an empty Diff. And vice versa.
	Equals(n JsonNode, options ...Option) bool

	// Diff produces a list of differences (Diff) between two
	// JsonNodes such that if the output Diff were applied to the
	// first JsonNode (Patch) then the two JsonNodes would be
	// Equal. The necessary Metadata is embeded in the Diff itself so
	// only the Diff is required to Patch a JsonNode.
	Diff(n JsonNode, options ...Option) Diff

	// Patch applies a Diff to a JsonNode. No Metadata is provided
	// because the original interpretation of the structure is
	// embedded in the Diff itself.
	Patch(d Diff) (JsonNode, error)

	jsonNodeInternals
}

type jsonNodeInternals interface {
	raw() interface{}
	hashCode(opts *Options) [8]byte
	equals(n JsonNode, o *Options) bool
	diff(n JsonNode, p Path, opts *Options, strategy patchStrategy) Diff
	patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error)
}

// NewJsonNode constructs a JsonNode from native Golang objects. See the
// function source for supported types and conversions. Slices are always
// placed into native JSON Arrays and interpretated as Lists, Sets or
// Multisets based on Metadata provided during Equals and Diff
// operations.
func NewJsonNode(n interface{}) (JsonNode, error) {
	switch t := n.(type) {
	case map[string]interface{}:
		m := newJsonObject()
		for k, v := range t {
			n, ok := v.(JsonNode)
			if !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				n = e
			}
			m[k] = n
		}
		return m, nil
	case map[interface{}]interface{}:
		m := newJsonObject()
		for k, v := range t {
			s, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("unsupported key type %T", k)
			}
			if _, ok := v.(JsonNode); !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				m[s] = e
			}
		}
		return m, nil
	case []interface{}:
		l := make(jsonArray, len(t))
		for i, v := range t {
			if _, ok := v.(JsonNode); !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				l[i] = e
			}
		}
		return l, nil
	case float64:
		return jsonNumber(t), nil
	case int:
		return jsonNumber(t), nil
	case uint8:
		return jsonNumber(float64(t)), nil
	case string:
		return jsonString(t), nil
	case bool:
		return jsonBool(t), nil
	case nil:
		return jsonNull{}, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", t)
	}
}

func nodeList(n ...JsonNode) []JsonNode {
	l := []JsonNode{}
	if len(n) == 0 {
		return l
	}
	if n[0].Equals(voidNode{}) {
		return l
	}
	return append(l, n...)
}

// Basic JSON types
type jsonBool bool
type jsonString string
type jsonNumber float64
type jsonNull struct{}
type jsonObject map[string]JsonNode
type jsonArray []JsonNode
type voidNode struct{}

// Exported aliases for main package compatibility
type JsonBool = jsonBool
type JsonString = jsonString
type JsonNumber = jsonNumber
type JsonNull = jsonNull
type JsonObject = jsonObject
type JsonArray = jsonArray
type VoidNodeType = voidNode

// Constructor functions
func newJsonObject() jsonObject {
	return make(jsonObject)
}

func NewJsonObject() jsonObject {
	return make(jsonObject)
}

func VoidNode() JsonNode {
	return voidNode{}
}

func NodeList(n ...JsonNode) []JsonNode {
	l := []JsonNode{}
	if len(n) == 0 {
		return l
	}
	if n[0].Equals(VoidNode()) {
		return l
	}
	return append(l, n...)
}

func ReadJsonString(s string) (JsonNode, error) {
	var n interface{}
	err := json.Unmarshal([]byte(s), &n)
	if err != nil {
		return nil, err
	}
	return NewJsonNode(n)
}

func ReadYamlString(s string) (JsonNode, error) {
	var n interface{}
	err := yaml.Unmarshal([]byte(s), &n)
	if err != nil {
		return nil, err
	}
	return NewJsonNode(n)
}

// Render functions  
func RenderJson(n interface{}) string {
	s, _ := json.Marshal(n)
	return string(s)
}

func RenderYaml(n interface{}) string {
	s, _ := yaml.Marshal(n)
	return strings.TrimSuffix(string(s), "\n")
}

// Hash function
func hash(b []byte) [8]byte {
	var result [8]byte
	for i := 0; i < 8; i++ {
		if i < len(b) {
			result[i] = b[i]
		}
	}
	return result
}

// Basic type implementations
func (b jsonBool) Json(_ ...Option) string {
	return RenderJson(b.raw())
}

func (b jsonBool) Yaml(_ ...Option) string {
	return RenderYaml(b.raw())
}

func (b jsonBool) raw() interface{} {
	return bool(b)
}

func (b1 jsonBool) Equals(n JsonNode, opts ...Option) bool {
	o := Refine(&Options{Retain: opts}, nil)
	return b1.equals(n, o)
}

func (b1 jsonBool) equals(n JsonNode, _ *Options) bool {
	b2, ok := n.(jsonBool)
	if !ok {
		return false
	}
	return b1 == b2
}

func (b jsonBool) hashCode(_ *Options) [8]byte {
	if b {
		return [8]byte{0x24, 0x6B, 0xE3, 0xE4, 0xAF, 0x59, 0xDC, 0x1C}
	} else {
		return [8]byte{0xC6, 0x38, 0x77, 0xD1, 0x0A, 0x7E, 0x1F, 0xBF}
	}
}

func (b jsonBool) Diff(n JsonNode, opts ...Option) Diff {
	o := Refine(&Options{Retain: opts}, nil)
	strategy := getPatchStrategy(o)
	return b.diff(n, make(Path, 0), o, strategy)
}

func (b jsonBool) diff(n JsonNode, path Path, opts *Options, strategy patchStrategy) Diff {
	return basicDiff(b, n, path, opts, strategy)
}

func (b jsonBool) Patch(d Diff) (JsonNode, error) {
	return patchAll(b, d)
}

func (b jsonBool) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {
	return patch(b, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

// String implementation
func (s jsonString) Json(_ ...Option) string {
	return RenderJson(s.raw())
}

func (s jsonString) Yaml(_ ...Option) string {
	return RenderYaml(s.raw())
}

func (s jsonString) raw() interface{} {
	return string(s)
}

func (s1 jsonString) Equals(n JsonNode, opts ...Option) bool {
	o := Refine(&Options{Retain: opts}, nil)
	return s1.equals(n, o)
}

func (s1 jsonString) equals(n JsonNode, _ *Options) bool {
	s2, ok := n.(jsonString)
	if !ok {
		return false
	}
	return s1 == s2
}

func (s jsonString) hashCode(_ *Options) [8]byte {
	return hash([]byte(s))
}

func (s jsonString) Diff(n JsonNode, opts ...Option) Diff {
	o := Refine(&Options{Retain: opts}, nil)
	strategy := getPatchStrategy(o)
	return s.diff(n, make(Path, 0), o, strategy)
}

func (s jsonString) diff(n JsonNode, path Path, opts *Options, strategy patchStrategy) Diff {
	return basicDiff(s, n, path, opts, strategy)
}

func (s jsonString) Patch(d Diff) (JsonNode, error) {
	return patchAll(s, d)
}

func (s jsonString) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {
	return patch(s, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

// Number implementation
func (n jsonNumber) Json(_ ...Option) string {
	return RenderJson(n.raw())
}

func (n jsonNumber) Yaml(_ ...Option) string {
	return RenderYaml(n.raw())
}

func (n jsonNumber) raw() interface{} {
	return float64(n)
}

func (n1 jsonNumber) Equals(n JsonNode, opts ...Option) bool {
	o := Refine(&Options{Retain: opts}, nil)
	return n1.equals(n, o)
}

func (n1 jsonNumber) equals(n JsonNode, o *Options) bool {
	n2, ok := n.(jsonNumber)
	if !ok {
		return false
	}
	// Check for precision option
	if precision, ok := getOption[PrecisionOption](o); ok {
		diff := float64(n1) - float64(n2)
		if diff < 0 {
			diff = -diff
		}
		return diff <= precision.Precision
	}
	return n1 == n2
}

func (n jsonNumber) hashCode(opts *Options) [8]byte {
	return hash([]byte(strconv.FormatFloat(float64(n), 'g', -1, 64)))
}

func (n jsonNumber) Diff(node JsonNode, opts ...Option) Diff {
	o := Refine(&Options{Retain: opts}, nil)
	strategy := getPatchStrategy(o)
	return n.diff(node, make(Path, 0), o, strategy)
}

func (n jsonNumber) diff(node JsonNode, path Path, opts *Options, strategy patchStrategy) Diff {
	return basicDiff(n, node, path, opts, strategy)
}

func (n jsonNumber) Patch(d Diff) (JsonNode, error) {
	return patchAll(n, d)
}

func (n jsonNumber) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {
	return patch(n, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

// Null implementation
func (n jsonNull) Json(_ ...Option) string {
	return "null"
}

func (n jsonNull) Yaml(_ ...Option) string {
	return "null"
}

func (n jsonNull) raw() interface{} {
	return nil
}

func (n1 jsonNull) Equals(n JsonNode, opts ...Option) bool {
	o := Refine(&Options{Retain: opts}, nil)
	return n1.equals(n, o)
}

func (n1 jsonNull) equals(n JsonNode, _ *Options) bool {
	_, ok := n.(jsonNull)
	return ok
}

func (n jsonNull) hashCode(_ *Options) [8]byte {
	return [8]byte{0x84, 0x27, 0xAD, 0x5B, 0x6B, 0x31, 0x47, 0xF3}
}

func (n jsonNull) Diff(node JsonNode, opts ...Option) Diff {
	o := Refine(&Options{Retain: opts}, nil)
	strategy := getPatchStrategy(o)
	return n.diff(node, make(Path, 0), o, strategy)
}

func (n jsonNull) diff(node JsonNode, path Path, opts *Options, strategy patchStrategy) Diff {
	return basicDiff(n, node, path, opts, strategy)
}

func (n jsonNull) Patch(d Diff) (JsonNode, error) {
	return patchAll(n, d)
}

func (n jsonNull) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {
	return patch(n, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

// Void implementation
func (v voidNode) Json(_ ...Option) string {
	return ""
}

func (v voidNode) Yaml(_ ...Option) string {
	return ""
}

func (v voidNode) raw() interface{} {
	return nil
}

func (v1 voidNode) Equals(n JsonNode, opts ...Option) bool {
	o := Refine(&Options{Retain: opts}, nil)
	return v1.equals(n, o)
}

func (v1 voidNode) equals(n JsonNode, _ *Options) bool {
	_, ok := n.(voidNode)
	return ok
}

func (v voidNode) hashCode(_ *Options) [8]byte {
	return [8]byte{0x56, 0x6F, 0x69, 0x64, 0x00, 0x00, 0x00, 0x00}
}

func (v voidNode) Diff(n JsonNode, opts ...Option) Diff {
	o := Refine(&Options{Retain: opts}, nil)
	strategy := getPatchStrategy(o)
	return v.diff(n, make(Path, 0), o, strategy)
}

func (v voidNode) diff(n JsonNode, path Path, opts *Options, strategy patchStrategy) Diff {
	return basicDiff(v, n, path, opts, strategy)
}

func (v voidNode) Patch(d Diff) (JsonNode, error) {
	return patchAll(v, d)
}

func (v voidNode) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {
	return patch(v, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

// Helper functions that need to be defined
func basicDiff(a, b JsonNode, path Path, opts *Options, strategy patchStrategy) Diff {
	if a.equals(b, opts) {
		return Diff{}
	}
	
	var e DiffElement
	switch strategy {
	case mergePatchStrategy:
		e = DiffElement{
			Path: path,
			Add:  []JsonNode{b},
			Metadata: Metadata{
				Merge: true,
			},
		}
	default:
		e = DiffElement{
			Path:   path,
			Remove: []JsonNode{a},
			Add:    []JsonNode{b},
		}
	}
	return Diff{e}
}

func patchAll(n JsonNode, d Diff) (JsonNode, error) {
	for _, element := range d {
		result, err := n.patch(Path{}, element.Path, []JsonNode{}, element.Remove, element.Add, []JsonNode{}, getPatchStrategyFromMetadata(element.Metadata))
		if err != nil {
			return nil, err
		}
		n = result
	}
	return n, nil
}

func patch(n JsonNode, pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {
	if len(pathAhead) == 0 {
		if len(oldValues) > 1 || len(newValues) > 1 {
			return nil, fmt.Errorf("more than one element in root diff")
		}
		if len(oldValues) > 0 && !n.equals(oldValues[0], &Options{}) {
			return nil, fmt.Errorf("expected %v but found %v", oldValues[0].Json(), n.Json())
		}
		if len(newValues) == 0 {
			return voidNode{}, nil
		}
		return newValues[0], nil
	}
	return nil, fmt.Errorf("patch not supported for this node type")
}

func getPatchStrategyFromMetadata(m Metadata) patchStrategy {
	if m.Merge {
		return mergePatchStrategy
	}
	return strictPatchStrategy
}

// Object implementation
func (o jsonObject) Json(_ ...Option) string {
	return RenderJson(o.raw())
}

func (o jsonObject) Yaml(_ ...Option) string {
	return RenderYaml(o.raw())
}

func (o jsonObject) raw() interface{} {
	m := make(map[string]interface{})
	for k, v := range o {
		m[k] = v.raw()
	}
	return m
}

func (o1 jsonObject) Equals(n JsonNode, opts ...Option) bool {
	o := Refine(&Options{Retain: opts}, nil)
	return o1.equals(n, o)
}

func (o1 jsonObject) equals(n JsonNode, opts *Options) bool {
	o2, ok := n.(jsonObject)
	if !ok {
		return false
	}
	if len(o1) != len(o2) {
		return false
	}
	for k, v1 := range o1 {
		v2, ok := o2[k]
		if !ok {
			return false
		}
		keyOpts := Refine(opts, PathKey(k))
		if !v1.equals(v2, keyOpts) {
			return false
		}
	}
	return true
}

func (o jsonObject) hashCode(opts *Options) [8]byte {
	h := make([]byte, 0)
	for k, v := range o {
		keyOpts := Refine(opts, PathKey(k))
		keyHash := hash([]byte(k))
		valueHash := v.hashCode(keyOpts)
		h = append(h, keyHash[:]...)
		h = append(h, valueHash[:]...)
	}
	return hash(h)
}

func (o jsonObject) Diff(n JsonNode, opts ...Option) Diff {
	opts2 := Refine(&Options{Retain: opts}, nil)
	strategy := getPatchStrategy(opts2)
	return o.diff(n, make(Path, 0), opts2, strategy)
}

func (o jsonObject) diff(n JsonNode, path Path, opts *Options, strategy patchStrategy) Diff {
	o2, ok := n.(jsonObject)
	if !ok {
		return basicDiff(o, n, path, opts, strategy)
	}
	
	d := make(Diff, 0)
	for k, v1 := range o {
		keyPath := append(path, PathKey(k))
		keyOpts := Refine(opts, PathKey(k))
		v2, ok := o2[k]
		if !ok {
			// Key removed
			e := DiffElement{
				Path:   keyPath,
				Remove: []JsonNode{v1},
			}
			d = append(d, e)
		} else {
			// Key exists in both - recurse
			subDiff := v1.diff(v2, keyPath, keyOpts, strategy)
			d = append(d, subDiff...)
		}
	}
	for k, v2 := range o2 {
		if _, ok := o[k]; !ok {
			// Key added
			keyPath := append(path, PathKey(k))
			e := DiffElement{
				Path: keyPath,
				Add:  []JsonNode{v2},
			}
			d = append(d, e)
		}
	}
	return d
}

func (o jsonObject) Patch(d Diff) (JsonNode, error) {
	return patchAll(o, d)
}

func (o jsonObject) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {
	if len(pathAhead) == 0 {
		return patch(o, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
	}
	
	key, ok := pathAhead[0].(PathKey)
	if !ok {
		return nil, fmt.Errorf("expected PathKey, got %T", pathAhead[0])
	}
	
	rest := pathAhead[1:]
	currentValue, exists := o[string(key)]
	if !exists {
		currentValue = voidNode{}
	}
	
	newValue, err := currentValue.patch(append(pathBehind, key), rest, before, oldValues, newValues, after, strategy)
	if err != nil {
		return nil, err
	}
	
	newObject := make(jsonObject)
	for k, v := range o {
		newObject[k] = v
	}
	
	if _, isVoid := newValue.(voidNode); isVoid {
		delete(newObject, string(key))
	} else {
		newObject[string(key)] = newValue
	}
	
	return newObject, nil
}

// Array implementation - basic dispatch version  
func (a jsonArray) Json(_ ...Option) string {
	return RenderJson(a.raw())
}

func (a jsonArray) Yaml(_ ...Option) string {
	return RenderYaml(a.raw())
}

func (a jsonArray) raw() interface{} {
	r := make([]interface{}, len(a))
	for i, n := range a {
		r[i] = n.raw()
	}
	return r
}

func (a1 jsonArray) Equals(n JsonNode, opts ...Option) bool {
	o := Refine(&Options{Retain: opts}, nil)
	return a1.equals(n, o)
}

func (a1 jsonArray) equals(n JsonNode, o *Options) bool {
	// This is basic implementation - real dispatch logic would be in array package
	a2, ok := n.(jsonArray)
	if !ok {
		return false
	}
	if len(a1) != len(a2) {
		return false
	}
	for i, v1 := range a1 {
		v2 := a2[i]
		indexOpts := Refine(o, PathIndex(i))
		if !v1.equals(v2, indexOpts) {
			return false
		}
	}
	return true
}

func (a jsonArray) hashCode(opts *Options) [8]byte {
	b := []byte{0xF5, 0x18, 0x0A, 0x71, 0xA4, 0xC4, 0x03, 0xF3}
	for i, n := range a {
		indexOpts := Refine(opts, PathIndex(i))
		h := n.hashCode(indexOpts)
		b = append(b, h[:]...)
	}
	return hash(b)
}

func (a jsonArray) Diff(n JsonNode, opts ...Option) Diff {
	o := Refine(&Options{Retain: opts}, nil)
	strategy := getPatchStrategy(o)
	return a.diff(n, make(Path, 0), o, strategy)
}

func (a jsonArray) diff(n JsonNode, path Path, opts *Options, strategy patchStrategy) Diff {
	// Basic implementation - real array dispatch would be in array package
	b, ok := n.(jsonArray)
	if !ok {
		return basicDiff(a, n, path, opts, strategy)
	}
	
	// Simple implementation for now
	if a.equals(b, opts) {
		return Diff{}
	}
	
	return basicDiff(a, b, path, opts, strategy)
}

func (a jsonArray) Patch(d Diff) (JsonNode, error) {
	return patchAll(a, d)
}

func (a jsonArray) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {
	return patch(a, pathBehind, pathAhead, before, oldValues, newValues, after, strategy)
}

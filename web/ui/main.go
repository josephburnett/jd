package main

import (
	"fmt"
	"runtime/debug"
	"sync"
	"syscall/js"

	jd "github.com/josephburnett/jd/lib"
)

const (
	commandId               = "command"
	aLabelId                = "a-label"
	aJsonId                 = "a-json"
	aErrorId                = "a-error"
	bLabelId                = "b-label"
	bJsonId                 = "b-json"
	bErrorId                = "b-error"
	diffLabelId             = "diff-label"
	diffId                  = "diff"
	diffErrorId             = "diff-error"
	modeDiffId              = "mode-diff"
	modePatchId             = "mode-patch"
	formatJsonId            = "format-json"
	formatYamlId            = "format-yaml"
	diffFormatJdId          = "diff-format-jd"
	diffFormatPatchId       = "diff-format-patch"
	arrayListId             = "array-list"
	arraySetId              = "array-set"
	arrayMsetId             = "array-mset"
	focusStyle              = "border:solid 3px #080"
	unfocusStyle            = "border:solid 3px #ccc"
	halfWidthStyle          = "width:97%"
	fullWidthStyle          = "width:98.5%"
	diffFormatJdPlaceholder = `@ ["foo"]
- "bar"
+ "baz"
`
	diffFormatPatchPlaceholder = `[{"op":"test","path":"/foo","value":"bar"},{"op":"remove","path":"/foo","value":"bar"},{"op":"add","path":"/foo","value":"baz"}]`
)

func main() {
	if _, err := newApp(); err != nil {
		panic(err)
	}
	select {}
}

type app struct {
	mux        sync.Mutex
	doc        js.Value
	changeCh   chan struct{}
	mode       string
	format     string
	diffFormat string
	array      string
}

func newApp() (*app, error) {
	a := &app{
		changeCh:   make(chan struct{}, 10),
		doc:        js.Global().Get("document"),
		mode:       modeDiffId,
		format:     formatJsonId,
		diffFormat: diffFormatJdId,
		array:      arrayListId,
	}
	for _, id := range []string{
		aJsonId,
		bJsonId,
		diffId,
	} {
		err := a.watchInput(id)
		if err != nil {
			return nil, err
		}
	}
	for _, id := range []string{
		modeDiffId,
		modePatchId,
	} {
		err := a.watchChange(id, &a.mode)
		if err != nil {
			return nil, err
		}
	}
	for _, id := range []string{
		formatJsonId,
		formatYamlId,
	} {
		err := a.watchChange(id, &a.format)
		if err != nil {
			return nil, err
		}
	}
	for _, id := range []string{
		diffFormatJdId,
		diffFormatPatchId,
	} {
		err := a.watchChange(id, &a.diffFormat)
		if err != nil {
			return nil, err
		}
	}
	for _, id := range []string{
		arrayListId,
		arraySetId,
		arrayMsetId,
	} {
		err := a.watchChange(id, &a.array)
		if err != nil {
			return nil, err
		}
	}
	go a.handleChange()
	a.changeCh <- struct{}{}
	return a, nil
}

func (a *app) watchInput(id string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
		defer a.catchPanic()
		a.changeCh <- struct{}{}
		return nil
	}
	element := a.getElementById(id)
	if element.IsNull() {
		return fmt.Errorf("id %v not found", id)
	}
	element.Call("addEventListener", "input", js.FuncOf(listener))
	return nil
}

func (a *app) watchChange(id string, s *string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
		defer a.catchPanic()
		a.mux.Lock()
		defer a.mux.Unlock()
		*s = id
		a.changeCh <- struct{}{}
		return nil
	}
	element := a.getElementById(id)
	if element.IsNull() {
		return fmt.Errorf("id %v not found", id)
	}
	element.Call("addEventListener", "change", js.FuncOf(listener))
	return nil
}

func (a *app) handleChange() {
	defer a.catchPanic()
	for {
		select {
		case <-a.changeCh:
			a.reconcile()
		}
	}
}

func (a *app) reconcile() {
	a.mux.Lock()
	defer a.mux.Unlock()

	a.setDerived()
	a.setPlaceholder()
	a.setCommandLabel()
	a.setInputLabels()
	a.setInputsEnabled()

	switch a.mode {
	case modeDiffId:
		a.printDiff()
	case modePatchId:
		a.printPatch()
	}
}

func (a *app) setDerived() {
	if a.diffFormat == diffFormatPatchId {
		a.array = arrayListId
		for id, val := range map[string]bool{
			arrayListId: true,
			arraySetId:  false,
			arrayMsetId: false,
		} {
			e := a.getElementById(id)
			e.Set("checked", val)
		}
	}
}

func (a *app) setPlaceholder() {
	d := a.getElementById(diffId)
	switch a.diffFormat {
	case diffFormatJdId:
		d.Set("placeholder", diffFormatJdPlaceholder)
	case diffFormatPatchId:
		d.Set("placeholder", diffFormatPatchPlaceholder)
	}
}

func (a *app) setCommandLabel() {
	command := "jd"
	switch a.mode {
	case modePatchId:
		command += " -p"
	default:
	}
	switch a.format {
	case formatYamlId:
		command += " -yaml"
	default:
	}
	switch a.diffFormat {
	case diffFormatPatchId:
		command += " -f patch"
	}
	switch a.array {
	case arraySetId:
		command += " -set"
	case arrayMsetId:
		command += " -mset"
	default:
	}
	switch a.mode {
	case modeDiffId:
		command += " a.json b.json"
	case modePatchId:
		command += " diff a.json"
	default:
	}
	a.setLabel(commandId, command)
}

func (a *app) setInputLabels() {
	aLabel := a.getElementById(aLabelId)
	bLabel := a.getElementById(bLabelId)
	if a.format == formatJsonId {
		aLabel.Set("innerHTML", "a.json")
		bLabel.Set("innerHTML", "b.json")
	} else {
		aLabel.Set("innerHTML", "a.yaml")
		bLabel.Set("innerHTML", "b.yaml")
	}
}

func (a *app) setInputsEnabled() {
	aJson := a.getElementById(aJsonId)
	bJson := a.getElementById(bJsonId)
	diffText := a.getElementById(diffId)
	switch a.mode {
	case modeDiffId:
		aJson.Set("style", focusStyle+";"+halfWidthStyle)
		bJson.Set("readonly", js.ValueOf(false))
		bJson.Set("style", focusStyle+";"+halfWidthStyle)
		diffText.Set("readonly", js.ValueOf(true))
		diffText.Set("style", unfocusStyle+";"+fullWidthStyle)
	case modePatchId:
		aJson.Set("style", focusStyle+";"+halfWidthStyle)
		bJson.Set("readonly", js.ValueOf(true))
		bJson.Set("style", unfocusStyle+";"+halfWidthStyle)
		diffText.Set("readonly", js.ValueOf(false))
		diffText.Set("style", focusStyle+";"+fullWidthStyle)
	default:
	}
	buttons := []string{
		arrayListId,
		arraySetId,
		arrayMsetId,
	}
	for _, id := range buttons {
		e := a.getElementById(id)
		switch a.diffFormat {
		case diffFormatJdId:
			e.Set("disabled", js.ValueOf(false))
		case diffFormatPatchId:
			e.Set("disabled", js.ValueOf(true))
		}
	}
}

func (a *app) getMetadata() []jd.Metadata {
	metadata := []jd.Metadata{}
	switch a.array {
	case arraySetId:
		metadata = append(metadata, jd.SET)
	case arrayMsetId:
		metadata = append(metadata, jd.MULTISET)
	default:
	}
	return metadata
}

func (a *app) parseAndTranslate(id string) (jd.JsonNode, error) {
	value := a.getElementById(id)
	nodeJson, errJson := jd.ReadJsonString(value.Get("value").String())
	nodeYaml, errYaml := jd.ReadYamlString(value.Get("value").String())
	// Translate YAML to JSON.
	if a.format == formatJsonId && errJson != nil && errYaml == nil {
		a.setTextarea(id, nodeYaml.Json())
	}
	// Translate JSON to YAML.
	if a.format == formatYamlId && errJson == nil {
		a.setTextarea(id, nodeJson.Yaml())
	}
	// Return any good parsing results.
	if errJson == nil {
		return nodeJson, nil
	}
	if errYaml == nil {
		return nodeYaml, nil
	}
	// Return an error relevant to the desired format.
	if a.format == formatJsonId {
		return nil, errJson
	} else {
		return nil, errYaml
	}
}

func (a *app) parseAndTranslateDiff() (jd.Diff, error) {
	diffText := a.getElementById(diffId)
	diffJd, errJd := jd.ReadDiffString(diffText.Get("value").String())
	diffPatch, errPatch := jd.ReadPatchString(diffText.Get("value").String())
	// Translate jd to patch.
	if a.diffFormat == diffFormatPatchId && errPatch != nil && errJd == nil {
		patchString, err := diffJd.RenderPatch()
		if err != nil {
			return nil, err
		}
		if patchString == "[]" {
			patchString = ""
		}
		a.setTextarea(diffId, patchString)
	}
	// Translate patch to jd.
	if a.diffFormat == diffFormatJdId && errJd != nil && errPatch == nil {
		a.setTextarea(diffId, diffPatch.Render())
	}
	// Return any good parsing resuilts.
	if errJd == nil {
		return diffJd, nil
	}
	if errPatch == nil {
		return diffPatch, nil
	}
	// Return an error relevant to the desired format.
	if a.diffFormat == diffFormatJdId {
		return nil, errJd
	} else {
		return nil, errPatch
	}
}

func (a *app) printDiff() {
	metadata := a.getMetadata()
	var fail bool
	// Read a
	aNode, err := a.parseAndTranslate(aJsonId)
	if err != nil {
		a.setLabel(aErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(aErrorId, "")
	}
	// Read b
	bNode, err := a.parseAndTranslate(bJsonId)
	if err != nil {
		a.setLabel(bErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(bErrorId, "")
	}
	if fail {
		a.setTextarea(diffId, "")
		return
	}
	// Print diff
	diff := aNode.Diff(bNode, metadata...)
	var out string
	switch a.diffFormat {
	case diffFormatJdId:
		out = diff.Render()
	case diffFormatPatchId:
		out, err = diff.RenderPatch()
		if err != nil {
			a.setLabel(diffErrorId, err.Error())
		}
		if out == "[]" {
			out = ""
		}
	}
	a.setTextarea(diffId, out)
}

func (a *app) printPatch() {
	metadata := a.getMetadata()
	var fail bool
	// Read a
	aNode, err := a.parseAndTranslate(aJsonId)
	if err != nil {
		a.setLabel(aErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(aErrorId, "")
	}
	// Read diff
	diff, err := a.parseAndTranslateDiff()
	if err != nil {
		a.setLabel(diffErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(diffErrorId, "")
	}
	if err != nil {
		a.setLabel(diffErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(diffErrorId, "")
	}
	if fail {
		a.setTextarea(bJsonId, "")
		return
	}
	// Print patch
	bNode, err := aNode.Patch(diff)
	if err != nil {
		a.setLabel(diffErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(diffErrorId, "")
	}
	if fail {
		a.setTextarea(bJsonId, "")
		return
	}
	var out string
	if a.format == formatJsonId {
		out = bNode.Json(metadata...)
	} else {
		out = bNode.Yaml(metadata...)
	}
	a.setTextarea(bJsonId, out)
}

func (a *app) getElementById(id string) js.Value {
	return a.doc.Call("getElementById", id)
}

func (a *app) setLabel(id, msg string) {
	a.getElementById(id).Set("innerHTML", msg)
}

func (a *app) setTextarea(id, text string) {
	a.getElementById(id).Set("value", text)
}

func (a *app) catchPanic() {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		msg := fmt.Sprintf("%v\n\n<pre>%q at:\n\n %v</pre>", crashMessage, r, stack)
		value := a.getElementById("crash")
		value.Set("innerHTML", msg)
		panic(r)
	}
}

const crashMessage = `Jd has crashed. Please report the following error at <a href="https://github.com/josephburnett/jd/issues">https://github.com/josephburnett/jd/issues</a>.`

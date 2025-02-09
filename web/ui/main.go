package main

import (
	"fmt"
	"runtime/debug"
	"sync"
	"syscall/js"

	jd "github.com/josephburnett/jd/v2"
)

const (
	commandId         = "command"
	aLabelId          = "a-label"
	aJsonId           = "a-json"
	aErrorId          = "a-error"
	bLabelId          = "b-label"
	bJsonId           = "b-json"
	bErrorId          = "b-error"
	diffLabelId       = "diff-label"
	diffId            = "diff"
	diffErrorId       = "diff-error"
	modeDiffId        = "mode-diff"
	modePatchId       = "mode-patch"
	formatJsonId      = "format-json"
	formatYamlId      = "format-yaml"
	diffFormatJdId    = "diff-format-jd"
	diffFormatPatchId = "diff-format-patch"
	diffFormatMergeId = "diff-format-merge"
	arrayListId       = "array-list"
	arraySetId        = "array-set"
	arrayMsetId       = "array-mset"
	focusStyle        = "border:solid 3px #080"
	unfocusStyle      = "border:solid 3px #ccc"
	halfWidthStyle    = "width:97%"
	fullWidthStyle    = "width:98.5%"
	placeholderA      = `{"foo":["bar","baz"]}`
	placeholderB      = `{"foo":["bar","baz","bam"]}`
)

func main() {
	if _, err := newApp(); err != nil {
		panic(err)
	}
	select {}
}

type app struct {
	mux            sync.Mutex
	doc            js.Value
	changeCh       chan struct{}
	mode           string
	format         string
	formatLast     string
	diffFormat     string
	diffFormatLast string
	array          string
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
		diffFormatMergeId,
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
	switch a.diffFormat {
	case diffFormatPatchId, diffFormatMergeId:
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

var (
	jdPlaceholderList string
	jdPlaceholderSet  string
	jdPlaceholderMset string
	patchPlaceholder  string
	mergePlaceholder  string
	yamlPlaceholderA  string
	yamlPlaceholderB  string
)

func init() {
	var (
		a   jd.JsonNode
		b   jd.JsonNode
		err error
	)
	a, err = jd.ReadJsonString(placeholderA)
	b, err = jd.ReadJsonString(placeholderB)
	if err != nil {
		panic(err)
	}
	jdPlaceholderList = a.Diff(b).Render()
	jdPlaceholderSet = a.Diff(b, jd.SET).Render()
	jdPlaceholderMset = a.Diff(b, jd.MULTISET).Render()
	patchPlaceholder, err = a.Diff(b).RenderPatch()
	mergePlaceholder, err = a.Diff(b, jd.MERGE).RenderMerge()
	if err != nil {
		panic(err)
	}
	yamlPlaceholderA = a.Yaml()
	yamlPlaceholderB = b.Yaml()
}

func (a *app) setPlaceholder() {
	d := a.getElementById(diffId)
	switch a.diffFormat {
	case diffFormatJdId:
		switch a.array {
		case arrayListId:
			d.Set("placeholder", jdPlaceholderList)
		case arraySetId:
			d.Set("placeholder", jdPlaceholderSet)
		case arrayMsetId:
			d.Set("placeholder", jdPlaceholderMset)
		}
	case diffFormatPatchId:
		d.Set("placeholder", patchPlaceholder)
	case diffFormatMergeId:
		d.Set("placeholder", mergePlaceholder)
	}
	aInput := a.getElementById(aJsonId)
	bInput := a.getElementById(bJsonId)
	switch a.format {
	case formatJsonId:
		aInput.Set("placeholder", placeholderA)
		bInput.Set("placeholder", placeholderB)
	case formatYamlId:
		aInput.Set("placeholder", yamlPlaceholderA)
		bInput.Set("placeholder", yamlPlaceholderB)
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
	case diffFormatMergeId:
		command += " -f merge"
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
		if a.format == formatJsonId {
			command += " a.json b.json"
		} else {
			command += " a.yaml b.yaml"
		}
	case modePatchId:
		if a.format == formatJsonId {
			command += " diff a.json"
		} else {
			command += " diff a.yaml"
		}
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
		case diffFormatPatchId, diffFormatMergeId:
			e.Set("disabled", js.ValueOf(true))
		}
	}
}

func (a *app) getMetadata() []jd.Option {
	options := []jd.Option{}
	switch a.array {
	case arraySetId:
		options = append(options, jd.SET)
	case arrayMsetId:
		options = append(options, jd.MULTISET)
	default:
	}
	switch a.diffFormat {
	case diffFormatMergeId:
		options = append(options, jd.MERGE)
	}
	return options
}

func (a *app) parseAndTranslate(id, formatLast string) (jd.JsonNode, error) {
	change := false
	if a.format != formatLast {
		change = true
	}
	value := a.getElementById(id)
	nodeJson, errJson := jd.ReadJsonString(value.Get("value").String())
	nodeYaml, errYaml := jd.ReadYamlString(value.Get("value").String())
	// Translate YAML to JSON.
	if change && a.format == formatJsonId && errJson != nil && errYaml == nil {
		a.setTextarea(id, nodeYaml.Json())
		return nodeYaml, nil
	}
	// Translate JSON to YAML.
	if change && a.format == formatYamlId && errJson == nil {
		a.setTextarea(id, nodeJson.Yaml())
		return nodeJson, nil
	}
	// Return good parsing results.
	if a.format == formatJsonId && errJson == nil {
		return nodeJson, nil
	}
	if a.format == formatYamlId && errYaml == nil {
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
	change := false
	if a.diffFormat != a.diffFormatLast {
		change = true
		a.diffFormatLast = a.diffFormat
	}
	diffText := a.getElementById(diffId)
	diffJd, errJd := jd.ReadDiffString(diffText.Get("value").String())
	diffPatch, errPatch := jd.ReadPatchString(diffText.Get("value").String())
	diffMerge, errMerge := jd.ReadMergeString(diffText.Get("value").String())
	// Translate jd to patch.
	if change && a.diffFormat == diffFormatPatchId && (errPatch != nil || errMerge != nil) && errJd == nil {
		patchString, err := diffJd.RenderPatch()
		if err != nil {
			return nil, err
		}
		if patchString == "[]" {
			patchString = ""
		}
		a.setTextarea(diffId, patchString)
		return diffJd, nil
	}
	// Translate jd to merge:
	if change && a.diffFormat == diffFormatMergeId && (errPatch != nil || errJd != nil) && errMerge == nil {
		mergeString, err := diffJd.RenderMerge()
		if err != nil {
			return nil, err
		}
		a.setTextarea(diffId, mergeString)
		return diffJd, nil
	}
	// Translate patch to jd.
	if change && a.diffFormat == diffFormatJdId && (errJd != nil || errMerge != nil) && errPatch == nil {
		a.setTextarea(diffId, diffPatch.Render())
		return diffPatch, nil
	}
	// Translate merge to jd.
	if change && a.diffFormat == diffFormatJdId && (errJd != nil || errPatch != nil) && errMerge == nil {
		a.setTextarea(diffId, diffMerge.Render())
		return diffMerge, nil
	}
	// Return good parsing results.
	if a.diffFormat == diffFormatJdId && errJd == nil {
		return diffJd, nil
	}
	if a.diffFormat == diffFormatPatchId && errPatch == nil {
		return diffPatch, nil
	}
	if a.diffFormat == diffFormatMergeId && errMerge == nil {
		return diffMerge, nil
	}
	// Return an error relevant to the desired format.
	switch a.diffFormat {
	case diffFormatJdId:
		return nil, errJd
	case diffFormatPatchId:
		return nil, errPatch
	case diffFormatMergeId:
		return nil, errMerge
	}
	return nil, fmt.Errorf("unsupported diff format: %v", a.diffFormat)
}

func (a *app) printDiff() {
	metadata := a.getMetadata()
	var fail bool
	// Read a
	aNode, err := a.parseAndTranslate(aJsonId, a.formatLast)
	if err != nil {
		a.setLabel(aErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(aErrorId, "")
	}
	// Read b
	bNode, err := a.parseAndTranslate(bJsonId, a.formatLast)
	if err != nil {
		a.setLabel(bErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(bErrorId, "")
	}
	// Mark clean translation
	a.formatLast = a.format
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
	case diffFormatMergeId:
		out, err = diff.RenderMerge()
		if err != nil {
			a.setLabel(diffErrorId, err.Error())
		}
	}
	a.setTextarea(diffId, out)
}

func (a *app) printPatch() {
	metadata := a.getMetadata()
	var fail bool
	// Read a
	aNode, err := a.parseAndTranslate(aJsonId, a.formatLast)
	if err != nil {
		a.setLabel(aErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(aErrorId, "")
	}
	a.formatLast = a.format
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

package main

import (
	"fmt"
	"sync"
	"syscall/js"

	jd "github.com/josephburnett/jd/lib"
)

const (
	commandId   = "command"
	aLabelId    = "a-label"
	aJsonId     = "a-json"
	aErrorId    = "a-error"
	bLabelId    = "b-label"
	bJsonId     = "b-json"
	bErrorId    = "b-error"
	diffLabelId = "diff-label"
	diffId      = "diff"
	diffErrorId = "diff-error"
	modeDiffId  = "mode-diff"
	modePatchId = "mode-patch"
	arrayListId = "array-list"
	arraySetId  = "array-set"
	arrayMsetId = "array-mset"
)

func main() {
	if _, err := newApp(); err != nil {
		panic(err)
	}
	select {}
}

type app struct {
	mux      sync.Mutex
	doc      js.Value
	changeCh chan struct{}
	mode     string
	array    string
}

func newApp() (*app, error) {
	a := &app{
		changeCh: make(chan struct{}, 10),
		doc:      js.Global().Get("document"),
		mode:     modeDiffId,
		array:    arrayListId,
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
		err := a.watchMode(id)
		if err != nil {
			return nil, err
		}
	}
	for _, id := range []string{
		arrayListId,
		arraySetId,
		arrayMsetId,
	} {
		err := a.watchArray(id)
		if err != nil {
			return nil, err
		}
	}
	go a.handleChange()
	return a, nil
}

func (a *app) watchInput(id string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
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

func (a *app) watchArray(id string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
		a.mux.Lock()
		defer a.mux.Unlock()
		a.array = id
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

func (a *app) watchMode(id string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
		a.mux.Lock()
		defer a.mux.Unlock()
		a.mode = id
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

	// set the command label
	command := "jd"
	switch a.mode {
	case modePatchId:
		command += " -p"
	default:
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
		command += " a.json diff"
	default:
	}
	a.setLabel(commandId, command)

	// Enable / disable inputs based on mode
	aJson := a.getElementById(aJsonId)
	bJson := a.getElementById(bJsonId)
	diffText := a.getElementById(diffId)
	switch a.mode {
	case modeDiffId:
		bJson.Set("readonly", js.ValueOf(false))
		diffText.Set("readonly", js.ValueOf(true))
	case modePatchId:
		bJson.Set("readonly", js.ValueOf(true))
		diffText.Set("readonly", js.ValueOf(false))
	default:
	}

	// Chose array semantic metadata
	metadata := []jd.Metadata{}
	switch a.array {
	case arraySetId:
		metadata = append(metadata, jd.SET)
	case arrayMsetId:
		metadata = append(metadata, jd.MULTISET)
	default:
	}

	var fail bool

	// Parse a.json
	aNode, err := jd.ReadJsonString(aJson.Get("value").String())
	if err != nil {
		a.setLabel(aErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(aErrorId, "")
	}

	switch a.mode {
	case modeDiffId:
		// parse b.json
		bNode, err := jd.ReadJsonString(bJson.Get("value").String())
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

		// produce diff
		diff := aNode.Diff(bNode, metadata...)
		a.setTextarea(diffId, diff.Render())
	case modePatchId:
		// parse diff
		diff, err := jd.ReadDiffString(diffText.Get("value").String())
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

		// produce b.json
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

		a.setTextarea(bJsonId, bNode.Json(metadata...))
	}
	return
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

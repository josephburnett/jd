package main

import (
	"fmt"
	"sync"
	"syscall/js"

	jd "github.com/josephburnett/jd/lib"
)

const (
	aLabelId    = "a-label"
	aJsonId     = "a-json"
	aErrorId    = "a-error"
	bLabelId    = "b-label"
	bJsonId     = "b-json"
	bErrorId    = "b-error"
	diffLabelId = "diff-label"
	diffId      = "diff"
	diffErrorId = "diff-error"
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
	patch    bool
	set      bool
	mset     bool
}

func newApp() (*app, error) {
	a := &app{
		changeCh: make(chan struct{}, 10),
		doc:      js.Global().Get("document"),
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
	go a.handleInput()
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

func (a *app) handleInput() {
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

	var fail bool

	aJson := a.getElementById(aJsonId)
	aNode, err := jd.ReadJsonString(aJson.Get("value").String())
	if err != nil {
		a.setLabel(aErrorId, err.Error())
		fail = true
	} else {
		a.setLabel(aErrorId, "")
	}

	bJson := a.getElementById(bJsonId)
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

	diff := aNode.Diff(bNode)
	a.setTextarea(diffId, diff.Render())
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

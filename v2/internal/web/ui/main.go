package main

import (
	"fmt"
	"runtime/debug"
	"sync"
	"syscall/js"

	"github.com/josephburnett/jd/v2/internal/web/presenter"
)

func main() {
	if _, err := newApp(); err != nil {
		panic(err)
	}
	select {}
}

type app struct {
	mux       sync.Mutex
	changeCh  chan struct{}
	presenter *presenter.Presenter
	view      *domView
}

func newApp() (*app, error) {
	view := newDomView()
	pres, err := presenter.New(view)
	if err != nil {
		return nil, err
	}

	a := &app{
		changeCh:  make(chan struct{}, 10),
		presenter: pres,
		view:      view.(*domView),
	}

	// Set up event listeners
	for _, id := range []string{
		presenter.AJsonId,
		presenter.BJsonId,
		presenter.DiffId,
	} {
		err := a.watchInput(id)
		if err != nil {
			return nil, err
		}
	}

	// Set up mode click handlers
	for _, id := range []string{
		presenter.ModeDiffId,
		presenter.ModePatchId,
	} {
		err := a.watchModeClick(id)
		if err != nil {
			return nil, err
		}
	}

	// Set up format click handlers
	for _, id := range []string{
		presenter.FormatJsonId,
		presenter.FormatYamlId,
	} {
		err := a.watchFormatClick(id)
		if err != nil {
			return nil, err
		}
	}

	// Set up diff format click handlers
	for _, id := range []string{
		presenter.DiffFormatJdId,
		presenter.DiffFormatPatchId,
		presenter.DiffFormatMergeId,
	} {
		err := a.watchDiffFormatClick(id)
		if err != nil {
			return nil, err
		}
	}

	// Watch options JSON textarea
	err = a.watchOptionsChange(presenter.OptionsJsonId)
	if err != nil {
		return nil, err
	}

	// Set up example buttons
	err = a.setupExampleButtons()
	if err != nil {
		return nil, err
	}

	go a.handleChange()
	a.changeCh <- struct{}{} // Initial reconcile
	return a, nil
}

func (a *app) watchInput(id string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
		defer a.catchPanic()
		a.changeCh <- struct{}{}
		return nil
	}
	element := a.view.getElementById(id)
	if element.IsNull() {
		return fmt.Errorf("id %v not found", id)
	}
	element.Call("addEventListener", "input", js.FuncOf(listener))
	return nil
}

func (a *app) watchModeClick(id string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
		defer a.catchPanic()
		a.mux.Lock()
		defer a.mux.Unlock()

		// Update styling for mode selection
		a.updateModeStyles(id)
		a.updatePresenterState()
		a.changeCh <- struct{}{}
		return nil
	}
	element := a.view.getElementById(id)
	if element.IsNull() {
		return fmt.Errorf("id %v not found", id)
	}
	element.Call("addEventListener", "click", js.FuncOf(listener))
	return nil
}

func (a *app) watchFormatClick(id string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
		defer a.catchPanic()
		a.mux.Lock()
		defer a.mux.Unlock()

		// Update styling for format selection
		a.updateFormatStyles(id)
		a.updatePresenterState()
		a.changeCh <- struct{}{}
		return nil
	}
	element := a.view.getElementById(id)
	if element.IsNull() {
		return fmt.Errorf("id %v not found", id)
	}
	element.Call("addEventListener", "click", js.FuncOf(listener))
	return nil
}

func (a *app) watchDiffFormatClick(id string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
		defer a.catchPanic()
		a.mux.Lock()
		defer a.mux.Unlock()

		// Update styling for diff format selection
		a.updateDiffFormatStyles(id)
		a.updatePresenterState()
		a.changeCh <- struct{}{}
		return nil
	}
	element := a.view.getElementById(id)
	if element.IsNull() {
		return fmt.Errorf("id %v not found", id)
	}
	element.Call("addEventListener", "click", js.FuncOf(listener))
	return nil
}

func (a *app) watchOptionsChange(id string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
		defer a.catchPanic()
		a.mux.Lock()
		defer a.mux.Unlock()
		a.updatePresenterState()
		a.changeCh <- struct{}{}
		return nil
	}
	element := a.view.getElementById(id)
	if element.IsNull() {
		return fmt.Errorf("id %v not found", id)
	}
	element.Call("addEventListener", "input", js.FuncOf(listener))
	return nil
}

func (a *app) setupExampleButtons() error {
	examples := map[string]string{
		"example-set":            `["SET"]`,
		"example-multiset":       `["MULTISET"]`,
		"example-precision":      `[{"precision": 0.1}]`,
		"example-setkeys":        `[{"setkeys": ["id", "name"]}]`,
		"example-path-set":       `[{"@": ["users"], "^": ["SET"]}]`,
		"example-path-precision": `[{"@": ["scores", 0], "^": [{"precision": 0.01}]}]`,
		"example-selective-diff": `["DIFF_OFF", {"@": ["important"], "^": ["DIFF_ON"]}, {"@": ["metadata", "timestamps"], "^": ["DIFF_ON"]}]`,
		"example-mixed":          `["SET", {"@": ["users"], "^": ["MULTISET"]}, {"@": ["temperature"], "^": [{"precision": 0.1}]}]`,
		"example-clear":          `[]`,
	}

	for buttonId, jsonExample := range examples {
		err := a.setupExampleButton(buttonId, jsonExample)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *app) setupExampleButton(buttonId, jsonExample string) error {
	listener := func(_ js.Value, _ []js.Value) interface{} {
		defer a.catchPanic()
		a.mux.Lock()
		defer a.mux.Unlock()

		// Set the options JSON textarea value
		a.view.SetTextarea(presenter.OptionsJsonId, jsonExample)

		// Update presenter state
		a.updatePresenterState()
		a.changeCh <- struct{}{}
		return nil
	}

	element := a.view.getElementById(buttonId)
	if element.IsNull() {
		return fmt.Errorf("example button id %v not found", buttonId)
	}
	element.Call("addEventListener", "click", js.FuncOf(listener))
	return nil
}

func (a *app) updatePresenterState() {
	mode := a.getCurrentMode()
	format := a.getCurrentFormat()
	diffFormat := a.getCurrentDiffFormat()
	optionsJSON := a.view.GetValue(presenter.OptionsJsonId)
	a.presenter.UpdateState(mode, format, diffFormat, optionsJSON)
}

func (a *app) updateModeStyles(selectedId string) {
	// Reset all mode styles
	a.view.getElementById(presenter.ModeDiffId).Set("style", "cursor: pointer; margin: 0 5px;")
	a.view.getElementById(presenter.ModePatchId).Set("style", "cursor: pointer; margin: 0 5px;")

	// Highlight selected
	selectedElement := a.view.getElementById(selectedId)
	selectedElement.Set("style", "cursor: pointer; color: #080; text-decoration: underline; margin: 0 5px;")
}

func (a *app) updateFormatStyles(selectedId string) {
	// Reset all format styles
	a.view.getElementById(presenter.FormatJsonId).Set("style", "cursor: pointer; margin: 0 5px;")
	a.view.getElementById(presenter.FormatYamlId).Set("style", "cursor: pointer; margin: 0 5px;")

	// Highlight selected
	selectedElement := a.view.getElementById(selectedId)
	selectedElement.Set("style", "cursor: pointer; color: #080; text-decoration: underline; margin: 0 5px;")
}

func (a *app) updateDiffFormatStyles(selectedId string) {
	// Reset all diff format styles
	a.view.getElementById(presenter.DiffFormatJdId).Set("style", "cursor: pointer; margin: 0 3px;")
	a.view.getElementById(presenter.DiffFormatPatchId).Set("style", "cursor: pointer; margin: 0 3px;")
	a.view.getElementById(presenter.DiffFormatMergeId).Set("style", "cursor: pointer; margin: 0 3px;")

	// Highlight selected
	selectedElement := a.view.getElementById(selectedId)
	selectedElement.Set("style", "cursor: pointer; color: #080; text-decoration: underline; margin: 0 3px;")
}

func (a *app) getCurrentMode() presenter.Mode {
	// Check which mode element has the underline style (is selected)
	diffElement := a.view.getElementById(presenter.ModeDiffId)
	if diffElement.Get("style").Get("textDecoration").String() == "underline" {
		return presenter.ModeDiff
	}
	return presenter.ModePatch
}

func (a *app) getCurrentFormat() presenter.Format {
	// Check which format element has the underline style (is selected)
	jsonElement := a.view.getElementById(presenter.FormatJsonId)
	if jsonElement.Get("style").Get("textDecoration").String() == "underline" {
		return presenter.FormatJSON
	}
	return presenter.FormatYAML
}

func (a *app) getCurrentDiffFormat() presenter.DiffFormat {
	// Check which diff format element has the underline style (is selected)
	jdElement := a.view.getElementById(presenter.DiffFormatJdId)
	patchElement := a.view.getElementById(presenter.DiffFormatPatchId)

	if jdElement.Get("style").Get("textDecoration").String() == "underline" {
		return presenter.DiffFormatJd
	}
	if patchElement.Get("style").Get("textDecoration").String() == "underline" {
		return presenter.DiffFormatPatch
	}
	return presenter.DiffFormatMerge
}

func (a *app) handleChange() {
	defer a.catchPanic()
	for {
		select {
		case <-a.changeCh:
			a.presenter.Reconcile()
		}
	}
}

func (a *app) catchPanic() {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		msg := fmt.Sprintf("%v\n\n<pre>%q at:\n\n %v</pre>", crashMessage, r, stack)
		crashElement := a.view.getElementById("crash")
		crashElement.Set("innerHTML", msg)
		panic(r)
	}
}

const crashMessage = `Jd has crashed. Please report the following error at <a href="https://github.com/josephburnett/jd/issues">https://github.com/josephburnett/jd/issues</a>.`

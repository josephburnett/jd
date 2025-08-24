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

	for _, id := range []string{
		presenter.ModeDiffId,
		presenter.ModePatchId,
	} {
		err := a.watchModeChange(id)
		if err != nil {
			return nil, err
		}
	}

	for _, id := range []string{
		presenter.FormatJsonId,
		presenter.FormatYamlId,
	} {
		err := a.watchFormatChange(id)
		if err != nil {
			return nil, err
		}
	}

	for _, id := range []string{
		presenter.DiffFormatJdId,
		presenter.DiffFormatPatchId,
		presenter.DiffFormatMergeId,
	} {
		err := a.watchDiffFormatChange(id)
		if err != nil {
			return nil, err
		}
	}

	for _, id := range []string{
		presenter.ArrayListId,
		presenter.ArraySetId,
		presenter.ArrayMsetId,
	} {
		err := a.watchArrayChange(id)
		if err != nil {
			return nil, err
		}
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

func (a *app) watchModeChange(id string) error {
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
	element.Call("addEventListener", "change", js.FuncOf(listener))
	return nil
}

func (a *app) watchFormatChange(id string) error {
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
	element.Call("addEventListener", "change", js.FuncOf(listener))
	return nil
}

func (a *app) watchDiffFormatChange(id string) error {
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
	element.Call("addEventListener", "change", js.FuncOf(listener))
	return nil
}

func (a *app) watchArrayChange(id string) error {
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
	element.Call("addEventListener", "change", js.FuncOf(listener))
	return nil
}

func (a *app) updatePresenterState() {
	mode := a.getCurrentMode()
	format := a.getCurrentFormat()
	diffFormat := a.getCurrentDiffFormat()
	array := a.getCurrentArray()
	a.presenter.UpdateState(mode, format, diffFormat, array)
}

func (a *app) getCurrentMode() presenter.Mode {
	if a.view.IsChecked(presenter.ModeDiffId) {
		return presenter.ModeDiff
	}
	return presenter.ModePatch
}

func (a *app) getCurrentFormat() presenter.Format {
	if a.view.IsChecked(presenter.FormatJsonId) {
		return presenter.FormatJSON
	}
	return presenter.FormatYAML
}

func (a *app) getCurrentDiffFormat() presenter.DiffFormat {
	if a.view.IsChecked(presenter.DiffFormatJdId) {
		return presenter.DiffFormatJd
	}
	if a.view.IsChecked(presenter.DiffFormatPatchId) {
		return presenter.DiffFormatPatch
	}
	return presenter.DiffFormatMerge
}

func (a *app) getCurrentArray() presenter.ArrayType {
	if a.view.IsChecked(presenter.ArrayListId) {
		return presenter.ArrayList
	}
	if a.view.IsChecked(presenter.ArraySetId) {
		return presenter.ArraySet
	}
	return presenter.ArrayMset
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
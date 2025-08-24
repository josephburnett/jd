package main

import (
	"syscall/js"

	"github.com/josephburnett/jd/v2/internal/web/view"
)

// domView implements the View interface using DOM operations
type domView struct {
	doc js.Value
}

// newDomView creates a new DOM view
func newDomView() view.View {
	return &domView{
		doc: js.Global().Get("document"),
	}
}

// GetValue returns the value of an input element
func (d *domView) GetValue(id string) string {
	return d.getElementById(id).Get("value").String()
}

// SetTextarea sets the value of a textarea element
func (d *domView) SetTextarea(id, text string) {
	d.getElementById(id).Set("value", text)
}

// SetLabel sets the innerHTML of a label element
func (d *domView) SetLabel(id, text string) {
	d.getElementById(id).Set("innerHTML", text)
}

// IsChecked returns whether a checkbox/radio is checked
func (d *domView) IsChecked(id string) bool {
	return d.getElementById(id).Get("checked").Bool()
}

// SetChecked sets the checked state of a checkbox/radio
func (d *domView) SetChecked(id string, checked bool) {
	d.getElementById(id).Set("checked", js.ValueOf(checked))
}

// SetDisabled sets the disabled state of an element
func (d *domView) SetDisabled(id string, disabled bool) {
	d.getElementById(id).Set("disabled", js.ValueOf(disabled))
}

// SetReadonly sets the readonly state of an input element
func (d *domView) SetReadonly(id string, readonly bool) {
	d.getElementById(id).Set("readonly", js.ValueOf(readonly))
}

// SetStyle sets the style attribute of an element
func (d *domView) SetStyle(id, style string) {
	d.getElementById(id).Set("style", style)
}

// SetPlaceholder sets the placeholder text of an input element
func (d *domView) SetPlaceholder(id, text string) {
	d.getElementById(id).Set("placeholder", text)
}

// getElementById is a helper to get an element by ID
func (d *domView) getElementById(id string) js.Value {
	return d.doc.Call("getElementById", id)
}

package presenter

import "github.com/josephburnett/jd/v2/internal/web/view"

// mockView implements view.View for testing
type mockView struct {
	values       map[string]string
	labels       map[string]string
	placeholders map[string]string
	checked      map[string]bool
	disabled     map[string]bool
	readonly     map[string]bool
	styles       map[string]string
}

// newMockView creates a new mock view for testing
func newMockView() view.View {
	return &mockView{
		values:       make(map[string]string),
		labels:       make(map[string]string),
		placeholders: make(map[string]string),
		checked:      make(map[string]bool),
		disabled:     make(map[string]bool),
		readonly:     make(map[string]bool),
		styles:       make(map[string]string),
	}
}

func (m *mockView) GetValue(id string) string {
	return m.values[id]
}

func (m *mockView) SetTextarea(id, text string) {
	m.values[id] = text
}

func (m *mockView) SetLabel(id, text string) {
	m.labels[id] = text
}

func (m *mockView) IsChecked(id string) bool {
	return m.checked[id]
}

func (m *mockView) SetChecked(id string, checked bool) {
	m.checked[id] = checked
}

func (m *mockView) SetDisabled(id string, disabled bool) {
	m.disabled[id] = disabled
}

func (m *mockView) SetReadonly(id string, readonly bool) {
	m.readonly[id] = readonly
}

func (m *mockView) SetStyle(id, style string) {
	m.styles[id] = style
}

func (m *mockView) SetPlaceholder(id, text string) {
	m.placeholders[id] = text
}

package view

// View defines the interface for DOM operations that the presenter needs.
// This abstraction allows the presenter logic to be tested without browser dependencies.
type View interface {
	// Input/Output operations
	GetValue(id string) string
	SetTextarea(id, text string)
	SetLabel(id, text string)

	// State management
	IsChecked(id string) bool
	SetChecked(id string, checked bool)
	SetDisabled(id string, disabled bool)
	SetReadonly(id string, readonly bool)
	SetStyle(id, style string)

	// Placeholder management
	SetPlaceholder(id, text string)
}

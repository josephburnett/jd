package presenter

// State represents the current application state
type State struct {
	Mode            Mode
	Format          Format
	DiffFormat      DiffFormat
	OptionsJSON     string // JSON string for advanced options
	FormatLast      Format
	DiffFormatLast  DiffFormat
	ValidationError string // Error from parsing options JSON
}

// Mode represents the application operation mode
type Mode string

const (
	ModeDiff  Mode = "mode-diff"
	ModePatch Mode = "mode-patch"
)

// Format represents the data format for inputs
type Format string

const (
	FormatJSON Format = "format-json"
	FormatYAML Format = "format-yaml"
)

// DiffFormat represents the output diff format
type DiffFormat string

const (
	DiffFormatJd    DiffFormat = "diff-format-jd"
	DiffFormatPatch DiffFormat = "diff-format-patch"
	DiffFormatMerge DiffFormat = "diff-format-merge"
)

// ArrayType represents how arrays should be treated
type ArrayType string

const (
	ArrayList ArrayType = "array-list"
	ArraySet  ArrayType = "array-set"
	ArrayMset ArrayType = "array-mset"
)

// NewState creates a new state with default values
func NewState() *State {
	return &State{
		Mode:        ModeDiff,
		Format:      FormatJSON,
		DiffFormat:  DiffFormatJd,
		OptionsJSON: "[]",
	}
}

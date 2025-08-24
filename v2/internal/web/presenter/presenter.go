package presenter

import (
	"fmt"

	jd "github.com/josephburnett/jd/v2"
	"github.com/josephburnett/jd/v2/internal/web/view"
)

const (
	// Element IDs
	CommandId         = "command"
	ALabelId          = "a-label"
	AJsonId           = "a-json"
	AErrorId          = "a-error"
	BLabelId          = "b-label"
	BJsonId           = "b-json"
	BErrorId          = "b-error"
	DiffLabelId       = "diff-label"
	DiffId            = "diff"
	DiffErrorId       = "diff-error"
	ModeDiffId        = "mode-diff"
	ModePatchId       = "mode-patch"
	FormatJsonId      = "format-json"
	FormatYamlId      = "format-yaml"
	DiffFormatJdId    = "diff-format-jd"
	DiffFormatPatchId = "diff-format-patch"
	DiffFormatMergeId = "diff-format-merge"
	ArrayListId       = "array-list"
	ArraySetId        = "array-set"
	ArrayMsetId       = "array-mset"

	// Styles
	FocusStyle     = "border:solid 3px #080"
	UnfocusStyle   = "border:solid 3px #ccc"
	HalfWidthStyle = "width:97%"
	FullWidthStyle = "width:98.5%"

	// Placeholders
	PlaceholderA = `{"foo":["bar","baz"]}`
	PlaceholderB = `{"foo":["bar","baz","bam"]}`
)

// Placeholders holds pre-computed placeholder values
type Placeholders struct {
	JdPlaceholderList string
	JdPlaceholderSet  string
	JdPlaceholderMset string
	PatchPlaceholder  string
	MergePlaceholder  string
	YamlPlaceholderA  string
	YamlPlaceholderB  string
}

// Presenter contains the core UI logic without DOM dependencies
type Presenter struct {
	state        *State
	placeholders *Placeholders
	view         view.View
}

// New creates a new presenter with the given view
func New(v view.View) (*Presenter, error) {
	placeholders, err := initPlaceholders()
	if err != nil {
		return nil, err
	}

	return &Presenter{
		state:        NewState(),
		placeholders: placeholders,
		view:         v,
	}, nil
}

// initPlaceholders pre-computes all placeholder values
func initPlaceholders() (*Placeholders, error) {
	a, err := jd.ReadJsonString(PlaceholderA)
	if err != nil {
		return nil, err
	}
	b, err := jd.ReadJsonString(PlaceholderB)
	if err != nil {
		return nil, err
	}

	jdPlaceholderList := a.Diff(b).Render()
	jdPlaceholderSet := a.Diff(b, jd.SET).Render()
	jdPlaceholderMset := a.Diff(b, jd.MULTISET).Render()

	patchPlaceholder, err := a.Diff(b).RenderPatch()
	if err != nil {
		return nil, err
	}

	mergePlaceholder, err := a.Diff(b, jd.MERGE).RenderMerge()
	if err != nil {
		return nil, err
	}

	return &Placeholders{
		JdPlaceholderList: jdPlaceholderList,
		JdPlaceholderSet:  jdPlaceholderSet,
		JdPlaceholderMset: jdPlaceholderMset,
		PatchPlaceholder:  patchPlaceholder,
		MergePlaceholder:  mergePlaceholder,
		YamlPlaceholderA:  a.Yaml(),
		YamlPlaceholderB:  b.Yaml(),
	}, nil
}

// UpdateState updates the presenter state and triggers reconciliation
func (p *Presenter) UpdateState(mode Mode, format Format, diffFormat DiffFormat, array ArrayType) {
	p.state.Mode = mode
	p.state.Format = format
	p.state.DiffFormat = diffFormat
	p.state.Array = array
	p.Reconcile()
}

// Reconcile updates the UI based on current state
func (p *Presenter) Reconcile() {
	p.setDerived()
	p.setPlaceholder()
	p.setCommandLabel()
	p.setInputLabels()
	p.setInputsEnabled()

	switch p.state.Mode {
	case ModeDiff:
		p.printDiff()
	case ModePatch:
		p.printPatch()
	}
}

// setDerived sets derived state based on current selections
func (p *Presenter) setDerived() {
	switch p.state.DiffFormat {
	case DiffFormatPatch, DiffFormatMerge:
		p.state.Array = ArrayList
		p.view.SetChecked(string(ArrayList), true)
		p.view.SetChecked(string(ArraySet), false)
		p.view.SetChecked(string(ArrayMset), false)
	}
}

// setPlaceholder updates placeholder text based on current format
func (p *Presenter) setPlaceholder() {
	switch p.state.DiffFormat {
	case DiffFormatJd:
		switch p.state.Array {
		case ArrayList:
			p.view.SetPlaceholder(DiffId, p.placeholders.JdPlaceholderList)
		case ArraySet:
			p.view.SetPlaceholder(DiffId, p.placeholders.JdPlaceholderSet)
		case ArrayMset:
			p.view.SetPlaceholder(DiffId, p.placeholders.JdPlaceholderMset)
		}
	case DiffFormatPatch:
		p.view.SetPlaceholder(DiffId, p.placeholders.PatchPlaceholder)
	case DiffFormatMerge:
		p.view.SetPlaceholder(DiffId, p.placeholders.MergePlaceholder)
	}

	switch p.state.Format {
	case FormatJSON:
		p.view.SetPlaceholder(AJsonId, PlaceholderA)
		p.view.SetPlaceholder(BJsonId, PlaceholderB)
	case FormatYAML:
		p.view.SetPlaceholder(AJsonId, p.placeholders.YamlPlaceholderA)
		p.view.SetPlaceholder(BJsonId, p.placeholders.YamlPlaceholderB)
	}
}

// setCommandLabel updates the command label based on current settings
func (p *Presenter) setCommandLabel() {
	command := "jd"

	switch p.state.Mode {
	case ModePatch:
		command += " -p"
	}

	switch p.state.Format {
	case FormatYAML:
		command += " -yaml"
	}

	switch p.state.DiffFormat {
	case DiffFormatPatch:
		command += " -f patch"
	case DiffFormatMerge:
		command += " -f merge"
	}

	switch p.state.Array {
	case ArraySet:
		command += " -set"
	case ArrayMset:
		command += " -mset"
	}

	switch p.state.Mode {
	case ModeDiff:
		if p.state.Format == FormatJSON {
			command += " a.json b.json"
		} else {
			command += " a.yaml b.yaml"
		}
	case ModePatch:
		if p.state.Format == FormatJSON {
			command += " diff a.json"
		} else {
			command += " diff a.yaml"
		}
	}

	p.view.SetLabel(CommandId, command)
}

// setInputLabels updates input field labels based on format
func (p *Presenter) setInputLabels() {
	if p.state.Format == FormatJSON {
		p.view.SetLabel(ALabelId, "a.json")
		p.view.SetLabel(BLabelId, "b.json")
	} else {
		p.view.SetLabel(ALabelId, "a.yaml")
		p.view.SetLabel(BLabelId, "b.yaml")
	}
}

// setInputsEnabled manages which inputs are enabled based on mode
func (p *Presenter) setInputsEnabled() {
	switch p.state.Mode {
	case ModeDiff:
		p.view.SetStyle(AJsonId, FocusStyle+";"+HalfWidthStyle)
		p.view.SetReadonly(BJsonId, false)
		p.view.SetStyle(BJsonId, FocusStyle+";"+HalfWidthStyle)
		p.view.SetReadonly(DiffId, true)
		p.view.SetStyle(DiffId, UnfocusStyle+";"+FullWidthStyle)
	case ModePatch:
		p.view.SetStyle(AJsonId, FocusStyle+";"+HalfWidthStyle)
		p.view.SetReadonly(BJsonId, true)
		p.view.SetStyle(BJsonId, UnfocusStyle+";"+HalfWidthStyle)
		p.view.SetReadonly(DiffId, false)
		p.view.SetStyle(DiffId, FocusStyle+";"+FullWidthStyle)
	}

	buttons := []string{ArrayListId, ArraySetId, ArrayMsetId}
	for _, id := range buttons {
		switch p.state.DiffFormat {
		case DiffFormatJd:
			p.view.SetDisabled(id, false)
		case DiffFormatPatch, DiffFormatMerge:
			p.view.SetDisabled(id, true)
		}
	}
}

// getMetadata returns jd options based on current state
func (p *Presenter) getMetadata() []jd.Option {
	options := []jd.Option{}
	switch p.state.Array {
	case ArraySet:
		options = append(options, jd.SET)
	case ArrayMset:
		options = append(options, jd.MULTISET)
	}
	switch p.state.DiffFormat {
	case DiffFormatMerge:
		options = append(options, jd.MERGE)
	}
	return options
}

// parseAndTranslate parses input and handles format conversion
func (p *Presenter) parseAndTranslate(id string, formatLast Format) (jd.JsonNode, error) {
	change := p.state.Format != formatLast
	value := p.view.GetValue(id)

	nodeJson, errJson := jd.ReadJsonString(value)
	nodeYaml, errYaml := jd.ReadYamlString(value)

	// Translate YAML to JSON
	if change && p.state.Format == FormatJSON && errJson != nil && errYaml == nil {
		p.view.SetTextarea(id, nodeYaml.Json())
		return nodeYaml, nil
	}
	
	// Translate JSON to YAML
	if change && p.state.Format == FormatYAML && errJson == nil {
		p.view.SetTextarea(id, nodeJson.Yaml())
		return nodeJson, nil
	}
	
	// Return good parsing results
	if p.state.Format == FormatJSON && errJson == nil {
		return nodeJson, nil
	}
	if p.state.Format == FormatYAML && errYaml == nil {
		return nodeYaml, nil
	}
	
	// Return error relevant to desired format
	if p.state.Format == FormatJSON {
		return nil, errJson
	}
	return nil, errYaml
}

// parseAndTranslateDiff parses diff input and handles format conversion
func (p *Presenter) parseAndTranslateDiff() (jd.Diff, error) {
	change := p.state.DiffFormat != p.state.DiffFormatLast
	if change {
		p.state.DiffFormatLast = p.state.DiffFormat
	}

	diffText := p.view.GetValue(DiffId)
	diffJd, errJd := jd.ReadDiffString(diffText)
	diffPatch, errPatch := jd.ReadPatchString(diffText)
	diffMerge, errMerge := jd.ReadMergeString(diffText)

	// Handle format conversions
	if change && p.state.DiffFormat == DiffFormatPatch && (errPatch != nil || errMerge != nil) && errJd == nil {
		patchString, err := diffJd.RenderPatch()
		if err != nil {
			return nil, err
		}
		if patchString == "[]" {
			patchString = ""
		}
		p.view.SetTextarea(DiffId, patchString)
		return diffJd, nil
	}

	if change && p.state.DiffFormat == DiffFormatMerge && (errPatch != nil || errJd != nil) && errMerge == nil {
		mergeString, err := diffJd.RenderMerge()
		if err != nil {
			return nil, err
		}
		p.view.SetTextarea(DiffId, mergeString)
		return diffJd, nil
	}

	if change && p.state.DiffFormat == DiffFormatJd && (errJd != nil || errMerge != nil) && errPatch == nil {
		p.view.SetTextarea(DiffId, diffPatch.Render())
		return diffPatch, nil
	}

	if change && p.state.DiffFormat == DiffFormatJd && (errJd != nil || errPatch != nil) && errMerge == nil {
		p.view.SetTextarea(DiffId, diffMerge.Render())
		return diffMerge, nil
	}

	// Return good parsing results
	if p.state.DiffFormat == DiffFormatJd && errJd == nil {
		return diffJd, nil
	}
	if p.state.DiffFormat == DiffFormatPatch && errPatch == nil {
		return diffPatch, nil
	}
	if p.state.DiffFormat == DiffFormatMerge && errMerge == nil {
		return diffMerge, nil
	}

	// Return error relevant to desired format
	switch p.state.DiffFormat {
	case DiffFormatJd:
		return nil, errJd
	case DiffFormatPatch:
		return nil, errPatch
	case DiffFormatMerge:
		return nil, errMerge
	}
	return nil, fmt.Errorf("unsupported diff format: %v", p.state.DiffFormat)
}

// printDiff handles diff mode computation and output
func (p *Presenter) printDiff() {
	metadata := p.getMetadata()
	var fail bool

	// Read a
	aNode, err := p.parseAndTranslate(AJsonId, p.state.FormatLast)
	if err != nil {
		p.view.SetLabel(AErrorId, err.Error())
		fail = true
	} else {
		p.view.SetLabel(AErrorId, "")
	}

	// Read b
	bNode, err := p.parseAndTranslate(BJsonId, p.state.FormatLast)
	if err != nil {
		p.view.SetLabel(BErrorId, err.Error())
		fail = true
	} else {
		p.view.SetLabel(BErrorId, "")
	}

	// Mark clean translation
	p.state.FormatLast = p.state.Format
	if fail {
		p.view.SetTextarea(DiffId, "")
		return
	}

	// Print diff
	diff := aNode.Diff(bNode, metadata...)
	var out string
	switch p.state.DiffFormat {
	case DiffFormatJd:
		out = diff.Render()
	case DiffFormatPatch:
		out, err = diff.RenderPatch()
		if err != nil {
			p.view.SetLabel(DiffErrorId, err.Error())
			return
		}
		if out == "[]" {
			out = ""
		}
	case DiffFormatMerge:
		out, err = diff.RenderMerge()
		if err != nil {
			p.view.SetLabel(DiffErrorId, err.Error())
			return
		}
	}
	p.view.SetTextarea(DiffId, out)
	p.view.SetLabel(DiffErrorId, "")
}

// printPatch handles patch mode computation and output
func (p *Presenter) printPatch() {
	metadata := p.getMetadata()
	var fail bool

	// Read a
	aNode, err := p.parseAndTranslate(AJsonId, p.state.FormatLast)
	if err != nil {
		p.view.SetLabel(AErrorId, err.Error())
		fail = true
	} else {
		p.view.SetLabel(AErrorId, "")
	}
	p.state.FormatLast = p.state.Format

	// Read diff
	diff, err := p.parseAndTranslateDiff()
	if err != nil {
		p.view.SetLabel(DiffErrorId, err.Error())
		fail = true
	} else {
		p.view.SetLabel(DiffErrorId, "")
	}

	if fail {
		p.view.SetTextarea(BJsonId, "")
		return
	}

	// Print patch
	bNode, err := aNode.Patch(diff)
	if err != nil {
		p.view.SetLabel(DiffErrorId, err.Error())
		p.view.SetTextarea(BJsonId, "")
		return
	}

	var out string
	if p.state.Format == FormatJSON {
		out = bNode.Json(metadata...)
	} else {
		out = bNode.Yaml(metadata...)
	}
	p.view.SetTextarea(BJsonId, out)
	p.view.SetLabel(DiffErrorId, "")
}
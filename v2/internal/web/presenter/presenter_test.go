package presenter

import (
	"testing"
)

func TestNew(t *testing.T) {
	view := newMockView()
	presenter, err := New(view)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if presenter == nil {
		t.Fatal("New() returned nil presenter")
	}
	
	// Verify initial state
	if presenter.state.Mode != ModeDiff {
		t.Errorf("Initial mode = %v, want %v", presenter.state.Mode, ModeDiff)
	}
	if presenter.state.Format != FormatJSON {
		t.Errorf("Initial format = %v, want %v", presenter.state.Format, FormatJSON)
	}
	if presenter.state.DiffFormat != DiffFormatJd {
		t.Errorf("Initial diff format = %v, want %v", presenter.state.DiffFormat, DiffFormatJd)
	}
	if presenter.state.Array != ArrayList {
		t.Errorf("Initial array type = %v, want %v", presenter.state.Array, ArrayList)
	}
}

func TestSetCommandLabel(t *testing.T) {
	tests := []struct {
		name       string
		mode       Mode
		format     Format
		diffFormat DiffFormat
		array      ArrayType
		want       string
	}{
		{
			name:       "basic diff json",
			mode:       ModeDiff,
			format:     FormatJSON,
			diffFormat: DiffFormatJd,
			array:      ArrayList,
			want:       "jd a.json b.json",
		},
		{
			name:       "diff yaml with set",
			mode:       ModeDiff,
			format:     FormatYAML,
			diffFormat: DiffFormatJd,
			array:      ArraySet,
			want:       "jd -yaml -set a.yaml b.yaml",
		},
		{
			name:       "patch mode with merge format",
			mode:       ModePatch,
			format:     FormatJSON,
			diffFormat: DiffFormatMerge,
			array:      ArrayList,
			want:       "jd -p -f merge diff a.json",
		},
		{
			name:       "patch yaml with patch format and multiset",
			mode:       ModePatch,
			format:     FormatYAML,
			diffFormat: DiffFormatPatch,
			array:      ArrayMset,
			want:       "jd -p -yaml -f patch -mset diff a.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newMockView()
			presenter, _ := New(view)
			
			presenter.state.Mode = tt.mode
			presenter.state.Format = tt.format
			presenter.state.DiffFormat = tt.diffFormat
			presenter.state.Array = tt.array
			
			presenter.setCommandLabel()
			
			mock := view.(*mockView)
			got := mock.labels[CommandId]
			if got != tt.want {
				t.Errorf("setCommandLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetInputLabels(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		wantA  string
		wantB  string
	}{
		{
			name:   "json format",
			format: FormatJSON,
			wantA:  "a.json",
			wantB:  "b.json",
		},
		{
			name:   "yaml format",
			format: FormatYAML,
			wantA:  "a.yaml",
			wantB:  "b.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newMockView()
			presenter, _ := New(view)
			presenter.state.Format = tt.format
			
			presenter.setInputLabels()
			
			mock := view.(*mockView)
			if got := mock.labels[ALabelId]; got != tt.wantA {
				t.Errorf("A label = %v, want %v", got, tt.wantA)
			}
			if got := mock.labels[BLabelId]; got != tt.wantB {
				t.Errorf("B label = %v, want %v", got, tt.wantB)
			}
		})
	}
}

func TestSetInputsEnabled(t *testing.T) {
	tests := []struct {
		name           string
		mode           Mode
		wantAStyle     string
		wantBReadonly  bool
		wantBStyle     string
		wantDiffReadonly bool
		wantDiffStyle  string
	}{
		{
			name:           "diff mode",
			mode:           ModeDiff,
			wantAStyle:     FocusStyle + ";" + HalfWidthStyle,
			wantBReadonly:  false,
			wantBStyle:     FocusStyle + ";" + HalfWidthStyle,
			wantDiffReadonly: true,
			wantDiffStyle:  UnfocusStyle + ";" + FullWidthStyle,
		},
		{
			name:           "patch mode",
			mode:           ModePatch,
			wantAStyle:     FocusStyle + ";" + HalfWidthStyle,
			wantBReadonly:  true,
			wantBStyle:     UnfocusStyle + ";" + HalfWidthStyle,
			wantDiffReadonly: false,
			wantDiffStyle:  FocusStyle + ";" + FullWidthStyle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newMockView()
			presenter, _ := New(view)
			presenter.state.Mode = tt.mode
			
			presenter.setInputsEnabled()
			
			mock := view.(*mockView)
			if got := mock.styles[AJsonId]; got != tt.wantAStyle {
				t.Errorf("A style = %v, want %v", got, tt.wantAStyle)
			}
			if got := mock.readonly[BJsonId]; got != tt.wantBReadonly {
				t.Errorf("B readonly = %v, want %v", got, tt.wantBReadonly)
			}
			if got := mock.styles[BJsonId]; got != tt.wantBStyle {
				t.Errorf("B style = %v, want %v", got, tt.wantBStyle)
			}
			if got := mock.readonly[DiffId]; got != tt.wantDiffReadonly {
				t.Errorf("Diff readonly = %v, want %v", got, tt.wantDiffReadonly)
			}
			if got := mock.styles[DiffId]; got != tt.wantDiffStyle {
				t.Errorf("Diff style = %v, want %v", got, tt.wantDiffStyle)
			}
		})
	}
}

func TestSetDerived(t *testing.T) {
	tests := []struct {
		name       string
		diffFormat DiffFormat
		wantArray  ArrayType
	}{
		{
			name:       "jd format preserves array setting",
			diffFormat: DiffFormatJd,
			wantArray:  ArraySet, // Should remain unchanged
		},
		{
			name:       "patch format forces list",
			diffFormat: DiffFormatPatch,
			wantArray:  ArrayList,
		},
		{
			name:       "merge format forces list",
			diffFormat: DiffFormatMerge,
			wantArray:  ArrayList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newMockView()
			presenter, _ := New(view)
			presenter.state.Array = ArraySet // Start with set
			presenter.state.DiffFormat = tt.diffFormat
			
			presenter.setDerived()
			
			if presenter.state.Array != tt.wantArray {
				t.Errorf("Array type = %v, want %v", presenter.state.Array, tt.wantArray)
			}
			
			// Check that checkboxes are updated correctly for patch/merge formats
			if tt.diffFormat == DiffFormatPatch || tt.diffFormat == DiffFormatMerge {
				mock := view.(*mockView)
				if !mock.checked[string(ArrayList)] {
					t.Error("ArrayList should be checked")
				}
				if mock.checked[string(ArraySet)] {
					t.Error("ArraySet should not be checked")
				}
				if mock.checked[string(ArrayMset)] {
					t.Error("ArrayMset should not be checked")
				}
			}
		})
	}
}

func TestParseAndTranslateJSON(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		format     Format
		formatLast Format
		wantError  bool
		wantOutput string // Expected output in textarea if format conversion happens
	}{
		{
			name:       "valid json, json format",
			input:      `{"foo": "bar"}`,
			format:     FormatJSON,
			formatLast: FormatJSON,
			wantError:  false,
		},
		{
			name:       "invalid json, json format",
			input:      `{"foo": bar}`,
			format:     FormatJSON,
			formatLast: FormatJSON,
			wantError:  true,
		},
		{
			name:       "valid yaml, yaml format",
			input:      `foo: bar`,
			format:     FormatYAML,
			formatLast: FormatYAML,
			wantError:  false,
		},
		{
			name:       "json to yaml conversion",
			input:      `{"foo": "bar"}`,
			format:     FormatYAML,
			formatLast: FormatJSON,
			wantError:  false,
			wantOutput: "foo: bar\n",
		},
		{
			name:       "yaml to json conversion",
			input:      `foo: bar`,
			format:     FormatJSON,
			formatLast: FormatYAML,
			wantError:  false,
			wantOutput: `{"foo":"bar"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newMockView()
			presenter, _ := New(view)
			presenter.state.Format = tt.format
			
			// Set up input value
			mock := view.(*mockView)
			mock.values[AJsonId] = tt.input
			
			_, err := presenter.parseAndTranslate(AJsonId, tt.formatLast)
			
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if tt.wantOutput != "" {
				if got := mock.values[AJsonId]; got != tt.wantOutput {
					t.Errorf("Output = %v, want %v", got, tt.wantOutput)
				}
			}
		})
	}
}

func TestPrintDiff(t *testing.T) {
	tests := []struct {
		name       string
		inputA     string
		inputB     string
		format     Format
		diffFormat DiffFormat
		array      ArrayType
		wantOutput string
		wantError  string
	}{
		{
			name:       "basic json diff",
			inputA:     `{"foo": "bar"}`,
			inputB:     `{"foo": "baz"}`,
			format:     FormatJSON,
			diffFormat: DiffFormatJd,
			array:      ArrayList,
			wantOutput: `@ ["foo"]
- "bar"
+ "baz"
`,
		},
		{
			name:       "empty diff",
			inputA:     `{"foo": "bar"}`,
			inputB:     `{"foo": "bar"}`,
			format:     FormatJSON,
			diffFormat: DiffFormatJd,
			array:      ArrayList,
			wantOutput: "",
		},
		{
			name:       "invalid json A",
			inputA:     `{"foo": "bar"`,
			inputB:     `{"foo": "baz"}`,
			format:     FormatJSON,
			diffFormat: DiffFormatJd,
			array:      ArrayList,
			wantError:  "A input error expected",
		},
		{
			name:       "invalid json B",
			inputA:     `{"foo": "bar"}`,
			inputB:     `{"foo": "baz"`,
			format:     FormatJSON,
			diffFormat: DiffFormatJd,
			array:      ArrayList,
			wantError:  "B input error expected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newMockView()
			presenter, _ := New(view)
			presenter.state.Format = tt.format
			presenter.state.DiffFormat = tt.diffFormat
			presenter.state.Array = tt.array
			presenter.state.Mode = ModeDiff
			
			// Set up input values
			mock := view.(*mockView)
			mock.values[AJsonId] = tt.inputA
			mock.values[BJsonId] = tt.inputB
			
			presenter.printDiff()
			
			if tt.wantError != "" {
				if mock.labels[AErrorId] == "" && mock.labels[BErrorId] == "" {
					t.Errorf("Expected error but got none. AError=%q, BError=%q, DiffError=%q", 
						mock.labels[AErrorId], mock.labels[BErrorId], mock.labels[DiffErrorId])
				}
				return
			}
			
			if got := mock.values[DiffId]; got != tt.wantOutput {
				t.Errorf("Diff output = %q, want %q", got, tt.wantOutput)
			}
			
			// Verify no errors
			if mock.labels[AErrorId] != "" {
				t.Errorf("Unexpected A error: %v", mock.labels[AErrorId])
			}
			if mock.labels[BErrorId] != "" {
				t.Errorf("Unexpected B error: %v", mock.labels[BErrorId])
			}
			if mock.labels[DiffErrorId] != "" {
				t.Errorf("Unexpected diff error: %v", mock.labels[DiffErrorId])
			}
		})
	}
}

func TestPrintPatch(t *testing.T) {
	tests := []struct {
		name       string
		inputA     string
		inputDiff  string
		format     Format
		diffFormat DiffFormat
		wantOutput string
		wantError  string
	}{
		{
			name:       "basic patch",
			inputA:     `{"foo": "bar"}`,
			inputDiff:  `@ ["foo"]` + "\n" + `- "bar"` + "\n" + `+ "baz"` + "\n",
			format:     FormatJSON,
			diffFormat: DiffFormatJd,
			wantOutput: `{"foo":"baz"}`,
		},
		{
			name:       "empty patch",
			inputA:     `{"foo": "bar"}`,
			inputDiff:  ``,
			format:     FormatJSON,
			diffFormat: DiffFormatJd,
			wantOutput: `{"foo":"bar"}`,
		},
		{
			name:       "invalid input A",
			inputA:     `{"foo": "bar"`,
			inputDiff:  ``,
			format:     FormatJSON,
			diffFormat: DiffFormatJd,
			wantError:  "A input error expected",
		},
		{
			name:       "invalid diff",
			inputA:     `{"foo": "bar"}`,
			inputDiff:  `invalid diff`,
			format:     FormatJSON,
			diffFormat: DiffFormatJd,
			wantError:  "diff error expected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newMockView()
			presenter, _ := New(view)
			presenter.state.Format = tt.format
			presenter.state.DiffFormat = tt.diffFormat
			presenter.state.Mode = ModePatch
			
			// Set up input values
			mock := view.(*mockView)
			mock.values[AJsonId] = tt.inputA
			mock.values[DiffId] = tt.inputDiff
			
			presenter.printPatch()
			
			if tt.wantError != "" {
				if mock.labels[AErrorId] == "" && mock.labels[DiffErrorId] == "" {
					t.Errorf("Expected error but got none. AError=%q, DiffError=%q", 
						mock.labels[AErrorId], mock.labels[DiffErrorId])
				}
				return
			}
			
			if got := mock.values[BJsonId]; got != tt.wantOutput {
				t.Errorf("Patch output = %q, want %q", got, tt.wantOutput)
			}
			
			// Verify no errors
			if mock.labels[AErrorId] != "" {
				t.Errorf("Unexpected A error: %v", mock.labels[AErrorId])
			}
			if mock.labels[DiffErrorId] != "" {
				t.Errorf("Unexpected diff error: %v", mock.labels[DiffErrorId])
			}
		})
	}
}

func TestUpdateState(t *testing.T) {
	view := newMockView()
	presenter, _ := New(view)
	
	// Update state
	presenter.UpdateState(ModePatch, FormatYAML, DiffFormatPatch, ArraySet)
	
	// Verify state was updated
	if presenter.state.Mode != ModePatch {
		t.Errorf("Mode = %v, want %v", presenter.state.Mode, ModePatch)
	}
	if presenter.state.Format != FormatYAML {
		t.Errorf("Format = %v, want %v", presenter.state.Format, FormatYAML)
	}
	if presenter.state.DiffFormat != DiffFormatPatch {
		t.Errorf("DiffFormat = %v, want %v", presenter.state.DiffFormat, DiffFormatPatch)
	}
	
	// Note: Array should be forced to ArrayList due to patch format
	if presenter.state.Array != ArrayList {
		t.Errorf("Array = %v, want %v (should be forced to list)", presenter.state.Array, ArrayList)
	}
}
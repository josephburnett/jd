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
	if presenter.state.OptionsJSON != "[]" {
		t.Errorf("Initial options JSON = %v, want %v", presenter.state.OptionsJSON, "[]")
	}
}

func TestSetCommandLabel(t *testing.T) {
	tests := []struct {
		name        string
		mode        Mode
		format      Format
		diffFormat  DiffFormat
		optionsJSON string
		want        string
	}{
		{
			name:        "basic diff json",
			mode:        ModeDiff,
			format:      FormatJSON,
			diffFormat:  DiffFormatJd,
			optionsJSON: "[]",
			want:        "jd a.json b.json",
		},
		{
			name:        "diff yaml with set option",
			mode:        ModeDiff,
			format:      FormatYAML,
			diffFormat:  DiffFormatJd,
			optionsJSON: `["SET"]`,
			want:        `jd -yaml -opts='["SET"]' a.yaml b.yaml`,
		},
		{
			name:        "patch mode with merge format",
			mode:        ModePatch,
			format:      FormatJSON,
			diffFormat:  DiffFormatMerge,
			optionsJSON: "[]",
			want:        "jd -p -f merge diff a.json",
		},
		{
			name:        "patch yaml with options",
			mode:        ModePatch,
			format:      FormatYAML,
			diffFormat:  DiffFormatPatch,
			optionsJSON: `[{"precision": 0.1}]`,
			want:        `jd -p -yaml -f patch -opts='[{"precision": 0.1}]' diff a.yaml`,
		},
		{
			name:        "path options example",
			mode:        ModeDiff,
			format:      FormatJSON,
			diffFormat:  DiffFormatJd,
			optionsJSON: `[{"@": ["users"], "^": ["SET"]}]`,
			want:        `jd -opts='[{"@": ["users"], "^": ["SET"]}]' a.json b.json`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newMockView()
			presenter, _ := New(view)

			presenter.state.Mode = tt.mode
			presenter.state.Format = tt.format
			presenter.state.DiffFormat = tt.diffFormat
			presenter.state.OptionsJSON = tt.optionsJSON
			presenter.validateAndParseOptions()

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
		name             string
		mode             Mode
		wantAStyle       string
		wantBReadonly    bool
		wantBStyle       string
		wantDiffReadonly bool
		wantDiffStyle    string
	}{
		{
			name:             "diff mode",
			mode:             ModeDiff,
			wantAStyle:       FocusStyle + ";" + HalfWidthStyle,
			wantBReadonly:    false,
			wantBStyle:       FocusStyle + ";" + HalfWidthStyle,
			wantDiffReadonly: true,
			wantDiffStyle:    UnfocusStyle + ";" + FullWidthStyle,
		},
		{
			name:             "patch mode",
			mode:             ModePatch,
			wantAStyle:       FocusStyle + ";" + HalfWidthStyle,
			wantBReadonly:    true,
			wantBStyle:       UnfocusStyle + ";" + HalfWidthStyle,
			wantDiffReadonly: false,
			wantDiffStyle:    FocusStyle + ";" + FullWidthStyle,
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

func TestValidateAndParseOptions(t *testing.T) {
	tests := []struct {
		name        string
		optionsJSON string
		wantError   bool
	}{
		{
			name:        "empty options",
			optionsJSON: "[]",
			wantError:   false,
		},
		{
			name:        "valid SET option",
			optionsJSON: `["SET"]`,
			wantError:   false,
		},
		{
			name:        "valid path option",
			optionsJSON: `[{"@": ["users"], "^": ["SET"]}]`,
			wantError:   false,
		},
		{
			name:        "invalid JSON",
			optionsJSON: `[invalid`,
			wantError:   true,
		},
		{
			name:        "invalid option",
			optionsJSON: `["INVALID_OPTION"]`,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newMockView()
			presenter, _ := New(view)
			presenter.state.OptionsJSON = tt.optionsJSON

			presenter.validateAndParseOptions()

			hasError := presenter.state.ValidationError != ""
			if hasError != tt.wantError {
				t.Errorf("validateAndParseOptions() error = %v, wantError %v. Error: %v",
					hasError, tt.wantError, presenter.state.ValidationError)
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
		name        string
		inputA      string
		inputB      string
		format      Format
		diffFormat  DiffFormat
		optionsJSON string
		wantOutput  string
		wantError   string
	}{
		{
			name:        "basic json diff",
			inputA:      `{"foo": "bar"}`,
			inputB:      `{"foo": "baz"}`,
			format:      FormatJSON,
			diffFormat:  DiffFormatJd,
			optionsJSON: "[]",
			wantOutput: `@ ["foo"]
- "bar"
+ "baz"
`,
		},
		{
			name:        "empty diff",
			inputA:      `{"foo": "bar"}`,
			inputB:      `{"foo": "bar"}`,
			format:      FormatJSON,
			diffFormat:  DiffFormatJd,
			optionsJSON: "[]",
			wantOutput:  "",
		},
		{
			name:        "invalid json A",
			inputA:      `{"foo": "bar"`,
			inputB:      `{"foo": "baz"}`,
			format:      FormatJSON,
			diffFormat:  DiffFormatJd,
			optionsJSON: "[]",
			wantError:   "A input error expected",
		},
		{
			name:        "invalid json B",
			inputA:      `{"foo": "bar"}`,
			inputB:      `{"foo": "baz"`,
			format:      FormatJSON,
			diffFormat:  DiffFormatJd,
			optionsJSON: "[]",
			wantError:   "B input error expected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newMockView()
			presenter, _ := New(view)
			presenter.state.Format = tt.format
			presenter.state.DiffFormat = tt.diffFormat
			presenter.state.OptionsJSON = tt.optionsJSON
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
	presenter.UpdateState(ModePatch, FormatYAML, DiffFormatPatch, `["SET"]`)

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
	if presenter.state.OptionsJSON != `["SET"]` {
		t.Errorf("OptionsJSON = %v, want %v", presenter.state.OptionsJSON, `["SET"]`)
	}
}

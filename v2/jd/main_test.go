package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const (
	jdFlags = "JD_FLAGS"
)

func TestMain(t *testing.T) {

	if flags := os.Getenv(jdFlags); flags != "" {
		args := strings.Split(flags, " ")
		os.Args = append([]string{os.Args[0]}, args...)
		main()
	}

	ref := func(s string) *string {
		return &s
	}

	testCases := []struct {
		name           string
		files          map[string]string
		args           []string
		exitCode       int
		out            *string
		outFile        string
		wantFileHeader string // if set, verify output starts with ^ {"file":"...<this>"}
	}{{
		name: "no diff",
		files: map[string]string{
			"a.json": `{"foo":"bar"}`,
			"b.json": `{"foo":"bar"}`,
		},
		args:     []string{"a.json", "b.json"},
		out:      ref(""),
		exitCode: 0,
	}, {
		name: "diff",
		files: map[string]string{
			"a.json": `{"foo":"bar"}`,
			"b.json": `{"foo":"baz"}`,
		},
		args: []string{"a.json", "b.json"},
		out: ref(s(
			`@ ["foo"]`,
			`- "bar"`,
			`+ "baz"`,
		)),
		exitCode:       1,
		wantFileHeader: "a.json",
	}, {
		name: "no diff in patch mode",
		files: map[string]string{
			"a.json": `{}`,
			"b.json": `{}`,
		},
		args:     []string{"-f", "patch", "a.json", "b.json"},
		out:      ref(`[]`),
		exitCode: 0,
	}, {
		name: "no diff in merge mode",
		files: map[string]string{
			"a.json": `{}`,
			"b.json": `{}`,
		},
		args:     []string{"-f", "merge", "a.json", "b.json"},
		out:      ref(`{}`),
		exitCode: 0,
	}, {
		name: "diff in patch mode",
		files: map[string]string{
			"a.json": `{"foo":"bar"}`,
			"b.json": `{"foo":"baz"}`,
		},
		args:     []string{"-f", "patch", "a.json", "b.json"},
		out:      ref(`[{"op":"test","path":"/foo","value":"bar"},{"op":"remove","path":"/foo","value":"bar"},{"op":"add","path":"/foo","value":"baz"}]`),
		exitCode: 1,
	}, {
		name: "diff in merge mode",
		files: map[string]string{
			"a.json": `{"foo":"bar"}`,
			"b.json": `{"foo":"baz"}`,
		},
		args:     []string{"-f", "merge", "a.json", "b.json"},
		out:      ref(`{"foo":"baz"}`),
		exitCode: 1,
	}, {
		name: "exit 0 on successful patch",
		files: map[string]string{
			"a.json": `{"foo":"bar"}`,
			"patch": `
@ ["foo"]
- "bar"
+ "baz"
`,
		},
		args:     []string{"-p", "patch", "a.json"},
		out:      ref(`{"foo":"baz"}`),
		exitCode: 0,
	}, {
		name: "exit 1 on unsuccessful patch",
		files: map[string]string{
			"a.json": `{"foo":"bar"}`,
			"patch": `
@ ["foo"]
- "zap"
+ "baz"
`,
		},
		args:     []string{"-p", "patch", "a.json"},
		exitCode: 2,
	}}

	testName := t.Name()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			temp, err := os.MkdirTemp("", testName)
			if err != nil {
				t.Fatalf("error creating temp directory: %v", err)
			}
			defer func() {
				os.RemoveAll(temp)
			}()
			files := map[string]struct{}{}
			for filename, content := range tc.files {
				files[filename] = struct{}{}
				err := os.WriteFile(fmt.Sprintf("%v%v%v", temp, os.PathSeparator, filename), []byte(content), 0666)
				if err != nil {
					t.Fatalf("error writing temp file: %v", err)
				}
			}
			args := make([]string, len(tc.args))
			for i, arg := range tc.args {
				if _, isFile := files[arg]; isFile {
					args[i] = fmt.Sprintf("%v%v%v", temp, os.PathSeparator, arg)
				} else {
					args[i] = arg
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run", testName)
			cmd.Env = append(os.Environ(), jdFlags+"="+strings.Join(args, " "))
			out, _ := cmd.CombinedOutput()
			outStr := string(out)
			if tc.wantFileHeader != "" {
				prefix := `^ {"file":"`
				if !strings.HasPrefix(outStr, prefix) {
					t.Errorf("wanted file header prefix %q. got %q", prefix, outStr)
				} else {
					headerEnd := strings.Index(outStr, "\n")
					header := outStr[:headerEnd]
					if !strings.HasSuffix(header, tc.wantFileHeader+`"}`) {
						t.Errorf("wanted file header ending with %q. got %q", tc.wantFileHeader, header)
					}
					outStr = outStr[headerEnd+1:]
				}
			}
			if tc.out != nil {
				if outStr != *tc.out {
					t.Errorf("wanted out %q. got %q", *tc.out, outStr)
				}
			}
			if exitCode := cmd.ProcessState.ExitCode(); exitCode != tc.exitCode {
				t.Errorf("wanted exit code %v. got %v", tc.exitCode, exitCode)
			}
		})
	}
}

func s(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}

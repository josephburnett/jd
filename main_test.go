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

	testCases := []struct {
		name     string
		files    map[string]string
		args     []string
		exitCode int
		out      string
		outFile  string
	}{{
		name: "no diff",
		files: map[string]string{
			"a.json": `{"foo":"bar"}`,
			"b.json": `{"foo":"bar"}`,
		},
		args:     []string{"a.json", "b.json"},
		out:      "",
		exitCode: 0,
	}, {
		name: "diff",
		files: map[string]string{
			"a.json": `{"foo":"bar"}`,
			"b.json": `{"foo":"baz"}`,
		},
		args: []string{"a.json", "b.json"},
		out: s(
			`@ ["foo"]`,
			`- "bar"`,
			`+ "baz"`,
		),
		exitCode: 1,
	}, {
		name: "no diff in patch mode",
		files: map[string]string{
			"a.json": `{}`,
			"b.json": `{}`,
		},
		args:     []string{"-f", "patch", "a.json", "b.json"},
		out:      `[]`,
		exitCode: 0,
	}, {
		name: "no diff in merge mode",
		files: map[string]string{
			"a.json": `{}`,
			"b.json": `{}`,
		},
		args:     []string{"-f", "merge", "a.json", "b.json"},
		out:      `{}`,
		exitCode: 0,
	}, {
		name: "diff in patch mode",
		files: map[string]string{
			"a.json": `{"foo":"bar"}`,
			"b.json": `{"foo":"baz"}`,
		},
		args:     []string{"-f", "patch", "a.json", "b.json"},
		out:      `[{"op":"test","path":"/foo","value":"bar"},{"op":"remove","path":"/foo","value":"bar"},{"op":"add","path":"/foo","value":"baz"}]`,
		exitCode: 1,
	}, {
		name: "diff in merge mode",
		files: map[string]string{
			"a.json": `{"foo":"bar"}`,
			"b.json": `{"foo":"baz"}`,
		},
		args:     []string{"-f", "merge", "a.json", "b.json"},
		out:      `{"foo":"baz"}`,
		exitCode: 1,
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
			if string(out) != tc.out {
				t.Errorf("wanted out %q. got %q", tc.out, string(out))
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

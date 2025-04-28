package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/josephburnett/jd/v2"
	"github.com/josephburnett/jd/v2/web/serve"
)

const version = "HEAD"

var (
	color         = flag.Bool("color", false, "Print color diff")
	format        = flag.String("f", "", "Diff format (jd, patch, merge)")
	gitDiffDriver = flag.Bool("git-diff-driver", false, "Use jd as a git diff driver.")
	mset          = flag.Bool("mset", false, "Arrays as multisets")
	opts          = flag.String("opts", "[]", "JSON array of options")
	output        = flag.String("o", "", "Output file")
	patch         = flag.Bool("p", false, "Patch mode")
	port          = flag.Int("port", 0, "Serve web UI on port")
	precision     = flag.Float64("precision", 0, "Maximum absolute difference for numbers to be equal")
	set           = flag.Bool("set", false, "Arrays as sets")
	setkeys       = flag.String("setkeys", "", "Keys to identify set objects")
	translate     = flag.String("t", "", "Translate mode")
	ver           = flag.Bool("version", false, "Print version and exit")
	yaml          = flag.Bool("yaml", false, "Read and write YAML")

	// This is here so that existing user commands that provide -v2 don't fail.
	_ = flag.Bool("v2", true, "Use the jd v2 library (deprecated, has no effect)")
)

func main() {
	if filepath.Base(os.Args[0]) == "jd-github-action" {
		fmt.Println("Running as GitHub Action...")
		runAsGitHubAction()
		return
	}
	flag.Parse()
	if *ver {
		fmt.Printf("jd version %v\n", version)
		return
	}
	if *port != 0 {
		if len(flag.Args()) > 0 {
			errorfAndExit("The web UI (-port) does not support arguments")
		}
		err := serveWeb(strconv.Itoa(*port))
		if err != nil {
			errorAndExit(err)
		}
		return
	}
	var (
		options []jd.Option
		err     error
	)
	options, err = parseOptions()
	if err != nil {
		errorAndExit(err)
	}
	if *gitDiffDriver {
		err := printGitDiffDriver(options)
		if err != nil {
			errorAndExit(err)
		}
		os.Exit(0)
		return
	}
	mode := diffMode
	if *patch {
		mode = patchMode
	}
	if *translate != "" {
		mode = translateMode
	}
	if *patch && *translate != "" {
		errorfAndExit("Patch and translate modes cannot be used together.")
	}
	var a, b string
	switch mode {
	case diffMode, patchMode:
		switch len(flag.Args()) {
		case 1:
			a = readFile(flag.Arg(0))
			b = readStdin()
		case 2:
			a = readFile(flag.Arg(0))
			b = readFile(flag.Arg(1))
		default:
			printUsageAndExit()
		}
	case translateMode:
		switch len(flag.Args()) {
		case 0:
			a = readStdin()
		case 1:
			a = readFile(flag.Arg(0))
		default:
			printUsageAndExit()
		}
	}
	switch mode {
	case diffMode:
		printDiff(a, b, options)
	case patchMode:
		printPatch(a, b, options)
	case translateMode:
		printTranslation(a)
	}
}

type mode string

const (
	diffMode      mode = "diff"
	patchMode     mode = "patch"
	translateMode mode = "trans"
)

func serveWeb(port string) error {
	if serve.Handle == nil {
		return fmt.Errorf("the web UI wasn't include in this build: use `make build` to include it")
	}
	http.HandleFunc("/", serve.Handle)
	log.Printf("Listening on http://localhost:%v...", port)
	return http.ListenAndServe(":"+port, nil)
}

func parseOptions() ([]jd.Option, error) {
	if *precision != 0.0 && (*set || *mset) {
		return nil, fmt.Errorf("-precision cannot be used with -set or -mset because they use hashcodes")
	}
	options, err := jd.UnmarshalOptions([]byte(*opts))
	if err != nil {
		return nil, err
	}
	if *set {
		options = append(options, jd.SET)
	}
	if *mset {
		options = append(options, jd.MULTISET)
	}
	if *setkeys != "" {
		keys := make([]string, 0)
		ks := strings.Split(*setkeys, ",")
		for _, k := range ks {
			trimmed := strings.TrimSpace(k)
			if trimmed == "" {
				return nil, fmt.Errorf("invalid set key: %v", k)
			}
			keys = append(keys, trimmed)
		}
		options = append(options, jd.SetKeys(keys...))
	}
	if *format == "merge" {
		options = append(options, jd.MERGE)
	}
	options = append(options, jd.Precision(*precision))
	return options, nil
}

func printUsageAndExit() {
	for _, line := range []string{
		``,
		`Usage: jd [OPTION]... FILE1 [FILE2]`,
		`Diff and patch JSON files.`,
		``,
		`Prints the diff of FILE1 and FILE2 to STDOUT.`,
		`When FILE2 is omitted the second input is read from STDIN.`,
		`When patching (-p) FILE1 is a diff.`,
		``,
		`Options:`,
		`  -color       Print color diff.`,
		`  -p           Apply patch FILE1 to FILE2 or STDIN.`,
		`  -o=FILE3     Write to FILE3 instead of STDOUT.`,
		`  -opts='[]'   JSON array of options.`,
		`  -set         Treat arrays as sets.`,
		`  -mset        Treat arrays as multisets (bags).`,
		`  -setkeys     Keys to identify set objects`,
		`  -yaml        Read and write YAML instead of JSON.`,
		`  -port=N      Serve web UI on port N`,
		`  -precision=N Maximum absolute difference for numbers to be equal.`,
		`               Example: -precision=0.00001`,
		`  -f=FORMAT    Read and write diff in FORMAT "jd" (default), "patch" (RFC 6902) or`,
		`               "merge" (RFC 7386)`,
		`  -t=FORMATS   Translate FILE1 between FORMATS. Supported formats are "jd",`,
		`               "patch" (RFC 6902), "merge" (RFC 7386), "json" and "yaml".`,
		`               FORMATS are provided as a pair separated by "2". E.g.`,
		`               "yaml2json" or "jd2patch".`,
		``,
		`Examples:`,
		`  jd a.json b.json`,
		`  cat b.json | jd a.json`,
		`  jd -o patch a.json b.json; jd patch a.json`,
		`  jd -set a.json b.json`,
		`  jd -f patch a.json b.json`,
		`  jd -f merge a.json b.json`,
		``,
		`Version: ` + version,
		``,
	} {
		fmt.Println(line)
	}
	os.Exit(2)
}

func printDiff(a, b string, options []jd.Option) {
	str, haveDiff, err := diff(a, b, options)
	if err != nil {
		errorAndExit(err)
	}
	if *output == "" {
		fmt.Print(str)
	} else {
		ioutil.WriteFile(*output, []byte(str), 0644)
	}
	if haveDiff {
		os.Exit(1)
	}
	os.Exit(0)
}

func printGitDiffDriver(options []jd.Option) error {
	if len(flag.Args()) != 7 {
		return fmt.Errorf("Git diff driver expects exactly 7 arguments.")
	}
	a := readFile(flag.Arg(1))
	b := readFile(flag.Arg(4))
	str, _, err := diff(a, b, options)
	if err != nil {
		return err
	}
	fmt.Print(str)
	os.Exit(0)
	return nil
}

func diff(a, b string, options []jd.Option) (string, bool, error) {
	var aNode, bNode jd.JsonNode
	var err error
	if *yaml {
		aNode, err = jd.ReadYamlString(a)
	} else {
		aNode, err = jd.ReadJsonString(a)
	}
	if err != nil {
		return "", false, err
	}
	if *yaml {
		bNode, err = jd.ReadYamlString(b)
	} else {
		bNode, err = jd.ReadJsonString(b)
	}
	if err != nil {
		return "", false, err
	}
	diff := aNode.Diff(bNode, options...)
	var renderOptions []jd.Option
	if *color {
		renderOptions = append(renderOptions, jd.COLOR)
	}
	var (
		str      string
		haveDiff bool
	)
	switch *format {
	case "", "jd":
		str = diff.Render(renderOptions...)
		if str != "" {
			haveDiff = true
		}
	case "patch":
		str, err = diff.RenderPatch()
		if err != nil {
			return "", false, err
		}
		if str != "[]" {
			haveDiff = true
		}
	case "merge":
		str, err = diff.RenderMerge()
		if err != nil {
			return "", false, err
		}
		if str != "{}" {
			haveDiff = true
		}
	default:
		return "", false, fmt.Errorf("Invalid format: %q", *format)
	}
	return str, haveDiff, nil
}

func printPatch(p, a string, options []jd.Option) {
	var diff jd.Diff
	var err error
	switch *format {
	case "", "jd":
		diff, err = jd.ReadDiffString(p)
	case "patch":
		diff, err = jd.ReadPatchString(p)
	case "merge":
		diff, err = jd.ReadMergeString(p)
	default:
		errorfAndExit("Invalid format: %q", *format)
	}
	if err != nil {
		errorAndExit(err)
	}
	var aNode jd.JsonNode
	if *yaml {
		aNode, err = jd.ReadYamlString(a)
	} else {
		aNode, err = jd.ReadJsonString(a)
	}
	if err != nil {
		errorAndExit(err)
	}
	bNode, err := aNode.Patch(diff)
	if err != nil {
		errorAndExit(err)
	}
	var out string
	if *yaml {
		out = bNode.Yaml(options...)
	} else {
		out = bNode.Json(options...)
	}
	if *output == "" {
		fmt.Print(out)
	} else {
		os.WriteFile(*output, []byte(out), 0644)
	}
	os.Exit(0)
}

func printTranslation(a string) {
	var out string
	switch *translate {
	case "jd2patch":
		diff, err := jd.ReadDiffString(a)
		if err != nil {
			errorAndExit(err)
		}
		out, err = diff.RenderPatch()
		if err != nil {
			errorAndExit(err)
		}
	case "patch2jd":
		patch, err := jd.ReadPatchString(a)
		if err != nil {
			errorAndExit(err)
		}
		out = patch.Render()
	case "jd2merge":
		diff, err := jd.ReadDiffString(a)
		if err != nil {
			errorAndExit(err)
		}
		out, err = diff.RenderMerge()
		if err != nil {
			errorAndExit(err)
		}
	case "merge2jd":
		patch, err := jd.ReadMergeString(a)
		if err != nil {
			errorAndExit(err)
		}
		out = patch.Render()
	case "json2yaml":
		node, err := jd.ReadJsonString(a)
		if err != nil {
			errorAndExit(err)
		}
		out = node.Yaml()
	case "yaml2json":
		node, err := jd.ReadYamlString(a)
		if err != nil {
			errorAndExit(err)
		}
		out = node.Json()
	default:
		errorfAndExit("unsupported translation: %q", *translate)
	}
	if *output == "" {
		fmt.Print(out)
	} else {
		ioutil.WriteFile(*output, []byte(out), 0644)
	}
	os.Exit(0)
}

func errorAndExit(err error) {
	errorfAndExit("%v", err.Error())
}

func errorfAndExit(msg string, args ...interface{}) {
	log.Printf(msg, args...)
	os.Exit(2)
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Print(err.Error())
		os.Exit(2)
	}
	return string(bytes)
}

func readStdin() string {
	r := bufio.NewReader(os.Stdin)
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		log.Print(err.Error())
		os.Exit(2)
	}
	return string(bytes)
}

func runAsGitHubAction() {
	gitHubOutput := os.Getenv("GITHUB_OUTPUT")
	if gitHubOutput == "" {
		errorfAndExit("GITHUB_OUTPUT no set. Are you running in GitHub CI?")
	}
	file, err := os.OpenFile(gitHubOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errorAndExit(err)
	}
	defer file.Close()
	if len(os.Args) < 2 {
		errorfAndExit("Running GitHub Action requires args")
	}
	// Actions do not accept list inputs so args must be string split.
	args := strings.Fields(os.Args[1])
	cmd := exec.Command("/jd", args...)
	out, _ := cmd.CombinedOutput()
	delimiter := strconv.Itoa(rand.Int())
	file.WriteString("output<<" + delimiter + "\n")
	file.WriteString(string(out))
	file.WriteString(delimiter + "\n")
	file.WriteString("exit_code=" + strconv.Itoa(cmd.ProcessState.ExitCode()) + "\n")
	os.Exit(0)
}

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
	"strconv"
	"strings"

	jd "github.com/josephburnett/jd/lib"
	v2 "github.com/josephburnett/jd/v2"
	"github.com/josephburnett/jd/web/serve"
)

const version = "HEAD"

var (
	color         = flag.Bool("color", false, "Print color diff")
	format        = flag.String("f", "", "Diff format (jd, patch, merge)")
	gitDiffDriver = flag.Bool("git-diff-driver", false, "Use jd as a git diff driver.")
	mset          = flag.Bool("mset", false, "Arrays as multisets")
	output        = flag.String("o", "", "Output file")
	patch         = flag.Bool("p", false, "Patch mode")
	port          = flag.Int("port", 0, "Serve web UI on port")
	precision     = flag.Float64("precision", 0, "Maximum absolute difference for numbers to be equal")
	set           = flag.Bool("set", false, "Arrays as sets")
	setkeys       = flag.String("setkeys", "", "Keys to identify set objects")
	translate     = flag.String("t", "", "Translate mode")
	ver           = flag.Bool("version", false, "Print version and exit")
	libv2         = flag.Bool("v2", true, "Use the jd v2 library")
	yaml          = flag.Bool("yaml", false, "Read and write YAML")
)

func main() {
	if os.Args[0] == "jd-github-action" {
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
			errorAndExit("The web UI (-port) does not support arguments")
		}
		err := serveWeb(strconv.Itoa(*port))
		if err != nil {
			errorAndExit(err.Error())
		}
		return
	}
	var (
		metadata []jd.Metadata
		options  []v2.Option
		err      error
	)
	if *libv2 {
		options, err = parseMetadataV2()
	} else {
		metadata, err = parseMetadata()
	}
	if err != nil {
		errorAndExit(err.Error())
	}
	if *gitDiffDriver {
		err := printGitDiffDriver(options)
		if err != nil {
			errorAndExit(err.Error())
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
		errorAndExit("Patch and translate modes cannot be used together.")
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
		if *libv2 {
			printDiffV2(a, b, options)
		} else {
			printDiff(a, b, metadata)
		}
	case patchMode:
		if *libv2 {
			printPatchV2(a, b, options)
		} else {
			printPatch(a, b, metadata)
		}
	case translateMode:
		if *libv2 {
			printTranslationV2(a)
		} else {
			printTranslation(a)
		}
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

func parseMetadata() ([]jd.Metadata, error) {
	if *precision != 0.0 && (*set || *mset) {
		return nil, fmt.Errorf("-precision cannot be used with -set or -mset because they use hashcodes")
	}
	metadata := make([]jd.Metadata, 0)
	if *set {
		metadata = append(metadata, jd.SET)
	}
	if *mset {
		metadata = append(metadata, jd.MULTISET)
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
		metadata = append(metadata, jd.Setkeys(keys...))
	}
	if *format == "merge" {
		metadata = append(metadata, jd.MERGE)
	}
	metadata = append(metadata, jd.SetPrecision(*precision))
	return metadata, nil
}

func parseMetadataV2() ([]v2.Option, error) {
	if *precision != 0.0 && (*set || *mset) {
		return nil, fmt.Errorf("-precision cannot be used with -set or -mset because they use hashcodes")
	}
	options := make([]v2.Option, 0)
	if *set {
		options = append(options, v2.SET)
	}
	if *mset {
		options = append(options, v2.MULTISET)
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
		options = append(options, v2.SetKeys(keys...))
	}
	if *format == "merge" {
		options = append(options, v2.MERGE)
	}
	options = append(options, v2.Precision(*precision))
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
		`  -v2          Use the JD v2 library and format (defaults true).`,
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

func printDiff(a, b string, metadata []jd.Metadata) {
	str, haveDiff, err := diff(a, b, metadata)
	if err != nil {
		errorAndExit(err.Error())
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

func printDiffV2(a, b string, options []v2.Option) {
	str, haveDiff, err := diffV2(a, b, options)
	if err != nil {
		errorAndExit(err.Error())
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

func printGitDiffDriver(options []v2.Option) error {
	if len(flag.Args()) != 7 {
		return fmt.Errorf("Git diff driver expects exactly 7 arguments.")
	}
	a := readFile(flag.Arg(1))
	b := readFile(flag.Arg(4))
	str, _, err := diffV2(a, b, options)
	if err != nil {
		return err
	}
	fmt.Print(str)
	os.Exit(0)
	return nil
}

func diff(a, b string, metadata []jd.Metadata) (string, bool, error) {
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
	diff := aNode.Diff(bNode, metadata...)
	var renderOptions []jd.RenderOption
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

func diffV2(a, b string, options []v2.Option) (string, bool, error) {
	var aNode, bNode v2.JsonNode
	var err error
	if *yaml {
		aNode, err = v2.ReadYamlString(a)
	} else {
		aNode, err = v2.ReadJsonString(a)
	}
	if err != nil {
		return "", false, err
	}
	if *yaml {
		bNode, err = v2.ReadYamlString(b)
	} else {
		bNode, err = v2.ReadJsonString(b)
	}
	if err != nil {
		return "", false, err
	}
	diff := aNode.Diff(bNode, options...)
	var renderOptions []v2.Option
	if *color {
		renderOptions = append(renderOptions, v2.COLOR)
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

func printPatch(p, a string, metadata []jd.Metadata) {
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
		errorAndExit(fmt.Sprintf("Invalid format: %q", *format))
	}
	if err != nil {
		errorAndExit(err.Error())
	}
	var aNode jd.JsonNode
	if *yaml {
		aNode, err = jd.ReadYamlString(a)
	} else {
		aNode, err = jd.ReadJsonString(a)
	}
	if err != nil {
		errorAndExit(err.Error())
	}
	bNode, err := aNode.Patch(diff)
	if err != nil {
		errorAndExit(err.Error())
	}
	var out string
	if *yaml {
		out = bNode.Yaml(metadata...)
	} else {
		out = bNode.Json(metadata...)
	}
	if *output == "" {
		fmt.Print(out)
	} else {
		ioutil.WriteFile(*output, []byte(out), 0644)
	}
	os.Exit(0)
}

func printPatchV2(p, a string, options []v2.Option) {
	var diff v2.Diff
	var err error
	switch *format {
	case "", "jd":
		diff, err = v2.ReadDiffString(p)
	case "patch":
		diff, err = v2.ReadPatchString(p)
	case "merge":
		diff, err = v2.ReadMergeString(p)
	default:
		errorAndExit(fmt.Sprintf("Invalid format: %q", *format))
	}
	if err != nil {
		errorAndExit(err.Error())
	}
	var aNode v2.JsonNode
	if *yaml {
		aNode, err = v2.ReadYamlString(a)
	} else {
		aNode, err = v2.ReadJsonString(a)
	}
	if err != nil {
		errorAndExit(err.Error())
	}
	bNode, err := aNode.Patch(diff)
	if err != nil {
		errorAndExit(err.Error())
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
		ioutil.WriteFile(*output, []byte(out), 0644)
	}
	os.Exit(0)
}

func printTranslation(a string) {
	var out string
	switch *translate {
	case "jd2patch":
		diff, err := jd.ReadDiffString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out, err = diff.RenderPatch()
		if err != nil {
			errorAndExit(err.Error())
		}
	case "patch2jd":
		patch, err := jd.ReadPatchString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out = patch.Render()
	case "jd2merge":
		diff, err := jd.ReadDiffString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out, err = diff.RenderMerge()
		if err != nil {
			errorAndExit(err.Error())
		}
	case "merge2jd":
		patch, err := jd.ReadMergeString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out = patch.Render()
	case "json2yaml":
		node, err := jd.ReadJsonString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out = node.Yaml()
	case "yaml2json":
		node, err := jd.ReadYamlString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out = node.Json()
	default:
		errorAndExit("unsupported translation: %q", *translate)
	}
	if *output == "" {
		fmt.Print(out)
	} else {
		ioutil.WriteFile(*output, []byte(out), 0644)
	}
	os.Exit(0)
}

func printTranslationV2(a string) {
	var out string
	switch *translate {
	case "jd2patch":
		diff, err := v2.ReadDiffString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out, err = diff.RenderPatch()
		if err != nil {
			errorAndExit(err.Error())
		}
	case "patch2jd":
		patch, err := v2.ReadPatchString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out = patch.Render()
	case "jd2merge":
		diff, err := v2.ReadDiffString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out, err = diff.RenderMerge()
		if err != nil {
			errorAndExit(err.Error())
		}
	case "merge2jd":
		patch, err := v2.ReadMergeString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out = patch.Render()
	case "json2yaml":
		node, err := v2.ReadJsonString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out = node.Yaml()
	case "yaml2json":
		node, err := v2.ReadYamlString(a)
		if err != nil {
			errorAndExit(err.Error())
		}
		out = node.Json()
	default:
		errorAndExit("unsupported translation: %q", *translate)
	}
	if *output == "" {
		fmt.Print(out)
	} else {
		ioutil.WriteFile(*output, []byte(out), 0644)
	}
	os.Exit(0)
}

func errorAndExit(msg string, args ...interface{}) {
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
		errorAndExit("GITHUB_OUTPUT no set. Are you running in GitHub CI?")
	}
	file, err := os.OpenFile(gitHubOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errorAndExit(err.Error())
	}
	defer file.Close()
	cmd := exec.Command("jd", os.Args[2:]...)
	out, err := cmd.CombinedOutput()
	delimiter := strconv.Itoa(rand.Int())
	file.WriteString("output<<" + delimiter + "\n")
	if err != nil {
		file.WriteString(err.Error() + "\n")
	}
	file.WriteString(string(out))
	file.WriteString(delimiter)
	os.Exit(cmd.ProcessState.ExitCode())
}

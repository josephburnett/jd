package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	jd "github.com/josephburnett/jd/lib"
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
	precision     = flag.Float64("precision", 0, "Precision for numbers")
	set           = flag.Bool("set", false, "Arrays as sets")
	setkeys       = flag.String("setkeys", "", "Keys to identify set objects")
	translate     = flag.String("t", "", "Translate mode")
	ver           = flag.Bool("version", false, "Print version and exit")
	yaml          = flag.Bool("yaml", false, "Read and write YAML")
)

func main() {
	flag.Parse()
	if *ver {
		fmt.Printf("jd version %v\n", version)
		return
	}
	if *port != 0 {
		err := serveWeb(strconv.Itoa(*port))
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
	metadata, err := parseMetadata()
	if err != nil {
		errorAndExit(err.Error())
	}
	if *gitDiffDriver {
		err := printGitDiffDriver(metadata)
		if err != nil {
			panic(err)
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
		printDiff(a, b, metadata)
	case patchMode:
		printPatch(a, b, metadata)
	case translateMode:
		printTranslation(a, metadata)
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
		`  -precision=N Precision for numbers. Positive number for decimal places or`,
		`               negative for significant figures.`,
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

func printDiff(a, b string, metadata []jd.Metadata) {
	str, err := diff(a, b, metadata)
	if err != nil {
		errorAndExit(err.Error())
	}
	if *output == "" {
		if str == "" {
			os.Exit(0)
		}
		fmt.Print(str)
		os.Exit(1)
	} else {
		if str == "" {
			os.Exit(0)
		}
		ioutil.WriteFile(*output, []byte(str), 0644)
		os.Exit(1)
	}
}

func printGitDiffDriver(metadata []jd.Metadata) error {
	if len(flag.Args()) != 7 {
		return fmt.Errorf("Git diff driver expects exactly 7 arguments.")
	}
	a := readFile(flag.Arg(1))
	b := readFile(flag.Arg(4))
	str, err := diff(a, b, metadata)
	if err != nil {
		return err
	}
	fmt.Print(str)
	os.Exit(0)
	return nil
}

func diff(a, b string, metadata []jd.Metadata) (string, error) {
	var aNode, bNode jd.JsonNode
	var err error
	if *yaml {
		aNode, err = jd.ReadYamlString(a)
	} else {
		aNode, err = jd.ReadJsonString(a)
	}
	if err != nil {
		return "", err
	}
	if *yaml {
		bNode, err = jd.ReadYamlString(b)
	} else {
		bNode, err = jd.ReadJsonString(b)
	}
	if err != nil {
		return "", err
	}
	diff := aNode.Diff(bNode, metadata...)
	var renderOptions []jd.RenderOption
	if *color {
		renderOptions = append(renderOptions, jd.COLOR)
	}
	var str string
	switch *format {
	case "", "jd":
		str = diff.Render(renderOptions...)
	case "patch":
		str, err = diff.RenderPatch()
		if err != nil {
			return "", err
		}
	case "merge":
		str, err = diff.RenderMerge()
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("Invalid format: %q", *format)
	}
	return str, nil
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
		if out == "" {
			os.Exit(0)
		}
		fmt.Print(out)
		os.Exit(1)
	} else {
		if out == "" {
			os.Exit(0)
		}
		ioutil.WriteFile(*output, []byte(out), 0644)
		os.Exit(1)
	}
}

func printTranslation(a string, metadata []jd.Metadata) {
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

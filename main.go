package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	jd "github.com/josephburnett/jd/lib"
)

var patch = flag.Bool("p", false, "Patch mode")
var output = flag.String("o", "", "Output file")
var set = flag.Bool("set", false, "Arrays as sets")
var mset = flag.Bool("mset", false, "Arrays as multisets")
var setkeys = flag.String("setkeys", "", "Keys to identify set objects")

func main() {
	flag.Parse()
	options, err := parseOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}
	var a, b string
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
	if *patch {
		patchJson(a, b, options)
	} else {
		diffJson(a, b, options)
	}
}

func parseOptions() ([]jd.Option, error) {
	options := make([]jd.Option, 0)
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
				return nil, fmt.Errorf("Invalid set key: %v", k)
			}
			keys = append(keys, trimmed)
		}
		options = append(options, jd.SetkeysOption(keys))
	}
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
		`  -p        Apply patch FILE1 to FILE2 or STDIN.`,
		`  -o=FILE3  Write to FILE3 instead of STDOUT.`,
		`  -set      Treat arrays as sets.`,
		`  -mset     Treat arrays as multisets (bags).`,
		``,
		`Examples:`,
		`  jd a.json b.json`,
		`  cat b.json | jd a.json`,
		`  jd -o patch a.json b.json; jd patch a.json`,
		`  jd -set a.json b.json`,
		``,
	} {
		fmt.Println(line)
	}
	os.Exit(1)
}

func diffJson(a, b string, options []jd.Option) {
	// TODO: remove options from reading
	aNode, err := jd.ReadJsonString(a, options...)
	if err != nil {
		log.Fatalf(err.Error())
	}
	// TODO: remove options from reading
	bNode, err := jd.ReadJsonString(b, options...)
	if err != nil {
		log.Fatalf(err.Error())
	}
	diff := aNode.Diff(bNode, options...)
	if *output == "" {
		fmt.Print(diff.Render())
	} else {
		ioutil.WriteFile(*output, []byte(diff.Render()), 0644)
	}
}

func patchJson(p, a string, options []jd.Option) {
	diff, err := jd.ReadDiffString(p, options...)
	if err != nil {
		log.Fatalf(err.Error())
	}
	// TODO: remove options from reading
	aNode, err := jd.ReadJsonString(a, options...)
	if err != nil {
		log.Fatalf(err.Error())
	}
	bNode, err := aNode.Patch(diff, options...)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if *output == "" {
		fmt.Print(bNode.Json())
	} else {
		ioutil.WriteFile(*output, []byte(bNode.Json()), 0644)
	}
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}

func readStdin() string {
	r := bufio.NewReader(os.Stdin)
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}

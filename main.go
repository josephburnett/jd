package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	jd "github.com/josephburnett/jd/lib"
)

var patch = flag.Bool("p", false, "Patch mode")

func main() {
	flag.Parse()
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
		patchJson(a, b)
	} else {
		diffJson(a, b)
	}
}

func printUsageAndExit() {
	for _, line := range []string{
		`Usage: jd <json_file_a> <json_file_b>`,
		`       jd -p <patch_file> <json_file_a>`,
		`       jd <json_file_a>`,
		`       jd -p <patch_file>`,
	} {
		fmt.Print(line)
	}
	os.Exit(1)
}

func diffJson(a, b string) {
	aNode, err := jd.ReadJsonString(a)
	if err != nil {
		log.Fatalf(err.Error())
	}
	bNode, err := jd.ReadJsonString(b)
	if err != nil {
		log.Fatalf(err.Error())
	}
	diff := aNode.Diff(bNode)
	fmt.Print(diff.Render())
}

func patchJson(p, a string) {
	diff, err := jd.ReadDiffString(p)
	if err != nil {
		log.Fatalf(err.Error())
	}
	aNode, err := jd.ReadJsonString(a)
	if err != nil {
		log.Fatalf(err.Error())
	}
	bNode, err := aNode.Patch(diff)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Print(bNode.Json())
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

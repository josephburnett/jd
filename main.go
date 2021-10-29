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

var mset = flag.Bool("mset", false, "Arrays as multisets")
var output = flag.String("o", "", "Output file")
var patch = flag.Bool("p", false, "Patch mode")
var port = flag.Int("port", 0, "Serve web UI on port")
var set = flag.Bool("set", false, "Arrays as sets")
var setkeys = flag.String("setkeys", "", "Keys to identify set objects")
var ver = flag.Bool("version", false, "Print version and exit")
var yaml = flag.Bool("yaml", false, "Read and write YAML")

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
		printPatch(a, b, metadata)
	} else {
		printDiff(a, b, metadata)
	}
}

func serveWeb(port string) error {
	if serve.Handle == nil {
		return fmt.Errorf("The web UI wasn't include in this build. Use `make release` to include it.")
	}
	http.HandleFunc("/", serve.Handle)
	log.Printf("Listening on :%v...", port)
	return http.ListenAndServe(":"+port, nil)
}

func parseMetadata() ([]jd.Metadata, error) {
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
				return nil, fmt.Errorf("Invalid set key: %v", k)
			}
			keys = append(keys, trimmed)
		}
		metadata = append(metadata, jd.Setkeys(keys...))
	}
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
		`  -p        Apply patch FILE1 to FILE2 or STDIN.`,
		`  -o=FILE3  Write to FILE3 instead of STDOUT.`,
		`  -set      Treat arrays as sets.`,
		`  -mset     Treat arrays as multisets (bags).`,
		`  -setkeys  Keys to identify set objects`,
		`  -yaml     Read and write YAML instead of JSON.`,
		`  -port=N   Serve web UI on port N`,
		``,
		`Examples:`,
		`  jd a.json b.json`,
		`  cat b.json | jd a.json`,
		`  jd -o patch a.json b.json; jd patch a.json`,
		`  jd -set a.json b.json`,
		``,
		`Version: ` + version,
		``,
	} {
		fmt.Println(line)
	}
	os.Exit(1)
}

func printDiff(a, b string, metadata []jd.Metadata) {
	var aNode, bNode jd.JsonNode
	var err error
	if *yaml {
		aNode, err = jd.ReadYamlString(a)
	} else {
		aNode, err = jd.ReadJsonString(a)
	}
	if err != nil {
		log.Fatalf(err.Error())
	}
	if *yaml {
		bNode, err = jd.ReadYamlString(b)
	} else {
		bNode, err = jd.ReadJsonString(b)
	}
	if err != nil {
		log.Fatalf(err.Error())
	}
	diff := aNode.Diff(bNode, metadata...)
	if *output == "" {
		str := diff.Render()
		if  str == "" {
			os.Exit(0)
		}
		fmt.Print(str)
		os.Exit(1)
	} else {
		str := diff.Render()
		if str == "" {
			os.Exit(0)
		}
		ioutil.WriteFile(*output, []byte(str), 0644)
		os.Exit(1)
	}
}

func printPatch(p, a string, metadata []jd.Metadata) {
	diff, err := jd.ReadDiffString(p)
	if err != nil {
		log.Fatalf(err.Error())
	}
	var aNode jd.JsonNode
	if *yaml {
		aNode, err = jd.ReadYamlString(a)
	} else {
		aNode, err = jd.ReadJsonString(a)
	}
	if err != nil {
		log.Fatalf(err.Error())
	}
	bNode, err := aNode.Patch(diff)
	if err != nil {
		log.Fatalf(err.Error())
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

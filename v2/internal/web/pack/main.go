package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

var files = []string{
	"wasm_exec.js",
	"index.html",
	"jd.wasm",
}

func main() {
	pack := "// +build include_web\n"
	pack += "\n"
	pack += "package serve\n"
	pack += "\n"
	pack += "var base64EncodedFiles = map[string]string{\n"
	for _, f := range files {
		b, err := os.ReadFile(fmt.Sprintf("internal/web/assets/%v", f))
		if err != nil {
			panic(err)
		}
		s := base64.StdEncoding.EncodeToString(b)
		pack += fmt.Sprintf("\t%q: %q,\n", f, s)
	}
	pack += "}"
	err := os.WriteFile("internal/web/serve/files.go", []byte(pack), 0644)
	if err != nil {
		panic(err)
	}
}

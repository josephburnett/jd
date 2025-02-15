package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		fmt.Printf("GITHUB_OUTPUT not set. Are you running in GitHub CI?")
		os.Exit(2)
	}
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer file.Close()
	cmd := exec.Command("jd", os.Args[2:]...)
	data, _ := cmd.CombinedOutput()
	exitCode := cmd.ProcessState.ExitCode()
	delimiter := strconv.Itoa(rand.Int())
	file.WriteString("output<<" + delimiter + "\n")
	file.WriteString(string(data))
	file.WriteString("delimiter\n")
	os.Exit(exitCode)
}

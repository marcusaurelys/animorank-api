package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {

	code := os.Getenv("CODE")
	stdin := os.Getenv("STDIN")
	name := "program.c"
	err := os.WriteFile(name, []byte(code), 0644)
	if err != nil {
		fmt.Printf("Error writing code to file. %v", err)
		os.Exit(1)
	}

	compiled := "program.out"
	compile := exec.Command("gcc", name, "-o", compiled, "-Wall")
	compileOutput, err := compile.CombinedOutput()
	if err != nil {
		fmt.Printf("Error compiling code. %v", string(compileOutput))
		os.Exit(1)
	}

	run := exec.Command("./" + compiled)
	if stdin != "" {
		run.Stdin = strings.NewReader(stdin)
	}

	runOutput, err := run.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running: %v", string(runOutput))
		os.Exit(1)
	}

	fmt.Printf(string(runOutput))
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
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
	go func() {
		time.Sleep(time.Second * 3)
		fmt.Println("EXECUTION TIMED OUT -1")
		os.Exit(2) //TLE ERROR
	}()

	compile := exec.Command("gcc", name, "-o", compiled, "-Wall")
	compileOutput, err := compile.CombinedOutput()
	if err != nil {
		fmt.Printf("Error compiling code. %v", string(compileOutput))
		os.Exit(3) //Compile Error
	}

	run := exec.Command("./" + compiled)
	if stdin != "" {
		run.Stdin = strings.NewReader(stdin)
	}

	runOutput, err := run.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running: %v", string(runOutput))
		os.Exit(4) //runtime error
	}

	value := string(runOutput)
	fmt.Println(value)
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	arguments := os.Args
	if len(arguments) != 4 {
		fmt.Println("Usage: mark.exe program.exe input output")
		os.Exit(1)
	}

	program, input, output := arguments[1], arguments[2], arguments[3]

	out, err := os.Open(output)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(out)
	if err != nil {
		log.Fatal(err)
	}

	answer := buf.String()
	buf.Reset()

	cmd := exec.Command(program)
	in, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}

	cmd.Stdin = in
	cmd.Stdout = &buf

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	result := buf.String()
	if answer != result {
		fmt.Printf("expected: %s\n", answer)
		fmt.Printf("result: %s\n", result)
	} else {
		fmt.Printf("CORRECT!\n")
	}
}
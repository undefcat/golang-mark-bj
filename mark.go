package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type testCase struct {
	name   string
	result bool
}

var (
	inputPath   string
	outputPath  string
	programPath string
)

func execute(tc testCase, c chan <- testCase) {
	var buf bytes.Buffer
	cmd := exec.Command(programPath)

	in, err := os.Open(filepath.Join(inputPath, tc.name))
	if err != nil {
		log.Fatal(err)
	}

	answer, err := ioutil.ReadFile(filepath.Join(outputPath, tc.name))
	if err != nil {
		log.Fatal(err)
	}

	cmd.Stdin = in
	cmd.Stdout = &buf
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	answer = bytes.TrimSpace(answer)
	result := bytes.TrimSpace(buf.Bytes())

	strAns := string(answer)
	strRes := string(result)

	tc.result = strAns == strRes
	c <- tc
}

func main() {
	args := os.Args
	if len(args) != 4 {
		fmt.Println("Usage: mark program input output")
		os.Exit(1)
	}

	programPath, inputPath, outputPath = args[1], args[2], args[3]

	inputInfo, err := os.Lstat(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	if inputInfo.IsDir() {
		dir, err := os.Open(inputPath)
		if err != nil {
			log.Fatal(err)
		}

		// -1 means read all files in dir.
		inputFiles, err := dir.Readdir(-1)
		if err != nil {
			log.Fatal(err)
		}

		tcChan := make(chan testCase, 3)

		i := 0
		for _, inputFile := range inputFiles {
			tc := testCase{inputFile.Name(), false}
			go execute(tc, tcChan)
			i++
		}

		for {
			tc := <-tcChan
			i--

			if tc.result {
				fmt.Printf("case %s CORRECT!\n", tc.name)
			} else {
				fmt.Printf("case %s WRONG!\n", tc.name)
			}

			if i == 0 {
				close(tcChan)
				break
			}
		}
	}
}
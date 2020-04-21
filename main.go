package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	parser := Parser{}
	fi, err := os.Stat(os.Args[1])
	check(err)
	if fi.IsDir() {
		outputFile, err := os.Create(filepath.Join(os.Args[1], fi.Name()+".asm"))
		check(err)
		defer outputFile.Close()
		files, err := ioutil.ReadDir(os.Args[1])
		parser.init(outputFile, true)
		for _, f := range files {
			if strings.HasSuffix(f.Name(), "vm") {
				inputFile, err := os.Open(filepath.Join(os.Args[1], f.Name()))
				check(err)
				parser.nextFile(strings.Replace(f.Name(), ".vm", "", 1), inputFile)
				for parser.hasMoreCommands() {
					parser.advance()
				}
				inputFile.Close()
			}
		}
	} else {
		outputFile, err := os.Create(strings.Replace(os.Args[1], ".vm", ".asm", 1))
		check(err)
		defer outputFile.Close()
		parser.init(outputFile, false)
		inputFile, err := os.Open(os.Args[1])
		check(err)

		parser.nextFile(strings.Replace(fi.Name(), ".vm", "", 1), inputFile)
		for parser.hasMoreCommands() {
			parser.advance()
		}
		inputFile.Close()
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Parser ...
type Parser struct {
	scanner    *bufio.Scanner
	codeWriter *CodeWriter
	nextLine   string
}

func (p *Parser) init(outputFile *os.File, isDir bool) {
	p.codeWriter = &CodeWriter{}
	p.codeWriter.init(outputFile)
	if isDir {
		p.codeWriter.WriteBootStrap()
	}
}

func (p *Parser) nextFile(className string, inputFile *os.File) {
	fmt.Printf("Reading file from...: %s\n", inputFile.Name())
	p.scanner = bufio.NewScanner(inputFile)
	p.codeWriter.setClassName(className)
}

func (p *Parser) hasMoreCommands() bool {
	for p.scanner.Scan() {
		p.nextLine = p.scanner.Text()
		if p.nextLine != "" && !strings.HasPrefix(p.nextLine, "//") {
			return true
		}
	}
	return false
}

func (p *Parser) advance() {
	fmt.Print("Advancing... " + p.nextLine + "\n")
	if isPushPopCmd(p.nextLine) {
		p.codeWriter.WritePushPop(p.nextLine)
	} else if isArithmeticCmd(p.nextLine) {
		p.codeWriter.WriteArithmetic(p.nextLine)
	} else {
		p.codeWriter.WriteControl(p.nextLine)
	}
}

func isArithmeticCmd(cmd string) bool {
	return strings.HasPrefix(cmd, "and") ||
		strings.HasPrefix(cmd, "sub") ||
		strings.HasPrefix(cmd, "add") ||
		strings.HasPrefix(cmd, "or") ||
		strings.HasPrefix(cmd, "neg") ||
		strings.HasPrefix(cmd, "not") ||
		strings.HasPrefix(cmd, "eq") ||
		strings.HasPrefix(cmd, "lt") ||
		strings.HasPrefix(cmd, "gt")
}
func isPushPopCmd(cmd string) bool {
	return strings.HasPrefix(cmd, "push") || strings.HasPrefix(cmd, "pop")
}

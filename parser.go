package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Parser ...
type Parser struct {
	file       *os.File
	scanner    *bufio.Scanner
	codeWriter *CodeWriter
	nextLine   string
}

func (p *Parser) init(filePath string) {
	var err error
	p.file, err = os.Open(filePath)
	p.codeWriter = &CodeWriter{}
	fmt.Printf("Reading file from...: %s\n", filePath)
	check(err)
	p.scanner = bufio.NewScanner(p.file)
	p.codeWriter.init(strings.Replace(filePath, "vm", "asm", 1))
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
	}
}

func (p *Parser) close() error {
	if err := p.codeWriter.close(); err != nil {
		return err
	}
	return p.file.Close()
}

func isArithmeticCmd(cmd string) bool {
	return opMap[cmd] != ""
}

func isPushPopCmd(cmd string) bool {
	return strings.HasPrefix(cmd, "push") || strings.HasPrefix(cmd, "pop")

}

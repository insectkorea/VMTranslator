package main

import "os"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	parser := Parser{}
	parser.init(os.Args[1])

	for parser.hasMoreCommands() {
		parser.advance()
	}
	parser.close()
}

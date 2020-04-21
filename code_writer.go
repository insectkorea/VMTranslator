package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CodeWriter struct {
	file      *os.File
	className string
}

var (
	segMap = map[string]string{
		"local":    "LCL",
		"this":     "THIS",
		"that":     "THAT",
		"argument": "ARG",
		"temp":     "R5",
		"pointer":  "R3",
	}
	opMap = map[string]string{
		"and": "&",
		"sub": "-",
		"add": "+",
		"or":  "|",
		"neg": "-",
		"not": "!",
		"eq":  "JEQ",
		"lt":  "JLT",
		"gt":  "JGT",
	}
	count = 0
)

func (c *CodeWriter) init(outputFile *os.File) {
	var err error
	c.file = outputFile
	fmt.Printf("Writing file to...: %s\n", outputFile.Name())
	check(err)
}

func (c *CodeWriter) setClassName(className string) {
	c.className = className
}

func (c *CodeWriter) WriteBootStrap() {
	c.write(formatBootStrap())
}

func (c *CodeWriter) WriteArithmetic(cmd string) {
	cmds := strings.Split(cmd, " ")
	// eq, lt, gt, sub, add, neg, or, not, and
	switch cmds[0] {
	case "eq", "lt", "gt":
		c.write(formatComp(opMap[cmds[0]], count))
		count++
		break
	case "not", "neg":
		c.write(formatNegNot(opMap[cmds[0]]))
		break
	case "and", "or", "add", "sub":
		c.write(formatArithmeticLogical(opMap[cmds[0]]))
		break
	default:
		panic("Unidentified operation")
	}
}

func (c *CodeWriter) WriteControl(cmd string) {
	cmds := strings.Split(cmd, " ")
	switch cmds[0] {
	case "label":
		c.write(formatLabel(c.className, cmds[1]))
		break
	case "goto":
		c.write(formatGoto(c.className, cmds[1]))
		break
	case "if-goto":
		c.write(formatIf(c.className, cmds[1]))
		break
	case "function":
		nVar, _ := strconv.Atoi(cmds[2])
		c.write(formatFunction(cmds[1], nVar))
		break
	case "call":
		nArg, _ := strconv.Atoi(cmds[2])
		c.write(formatCall(cmds[1], nArg, cmds[1]+"$return."+strconv.Itoa(count)))
		count++
		break
	case "return":
		c.write(formatReturn())
		break
	}
}

func (c *CodeWriter) WritePushPop(cmd string) {
	cmds := strings.Split(cmd, " ")
	switch cmds[0] {
	case "push":
		switch cmds[1] {
		case "constant":
			c.write(formatPushConst(cmds[2]))
			break
		case "local", "argument", "this", "that":
			c.write(formatPushSeg(cmds[2], segMap[cmds[1]]))
			break
		case "temp", "pointer":
			c.write(formatPushTemp(cmds[2], segMap[cmds[1]]))
			break
		case "static":
			c.write(formatPushStatic(c.className, cmds[2]))
			break
		default:
			panic("Unidentified segment to push" + cmds[1])
		}
		break
	case "pop":
		switch cmds[1] {
		case "local", "argument", "this", "that":
			c.write(formatPopSeg(cmds[2], segMap[cmds[1]]))
			break
		case "temp", "pointer":
			c.write(formatPopTemp(cmds[2], segMap[cmds[1]]))
			break
		case "static":
			c.write(formatPopStatic(c.className, cmds[2], segMap[cmds[1]]))
			break
		default:
			panic("Unidentified segment to pop " + cmds[1])
		}
	default:
		panic("Unidentified operation " + cmds[1])
	}
}

func (c *CodeWriter) write(cmd string) {
	fmt.Print("Writing... " + cmd + "\n")
	c.file.WriteString(cmd)
}

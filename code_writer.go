package main

import (
	"fmt"
	"os"
	"strings"
)

type CodeWriter struct {
	file *os.File
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

func (c *CodeWriter) init(filePath string) {
	var err error
	c.file, err = os.Create(filePath)
	fmt.Printf("Writing file to...: %s\n", filePath)
	check(err)
}

func (c *CodeWriter) WriteArithmetic(cmd string) {
	// eq, lt, gt, sub, add, neg, or, not, and
	switch cmd {
	case "eq", "lt", "gt":
		c.write(formatComp(opMap[cmd], count))
		count++
		break
	case "not", "neg":
		c.write(formatNegNot(opMap[cmd]))
		break
	case "and", "or", "add", "sub":
		c.write(formatArithmeticLogical(opMap[cmd]))
		break
	default:
		panic("Unidentified operation")
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
			c.write(formatPushStatic(cmds[2]))
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
			c.write(formatPopStatic(cmds[2], segMap[cmds[1]]))
			break
		default:
			panic("Unidentified segment to pop " + cmds[1])
		}
	default:
		panic("Unidentified operation " + cmds[1])
	}
}

func formatComp(op string, pos int) string {
	return fmt.Sprintf(`@SP
AM=M-1
D=M
A=A-1
D=M-D
@COMP%d
D; %s
@SP
AM=M-1
M=0
@DONE%d
D; JMP
(COMP%d)
@SP
AM=M-1
M=-1
(DONE%d)
@SP
M=M+1
`, pos, op, pos, pos, pos)
}

func formatArithmeticLogical(op string) string {
	return fmt.Sprintf(`@SP
AM=M-1
D=M
A=A-1
M=M%sD
`, op)
}

func formatNegNot(op string) string {
	return fmt.Sprintf(`@SP
A=M-1
M=%sM
`, op)
}

func formatPushConst(cons string) string {
	return fmt.Sprintf(`@%s
D=A
@SP
A=M
M=D
@SP
M=M+1
`, cons)
}

func formatPushSeg(pos string, seg string) string {
	return fmt.Sprintf(`@%s
D=A
@%s
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
`, pos, seg)
}

func formatPushTemp(pos string, seg string) string {
	return fmt.Sprintf(`@%s
D=A
@%s
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
`, pos, seg)
}

func formatPushStatic(pos string) string {
	return fmt.Sprintf(`@STATIC%s
D=M
@SP
A=M
M=D
@SP
M=M+1
`, pos)
}

func formatPopSeg(pos string, seg string) string {
	return fmt.Sprintf(`@%s
D=A
@%s
D=D+M
@R13
M=D
@SP
AM=M-1
D=M
@R13
A=M
M=D
`, pos, seg)
}

func formatPopTemp(pos string, seg string) string {
	return fmt.Sprintf(`@%s
D=A
@%s
D=D+A
@R13
M=D
@SP
AM=M-1
D=M
@R13
A=M
M=D
`, pos, seg)
}

func formatPopStatic(pos string, seg string) string {
	return fmt.Sprintf(`@SP
AM=M-1
D=M
@STATIC%s
M=D
`, pos)
}

func (c *CodeWriter) write(cmd string) {
	fmt.Print("Writing... " + cmd + "\n")
	c.file.WriteString(cmd)
}

func (c *CodeWriter) close() error {
	return c.file.Close()
}

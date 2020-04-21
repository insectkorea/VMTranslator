package main

import "fmt"

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

func formatPushStatic(className string, pos string) string {
	return fmt.Sprintf(`@%s%s
D=M
@SP
A=M
M=D
@SP
M=M+1
`, className, pos)
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

func formatPopStatic(className string, pos string, seg string) string {
	return fmt.Sprintf(`@SP
AM=M-1
D=M
@%s%s
M=D
`, className, pos)
}

func formatLabel(filename string, label string) string {
	return fmt.Sprintf(`(%s.$%s)
`, filename, label)
}

func formatGoto(filename string, label string) string {
	return fmt.Sprintf(`@%s.$%s
D; JMP
`, filename, label)
}

func formatIf(filename string, label string) string {
	return fmt.Sprintf(`@SP
AM=M-1
D=M
@%s.$%s
D; JNE
`, filename, label)
}

func formatFunction(funcName string, nVar int) string {
	local := ""
	initializeLocal := `@SP
AM=M+1
M=0
`
	for i := 0; i < nVar; i++ {
		local += initializeLocal
	}

	return fmt.Sprintf(`(%s)
%s
`, funcName, local)

}

func formatReturn() string {
	return fmt.Sprintf(`@LCL
D=M
@R13
M=D
@5
D=D-A
A=D
D=M
@R14
M=D
@SP
AM=M-1
D=M
@ARG
A=M
M=D
@ARG
D=M
@SP
M=D+1
@R13
D=M
@1
A=D-A
D=M
@THAT
M=D
@R13
D=M
@2
A=D-A
D=M
@THIS
M=D
@R13
D=M
@3
A=D-A
D=M
@ARG
M=D
@R13
D=M
@4
A=D-A
D=M
@LCL
M=D
@R14
A=M
D; JMP
`)

}

func formatCall(fnc string, nArg int, retAdr string) string {
	return fmt.Sprintf(`@%s
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
D=M
@5
D=D-A
@%d
D=D-A
@ARG
M=D
@SP
D=M
@LCL
M=D
@%s
D; JMP
(%s)
`, retAdr, nArg, fnc, retAdr)
}

func formatBootStrap() string {
	return fmt.Sprintf(`@256
D=A
@SP
M=D
%s
`, formatCall("Sys.init", 0, "Bootstrap.$return"))
}

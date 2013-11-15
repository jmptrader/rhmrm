package rhmrm

import "fmt"

// Instruction is a unit of machine language
type Instruction Word

// MkInstruction2 creates 2-operand instruction.
func MkInstruction2(op, a, b int16) (i Instruction) {
	i = Instruction(op & 0x3f)
	i |= Instruction(a & 0x1f << 6)
	i |= Instruction(b & 0x1f << 11)
	return
}

// WMkInstruction2 is like MkInstruction2, but returns Word.
func WMkInstruction2(op, a, b int16) Word {
	return Word(MkInstruction2(op, a, b))
}

// MkInstruction1 creates 1-operand instruction.
func MkInstruction1(op, c int16) (i Instruction) {
	i = Instruction(op & 0x3f)
	i |= Instruction(c & 0x7ff << 6)
	return
}

// WMkInstruction1 is like MkInstruction1, but returns Word.
func WMkInstruction1(op, c int16) Word {
	return Word(MkInstruction1(op, c))
}

// sextend10 sign-extends its' 10 bit parameter
func sextend10(n Word) (r Word) {
	if n & (1 << 9) != 0 {
		r = n | 0xfc00
	} else {
		r = n &^ 0xfc00
	}
	return
}

// op returns lowest 6 bits of an instruction
func (i Instruction) op() Word { return Word(i & 0x3f) }

// opstring returns string representation of the opcode
func (i Instruction) opstring() string { return op_strings[i.op()] }

// a returns bits [9:5] of an instruction
func (i Instruction) a() Word { return Word(i >> 6 & 0x1f) }

// b returns bits [15:10] of an instruction
func (i Instruction) b() Word { return Word(i >> 11) }

// c returns bits [15:6] of an instruction
func (i Instruction) c() Word { return Word(i >> 6) }

// cs is like c, but sign-extends the result to 16 bits
func (i Instruction) cs() Word { return sextend10(i.c()) }

// decouple extracts operator and operands from an instruction
func (i Instruction) decouple() (op Word, args []Word) {
	op = i.op()
	if op >= OP_IMP && op <= OP_CMN {
		args = []Word{ i.a(), i.b() }
	} else {
		args = []Word{ i.c() }
	}
	return
}

/*
// opstring returns string for opcode
func opstring(op Word) string {
	if int(op) >= len(op_names) {
		panic(fmt.Sprint("Bad operator index:", op))
	}
	return op_names[op]
}
*/

// kstring returns string for a control register
func kstring(r Word) (str string) {
	mode := r >> 3
	kreg := r & 7

	str = fmt.Sprintf("%2i", kreg)

	switch mode {
	case AM_AND:
		str = fmt.Sprintf("and(%s)", str)
	case AM_IOR:
		str = fmt.Sprintf("ior(%s)", str)
	case AM_XOR:
		str = fmt.Sprintf("xor(%s)", str)
	}

	return str
}

func (i Instruction) String() (s string) {
	switch {
	case i.op() == OP_IMP:
		s = fmt.Sprintf("%3v %3v %2d",
			i.opstring(), imp_strings[i.a()], i.b())
	case i.op() == OP_MTC:
		s = fmt.Sprintf("%3v %s, %2d",
			i.opstring(), kstring(i.a()), i.b())
	case i.op() == OP_MFC:
		s = fmt.Sprintf("%3v %2d, %s",
			i.opstring(), i.a(), kstring(i.b()))
	case i.op() >= OP_MOV && i.op() <= OP_CMN:
		s = fmt.Sprintf("%3v %2d, %d", i.opstring(), i.a(), i.b())
	case i.op() >= OP_JMP && i.op() <= OP_JNE:
		s = fmt.Sprintf("%3v %2d", i.opstring(), int16(i.cs()))
	case i.op() == OP_SWI || i.op() == OP_HWI || i.op() == OP_IRE:
		s = fmt.Sprintf("%3v %2d", i.opstring(), i.c())
	}
	return
}

var op_strings = []string{
	OP_IMP: "imp",
	OP_MOV: "mov",
	OP_MTC: "mtc",
	OP_MFC: "mfc",

	OP_STR: "str",
	OP_PSH: "psh",
	OP_LOA: "loa",
	OP_POP: "pop",
	OP_MOM: "mom",

	OP_SRL: "srl",

	OP_ADD: "add",
	OP_ADX: "adx",
	OP_SUB: "sub",
	OP_SBX: "sbx",
	OP_MUL: "mul",
	OP_MLI: "mli",
	OP_DIV: "div",
	OP_DVI: "dvi",
	OP_MOD: "mod",
	OP_MDI: "mdi",
	OP_INC: "inc",

	OP_AND: "and",
	OP_IOR: "ior",
	OP_XOR: "xor",
	OP_BIC: "bic",
	OP_SHL: "shl",
	OP_ASR: "asr",
	OP_SHR: "shr",
	OP_ROL: "rol",
	OP_ROR: "ror",

	OP_TST: "tst",
	OP_TEQ: "teq",
	OP_CMP: "cmp",
	OP_CMN: "cmn",

	OP_JMP: "jmp",
	OP_JLT: "jlt",
	OP_JLE: "jle",
	OP_JGT: "jgt",
	OP_JGE: "jge",
	OP_JEQ: "jeq",
	OP_JNE: "jne",

	OP_SWI: "swi",
	OP_HWI: "hwi",
	OP_IRE: "ire",
}

var imp_strings = []string{
	IMP_BRK: "brk",
	IMP_MOV: "mov",
	IMP_MTC: "mtc",

	IMP_STR: "str",
	IMP_PSH: "psh",

	IMP_SRL: "srl",
	
	IMP_ADD: "add",
	IMP_ADX: "adx",
	IMP_SUB: "sub",
	IMP_SBX: "sbx",
	IMP_MUL: "mul",
	IMP_MLI: "mli",
	IMP_DIV: "div",
	IMP_DVI: "dvi",
	IMP_MOD: "mod",
	IMP_MDI: "mdi",
	IMP_INC: "inc",

	IMP_AND: "and",
	IMP_IOR: "ior",
	IMP_XOR: "xor",
	IMP_BIC: "bic",
	IMP_SHL: "shl",
	IMP_ASR: "asr",
	IMP_SHR: "shr",
	IMP_ROL: "rol",
	IMP_ROR: "ror",

	IMP_TST: "tst",
	IMP_TEQ: "teq",
	IMP_CMP: "cmp",
	IMP_CMN: "cmn",
}

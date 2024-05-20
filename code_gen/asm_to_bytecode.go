package code_gen

import (
	"go-compiler/utils"
	"strconv"
	"strings"
)

type Assembler struct {
	Asm       string    `json:"asm"`
	Memory    []int64   `json:"memory"`
	Pc        int64     `json:"pc"`
	CodeCount int64     `json:"code_count"`
	Register  *Register `json:"register"`
}

func NewAssembler() *Assembler {
	return &Assembler{
		Memory:   make([]int64, 0),
		Register: NewRegister(),
	}
}

func (this *Assembler) AddAsm(asm string) {
	this.Asm = asm
}

func (this *Assembler) AddMemory(address int64) {
	this.Memory = append(this.Memory, address)
}

func (this *Assembler) AddPc(pc int64) {
	this.Pc += pc
}

func (this *Assembler) AddCodeCount(count int64) {
	this.CodeCount += count
}

func (this *Assembler) turnCodeToNum(code string) int64 {
	// 数字、string、bool、null
	// 前 3 位是类型
	// 数字 000
	// string 001
	// bool 010
	// null 011
	var result int64
	if code == "true" {
		result = (2 << 13) + 1
	} else if code == "false" {
		result = (2 << 13) + 0
	} else if code == "null" {
		result = (3 << 13) + 0
	} else if strings.Contains("1234567890", string(code[0])) {
		result, _ = strconv.ParseInt(code, 10, 64)
	} else if strings.Contains("abcdefghijklmnopqrstuvwxyz", string(code[0])) {

	}
	return result
}

func (this *Assembler) opSet2(code []string) {
	// set a1 true
	this.Pc += 3
	reg1 := code[1]
	r1 := this.Register.ReturnRegByName(reg1)
	num := this.turnCodeToNum(code[2])
	this.Memory = append(this.Memory, InstructionSet2.Value(), r1, num)
	this.CodeCount += 3
}

func (this *Assembler) opSub2(code []string) {
	//sub2 a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionSubtract2.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opAdd2(code []string) {
	//add2 a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionAdd2.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opMul2(code []string) {
	//mul a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionMultiply2.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opDiv2(code []string) {
	//div2 a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionDiv2.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opPush(code []string) {
	//push reg
	this.Pc += 2
	reg1 := code[1]

	r1 := this.Register.ReturnRegByName(reg1)

	this.Memory = append(this.Memory, InstructionPush.Value(), r1)
	this.CodeCount += 2
}

func (this *Assembler) opPop(code []string) {
	//pop reg
	this.Pc += 2
	reg1 := code[1]

	r1 := this.Register.ReturnRegByName(reg1)

	this.Memory = append(this.Memory, InstructionPop.Value(), r1)
	this.CodeCount += 2
}

func (this *Assembler) opBoolAnd(code []string) {
	//bool_and a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionBoolAnd.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opBoolOr(code []string) {
	//bool_or a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionBoolOr.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opBoolLessThan(code []string) {
	//bool_lt a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionBoolLessThan.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opBoolLessThanEquals(code []string) {
	//bool_lte a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionBoolLessThanEquals.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opBoolGreaterThan(code []string) {
	//bool_gt a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionBoolGreaterThan.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opBoolGreaterThanEquals(code []string) {
	//bool_gte a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionBoolGreaterThanEquals.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opBoolEquals(code []string) {
	//bool_eq a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionBoolEquals.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opBoolNotEquals(code []string) {
	//bool_neq a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionBoolNotEquals.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opPlusAssign(code []string) {
	//num_plus_eq a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionPlusAssign.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opSubtractAssign(code []string) {
	//num_sub_eq a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionSubtractAssign.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opMultiplyAssign(code []string) {
	//num_mul_eq a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionMultiplyAssign.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opDivideAssign(code []string) {
	//num_div_eq a1 a2 a3
	this.Pc += 4
	reg1 := code[1]
	reg2 := code[2]
	reg3 := code[3]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)
	r3 := this.Register.ReturnRegByName(reg3)

	this.Memory = append(this.Memory, InstructionDivideAssign.Value(), r1, r2, r3)
	this.CodeCount += 4
}

func (this *Assembler) opSaveFromRegister2(code []string) {
	//save_from_register a1 a2
	this.Pc += 3
	reg1 := code[1]
	reg2 := code[2]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)

	this.Memory = append(this.Memory, InstructionSaveFromRegister2.Value(), r1, r2)
	this.CodeCount += 3
}

func (this *Assembler) opLoadFromRegister2(code []string) {
	//load_from_register a1 a2
	this.Pc += 3
	reg1 := code[1]
	reg2 := code[2]

	r1 := this.Register.ReturnRegByName(reg1)
	r2 := this.Register.ReturnRegByName(reg2)

	this.Memory = append(this.Memory, InstructionLoadFromRegister2.Value(), r1, r2)
	this.CodeCount += 3
}

func (this *Assembler) opJump(code []string) {
	//jump @fuck
	//this.Pc += 3
	//reg1 := code[1]
	//reg2 := code[2]
	//
	//r1 := this.Register.ReturnRegByName(reg1)
	//r2 := this.Register.ReturnRegByName(reg2)
	//
	//this.Memory = append(this.Memory, InstructionSaveFromRegister2.Value(), r1, r2)
	//this.CodeCount += 3
}

func (this *Assembler) opJumpFromRegister(code []string) {
	//jump @fuck
	//this.Pc += 3
	//reg1 := code[1]
	//reg2 := code[2]
	//
	//r1 := this.Register.ReturnRegByName(reg1)
	//r2 := this.Register.ReturnRegByName(reg2)
	//
	//this.Memory = append(this.Memory, InstructionSaveFromRegister2.Value(), r1, r2)
	//this.CodeCount += 3
}

func (this *Assembler) opHalt(code []string) {
	//halt
	this.Pc += 1
	this.Memory = append(this.Memory, InstructionHalt.Value())
	this.CodeCount += 1
}

func (this *Assembler) opFuncInfo(op string) func(code []string) {
	fun := map[string]func(asm []string){
		InstructionSet2.Name(): this.opSet2,
		// +-*/
		InstructionAdd2.Name():      this.opAdd2,
		InstructionSubtract2.Name(): this.opSub2,
		InstructionMultiply2.Name(): this.opMul2,
		InstructionDiv2.Name():      this.opDiv2,
		// save \load
		InstructionSaveFromRegister2.Name(): this.opSaveFromRegister2,
		InstructionLoadFromRegister2.Name(): this.opLoadFromRegister2,
		// bool
		InstructionBoolAnd.Name():               this.opBoolAnd,
		InstructionBoolOr.Name():                this.opBoolOr,
		InstructionBoolGreaterThan.Name():       this.opBoolGreaterThan,
		InstructionBoolGreaterThanEquals.Name(): this.opBoolGreaterThanEquals,
		InstructionBoolLessThan.Name():          this.opBoolLessThan,
		InstructionBoolLessThanEquals.Name():    this.opBoolLessThanEquals,
		InstructionBoolEquals.Name():            this.opBoolEquals,
		InstructionBoolNotEquals.Name():         this.opBoolNotEquals,
		// push pop
		InstructionPush.Name(): this.opPush,
		InstructionPop.Name():  this.opPop,
		// += -= *= /=
		InstructionPlusAssign.Name():     this.opPlusAssign,
		InstructionSubtractAssign.Name(): this.opSubtractAssign,
		InstructionMultiplyAssign.Name(): this.opMultiplyAssign,
		InstructionDivideAssign.Name():   this.opDivideAssign,
		// jump
		InstructionJump.Name():             this.opJump,
		InstructionJumpFromRegister.Name(): this.opJumpFromRegister,
		// halt
		InstructionHalt.Name(): this.opHalt,
	}
	return fun[op]
}

// todo
func (this *Assembler) opPseudoFunction(code []string) {
}

// todo
func (this *Assembler) opPseudoReturn(code []string) {

}

// todo
func (this *Assembler) opPseudoFuncVar(code []string) {

}

func (this *Assembler) opPseudoFuncInfo(op string) func(code []string) {
	fun := map[string]func(asm []string){
		InstructionPseudoFunction.Name(): this.opPseudoFunction,
		InstructionPseudoReturn.Name():   this.opPseudoReturn,
		InstructionPseudoFuncVar.Name():  this.opPseudoFuncVar,
	}
	return fun[op]
}

// 这里去掉注释啥的
func (this *Assembler) preProcessAsm(asm string) string {
	return asm
}

func (this *Assembler) compileMachineCode(asm string) []int64 {
	memory := make([]int64, 0)

	asm = this.preProcessAsm(asm)
	lines := strings.Split(asm, "\n")
	var i = 0
	for i < len(lines) {
		lines[i] = strings.TrimSpace(lines[i])
		// 跳过空行
		if len(lines[i]) == 0 {
			i += 1
			continue
		}
		// 将指令分割成数组
		var line = strings.Split(lines[i], " ")
		op := string(line[0])
		utils.LogInfo("op", op)
		op = strings.TrimSpace(op)
		if len(op) == 0 {
			i += 1
			continue
		}
		if string(op[0]) == "@" {
			// 处理地址
		} else if string(op[0]) == "." {
			// 处理伪指令
			this.opPseudoFuncInfo(op)(line)
		} else {
			// 处理指令
			this.opFuncInfo(op)(line)
		}
		i += 1
	}
	memory = this.Memory
	return memory
}

func (this *Assembler) Compile() []int64 {
	this.Memory = this.compileMachineCode(this.Asm)
	return this.Memory
}

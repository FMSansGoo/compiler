package code_gen

import (
	"go-compiler/utils"
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
	// 数字、string、bool、null、void
	return 0
}

func (this *Assembler) opSet2(code []string) {
	this.Pc += 3
	reg1 := code[1]
	r1 := this.Register.ReturnRegByName(reg1)
	num := this.turnCodeToNum(code[2])
	this.Memory = append(this.Memory, 0, r1, num)
	this.CodeCount += 3
}

func (this *Assembler) opFuncInfo(op string) func(code []string) {
	fun := map[string]func(asm []string){
		InstructionSet2.Name(): this.opSet2,
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
		// 处理伪指令
		if string(op[0]) == "." {

		} else {
			// 利用 表驱动法
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

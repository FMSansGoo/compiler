package asm_vm_stack_base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go-compiler/utils"
)

type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(GetOpCodeFromValue(ins[i]).Name())
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		utils.LogInfo("def   ", def)
		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := def.OperandWidths

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, def.OperandWidths*def.OperandNums)
	offset := 0

	switch def.OperandWidths {
	case 2:
		operands[offset] = int(ReadUint16(ins[offset:]))
	}

	offset += def.OperandWidths

	return operands, offset
}

type Definition struct {
	Name          string // 指令名称
	OperandWidths int    // 操作数有多少个字节宽度
	OperandNums   int    // 除了指令外有多少个操作数

}

var definitions = map[OpCode]*Definition{
	//OperandWidths 包含每个操作数占用的字节数。
	// const 2 字节，16位，1 个操作数
	OpCodeConstant:          {OpCodeConstant.Name(), 2, 1},
	OpCodeAdd:               {OpCodeAdd.Name(), 0, 0},
	OpCodeSub:               {OpCodeSub.Name(), 0, 0},
	OpCodeMul:               {OpCodeMul.Name(), 0, 0},
	OpCodeDiv:               {OpCodeDiv.Name(), 0, 0},
	OpCodePop:               {OpCodePop.Name(), 0, 0},
	OpCodeTrue:              {OpCodeTrue.Name(), 0, 0},
	OpCodeFalse:             {OpCodeFalse.Name(), 0, 0},
	OpCodeNull:              {OpCodeNull.Name(), 0, 0},
	OpCodeEquals:            {OpCodeEquals.Name(), 0, 0},
	OpCodeNotEquals:         {OpCodeNotEquals.Name(), 0, 0},
	OpCodeGreaterThan:       {OpCodeGreaterThan.Name(), 0, 0},
	OpCodeGreaterThanEquals: {OpCodeGreaterThanEquals.Name(), 0, 0},
	OpCodeLessThan:          {OpCodeLessThan.Name(), 0, 0},
	OpCodeLessThanEquals:    {OpCodeLessThanEquals.Name(), 0, 0},

	OpCodeNot:        {OpCodeNot.Name(), 0, 0},
	OpCodeMinus:      {OpCodeMinus.Name(), 0, 0},
	OpCodeObjectCall: {OpCodeObjectCall.Name(), 0, 0},
	//OpCodeBreak:       {OpCodeBreak.Name(), 0, 0},

	OpCodeJump:          {OpCodeJump.Name(), 2, 1},
	OpCodeJumpNotTruthy: {OpCodeJumpNotTruthy.Name(), 2, 1},
	// 全局变量能占有 65536 字节
	OpCodeSetGlobal: {OpCodeSetGlobal.Name(), 2, 1},
	OpCodeGetGlobal: {OpCodeGetGlobal.Name(), 2, 1},
	// 局部变量占有 256 字节就行
	OpCodeSetLocal: {OpCodeSetLocal.Name(), 1, 1},
	OpCodeGetLocal: {OpCodeGetLocal.Name(), 1, 1},
	//
	OpCodeArray:  {OpCodeArray.Name(), 2, 1},
	OpCodeDict:   {OpCodeDict.Name(), 2, 1},
	OpCodeReturn: {OpCodeReturn.Name(), 0, 0},
	// closure，第一个数函数在常量池的索引，第二个数用于指定栈中有多少自由变量需要转移到即将创建的闭包中
	OpCodeClosure:      {OpCodeClosure.Name(), 2, 2},
	OpCodeFunctionCall: {OpCodeFunctionCall.Name(), 2, 1},
	// 获取自由变量
	OpCodeGetFree: {OpCodeGetFree.Name(), 1, 1},
	// 获取内置函数
	OpCodeGetBuiltin: {OpCodeGetBuiltin.Name(), 1, 1},
}

func Lookup(op string) (*Definition, error) {
	opCode := GetOpCodeFromName(op)
	def, ok := definitions[opCode]
	if !ok {
		return nil, fmt.Errorf("opcode %s undefined", op)
	}
	return def, nil
}

// Make 编译指令和操作数为字节码[]byte
func GenerateByte(op OpCode, operands ...int) []byte {
	// 读取指令
	opInfo, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	// 计算指令字节码长度
	instructionLen := 1
	instructionLen += opInfo.OperandNums * opInfo.OperandWidths

	// 字节码第一个元素为指令
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op.Value())

	// 将操作数依次加入字节码中
	offset := 1
	for _, o := range operands {
		// 读取操作数的长度
		width := opInfo.OperandWidths
		switch width {
		case 2:
			// 将操作数转化为大端存储的字节，存入instruction，一次占据同长度相同的数组单元
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		// 计算偏移量
		offset += width
	}

	return instruction
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}

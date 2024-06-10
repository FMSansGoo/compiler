package asm_vm_stack_base

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Definition struct {
	Name          string // 指令名称
	OperandWidths int    // 操作数有多少个字节宽度
	OperandNums   int    // 除了指令外有多少个操作数

}

var definitions = map[OpCode]*Definition{
	//OperandWidths 包含每个操作数占用的字节数。
	// const 2 字节，16位，1 个操作数
	OpCodeConstant:    {OpCodeConstant.Name(), 2, 1},
	OpCodeAdd:         {OpCodeAdd.Name(), 0, 0},
	OpCodeSub:         {OpCodeSub.Name(), 0, 0},
	OpCodeMul:         {OpCodeMul.Name(), 0, 0},
	OpCodeDiv:         {OpCodeDiv.Name(), 0, 0},
	OpCodePop:         {OpCodePop.Name(), 0, 0},
	OpCodeTrue:        {OpCodeTrue.Name(), 0, 0},
	OpCodeFalse:       {OpCodeFalse.Name(), 0, 0},
	OpCodeNull:        {OpCodeNull.Name(), 0, 0},
	OpCodeEquals:      {OpCodeEquals.Name(), 0, 0},
	OpCodeNotEquals:   {OpCodeNotEquals.Name(), 0, 0},
	OpCodeGreaterThan: {OpCodeGreaterThan.Name(), 0, 0},
}

//OpEqual:         {"OpEqual", []int{}},
//OpNotEqual:      {"OpNotEqual", []int{}},
//OpGreaterThan:   {"OpGreaterThan", []int{}},

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

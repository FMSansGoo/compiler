package code_gen

import "fmt"

//op值	指令	    说明
//255	halt	终止程序
//0	set	将8位值设置到寄存器低8位
//1	load	从内存加载8位值到寄存器
//2	add	将两个8位寄存器值相加
//3	save	将寄存器值保存到8位内存地址
//4	compare	比较两个寄存器值大小
//5	jump_if_less	条件跳转到指定16位内存地址
//6	jump	无条件跳转到指定16位内存地址
//7	save_from_register	将一个寄存器8位值保存到另一寄存器指定的8位内存地址
//8	set2	与set类似,但操作数为16位
//9	load2	从16位内存地址加载值到寄存器
//10	add2	与add类似,但操作数为16位
//11	save2	将寄存器值保存到16位内存地址
//12	sub2	16位减法
//13	load_from_register	从一寄存器指定的8位内存地址加载值到另一寄存器
//14	load_from_register2	与load_from_register类似,但加载的是16位值
//15	save_from_register2	将16位寄存器值保存到另一寄存器指定的16位内存地址
//16	jump_from_register	跳转到一个寄存器指定的内存地址
//17	shift_right	对一个寄存器值执行逻辑右移
//19	bitAnd	对两个寄存器值执行位与操作
//20	multiply2	将两个16位寄存器值相乘
//21    push reg 把数据推到栈上，栈指针 + 2
//22    pop reg 把数据从栈上退出来，栈指针 - 2
// Instruction

type Instruction struct {
	name  string
	value int64
}

var (
	validInstructions = []Instruction{}

	InstructionHalt                  = newInstruction("halt", 255)
	InstructionSet                   = newInstruction("set", 0)
	InstructionLoad                  = newInstruction("load", 1)
	InstructionAdd                   = newInstruction("add", 2)
	InstructionSave                  = newInstruction("save", 3)
	InstructionCompare               = newInstruction("compare", 4)
	InstructionJumpIfLess            = newInstruction("jump_if_less", 5)
	InstructionJump                  = newInstruction("jump", 6)
	InstructionSaveFromRegister      = newInstruction("save_from_register", 7)
	InstructionSet2                  = newInstruction("set2", 8)
	InstructionLoad2                 = newInstruction("load2", 9)
	InstructionAdd2                  = newInstruction("add2", 10)
	InstructionSave2                 = newInstruction("save2", 11)
	InstructionSubtract2             = newInstruction("sub2", 12)
	InstructionLoadFromRegister      = newInstruction("load_from_register", 13)
	InstructionLoadFromRegister2     = newInstruction("load_from_register2", 14)
	InstructionSaveFromRegister2     = newInstruction("save_from_register2", 15)
	InstructionJumpFromRegister      = newInstruction("jump_from_register", 16)
	InstructionShiftRight            = newInstruction("shift_right", 17)
	InstructionBitAnd                = newInstruction("bit_and", 19)
	InstructionMultiply2             = newInstruction("mul2", 20)
	InstructionDiv2                  = newInstruction("div2", 21)
	InstructionPush                  = newInstruction("push", 22)
	InstructionPop                   = newInstruction("pop", 23)
	InstructionBoolAnd               = newInstruction("bool_and", 24)
	InstructionBoolLessThan          = newInstruction("bool_lt", 25)
	InstructionBoolGreaterThan       = newInstruction("bool_gt", 26)
	InstructionBoolLessThanEquals    = newInstruction("bool_lte", 27)
	InstructionBoolGreaterThanEquals = newInstruction("bool_gte", 28)
	InstructionBoolEquals            = newInstruction("bool_eq", 29)
	InstructionBoolNotEquals         = newInstruction("bool_neq", 30)
	InstructionBoolOr                = newInstruction("bool_or", 31)
	InstructionBoolNot               = newInstruction("bool_not", 32)
	InstructionPlusAssign            = newInstruction("num_plus_eq", 33)
	InstructionSubtractAssign        = newInstruction("num_sub_eq", 35)
	InstructionMultiplyAssign        = newInstruction("num_mul_eq", 36)
	InstructionDivideAssign          = newInstruction("num_div_eq", 37)

	InstructionIf = newInstruction("if", 50)
)

func newInstruction(name string, value int64) Instruction {
	o := Instruction{name: name, value: value}
	validInstructions = append(validInstructions, o)
	return o
}

func InstructionAll() []Instruction {
	return validInstructions
}

func (t Instruction) valid() bool {
	for _, v := range InstructionAll() {
		if v == t {
			return true
		}
	}
	return false
}

func (t Instruction) Value() int64 {
	if !t.valid() {
		panic(fmt.Errorf("invalid Instruction: (%+v)", t))
	}
	return t.value
}

func (t Instruction) ValuePtr() *int64 {
	v := t.Value()
	return &v
}

func (t Instruction) Name() string {
	if !t.valid() {
		panic(fmt.Errorf("invalid Instruction: (%+v)", t))
	}
	return t.name
}

func (t Instruction) NamePtr() *string {
	n := t.Name()
	return &n
}

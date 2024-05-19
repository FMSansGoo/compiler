package code_gen

import "fmt"

// pa: 0b00000000,
// a1: 0b00010000,
// a2: 0b00100000,
// a3: 0b00110000,
// c1: 0b01000000,
// f1: 0b01010000,

type Register struct {
	RegisterPointer int
}

func NewRegister() *Register {
	return &Register{
		RegisterPointer: 1,
	}
}

// 分配和回收寄存器
func (this *Register) ReturnRegAlloc() string {
	s := fmt.Sprintf("a%d", this.RegisterPointer)
	this.RegisterPointer += 1
	return s
}

func (this *Register) ReturnRegPop() string {
	this.RegisterPointer -= 1
	s := fmt.Sprintf("a%d", this.RegisterPointer)
	return s
}

func (this *Register) ReturnRegByName(name string) int64 {
	registerMap := map[string]int{
		"a1": 0b00010000,
		"a2": 0b00100000,
		"a3": 0b00110000,
		"c1": 0b01000000,
		"f1": 0b01010000,
	}
	return int64(registerMap[name])
}

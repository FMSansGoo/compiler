package code_gen

// pa: 0b00000000,
// a1: 0b00010000,
// a2: 0b00100000,
// a3: 0b00110000,
// c1: 0b01000000,
// f1: 0b01010000,
type AsmRegister string

const (
	A1 AsmRegister = "a1"
	A2 AsmRegister = "a2"
	A3 AsmRegister = "a3"
	C1 AsmRegister = "c1"
	F1 AsmRegister = "f1"
	PA AsmRegister = "pa"
)

var AsmRegisterToNum = map[AsmRegister]int{
	PA: 0b00000000,
	A1: 0b00010000,
	A2: 0b00100000,
	A3: 0b00110000,
	C1: 0b01000000,
	F1: 0b01010000,
}

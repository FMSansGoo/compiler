package asm_vm_stack_base

import "fmt"

// 纯数字函数的自有指令。由 Go 循环执行，不是 CPU 机器码。VM.JitEnabled 默认关闭。
type NumOp int

const (
	NumLoadConst NumOp = iota
	NumLoadLocal
	NumStoreLocal
	NumAdd
	NumSub
	NumMul
	NumDiv
	NumNeg
	NumPop
	NumReturn
)

func (op NumOp) String() string {
	names := [...]string{
		"LoadConst",
		"LoadLocal",
		"StoreLocal",
		"Add",
		"Sub",
		"Mul",
		"Div",
		"Neg",
		"Pop",
		"Return",
	}
	if int(op) < 0 || int(op) >= len(names) {
		return fmt.Sprintf("NumOp(%d)", int(op))
	}
	return names[op]
}

type NumInstr struct {
	Op  NumOp
	Arg int
}

type NumericCode struct {
	Instrs    []NumInstr
	Consts    []float64
	NumLocals int
}

func compileNumeric(fn *CompiledFunctionObject) *NumericCode {
	if fn == nil || len(fn.Instructions) == 0 {
		return nil
	}
	code := &NumericCode{NumLocals: fn.NumLocals}
	ins := fn.Instructions
	for i := 0; i < len(ins); {
		op := GetOpCodeFromValue(ins[i])
		i++
		switch op {
		case OpCodeConstant:
			index := int(ReadUint16(ins[i:]))
			i += 2
			if index >= len(fn.Constants) {
				return nil
			}
			number, ok := numberValue(fn.Constants[index])
			if !ok {
				return nil
			}
			code.Instrs = append(code.Instrs, NumInstr{Op: NumLoadConst, Arg: len(code.Consts)})
			code.Consts = append(code.Consts, number)
		case OpCodeAdd:
			code.Instrs = append(code.Instrs, NumInstr{Op: NumAdd})
		case OpCodeSub:
			code.Instrs = append(code.Instrs, NumInstr{Op: NumSub})
		case OpCodeMul:
			code.Instrs = append(code.Instrs, NumInstr{Op: NumMul})
		case OpCodeDiv:
			code.Instrs = append(code.Instrs, NumInstr{Op: NumDiv})
		case OpCodeMinus:
			code.Instrs = append(code.Instrs, NumInstr{Op: NumNeg})
		case OpCodePop:
			code.Instrs = append(code.Instrs, NumInstr{Op: NumPop})
		case OpCodeGetLocal:
			code.Instrs = append(code.Instrs, NumInstr{Op: NumLoadLocal, Arg: int(ReadUint8(ins[i:]))})
			i++
		case OpCodeSetLocal:
			code.Instrs = append(code.Instrs, NumInstr{Op: NumStoreLocal, Arg: int(ReadUint8(ins[i:]))})
			i++
		case OpCodeReturn:
			code.Instrs = append(code.Instrs, NumInstr{Op: NumReturn})
		default:
			return nil
		}
	}
	return code
}

func numberValue(obj Object) (float64, bool) {
	switch value := obj.(type) {
	case *NumberObject:
		return value.Value, true
	case NumberObject:
		return value.Value, true
	default:
		return 0, false
	}
}

func runNumeric(code *NumericCode, args []float64) (float64, error) {
	slots := make([]float64, code.NumLocals+len(code.Instrs)+1)
	copy(slots, args)
	sp := code.NumLocals
	for _, instr := range code.Instrs {
		switch instr.Op {
		case NumLoadConst:
			slots[sp] = code.Consts[instr.Arg]
			sp++
		case NumLoadLocal:
			slots[sp] = slots[instr.Arg]
			sp++
		case NumStoreLocal:
			sp--
			slots[instr.Arg] = slots[sp]
		case NumAdd:
			sp--
			slots[sp-1] += slots[sp]
		case NumSub:
			sp--
			slots[sp-1] -= slots[sp]
		case NumMul:
			sp--
			slots[sp-1] *= slots[sp]
		case NumDiv:
			sp--
			slots[sp-1] /= slots[sp]
		case NumNeg:
			slots[sp-1] = -slots[sp-1]
		case NumPop:
			sp--
		case NumReturn:
			sp--
			return slots[sp], nil
		default:
			return 0, fmt.Errorf("unknown numeric op %s", instr.Op)
		}
	}
	return 0, fmt.Errorf("numeric code has no return")
}

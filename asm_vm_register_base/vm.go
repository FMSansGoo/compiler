package asm_vm_register_base

import "fmt"

// Machine 解释 Assembler 产出的 []int64。
// 第一步只执行 set2、加减乘除和 halt。跳转、比较、函数和内存访问还没跑。
type Machine struct {
	Memory []int64
	PC     int
	Regs   map[int64]int64
}

func Run(memory []int64) (*Machine, error) {
	machine := &Machine{
		Memory: memory,
		Regs:   map[int64]int64{},
	}
	for machine.PC < len(machine.Memory) {
		op := machine.Memory[machine.PC]
		switch op {
		case InstructionHalt.Value():
			machine.PC++
			return machine, nil
		case InstructionSet2.Value():
			dest, value, err := machine.readRegImm()
			if err != nil {
				return nil, err
			}
			machine.Regs[dest] = value
		case InstructionAdd2.Value():
			if err := machine.arith(func(left, right int64) (int64, error) {
				return left + right, nil
			}); err != nil {
				return nil, err
			}
		case InstructionSubtract2.Value():
			if err := machine.arith(func(left, right int64) (int64, error) {
				return left - right, nil
			}); err != nil {
				return nil, err
			}
		case InstructionMultiply2.Value():
			if err := machine.arith(func(left, right int64) (int64, error) {
				return left * right, nil
			}); err != nil {
				return nil, err
			}
		case InstructionDiv2.Value():
			if err := machine.arith(func(left, right int64) (int64, error) {
				if right == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				return left / right, nil
			}); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported register opcode %d at %d", op, machine.PC)
		}
	}
	return nil, fmt.Errorf("register code has no halt")
}

func (m *Machine) readRegImm() (int64, int64, error) {
	if m.PC+2 >= len(m.Memory) {
		return 0, 0, fmt.Errorf("truncated set2 at %d", m.PC)
	}
	dest := m.Memory[m.PC+1]
	value := m.Memory[m.PC+2]
	m.PC += 3
	return dest, value, nil
}

func (m *Machine) arith(op func(left, right int64) (int64, error)) error {
	if m.PC+3 >= len(m.Memory) {
		return fmt.Errorf("truncated arithmetic at %d", m.PC)
	}
	left := m.Regs[m.Memory[m.PC+1]]
	right := m.Regs[m.Memory[m.PC+2]]
	dest := m.Memory[m.PC+3]
	value, err := op(left, right)
	if err != nil {
		return err
	}
	m.Regs[dest] = value
	m.PC += 4
	return nil
}

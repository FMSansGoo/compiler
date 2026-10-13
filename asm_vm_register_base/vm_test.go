package asm_vm_register_base

import "testing"

func TestRunArithmetic(t *testing.T) {
	regs := NewRegister()
	cases := []struct {
		asm  string
		reg  string
		want int64
	}{
		{"set2 a1 1\nset2 a2 2\nadd2 a1 a2 a3\nhalt\n", "a3", 3},
		{"set2 a1 5\nset2 a2 2\nsub2 a1 a2 a3\nhalt\n", "a3", 3},
		{"set2 a1 3\nset2 a2 4\nmul2 a1 a2 a3\nhalt\n", "a3", 12},
		{"set2 a1 8\nset2 a2 2\ndiv2 a1 a2 a3\nhalt\n", "a3", 4},
	}
	for _, item := range cases {
		assembler := NewAssembler()
		assembler.AddAsm(item.asm)
		machine, err := Run(assembler.Compile())
		if err != nil {
			t.Fatalf("%s: %v", item.asm, err)
		}
		got := machine.Regs[regs.ReturnRegByName(item.reg)]
		if got != item.want {
			t.Fatalf("%s: got %d, want %d", item.asm, got, item.want)
		}
	}
}

func TestRunRejectsMissingHalt(t *testing.T) {
	assembler := NewAssembler()
	assembler.AddAsm("set2 a1 1\n")
	if _, err := Run(assembler.Compile()); err == nil {
		t.Fatal("expected missing halt")
	}
}

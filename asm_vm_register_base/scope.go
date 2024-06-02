package asm_vm_register_base

type VariableInfo struct {
	Name    string `json:"name"`
	Address int64  `json:"address"`
	InStack bool   `json:"in_stack"`
}

type SymbolTable struct {
	Table  map[string]VariableInfo
	Parent *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{Table: make(map[string]VariableInfo), Parent: nil}
}

func (s *SymbolTable) SetParent(parent *SymbolTable) {
	if parent != nil && s.Parent == nil {
		s.Parent = parent
		return
	}
}

// 向作用域中添加某一变量的偏移量
func (s *SymbolTable) AddVariableInfo(name string, Address int64, InStack bool) {
	s.Table[name] = VariableInfo{
		Name:    name,
		Address: Address,
		InStack: InStack,
	}
}

// 查找符号
func (s *SymbolTable) LookupVariableInfo(name string) (VariableInfo, bool) {
	signature, ok := s.Table[name]
	if ok {
		return signature, true
	}

	if s.Parent != nil {
		return s.Parent.LookupVariableInfo(name)
	}

	return VariableInfo{}, false
}

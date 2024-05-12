package main

// 符号表项
type Symbol struct {
	VarType   string    // const var
	Name      string    // 变量名
	Value     ValueType // 变量类型
	ExtraInfo ValueType // 变量额外信息（现在只用在1.函数的返回值类型）
}

// 作用域结构
type Scope struct {
	Parent    *Scope
	ScopeName string
	Table     map[string]Symbol
	SubScopes []*Scope
}

// 创建新作用域
func NewScope() *Scope {
	return &Scope{
		Parent:    nil,
		Table:     make(map[string]Symbol),
		SubScopes: make([]*Scope, 0),
	}
}

func (s *Scope) SetParent(parent *Scope) {
	if parent != nil && s.Parent == nil {
		s.Parent = parent
		s.Parent.SubScopes = append(s.Parent.SubScopes, s)
		return
	}
}

// 向作用域中添加符号
func (s *Scope) AddSymbol(varType string, name string, value ValueType, extraInfo ValueType) {
	// func or class 作用域中已经存在该符号
	if symbol, ok := s.Table[name]; ok && symbol.Name == name && (symbol.Value == ValueTypeFunctionExpression || symbol.Value == ValueTypeClassExpression) {
		logError("not allow the same function name or class name", name, s.ScopeName)
	}

	s.Table[name] = Symbol{VarType: varType, Name: name, Value: value, ExtraInfo: extraInfo}
}

// 查找符号
func (s *Scope) LookupSymbol(name string) (Symbol, bool) {
	symbol, ok := s.Table[name]
	if ok {
		return symbol, true
	}

	if s.Parent != nil {
		return s.Parent.LookupSymbol(name)
	}

	return Symbol{}, false
}

func (s *Scope) LookupScope(name string) (*Scope, bool) {
	_, ok := s.Table[name]
	if ok {
		return s, true
	}

	if s.Parent != nil {
		return s.Parent.LookupScope(name)
	}

	return nil, false
}

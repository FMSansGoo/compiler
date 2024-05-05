package main

// 符号表项
type Symbol struct {
	VarType string
	Name    string
	Value   ValueType
}

// 作用域结构
type Scope struct {
	Parent    *Scope
	ScopeName string
	Table     map[string]Symbol
}

// 创建新作用域
func NewScope() *Scope {
	return &Scope{
		Parent: nil,
		Table:  make(map[string]Symbol),
	}
}

func (s *Scope) SetParent(parent *Scope) {
	if parent != nil && s.Parent == nil {
		s.Parent = parent
		return
	}

}

// 向作用域中添加符号
func (s *Scope) AddSymbol(varType string, name string, value ValueType) {
	// 作用域中已经存在该符号
	if symbol, ok := s.Table[name]; ok && symbol.Name == name && symbol.Value == ValueTypeFunctionExpression {
		logError("not allow the same function", name, s.ScopeName)
	}
	s.Table[name] = Symbol{VarType: varType, Name: name, Value: value}
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

// demo
//func main() {
// 创建全局作用域
//globalScope := NewScope(nil)
//
//// 添加全局符号 test
//testFunc := func() interface{} {
//	return 1
//}
//globalScope.AddSymbol("test", testFunc)
//
//// 添加全局符号 main
//mainFunc := func() {
//	localScope := NewScope(globalScope)
//
//	// 在 main 函数作用域中添加符号 a 和 b
//	localScope.AddSymbol("a", 1)
//	localScope.AddSymbol("b", 1)
//
//	// 执行赋值语句 a = test()
//	symbol, ok := localScope.LookupSymbol("a")
//	if ok {
//		if test, ok := globalScope.LookupSymbol("test"); ok {
//			result := test.Value.(func() interface{})()
//			symbol.Value = result
//		}
//	}
//}
//
//// 执行 main 函数
//mainFunc()
//
//// 查找符号 a 和 b
//symbol, ok := globalScope.LookupSymbol("a")
//if ok {
//	fmt.Println("a =", symbol.Value)
//} else {
//	fmt.Println("a 未定义")
//}
//
//symbol, ok = globalScope.LookupSymbol("b")
//if ok {
//	fmt.Println("b =", symbol.Value)
//} else {
//	fmt.Println("b 未定义")
//}
//}

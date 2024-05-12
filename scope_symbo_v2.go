package main

import "fmt"

type ScopeV2 struct {
	Table  map[string]Signature
	Parent *ScopeV2
}

func NewScopeV2() *ScopeV2 {
	return &ScopeV2{Table: make(map[string]Signature), Parent: nil}
}

func (s *ScopeV2) SetParent(parent *ScopeV2) {
	if parent != nil && s.Parent == nil {
		s.Parent = parent
		return
	}
}

// 向作用域中添加某一变量的返回值
func (s *ScopeV2) AddSignature(name string, ReturnType AllType, isStatic bool, varType string) {
	// func or class 作用域中已经存在该符号
	if signature, ok := s.Table[name]; ok && signature.Name == name &&
		(signature.ReturnType.ValueType() == ClassType{}.ValueType() || signature.ReturnType.ValueType() == FunctionType{}.ValueType()) {
		s.LogNowScope()
		logError("not allow the same function name or class name", name)
	}
	s.Table[name] = Signature{
		Name:       name,
		ReturnType: ReturnType,
		IsStatic:   isStatic,
		VarType:    varType,
	}

}

// 查找符号
func (s *ScopeV2) LookupSignature(name string) (Signature, bool) {
	signature, ok := s.Table[name]
	if ok {
		return signature, true
	}

	if s.Parent != nil {
		return s.Parent.LookupSignature(name)
	}

	return Signature{}, false
}

// 输出当前作用域所有符号 + symbol
func (s *ScopeV2) LogNowScope() {
	str := fmt.Sprintf("\nScopeV2:")
	if len(s.Table) == 0 {
		str += fmt.Sprintf("Empty\n")
	}
	for name, signature := range s.Table {
		retType := signature.ReturnType.ValueType()
		funType := FunctionType{}.ValueType()
		classType := ClassType{}.ValueType()
		if retType == funType {
			str += fmt.Sprintf("(name: %s, func signature RetType: %+v signature:%+v)\n", name, signature.ReturnType.(FunctionType).ReturnType.ValueType(), signature)
		} else if retType == classType {
			str += fmt.Sprintf("(name: %s, class memberType: %+v signature:%+v)\n", name, signature.ReturnType.(ClassType).MemberSignatures, signature)
		} else {
			str += fmt.Sprintf("(name: %s, signature RetType: %+v signature:%+v)\n", name, signature.ReturnType.ValueType(), signature)
		}
	}
	fmt.Printf("%v\n", str)
}

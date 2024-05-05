package main

import "fmt"

/*
	语义分析
	前两次分别对类的生命和成员的声明进行解析并构建符号表（类型和成员），第三次再对方法体进行解析。这样就可以方便地处理不同顺序定义的问题。总的来说，三次遍历的任务是：
	1	第一遍：扫描所有class定义，检查有无重名的情况。
	2	第二遍：检查类的基类是否存在，检测是否循环继承；检查所有字段的类型以及是否重名；检查所有方法参数和返回值的类型以及是否重复定义（签名完全一致的情况）。
	3	第三遍：检查所有方法体中语句和表达式的语义。

*/

type SemanticAnalysis struct {
	Ast          Program
	CurrentScope *Scope
}

func NewSemanticAnalysis(program Program) *SemanticAnalysis {
	s := &SemanticAnalysis{
		Ast:          program,
		CurrentScope: NewScope(),
	}
	s.CurrentScope.ScopeName = "global"
	return s
}

func (this *SemanticAnalysis) visit() {
	if this.Ast.Type() == AstTypeProgram.Name() {
		this.visitProgram(this.Ast.Body)
	}
}

func (this *SemanticAnalysis) visitProgram(body []Node) {
	for _, item := range body {
		logInfo("visitProgram visit item", item.Type())
		switch item.Type() {
		// 变量定义
		case AstTypeVariableDeclaration.Name():
			this.visitVariableDeclaration(item)
		// 赋值
		case AstTypeAssignmentExpression.Name():
			this.visitAssignmentExpression(item)
		// 访问 block
		case AstTypeBlockStatement.Name():
			this.visitBlockStatement(item)
		}
		logInfo("visitProgram visit item after currentScope", this.CurrentScope)
	}
}

func (this *SemanticAnalysis) visitAssignmentExpression(node Node) {
	left := node.(AssignmentExpression).Left
	var variableName string
	switch left.Type() {
	case AstTypeIdentifier.Name():
		_, variableName, _ = this.visitIdentifier(left)
	}
	symbol, ok := this.CurrentScope.LookupSymbol(variableName)
	if ok && symbol.VarType == TokenTypeConst.Name() {
		logError("const variable cannot be reassigned", variableName)
		return
	}

	varType := TokenTypeVar.Name()
	valueType := ValueTypeError
	right := node.(AssignmentExpression).Right
	switch right.Type() {
	case AstTypeBinaryExpression.Name():
		valueType = this.visitBinaryExpression(right)
	case AstTypeIdentifier.Name():
		valueType, _, varType = this.visitIdentifier(right)
	case AstTypeNumberLiteral.Name():
		valueType = this.visitNumberLiteral(right)
	case AstTypeStringLiteral.Name():
		valueType, _ = this.visitStringLiteral(right)
	case AstTypeBooleanLiteral.Name():
		valueType = this.visitBooleanLiteral(right)
	case AstTypeArrayLiteral.Name():
		valueType = this.visitArrayLiteral(right)
	case AstTypeDictLiteral.Name():
		valueType = this.visitDictLiteral(right)
	case AstTypeNullLiteral.Name():
		valueType = this.visitNullLiteral(right)
	}
	this.CurrentScope.AddSymbol(varType, variableName, valueType)
	return
}

func (this *SemanticAnalysis) visitFunctionExpression(node Node) ValueType {
	// 进入函数就是一个新的作用域
	funcScope := NewScope()
	funcScope.SetParent(this.CurrentScope)
	this.CurrentScope = funcScope
	// 开始进入作用域
	params := node.(FunctionExpression).Params
	body := node.(FunctionExpression).Body
	logInfo("visitFunctionExpression", params, body)
	if body.Type() == AstTypeBlockStatement.Name() {
		this.visitBlockStatement(body)
		// 回滚到上一个作用域中
		this.CurrentScope = this.CurrentScope.Parent
		return ValueTypeFunctionExpression
	}

	return ValueTypeError
}

func (this *SemanticAnalysis) visitBlockStatement(node Node) {
	// 开始新的作用域
	logInfo("visitBlockStatement visit node", node.Type())
	for _, item := range node.(BlockStatement).Body {
		logInfo("visitBlockStatement visit item", item.Type())
		switch item.Type() {
		// 变量定义
		case AstTypeVariableDeclaration.Name():
			this.visitVariableDeclaration(item)
		// 赋值
		case AstTypeAssignmentExpression.Name():
			this.visitAssignmentExpression(item)
		}
	}
	logInfo("visitBlockStatement current Scope", this.CurrentScope)
	return
}
func (this *SemanticAnalysis) visitVariableDeclaration(node Node) {
	//type VariableDeclaration struct {
	//	Kind  string // kind属性
	//	Name  Node   // name属性
	//	Value Node   // value属性
	//}
	varType := node.(VariableDeclaration).Kind
	left := node.(VariableDeclaration).Name
	var variableName string
	switch left.Type() {
	case AstTypeIdentifier.Name():
		_, variableName, _ = this.visitIdentifier(left)
	}

	valueType := ValueTypeError
	right := node.(VariableDeclaration).Value
	switch right.Type() {
	// 访问函数
	case AstTypeFunctionExpression.Name():
		valueType = this.visitFunctionExpression(right)
	case AstTypeBinaryExpression.Name():
		valueType = this.visitBinaryExpression(right)
	case AstTypeNumberLiteral.Name():
		valueType = this.visitNumberLiteral(right)
	case AstTypeNullLiteral.Name():
		valueType = this.visitNullLiteral(right)
	case AstTypeStringLiteral.Name():
		valueType, _ = this.visitStringLiteral(right)
	case AstTypeArrayLiteral.Name():
		valueType = this.visitArrayLiteral(right)
	case AstTypeDictLiteral.Name():
		valueType = this.visitDictLiteral(right)
	case AstTypeIdentifier.Name():
		var varName string
		_, varName, _ = this.visitIdentifier(right)
		symbol, ok := this.CurrentScope.LookupSymbol(varName)
		if ok {
			valueType = symbol.Value
		} else {
			logError("undeclared variable", varName)
			return
		}
	}
	this.CurrentScope.AddSymbol(varType, variableName, valueType)

	fmt.Printf("this.CurrentScope: %+v\n", this.CurrentScope)
	return
}

func (this *SemanticAnalysis) visitArrayLiteral(node Node) ValueType {
	// 做一下限制，不能多种类型混合在数组里面一起
	if node.Type() == AstTypeArrayLiteral.Name() {
		if len(node.(ArrayLiteral).Values) > 0 {
			checkType := node.(ArrayLiteral).Values[0].Type()
			// 数组
			for _, item := range node.(ArrayLiteral).Values {
				if item.Type() != checkType {
					logError("array literal type error", checkType, item.Type())
					return ValueTypeError
				}
			}
		}
		// 这里
		return ValueTypeArrayLiteral
	}
	return ValueTypeError
}

func (this *SemanticAnalysis) visitBinaryExpression(node Node) ValueType {
	// BinaryExpression节点结构
	//type BinaryExpression struct {
	//	Operator string // operator属性
	//	Left     Node   // left属性
	//	Right    Node   // right属性
	//}
	if node.Type() == AstTypeBinaryExpression.Name() {
		// op + - * /
		switch node.(BinaryExpression).Operator {
		case "+":
			// left
			left := node.(BinaryExpression).Left
			leftValueType := ValueTypeError
			switch left.Type() {
			case AstTypeBinaryExpression.Name():
				leftValueType = this.visitBinaryExpression(left)
			case AstTypeNumberLiteral.Name():
				leftValueType = this.visitNumberLiteral(left)
			case AstTypeStringLiteral.Name():
				leftValueType, _ = this.visitStringLiteral(left)
			case AstTypeIdentifier.Name():
				leftValueType, _, _ = this.visitIdentifier(left)
			default:
				logError("左值类型错误", node.(BinaryExpression).Left)
				return ValueTypeError
			}

			// right
			right := node.(BinaryExpression).Right
			rightValueType := ValueTypeError
			switch right.Type() {
			case AstTypeBinaryExpression.Name():
				rightValueType = this.visitBinaryExpression(right)
			case AstTypeNumberLiteral.Name():
				rightValueType = this.visitNumberLiteral(right)
			case AstTypeStringLiteral.Name():
				rightValueType, _ = this.visitStringLiteral(right)
			case AstTypeIdentifier.Name():
				rightValueType, _, _ = this.visitIdentifier(right)
			default:
				logError("右值类型错误", node.(BinaryExpression).Right)
				return ValueTypeError
			}
			logInfo("visitBinaryExpression", leftValueType, rightValueType)

			isString := false

			switch leftValueType {
			case ValueTypeNumber:
			case ValueTypeString:
				isString = true
			default:
				logError("左值类型错误", node.(BinaryExpression).Left)
				return ValueTypeError
			}

			switch rightValueType {
			case ValueTypeNumber:
			case ValueTypeString:
				isString = true
			default:
				logError("右值类型错误", node.(BinaryExpression).Right)
				return ValueTypeError
			}

			if isString {
				return ValueTypeString
			}
			return ValueTypeNumber
		case "-":
			fallthrough
		case "*":
			fallthrough
		case "/":
			// left
			left := node.(BinaryExpression).Left
			leftValueType := ValueTypeError
			switch left.Type() {
			case AstTypeBinaryExpression.Name():
				leftValueType = this.visitBinaryExpression(left)
			case AstTypeNumberLiteral.Name():
				leftValueType = this.visitNumberLiteral(left)
			default:
				logError("左值类型错误", node.(BinaryExpression).Left)
				return ValueTypeError
			}

			// right
			right := node.(BinaryExpression).Right
			rightValueType := ValueTypeError
			switch right.Type() {
			case AstTypeBinaryExpression.Name():
				rightValueType = this.visitBinaryExpression(right)
			case AstTypeNumberLiteral.Name():
				rightValueType = this.visitNumberLiteral(right)
			default:
				logError("右值类型错误", node.(BinaryExpression).Right)
				return ValueTypeError
			}
			switch leftValueType {
			case ValueTypeNumber:
			default:
				logError("左值类型错误", node.(BinaryExpression).Left)
				return ValueTypeError
			}

			switch rightValueType {
			case ValueTypeNumber:
			default:
				logError("右值类型错误", node.(BinaryExpression).Right)
				return ValueTypeError
			}
			return ValueTypeNumber
		}
	}

	return ValueTypeError
}

func (this *SemanticAnalysis) visitDictLiteral(node Node) ValueType {
	// 做一下限制，不能多种类型混合在dict
	// key 也不能重复
	if node.Type() == AstTypeDictLiteral.Name() {
		if len(node.(DictLiteral).Values) <= 0 {
			return ValueTypeDictLiteral
		}
		keyNames := []string{}
		_, firstDvType, _ := this.visitPropertyAssignment(node.(DictLiteral).Values[0])
		for _, value := range node.(DictLiteral).Values {
			_, dvType, keyName := this.visitPropertyAssignment(value)
			if dvType != firstDvType {
				logError("dict literal value type error", firstDvType, dvType)
				return ValueTypeError
			}
			inSlice := InStringSlice(keyNames, keyName)
			if inSlice {
				logError("dict literal duplicate key", keyName)
				return ValueTypeError
			}
			keyNames = append(keyNames, keyName)
		}
		return ValueTypeDictLiteral
	}
	return ValueTypeError
}

func (this *SemanticAnalysis) visitPropertyAssignment(node Node) (valueType ValueType, kvType ValueType, keyName string) {
	if node.Type() == AstTypePropertyAssignment.Name() {
		_, keyName = this.visitStringLiteral(node.(PropertyAssignment).Key)

		kv := node.(PropertyAssignment).Value
		switch kv.Type() {
		case AstTypeFunctionExpression.Name():
			kvType = this.visitFunctionExpression(kv)
		case AstTypeBinaryExpression.Name():
			kvType = this.visitBinaryExpression(kv)
		case AstTypeNumberLiteral.Name():
			kvType = this.visitNumberLiteral(kv)
		case AstTypeNullLiteral.Name():
			kvType = this.visitNullLiteral(kv)
		case AstTypeStringLiteral.Name():
			kvType, _ = this.visitStringLiteral(kv)
		case AstTypeArrayLiteral.Name():
			kvType = this.visitArrayLiteral(kv)
		case AstTypeDictLiteral.Name():
			kvType = this.visitDictLiteral(kv)
		case AstTypeIdentifier.Name():
			var varName string
			_, varName, _ = this.visitIdentifier(kv)
			symbol, ok := this.CurrentScope.LookupSymbol(varName)
			if ok {
				kvType = symbol.Value
			} else {
				logError("undeclared variable", varName)
				return
			}
		}
		return ValueTypePropertyAssignment, kvType, keyName
	}
	return ValueTypeError, ValueTypeError, ""
}

func (this *SemanticAnalysis) visitStringLiteral(node Node) (valueType ValueType, value string) {
	if node.Type() == AstTypeStringLiteral.Name() {
		value = node.(StringLiteral).Value
		return ValueTypeString, value
	}
	return ValueTypeError, ""
}

func (this *SemanticAnalysis) visitBooleanLiteral(node Node) ValueType {
	if node.Type() == AstTypeBooleanLiteral.Name() {
		return ValueTypeBoolean
	}
	return ValueTypeError
}

func (this *SemanticAnalysis) visitNullLiteral(node Node) ValueType {
	if node.Type() == AstTypeNullLiteral.Name() {
		return ValueTypeNull
	}
	return ValueTypeError
}

func (this *SemanticAnalysis) visitNumberLiteral(node Node) ValueType {
	if node.Type() == AstTypeNumberLiteral.Name() {
		return ValueTypeNumber
	}
	return ValueTypeError
}

// 这里要做一下类型判断
func (this *SemanticAnalysis) visitIdentifier(node Node) (valueType ValueType, variableName string, varType string) {
	if node.Type() == AstTypeIdentifier.Name() {
		symbol, ok := this.CurrentScope.LookupSymbol(node.(Identifier).Value)
		if ok {
			return symbol.Value, symbol.Name, symbol.VarType
		} else {
			return ValueTypeIdentifier, node.(Identifier).Value, ""
		}
	}
	return ValueTypeError, "", ""
}

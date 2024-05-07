package main

import "fmt"

/*
	todo 还没做
	语义分析
	前两次分别对类的生命和成员的声明进行解析并构建符号表（类型和成员），第2次再对方法体进行解析。这样就可以方便地处理不同顺序定义的问题。总的来说，2次遍历的任务是：
	1	第一遍：扫描所有class定义，检查有无重名的情况。检查类的基类是否存在，检测是否循环继承；检查所有字段的类型以及是否重名；检查所有方法参数和返回值的类型以及是否重复定义（签名完全一致的情况）。
	3	第二遍：检查所有方法体中语句和表达式的语义。

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
		// 访问 if
		case AstTypeIfStatement.Name():
			this.visitIfStatement(item)
		// 访问 while
		case AstTypeWhileStatement.Name():
			this.visitWhileStatement(item)
		// 访问 for
		case AstTypeForStatement.Name():
			this.visitForStatement(item)
		// 访问 block
		case AstTypeBlockStatement.Name():
			this.visitBlockStatement(item)
		// 访问 class
		case AstTypeClassExpression.Name():
			this.visitClassExpression(item)
		}
		logInfo("visitProgram visit item after currentScope", this.CurrentScope)
	}
}

func (this *SemanticAnalysis) visitClassExpression(node Node) ValueType {
	// 1.检查 super 的 class 是否存在
	logInfo("visitClassExpression", node.Type())
	if node.Type() == AstTypeClassExpression.Name() {

		classNameExp := node.(ClassExpression).Name
		className := ""
		classType := ValueTypeError
		classType, className, _ = this.visitIdentifier(classNameExp)

		// 表示现在找不到这个类
		if classType == ValueTypeIdentifier {
			classType = ValueTypeClassExpression
		}

		superClass := node.(ClassExpression).SuperClass
		superClassName := ""
		superClassType := ValueTypeClassExpression
		// 1.检查 super 的 class 是否存在
		if superClass != nil {
			logInfo("visitClassExpression superClass", superClass)
			superClassType, superClassName, _ = this.visitIdentifier(superClass)
			logInfo("visitClassExpression superClass", superClassType, superClassName)

			// 表示现在找不到这个类
			if superClassType == ValueTypeIdentifier {
				logError("super class not found", superClassName)
				return ValueTypeError
			}
		}
		// 现在开始处理这个类
		classBody := node.(ClassExpression).Body
		logInfo("visitClassExpression classBody", classBody.Type())
		if classBody.Type() == AstTypeClassBodyStatement.Name() {
			this.visitClassBodyStatement(classBody)
		}
		this.CurrentScope.AddSymbol("const", className, classType, ValueTypeError)
	}
	return ValueTypeError
}

func (this *SemanticAnalysis) visitClassBodyStatement(node Node) {
	// 开始新的作用域
	// 进入 block 就是一个新的作用域
	funcScope := NewScope()
	funcScope.SetParent(this.CurrentScope)
	this.CurrentScope = funcScope
	logInfo("visitClassBodyStatement visit node", node.Type())
	for _, item := range node.(ClassBodyStatement).Body {
		logInfo("visitClassBodyStatement visit item", item.Type())
		switch item.Type() {
		// 变量定义
		case AstTypeClassVariableDeclaration.Name():
			this.visitClassVariableDeclaration(item)
		// 赋值
		case AstTypeAssignmentExpression.Name():
			this.visitAssignmentExpression(item)
		}
	}
	// 出来 block 了，就要回滚到上一个作用域中
	this.CurrentScope = this.CurrentScope.Parent
	logInfo("visitClassBodyStatement current Scope", this.CurrentScope)
	return
}

func (this *SemanticAnalysis) visitClassVariableDeclaration(node Node) {
	//type VariableDeclaration struct {
	//	Kind  string // kind属性
	//	Name  Node   // name属性
	//	Value Node   // value属性
	//}
	varType := node.(ClassVariableDeclaration).Kind
	left := node.(ClassVariableDeclaration).Name
	var variableName string
	switch left.Type() {
	case AstTypeMemberExpression.Name():
		// 这里来做 . 语法和 this \ cls 的判别
		if left.(MemberExpression).ElementType == "dot" {
			if left.(MemberExpression).Object.Type() == AstTypeIdentifier.Name() {
				//todo 这里暂时用 property 当做变量名，其实应该用类变量名
				_, variableName, _ = this.visitIdentifier(left.(MemberExpression).Property)
				if variableName == TokenTypeCls.Name() || variableName == TokenTypeThis.Name() {

				} else {
					logError("invalid member expression", left)
					return
				}
			}
		} else {
			logError("invalid member expression", left)
			return
		}
	case AstTypeIdentifier.Name():
		_, variableName, _ = this.visitIdentifier(left)
	}

	valueType := ValueTypeError
	right := node.(ClassVariableDeclaration).Value
	switch right.Type() {
	// 访问函数
	case AstTypeFunctionExpression.Name():
		valueType, _ = this.visitFunctionExpression(right)
	case AstTypeBinaryExpression.Name():
		valueType = this.visitBinaryExpression(right)
	case AstTypeNumberLiteral.Name():
		valueType = this.visitNumberLiteral(right)
	case AstTypeNullLiteral.Name():
		valueType = this.visitNullLiteral(right)
	case AstTypeBooleanLiteral.Name():
		valueType = this.visitBooleanLiteral(right)
	case AstTypeStringLiteral.Name():
		valueType, _ = this.visitStringLiteral(right)
	case AstTypeArrayLiteral.Name():
		valueType = this.visitArrayLiteral(right)
	case AstTypeUnaryExpression.Name():
		valueType = this.visitUnaryExpression(right)
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
	this.CurrentScope.AddSymbol(varType, variableName, valueType, ValueTypeError)

	fmt.Printf("this.CurrentScope: %+v\n", this.CurrentScope)
	return
}

func (this *SemanticAnalysis) visitAssignmentExpression(node Node) {
	left := node.(AssignmentExpression).Left
	var variableName string
	switch left.Type() {
	case AstTypeIdentifier.Name():
		_, variableName, _ = this.visitIdentifier(left)
	}
	symbol, ok := this.CurrentScope.LookupSymbol(variableName)
	// const 检查
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
	case AstTypeUnaryExpression.Name():
		valueType = this.visitUnaryExpression(right)
	}
	// 强类型检查，如果右边的值不是同一个类型就报错
	if ok && symbol.Value != valueType {
		logError("variable cannot be reassigned to another type", variableName)
		return
	}
	this.CurrentScope.AddSymbol(varType, variableName, valueType, ValueTypeError)
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

	extraInfoValueType := ValueTypeError
	switch right.Type() {
	// 访问函数，函数的返回类型
	case AstTypeFunctionExpression.Name():
		valueType, extraInfoValueType = this.visitFunctionExpression(right)
	case AstTypeBinaryExpression.Name():
		valueType = this.visitBinaryExpression(right)
	case AstTypeNumberLiteral.Name():
		valueType = this.visitNumberLiteral(right)
	case AstTypeNullLiteral.Name():
		valueType = this.visitNullLiteral(right)
	case AstTypeBooleanLiteral.Name():
		valueType = this.visitBooleanLiteral(right)
	case AstTypeStringLiteral.Name():
		valueType, _ = this.visitStringLiteral(right)
	case AstTypeArrayLiteral.Name():
		valueType = this.visitArrayLiteral(right)
	case AstTypeUnaryExpression.Name():
		valueType = this.visitUnaryExpression(right)
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
	this.CurrentScope.AddSymbol(varType, variableName, valueType, extraInfoValueType)

	fmt.Printf("this.CurrentScope: %+v\n", this.CurrentScope)
	return
}

func (this *SemanticAnalysis) visitFunctionExpression(node Node) (ValueType, ValueType) {
	params := node.(FunctionExpression).Params
	body := node.(FunctionExpression).Body
	logInfo("visitFunctionExpression", params, body)
	funcReturnType := ValueTypeNull
	if body.Type() == AstTypeBlockStatement.Name() {
		var hasReturn bool
		hasReturn, funcReturnType = this.visitBlockStatement(body)
		if hasReturn {
			return ValueTypeFunctionExpression, funcReturnType
		}
	}

	return ValueTypeError, ValueTypeError
}

func (this *SemanticAnalysis) visitBlockStatement(node Node) (hasReturn bool, returnValueType ValueType) {
	// 开始新的作用域
	// 进入 block 就是一个新的作用域
	funcScope := NewScope()
	funcScope.SetParent(this.CurrentScope)
	this.CurrentScope = funcScope

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
		//return
		case AstTypeReturnStatement.Name():
			hasReturn = true
			returnValueType = this.visitReturnStatement(item)
		}
	}
	logInfo("visitBlockStatement before current Scope", this.CurrentScope)
	// 出来 block 了，就要回滚到上一个作用域中
	this.CurrentScope = this.CurrentScope.Parent
	logInfo("visitBlockStatement after current Scope", this.CurrentScope)
	return
}

func (this *SemanticAnalysis) visitReturnStatement(node Node) ValueType {
	if node.Type() == AstTypeReturnStatement.Name() {
		v := node.(ReturnStatement).Value

		var rightType ValueType
		switch v.Type() {
		case AstTypeBinaryExpression.Name():
			rightType = this.visitBinaryExpression(v)
		case AstTypeNumberLiteral.Name():
			rightType = this.visitNumberLiteral(v)
		case AstTypeNullLiteral.Name():
			rightType = this.visitNullLiteral(v)
		case AstTypeBooleanLiteral.Name():
			rightType = this.visitBooleanLiteral(v)
		case AstTypeStringLiteral.Name():
			rightType, _ = this.visitStringLiteral(v)
		case AstTypeArrayLiteral.Name():
			rightType = this.visitArrayLiteral(v)
		case AstTypeUnaryExpression.Name():
			rightType = this.visitUnaryExpression(v)
		case AstTypeDictLiteral.Name():
			rightType = this.visitDictLiteral(v)
		case AstTypeIdentifier.Name():
			var varName string
			_, varName, _ = this.visitIdentifier(v)
			symbol, ok := this.CurrentScope.LookupSymbol(varName)
			if ok {
				rightType = symbol.Value
			} else {
				logError("undeclared variable", varName)
				return ValueTypeError
			}
		default:
			logError("not support return value type", v.Type())
			return ValueTypeError
		}
		logInfo("rightType", rightType)
		return rightType
	}
	return ValueTypeError
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
		// op
		switch node.(BinaryExpression).Operator {
		case "+":
			// 这里接受多个类型
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
			fallthrough
		case "+=":
			fallthrough
		case "-=":
			fallthrough
		case "*=":
			fallthrough
		case "/=":
			// left
			left := node.(BinaryExpression).Left
			leftValueType := ValueTypeError
			switch left.Type() {
			case AstTypeIdentifier.Name():
				leftValueType, _, _ = this.visitIdentifier(left)
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
		case ">":
			fallthrough
		case "<":
			fallthrough
		case ">=":
			fallthrough
		case "<=":
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

			if leftValueType != rightValueType {
				logError("类型不匹配", node.(BinaryExpression).Left, node.(BinaryExpression).Right)
				return ValueTypeError
			}
			return ValueTypeBoolean
		case "==":
			fallthrough
		case "!=":
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
			case AstTypeBooleanLiteral.Name():
				leftValueType = this.visitBooleanLiteral(left)
			case AstTypeNullLiteral.Name():
				leftValueType = this.visitNullLiteral(left)
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
			case AstTypeBooleanLiteral.Name():
				rightValueType = this.visitBooleanLiteral(right)
			case AstTypeNullLiteral.Name():
				rightValueType = this.visitNullLiteral(right)
			default:
				logError("右值类型错误", node.(BinaryExpression).Right)
				return ValueTypeError
			}
			logInfo("visitBinaryExpression", leftValueType, rightValueType)

			if leftValueType != rightValueType {
				logError("类型不匹配", node.(BinaryExpression).Left, node.(BinaryExpression).Right)
				return ValueTypeError
			}
			return ValueTypeBoolean
		case "and":
			fallthrough
		case "or":
			// left
			left := node.(BinaryExpression).Left
			leftValueType := ValueTypeError
			switch left.Type() {
			case AstTypeBinaryExpression.Name():
				leftValueType = this.visitBinaryExpression(left)
			case AstTypeIdentifier.Name():
				leftValueType, _, _ = this.visitIdentifier(left)
			case AstTypeBooleanLiteral.Name():
				leftValueType = this.visitBooleanLiteral(left)
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
			case AstTypeIdentifier.Name():
				rightValueType, _, _ = this.visitIdentifier(right)
			case AstTypeBooleanLiteral.Name():
				rightValueType = this.visitBooleanLiteral(right)
			default:
				logError("右值类型错误", node.(BinaryExpression).Right)
				return ValueTypeError
			}
			logInfo("in and not visitBinaryExpression", leftValueType, rightValueType)

			if (leftValueType != ValueTypeBoolean) && (rightValueType != ValueTypeBoolean) && (leftValueType != rightValueType) {
				logError("类型不匹配", node.(BinaryExpression).Left, node.(BinaryExpression).Right)
				return ValueTypeError
			}
			return ValueTypeBoolean
		default:
			logError("not support binary expression operator", node.(BinaryExpression).Operator)
			return ValueTypeError
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

func (this *SemanticAnalysis) visitPropertyAssignment(node Node) (valueType ValueType, vType ValueType, keyName string) {
	if node.Type() == AstTypePropertyAssignment.Name() {
		_, keyName = this.visitStringLiteral(node.(PropertyAssignment).Key)

		v := node.(PropertyAssignment).Value
		switch v.Type() {
		case AstTypeFunctionExpression.Name():
			vType, _ = this.visitFunctionExpression(v)
		case AstTypeBinaryExpression.Name():
			vType = this.visitBinaryExpression(v)
		case AstTypeNumberLiteral.Name():
			vType = this.visitNumberLiteral(v)
		case AstTypeNullLiteral.Name():
			vType = this.visitNullLiteral(v)
		case AstTypeStringLiteral.Name():
			vType, _ = this.visitStringLiteral(v)
		case AstTypeArrayLiteral.Name():
			vType = this.visitArrayLiteral(v)
		case AstTypeDictLiteral.Name():
			vType = this.visitDictLiteral(v)
		case AstTypeIdentifier.Name():
			var varName string
			_, varName, _ = this.visitIdentifier(v)
			symbol, ok := this.CurrentScope.LookupSymbol(varName)
			if ok {
				vType = symbol.Value
			} else {
				logError("undeclared variable", varName)
				return
			}
		}
		return ValueTypePropertyAssignment, vType, keyName
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

// 变量，这里要做一下类型判断
func (this *SemanticAnalysis) visitIdentifier(node Node) (valueType ValueType, variableName string, varType string) {
	if node.Type() == AstTypeIdentifier.Name() {
		symbol, ok := this.CurrentScope.LookupSymbol(node.(Identifier).Value)
		if ok {
			return symbol.Value, symbol.Name, symbol.VarType
		} else {
			return ValueTypeIdentifier, node.(Identifier).Value, "var"
		}
	}
	return ValueTypeError, "", ""
}

// not
func (this *SemanticAnalysis) visitUnaryExpression(node Node) ValueType {
	if node.Type() == AstTypeUnaryExpression.Name() {
		op := node.(UnaryExpression).Operator
		switch op {
		case "not":
			vValueType := ValueTypeError
			// not 后面只能加 identifier / bool
			v := node.(UnaryExpression).Value
			switch v.Type() {
			case AstTypeBinaryExpression.Name():
				vValueType = this.visitBinaryExpression(v)
			case AstTypeIdentifier.Name():
				vValueType, _, _ = this.visitIdentifier(v)
			case AstTypeBooleanLiteral.Name():
				vValueType = this.visitBooleanLiteral(v)
			}
			logInfo("not value type", vValueType)
			if vValueType != ValueTypeBoolean {
				logError("not value type error", ValueTypeBoolean, vValueType)
				return ValueTypeError
			}
			return ValueTypeBoolean
		}

		return ValueTypeBoolean
	}
	return ValueTypeError
}

func (this *SemanticAnalysis) visitIfStatement(node Node) ValueType {
	if node.Type() == AstTypeIfStatement.Name() {
		condition := node.(IfStatement).Condition
		conditionType := ValueTypeError
		switch condition.Type() {
		case AstTypeBinaryExpression.Name():
			conditionType = this.visitBinaryExpression(condition)
		case AstTypeIdentifier.Name():
			conditionType, _, _ = this.visitIdentifier(condition)
		case AstTypeBooleanLiteral.Name():
			conditionType = this.visitBooleanLiteral(condition)
		default:
			logError("not support if statement condition type", conditionType)
			return ValueTypeError
		}
		logInfo("if condition type", conditionType)

		if conditionType != ValueTypeBoolean {
			logError("if condition type error", ValueTypeBoolean, conditionType)
			return ValueTypeError
		}

		consequent := node.(IfStatement).Consequent
		if consequent.Type() == AstTypeBlockStatement.Name() {
			this.visitBlockStatement(consequent)
		}

		alternate := node.(IfStatement).Alternate
		if alternate != nil {
			if alternate.Type() == AstTypeBlockStatement.Name() {
				this.visitBlockStatement(alternate)
			}
		}

		return ValueTypeIfStatement
	}
	return ValueTypeError
}
func (this *SemanticAnalysis) visitWhileStatement(node Node) ValueType {
	if node.Type() == AstTypeWhileStatement.Name() {
		condition := node.(WhileStatement).Condition
		conditionType := ValueTypeError
		switch condition.Type() {
		case AstTypeBinaryExpression.Name():
			conditionType = this.visitBinaryExpression(condition)
		case AstTypeIdentifier.Name():
			conditionType, _, _ = this.visitIdentifier(condition)
		case AstTypeBooleanLiteral.Name():
			conditionType = this.visitBooleanLiteral(condition)
		default:
			logError("not support while statement condition type", conditionType)
			return ValueTypeError
		}
		logInfo("while condition type", conditionType)

		if conditionType != ValueTypeBoolean {
			logError("while condition type error", ValueTypeBoolean, conditionType)
			return ValueTypeError
		}

		body := node.(WhileStatement).Body
		if body.Type() == AstTypeBlockStatement.Name() {
			this.visitBlockStatement(body)
		}

		return ValueTypeWhileStatement
	}
	return ValueTypeError
}

func (this *SemanticAnalysis) visitForStatement(node Node) ValueType {
	if node.Type() == AstTypeForStatement.Name() {
		init := node.(ForStatement).Init
		initType := ValueTypeError
		switch init.Type() {
		case AstTypeVariableDeclaration.Name():
			this.visitVariableDeclaration(init)
		case AstTypeAssignmentExpression.Name():
			this.visitAssignmentExpression(init)
		case AstTypeBooleanLiteral.Name():
			initType = this.visitBooleanLiteral(init)
		default:
			logError("not support for statement init type", init)
			return ValueTypeError
		}
		logInfo("for init type", initType)

		testExp := node.(ForStatement).Test
		testType := ValueTypeError
		switch testExp.Type() {
		case AstTypeBinaryExpression.Name():
			testType = this.visitBinaryExpression(testExp)
		default:
			logError("not support for statement test type", testExp)
			return ValueTypeError
		}
		logInfo("for testType type", testType)

		update := node.(ForStatement).Update
		updateType := ValueTypeError
		switch update.Type() {
		case AstTypeBinaryExpression.Name():
			testType = this.visitBinaryExpression(update)
		default:
			logError("not support for statement update type", update)
			return ValueTypeError
		}
		logInfo("for updateType type", updateType)

		body := node.(ForStatement).Body
		if body.Type() == AstTypeBlockStatement.Name() {
			this.visitBlockStatement(body)
		}

		return ValueTypeError
	}
	return ValueTypeError
}

// Todo array dot 语法 [] .
func (this *SemanticAnalysis) visitMemberExpression(node Node) ValueType {
	if node.Type() == AstTypeMemberExpression.Name() {
		et := node.(MemberExpression).ElementType
		switch et {
		case "dot":

		case "array":
		}

		return ValueTypeError
	}
	return ValueTypeError
}

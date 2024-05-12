package main

import "fmt"

type SemanticAnalysisV2 struct {
	Ast          Program
	CurrentScope *ScopeV2
}

func NewSemanticAnalysisV2(program Program) *SemanticAnalysisV2 {
	s := &SemanticAnalysisV2{
		Ast:          program,
		CurrentScope: NewScopeV2(),
	}
	return s
}

func (this *SemanticAnalysisV2) visit() {
	if this.Ast.Type() != AstTypeProgram.Name() {
		return
	}
	this.visitProgram(this.Ast.Body)

}

func (this *SemanticAnalysisV2) visitProgram(body []Node) {
	for _, item := range body {
		logInfo("visitProgram visit item", item.Type())
		switch item.Type() {
		// 变量定义
		case AstTypeVariableDeclaration.Name():
			this.visitVariableDeclaration(item)
		//// 赋值
		case AstTypeAssignmentExpression.Name():
			this.visitAssignmentExpression(item)
		//// 访问 if
		case AstTypeIfStatement.Name():
			this.visitIfStatement(item)
		// 访问 while
		case AstTypeWhileStatement.Name():
			this.visitWhileStatement(item)
		//// 访问 for
		case AstTypeForStatement.Name():
			this.visitForStatement(item)
		//// 访问 block
		case AstTypeBlockStatement.Name():
			this.visitBlockStatement(item)
		//// 访问 class
		//case AstTypeClassExpression.Name():
		//	this.visitClassExpression(item)
		case AstTypeCallExpression.Name():
			this.visitCallExpression(item)
		default:
			logError("visitProgram visit item default", item.Type())
		}
		logInfo("visitProgram visit item after currentScope")
		this.CurrentScope.LogNowScope()
	}
}

// 变量定义
func (this *SemanticAnalysisV2) visitVariableDeclaration(node Node) {
	////type VariableDeclaration struct {
	////	Kind  string // kind属性
	////	Name  Node   // name属性
	////	Value Node   // value属性
	////}
	varType := node.(VariableDeclaration).Kind
	left := node.(VariableDeclaration).Name
	var variableName string
	switch left.Type() {
	case AstTypeIdentifier.Name():
		_, variableName, _ = this.visitIdentifier(left)
	default:
		logError("visitVariableDeclaration invalid left variable declaration", left)
		return
	}
	// 做一下限制，变量名不为空
	if len(variableName) == 0 {
		logError("visitVariableDeclaration invalid left variable declaration", left)
		return
	}
	//
	var valueType AllType
	right := node.(VariableDeclaration).Value

	switch right.Type() {
	// 访问函数，函数的返回类型
	case AstTypeFunctionExpression.Name():
		valueType = this.visitFunctionExpression(right)
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
		signature, ok := this.CurrentScope.LookupSignature(varName)
		if ok {
			valueType = signature.ReturnType
		} else {
			logError("undeclared variable", varName)
			return
		}
	case AstTypeCallExpression.Name():
		valueType, _ = this.visitCallExpression(right)
	default:
		logError("visitVariableDeclaration invalid right variable declaration", right)
		return
	}
	// 先不处理常量方法
	this.CurrentScope.AddSignature(variableName, valueType, false, varType)

	fmt.Printf("visitVariableDeclaration this.CurrentScope")
	this.CurrentScope.LogNowScope()
	return
}

// 赋值
func (this *SemanticAnalysisV2) visitAssignmentExpression(node Node) {
	left := node.(AssignmentExpression).Left
	logInfo("visitClassVariableDeclaration visitAssignmentExpression", node.(AssignmentExpression))
	var variableName string
	switch left.Type() {
	case AstTypeIdentifier.Name():
		_, variableName, _ = this.visitIdentifier(left)
	//case AstTypeMemberExpression.Name():
	//	_, variableName = this.visitMemberExpression(left)
	default:
		logError("invalid assignment expression", left)
		return
	}
	fmt.Printf("visitAssignmentExpression variableName %v left:%v\n", variableName, left)
	varSignature, ok := this.CurrentScope.LookupSignature(variableName)
	// const 检查
	if ok && varSignature.VarType == TokenTypeConst.Name() {
		logError("const variable cannot be reassigned", variableName)
		return
	}

	varType := TokenTypeVar.Name()
	var valueType AllType
	valueType = UnKnownType{}
	//valueType := ValueTypeIdentifier
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
	case AstTypeCallExpression.Name():
		valueType, _ = this.visitCallExpression(right)
	default:
		logError("invalid assignment expression", right)
	}
	unknownValueType := UnKnownType{}
	// 强类型检查，如果右边的值不是同一个类型就报错
	if ok && varSignature.ReturnType != valueType && valueType.ValueType() != unknownValueType.ValueType() {
		logError("variable cannot be reassigned to another type", variableName, varSignature, valueType)
		return
	}
	// 如果判断不出类型，就用原有的类型
	if ok && valueType.ValueType() == unknownValueType.ValueType() {
		valueType = varSignature.ReturnType
	}
	this.CurrentScope.AddSignature(variableName, valueType, false, varType)
	//this.CurrentScope.AddClass(variableName, classScope)
	return
}

func (this *SemanticAnalysisV2) visitFunctionExpression(node Node) (functionType AllType) {
	if node.Type() != AstTypeFunctionExpression.Name() {
		return UnKnownType{}
	}
	params := node.(FunctionExpression).Params

	// 开始新的作用域
	funcScope := NewScopeV2()
	funcScope.SetParent(this.CurrentScope)
	this.CurrentScope = funcScope
	// 这里用 defer 主要是怕提前 return，作用域没回滚
	defer func() {
		this.CurrentScope = this.CurrentScope.Parent
	}()

	signatures := make([]Signature, 0)
	for _, param := range params {
		if param.Type() != AstTypeIdentifier.Name() {
			logError("param must be identifier", param.Type())
			return UnKnownType{}
		}
		valueType, variableName, varType := this.visitIdentifier(param)
		s := Signature{
			Name:       variableName,
			ReturnType: valueType,
			IsStatic:   false,
			VarType:    varType,
		}
		signatures = append(signatures, s)
		this.CurrentScope.AddSignature(variableName, valueType, false, varType)
	}
	body := node.(FunctionExpression).Body
	logInfo("visitFunctionExpression", params, body)

	var funcReturnType AllType
	funcReturnType = VoidType{}
	if body.Type() == AstTypeBlockStatement.Name() {
		funcReturnType = this.visitBlockStatement(body)
	}

	return FunctionType{
		Params:     signatures,
		ReturnType: funcReturnType,
	}
}

// 表达式
func (this *SemanticAnalysisV2) visitBinaryExpression(node Node) AllType {
	// BinaryExpression节点结构
	//type BinaryExpression struct {
	//	Operator string // operator属性
	//	Left     Node   // left属性
	//	Right    Node   // right属性
	//}
	if node.Type() != AstTypeBinaryExpression.Name() {
		return UnKnownType{}
	}
	// op
	switch node.(BinaryExpression).Operator {
	case "+":
		// 这里接受多个类型
		left := node.(BinaryExpression).Left
		var leftValueType AllType
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
			return UnKnownType{}
		}

		// right
		right := node.(BinaryExpression).Right
		var rightValueType AllType
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
			return UnKnownType{}
		}
		logInfo("visitBinaryExpression", leftValueType, rightValueType)

		isString := false

		switch leftValueType {
		case NumberType{}:
		case StringType{}:
			isString = true
		default:
			logError("左值类型错误", node.(BinaryExpression).Left)
			return UnKnownType{}
		}

		switch rightValueType {
		case NumberType{}:
		case StringType{}:
			isString = true
		default:
			logError("右值类型错误", node.(BinaryExpression).Right)
			return UnKnownType{}
		}

		if isString {
			return StringType{}
		}
		return NumberType{}
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
		var leftValueType AllType
		switch left.Type() {
		case AstTypeIdentifier.Name():
			leftValueType, _, _ = this.visitIdentifier(left)
		case AstTypeBinaryExpression.Name():
			leftValueType = this.visitBinaryExpression(left)
		case AstTypeNumberLiteral.Name():
			leftValueType = this.visitNumberLiteral(left)
		default:
			logError("左值类型错误", node.(BinaryExpression).Left)
			return UnKnownType{}
		}

		// right
		right := node.(BinaryExpression).Right
		var rightValueType AllType
		switch right.Type() {
		case AstTypeBinaryExpression.Name():
			rightValueType = this.visitBinaryExpression(right)
		case AstTypeNumberLiteral.Name():
			rightValueType = this.visitNumberLiteral(right)
		default:
			logError("右值类型错误", node.(BinaryExpression).Right)
			return UnKnownType{}
		}
		switch leftValueType {
		case NumberType{}:
		default:
			logError("左值类型错误", node.(BinaryExpression).Left)
			return UnKnownType{}
		}

		switch rightValueType {
		case NumberType{}:
		default:
			logError("右值类型错误", node.(BinaryExpression).Right)
			return UnKnownType{}
		}
		return NumberType{}
	case ">":
		fallthrough
	case "<":
		fallthrough
	case ">=":
		fallthrough
	case "<=":
		// left
		left := node.(BinaryExpression).Left
		var leftValueType AllType
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
			return UnKnownType{}
		}

		// right
		right := node.(BinaryExpression).Right
		var rightValueType AllType
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
			return UnKnownType{}
		}
		logInfo("visitBinaryExpression", leftValueType, rightValueType)

		if leftValueType != rightValueType {
			logError("类型不匹配", node.(BinaryExpression).Left, node.(BinaryExpression).Right)
			return UnKnownType{}
		}
		return BooleanType{}
	case "==":
		fallthrough
	case "!=":
		// left
		left := node.(BinaryExpression).Left
		var leftValueType AllType
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
			return UnKnownType{}
		}

		// right
		right := node.(BinaryExpression).Right
		var rightValueType AllType
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
			return UnKnownType{}
		}
		logInfo("visitBinaryExpression", leftValueType, rightValueType)

		if leftValueType != rightValueType {
			logError("类型不匹配", node.(BinaryExpression).Left, node.(BinaryExpression).Right)
			return UnKnownType{}
		}
		return BooleanType{}
	case "and":
		fallthrough
	case "or":
		// left
		left := node.(BinaryExpression).Left
		var leftValueType AllType
		switch left.Type() {
		case AstTypeBinaryExpression.Name():
			leftValueType = this.visitBinaryExpression(left)
		case AstTypeIdentifier.Name():
			leftValueType, _, _ = this.visitIdentifier(left)
		case AstTypeBooleanLiteral.Name():
			leftValueType = this.visitBooleanLiteral(left)
		default:
			logError("左值类型错误", node.(BinaryExpression).Left)
			return UnKnownType{}
		}

		// right
		right := node.(BinaryExpression).Right
		var rightValueType AllType
		switch right.Type() {
		case AstTypeBinaryExpression.Name():
			rightValueType = this.visitBinaryExpression(right)
		case AstTypeIdentifier.Name():
			rightValueType, _, _ = this.visitIdentifier(right)
		case AstTypeBooleanLiteral.Name():
			rightValueType = this.visitBooleanLiteral(right)
		default:
			logError("右值类型错误", node.(BinaryExpression).Right)
			return UnKnownType{}
		}
		logInfo("in and not visitBinaryExpression", leftValueType, rightValueType)

		if (leftValueType.ValueType() != BooleanType{}.ValueType()) && (rightValueType.ValueType() != BooleanType{}.ValueType()) {
			logError("类型不匹配", node.(BinaryExpression).Left, node.(BinaryExpression).Right)
			return UnKnownType{}
		}
		return BooleanType{}
	default:
		logError("not support binary expression operator", node.(BinaryExpression).Operator)
		return UnKnownType{}
	}
}

// not
func (this *SemanticAnalysisV2) visitUnaryExpression(node Node) AllType {
	if node.Type() != AstTypeUnaryExpression.Name() {
		return UnKnownType{}
	}
	op := node.(UnaryExpression).Operator
	switch op {
	case "not":
		var vValueType AllType
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
		boolType := BooleanType{}.ValueType()
		if vValueType.ValueType() != boolType {
			logError("not value type error", vValueType)
			return UnKnownType{}
		}
		return BooleanType{}
	default:
		logError("not support unary expression operator", node.(UnaryExpression).Operator)
		return UnKnownType{}
	}

}

func (this *SemanticAnalysisV2) visitIdentifier(node Node) (valueType AllType, variableName string, varType string) {
	if node.Type() != AstTypeIdentifier.Name() {
		return VoidType{}, "", varType
	}
	signature, ok := this.CurrentScope.LookupSignature(node.(Identifier).Value)
	if ok {
		return signature.ReturnType, signature.Name, varType
	}
	return UnKnownType{}, node.(Identifier).Value, varType
}

func (this *SemanticAnalysisV2) visitStringLiteral(node Node) (valueType AllType, value string) {
	if node.Type() != AstTypeStringLiteral.Name() {
		return VoidType{}, ""
	}
	value = node.(StringLiteral).Value
	return StringType{}, value
}

func (this *SemanticAnalysisV2) visitBooleanLiteral(node Node) AllType {
	if node.Type() != AstTypeBooleanLiteral.Name() {
		return VoidType{}
	}
	return BooleanType{}
}

func (this *SemanticAnalysisV2) visitNullLiteral(node Node) AllType {
	if node.Type() != AstTypeNullLiteral.Name() {
		return VoidType{}
	}
	return NullType{}
}

func (this *SemanticAnalysisV2) visitNumberLiteral(node Node) AllType {
	if node.Type() != AstTypeNumberLiteral.Name() {
		return VoidType{}
	}
	return NumberType{}
}

func (this *SemanticAnalysisV2) visitWhileStatement(node Node) AllType {
	if node.Type() != AstTypeWhileStatement.Name() {
		return UnKnownType{}
	}
	condition := node.(WhileStatement).Condition
	var conditionType AllType
	switch condition.Type() {
	case AstTypeBinaryExpression.Name():
		conditionType = this.visitBinaryExpression(condition)
	case AstTypeIdentifier.Name():
		conditionType, _, _ = this.visitIdentifier(condition)
	case AstTypeBooleanLiteral.Name():
		conditionType = this.visitBooleanLiteral(condition)
	default:
		logError("not support while statement condition type", conditionType)
		return UnKnownType{}
	}

	booleanTypeName := BooleanType{}.ValueType()
	if conditionType.ValueType() != booleanTypeName {
		logError("while condition type error", conditionType)
		return UnKnownType{}
	}

	body := node.(WhileStatement).Body
	if body.Type() == AstTypeBlockStatement.Name() {
		this.visitBlockStatement(body)
	}

	return VoidType{}
}

func (this *SemanticAnalysisV2) visitForStatement(node Node) AllType {
	if node.Type() != AstTypeForStatement.Name() {
		return UnKnownType{}
	} // 开始新的作用域

	funcScope := NewScopeV2()
	funcScope.SetParent(this.CurrentScope)
	this.CurrentScope = funcScope
	// 这里用 defer 主要是怕提前 return，作用域没回滚
	defer func() {
		this.CurrentScope = this.CurrentScope.Parent
	}()

	init := node.(ForStatement).Init
	var initType AllType
	initType = UnKnownType{}
	switch init.Type() {
	case AstTypeVariableDeclaration.Name():
		this.visitVariableDeclaration(init)
	case AstTypeAssignmentExpression.Name():
		this.visitAssignmentExpression(init)
	case AstTypeBooleanLiteral.Name():
		initType = this.visitBooleanLiteral(init)
	default:
		logError("not support for statement init type", init)
		return UnKnownType{}
	}
	logInfo("for init type", initType)

	testExp := node.(ForStatement).Test
	var testType AllType
	testType = UnKnownType{}
	switch testExp.Type() {
	case AstTypeBinaryExpression.Name():
		testType = this.visitBinaryExpression(testExp)
	default:
		logError("not support for statement test type", testExp)
		return UnKnownType{}
	}
	logInfo("for testType type", testType)

	update := node.(ForStatement).Update
	var updateType AllType
	updateType = UnKnownType{}
	switch update.Type() {
	case AstTypeBinaryExpression.Name():
		testType = this.visitBinaryExpression(update)
	default:
		logError("not support for statement update type", update)
		return UnKnownType{}
	}
	logInfo("for updateType type", updateType)

	body := node.(ForStatement).Body
	if body.Type() == AstTypeBlockStatement.Name() {
		this.visitBlockStatement(body)
	}

	return VoidType{}
}

func (this *SemanticAnalysisV2) visitIfStatement(node Node) AllType {
	if node.Type() == AstTypeIfStatement.Name() {
		return UnKnownType{}
	}
	condition := node.(IfStatement).Condition
	var conditionType AllType
	conditionType = UnKnownType{}
	switch condition.Type() {
	case AstTypeBinaryExpression.Name():
		conditionType = this.visitBinaryExpression(condition)
	case AstTypeIdentifier.Name():
		conditionType, _, _ = this.visitIdentifier(condition)
	case AstTypeBooleanLiteral.Name():
		conditionType = this.visitBooleanLiteral(condition)
	default:
		logError("not support if statement condition type", conditionType)
		return UnKnownType{}
	}
	logInfo("if condition type", conditionType)

	bValueType := BooleanType{}.ValueType()
	if conditionType.ValueType() != bValueType {
		logError("if condition type error", conditionType)
		return UnKnownType{}
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

	return VoidType{}
}

// block
func (this *SemanticAnalysisV2) visitBlockStatement(node Node) AllType {
	// 开始新的作用域
	// 进入 block 就是一个新的作用域
	funcScope := NewScopeV2()
	funcScope.SetParent(this.CurrentScope)
	this.CurrentScope = funcScope

	// 防止中途退出作用域没返回
	defer func() {
		logInfo("visitBlockStatement after current Scope", this.CurrentScope)
		this.CurrentScope = this.CurrentScope.Parent
	}()

	// 默认函数返回都是 null
	var retValueType AllType
	retValueType = UnKnownType{}
	logInfo("visitBlockStatement visit node before", node.Type())
	for _, item := range node.(BlockStatement).Body {
		logInfo("visitBlockStatement visit item", item.Type())
		switch item.Type() {
		// 变量定义
		case AstTypeVariableDeclaration.Name():
			this.visitVariableDeclaration(item)
		// 赋值
		case AstTypeAssignmentExpression.Name():
			this.visitAssignmentExpression(item)
		// for 循环 break
		case AstTypeBreakStatement.Name():
			this.visitBreakStatement(item)
		// for 循环 continue
		case AstTypeContinueStatement.Name():
			this.visitContinueStatement(item)
		// return
		case AstTypeReturnStatement.Name():
			retValueType = this.visitReturnStatement(item)
		// 函数调用
		case AstTypeCallExpression.Name():
			this.visitCallExpression(item)
		default:
			logError("not support block statement type", item)
		}
	}
	logInfo("visitBlockStatement after current Scope", retValueType)
	this.CurrentScope.LogNowScope()
	return retValueType
}

// break
func (this *SemanticAnalysisV2) visitBreakStatement(node Node) AllType {
	return VoidType{}
}

func (this *SemanticAnalysisV2) visitContinueStatement(node Node) AllType {
	return VoidType{}
}

func (this *SemanticAnalysisV2) visitReturnStatement(node Node) AllType {
	if node.Type() != AstTypeReturnStatement.Name() {
		return UnKnownType{}
	}
	v := node.(ReturnStatement).Value
	if v == nil {
		return VoidType{}
	}

	var rightType AllType
	rightType = VoidType{}
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
		symbol, ok := this.CurrentScope.LookupSignature(varName)
		if ok {
			rightType = symbol.ReturnType
		} else {
			logError("undeclared variable", varName)
			return UnKnownType{}
		}
	default:
		logError("not support return value type", v.Type())
		return UnKnownType{}
	}
	logInfo("rightType", rightType)
	return rightType
}

// array
func (this *SemanticAnalysisV2) visitArrayLiteral(node Node) AllType {
	if node.Type() != AstTypeArrayLiteral.Name() {
		return VoidType{}
	}
	// 做一下限制，不能多种类型混合在数组里面一起
	if len(node.(ArrayLiteral).Values) > 0 {
		firstElement := node.(ArrayLiteral).Values[0]
		var firstElementType AllType
		checkType := firstElement.Type()
		switch checkType {
		case AstTypeStringLiteral.Name():
			firstElementType, _ = this.visitStringLiteral(firstElement)
		case AstTypeNullLiteral.Name():
			firstElementType = this.visitNumberLiteral(firstElement)
		case AstTypeBooleanLiteral.Name():
			firstElementType = this.visitBooleanLiteral(firstElement)
		case AstTypeNumberLiteral.Name():
			firstElementType = this.visitNumberLiteral(firstElement)
		default:
			logError("array literal type error", checkType, node.(ArrayLiteral).Values[0].Type())
		}
		// 数组
		for _, item := range node.(ArrayLiteral).Values {
			if item.Type() != checkType {
				logError("array literal type error", checkType, item.Type())
				return VoidType{}
			}
		}
		return ArrayType{
			ElementType: firstElementType,
		}
	}
	// 这里
	return ArrayType{}
}

// dict
func (this *SemanticAnalysisV2) visitDictLiteral(node Node) AllType {
	if node.Type() != AstTypeDictLiteral.Name() {
		return VoidType{}
	}
	// 做一下限制，不能多种类型混合在dict
	// key 也不能重复
	if len(node.(DictLiteral).Values) <= 0 {
		return DictType{
			KeyType: VoidType{},
			VType:   VoidType{},
		}
	}
	keyNames := []string{}
	firstVType, _ := this.visitPropertyAssignment(node.(DictLiteral).Values[0])
	for _, value := range node.(DictLiteral).Values {
		dictVType, keyName := this.visitPropertyAssignment(value)
		if dictVType != firstVType {
			logError("dict literal value type not the same error", firstVType, dictVType)
			return VoidType{}
		}
		inSlice := InStringSlice(keyNames, keyName)
		if inSlice {
			logError("dict literal duplicate key", keyName)
			return VoidType{}
		}
		keyNames = append(keyNames, keyName)
	}
	return DictType{
		KeyType: StringType{},
		VType:   firstVType,
	}
}

// k:v
func (this *SemanticAnalysisV2) visitPropertyAssignment(node Node) (vType AllType, keyName string) {
	if node.Type() != AstTypePropertyAssignment.Name() {
		return VoidType{}, ""
	}
	_, keyName = this.visitStringLiteral(node.(PropertyAssignment).Key)

	v := node.(PropertyAssignment).Value
	switch v.Type() {
	// todo
	case AstTypeFunctionExpression.Name():
		vType = this.visitFunctionExpression(v)
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
		symbol, ok := this.CurrentScope.LookupSignature(varName)
		if ok {
			vType = symbol.ReturnType
		} else {
			logError("undeclared variable", varName)
			return
		}
	default:
		logError("property assignment type error", v.Type(), node.Type())
	}
	return vType, keyName
}

// 调用
func (this *SemanticAnalysisV2) visitCallExpression(node Node) (AllType, string) {
	if node.Type() != AstTypeCallExpression.Name() {
		return UnKnownType{}, ""
	}
	if node.(CallExpression).Object.Type() == AstTypeMemberExpression.Name() {
		return this.visitMemberExpression(node.(CallExpression).Object)
	} else if node.(CallExpression).Object.Type() == AstTypeIdentifier.Name() {
		valueType, variableName, _ := this.visitIdentifier(node.(CallExpression).Object)
		functionValueType := FunctionType{}
		if valueType.ValueType() == functionValueType.ValueType() {
			return valueType.(FunctionType).ReturnType, variableName
		}
		return valueType, variableName
	}
	logInfo("visitCallExpression", node.(CallExpression).Object.Type())
	return UnKnownType{}, ""
}

// Todo array dot 语法 [] .
func (this *SemanticAnalysisV2) visitMemberExpression(node Node) (AllType, string) {
	//var variableName string
	//if node.Type() == AstTypeMemberExpression.Name() {
	//	et := node.(MemberExpression).ElementType
	//	switch et {
	//	case "dot":
	//		logInfo("visitMemberExpression node.(before)", node.(MemberExpression))
	//		if node.(MemberExpression).Object.Type() == AstTypeIdentifier.Name() {
	//			var specificTypeName string
	//			_, specificTypeName, _ = this.visitIdentifier(node.(MemberExpression).Object)
	//			logInfo("visitMemberExpression node.(MemberExpression)", node.(MemberExpression))
	//			if specificTypeName == TokenTypeCls.Name() || specificTypeName == TokenTypeThis.Name() {
	//				//logInfo("visitMemberExpression node.(MemberExpression).Property", node.(MemberExpression).Property)
	//				//// classPropertyName
	//				//var variableValueType = ValueTypeIdentifier
	//				//variableValueType, variableName, _ = this.visitIdentifier(node.(MemberExpression).Property)
	//				//return variableValueType, variableName
	//			} else {
	//				//// 1.处理.new()
	//				//var variableValueType = ValueTypeIdentifier
	//				//variableValueType, variableName, _ = this.visitIdentifier(node.(MemberExpression).Property)
	//				//logInfo("visitMemberExpression variableName", variableValueType, variableName)
	//				//if variableName == TokenTypeNew.Name() {
	//				//	variableValueType, _, _ = this.visitIdentifier(node.(MemberExpression).Object)
	//				//	return variableValueType, ""
	//				//}
	//				// 2.处理类方法
	//				if classScope, ok := this.CurrentScope.LookupScope(specificTypeName); ok {
	//					for _, subScope := range classScope.SubScopes {
	//						if symbol, ok1 := subScope.LookupSymbol(variableName); ok1 {
	//							logInfo("visitMemberExpression subScope symbol", symbol)
	//							return symbol.ExtraInfo, variableName
	//						}
	//					}
	//				}
	//				return variableValueType, variableName
	//			}
	//		}
	//
	//	case "array":
	//	}
	//	return UnKnownType{}, ""
	//}
	return UnKnownType{}, ""
}

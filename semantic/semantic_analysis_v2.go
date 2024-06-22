package semantic

import (
	"fmt"
	"go-compiler/lexer"
	"go-compiler/parser"
	"go-compiler/utils"
)

type SemanticAnalysisV2 struct {
	Ast          parser.Program
	CurrentScope *ScopeV2
}

func NewSemanticAnalysisV2(program parser.Program) *SemanticAnalysisV2 {
	s := &SemanticAnalysisV2{
		Ast:          program,
		CurrentScope: NewScopeV2(),
	}
	return s
}

func (this *SemanticAnalysisV2) Visit() {
	if this.Ast.Type() != parser.AstTypeProgram.Name() {
		return
	}
	this.visitProgram(this.Ast.Body)

}

func (this *SemanticAnalysisV2) visitProgram(body []parser.Node) {
	for _, item := range body {
		utils.LogInfo("visitProgram visit item", item.Type())
		switch item.Type() {
		// 变量定义
		case parser.AstTypeVariableDeclaration.Name():
			this.visitVariableDeclaration(item)
		//// 赋值
		case parser.AstTypeAssignmentExpression.Name():
			this.visitAssignmentExpression(item)
		//// 访问 if
		case parser.AstTypeIfStatement.Name():
			this.visitIfStatement(item)
		// 访问 while
		case parser.AstTypeWhileStatement.Name():
			this.visitWhileStatement(item)
		//// 访问 for
		case parser.AstTypeForStatement.Name():
			this.visitForStatement(item)
		//// 访问 block
		case parser.AstTypeBlockStatement.Name():
			this.visitBlockStatement(item)
		// 访问 class
		case parser.AstTypeClassExpression.Name():
			this.visitClassExpression(item)
		// 调用函数
		case parser.AstTypeCallExpression.Name():
			this.visitCallExpression(item)
		default:
			utils.LogError("visitProgram visit item default", item.Type())
		}
		utils.LogInfo("visitProgram visit item after currentScope")
		this.CurrentScope.LogNowScope()
	}
}

// 变量定义
func (this *SemanticAnalysisV2) visitVariableDeclaration(node parser.Node) {
	////type VariableDeclaration struct {
	////	Kind  string // kind属性
	////	Name  Node   // name属性
	////	Value Node   // value属性
	////}
	varType := node.(parser.VariableDeclaration).Kind
	left := node.(parser.VariableDeclaration).Name
	var variableName string
	switch left.Type() {
	case parser.AstTypeIdentifier.Name():
		_, variableName, _ = this.visitIdentifier(left)
	default:
		utils.LogError("visitVariableDeclaration invalid left variable declaration", left)
		return
	}
	// 做一下限制，变量名不为空
	if len(variableName) == 0 {
		utils.LogError("visitVariableDeclaration invalid left variable declaration", left)
		return
	}
	//
	var valueType AllType
	right := node.(parser.VariableDeclaration).Value

	switch right.Type() {
	// 访问函数，函数的返回类型
	case parser.AstTypeFunctionExpression.Name():
		valueType = this.visitFunctionExpression(right)
	case parser.AstTypeBinaryExpression.Name():
		valueType = this.visitBinaryExpression(right)
	case parser.AstTypeNumberLiteral.Name():
		valueType = this.visitNumberLiteral(right)
	case parser.AstTypeNullLiteral.Name():
		valueType = this.visitNullLiteral(right)
	case parser.AstTypeBooleanLiteral.Name():
		valueType = this.visitBooleanLiteral(right)
	case parser.AstTypeStringLiteral.Name():
		valueType, _ = this.visitStringLiteral(right)
	case parser.AstTypeArrayLiteral.Name():
		valueType = this.visitArrayLiteral(right)
	case parser.AstTypeUnaryExpression.Name():
		valueType = this.visitUnaryExpression(right)
	case parser.AstTypeDictLiteral.Name():
		valueType = this.visitDictLiteral(right)
	case parser.AstTypeIdentifier.Name():
		var varName string
		_, varName, _ = this.visitIdentifier(right)
		signature, ok := this.CurrentScope.LookupSignature(varName)
		if ok {
			valueType = signature.ReturnType
		} else {
			utils.LogError("undeclared variable", varName)
			return
		}
	case parser.AstTypeCallExpression.Name():
		valueType, _ = this.visitCallExpression(right)
	default:
		utils.LogError("visitVariableDeclaration invalid right variable declaration", right)
		return
	}
	// 先不处理常量方法
	this.CurrentScope.AddSignature(variableName, valueType, false, varType)

	fmt.Printf("in visitVariableDeclaration this.CurrentScope")
	this.CurrentScope.LogNowScope()
	return
}

// 赋值
func (this *SemanticAnalysisV2) visitAssignmentExpression(node parser.Node) {
	left := node.(parser.AssignmentExpression).Left
	utils.LogInfo("visitClassVariableDeclaration visitAssignmentExpression", node.(parser.AssignmentExpression))
	var variableName string
	switch left.Type() {
	case parser.AstTypeIdentifier.Name():
		_, variableName, _ = this.visitIdentifier(left)
	//case AstTypeMemberExpression.Name():
	//	_, variableName = this.visitMemberExpression(left)
	default:
		utils.LogError("invalid assignment expression", left)
		return
	}
	fmt.Printf("visitAssignmentExpression variableName %v left:%v\n", variableName, left)
	varSignature, ok := this.CurrentScope.LookupSignature(variableName)
	// const 检查
	if ok && varSignature.VarType == lexer.TokenTypeConst.Name() {
		utils.LogError("const variable cannot be reassigned", variableName)
		return
	}

	varType := lexer.TokenTypeVar.Name()
	var valueType AllType
	valueType = UnKnownType{}
	//valueType := ValueTypeIdentifier
	right := node.(parser.AssignmentExpression).Right
	switch right.Type() {
	case parser.AstTypeBinaryExpression.Name():
		valueType = this.visitBinaryExpression(right)
	case parser.AstTypeIdentifier.Name():
		valueType, _, varType = this.visitIdentifier(right)
	case parser.AstTypeNumberLiteral.Name():
		valueType = this.visitNumberLiteral(right)
	case parser.AstTypeStringLiteral.Name():
		valueType, _ = this.visitStringLiteral(right)
	case parser.AstTypeBooleanLiteral.Name():
		valueType = this.visitBooleanLiteral(right)
	case parser.AstTypeArrayLiteral.Name():
		valueType = this.visitArrayLiteral(right)
	case parser.AstTypeDictLiteral.Name():
		valueType = this.visitDictLiteral(right)
	case parser.AstTypeNullLiteral.Name():
		valueType = this.visitNullLiteral(right)
	case parser.AstTypeUnaryExpression.Name():
		valueType = this.visitUnaryExpression(right)
	case parser.AstTypeCallExpression.Name():
		valueType, _ = this.visitCallExpression(right)
	default:
		utils.LogError("invalid assignment expression", right)
	}
	unknownValueType := UnKnownType{}
	// 强类型检查，如果右边的值不是同一个类型就报错
	if ok && varSignature.ReturnType != valueType && valueType.ValueType() != unknownValueType.ValueType() {
		utils.LogError("variable cannot be reassigned to another type", variableName, varSignature, valueType)
		return
	}
	// 如果判断不出类型，就用原有的类型
	if ok && valueType.ValueType() == unknownValueType.ValueType() {
		valueType = varSignature.ReturnType
	}
	this.CurrentScope.AddSignature(variableName, valueType, false, varType)
	return
}

func (this *SemanticAnalysisV2) visitFunctionExpression(node parser.Node) (functionType AllType) {
	if node.Type() != parser.AstTypeFunctionExpression.Name() {
		return UnKnownType{}
	}
	params := node.(parser.FunctionExpression).Params

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
		if param.Type() != parser.AstTypeIdentifier.Name() {
			utils.LogError("param must be identifier", param.Type())
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
	body := node.(parser.FunctionExpression).Body
	utils.LogInfo("visitFunctionExpression", params, body)

	var funcReturnType AllType
	funcReturnType = VoidType{}
	if body.Type() == parser.AstTypeBlockStatement.Name() {
		funcReturnType = this.visitBlockStatement(body)
	}

	return FunctionType{
		Params:     signatures,
		ReturnType: funcReturnType,
	}
}

// 表达式
func (this *SemanticAnalysisV2) visitBinaryExpression(node parser.Node) AllType {
	// BinaryExpression节点结构
	//type BinaryExpression struct {
	//	Operator string // operator属性
	//	Left     Node   // left属性
	//	Right    Node   // right属性
	//}
	if node.Type() != parser.AstTypeBinaryExpression.Name() {
		return UnKnownType{}
	}
	// op
	switch node.(parser.BinaryExpression).Operator {
	case "+":
		// 这里接受多个类型
		left := node.(parser.BinaryExpression).Left
		var leftValueType AllType
		switch left.Type() {
		case parser.AstTypeBinaryExpression.Name():
			leftValueType = this.visitBinaryExpression(left)
		case parser.AstTypeNumberLiteral.Name():
			leftValueType = this.visitNumberLiteral(left)
		case parser.AstTypeStringLiteral.Name():
			leftValueType, _ = this.visitStringLiteral(left)
		case parser.AstTypeIdentifier.Name():
			leftValueType, _, _ = this.visitIdentifier(left)
		default:
			utils.LogError("左值类型错误", node.(parser.BinaryExpression).Left)
			return UnKnownType{}
		}

		// right
		right := node.(parser.BinaryExpression).Right
		var rightValueType AllType
		switch right.Type() {
		case parser.AstTypeBinaryExpression.Name():
			rightValueType = this.visitBinaryExpression(right)
		case parser.AstTypeNumberLiteral.Name():
			rightValueType = this.visitNumberLiteral(right)
		case parser.AstTypeStringLiteral.Name():
			rightValueType, _ = this.visitStringLiteral(right)
		case parser.AstTypeIdentifier.Name():
			rightValueType, _, _ = this.visitIdentifier(right)
		default:
			utils.LogError("右值类型错误", node.(parser.BinaryExpression).Right)
			return UnKnownType{}
		}
		utils.LogInfo("visitBinaryExpression", leftValueType, rightValueType)

		isString := false

		switch leftValueType {
		case NumberType{}:
		case StringType{}:
			isString = true
		default:
			utils.LogError("左值类型错误", node.(parser.BinaryExpression).Left)
			return UnKnownType{}
		}

		switch rightValueType {
		case NumberType{}:
		case StringType{}:
			isString = true
		default:
			utils.LogError("右值类型错误", node.(parser.BinaryExpression).Right)
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
		left := node.(parser.BinaryExpression).Left
		var leftValueType AllType
		switch left.Type() {
		case parser.AstTypeIdentifier.Name():
			leftValueType, _, _ = this.visitIdentifier(left)
		case parser.AstTypeBinaryExpression.Name():
			leftValueType = this.visitBinaryExpression(left)
		case parser.AstTypeNumberLiteral.Name():
			leftValueType = this.visitNumberLiteral(left)
		default:
			utils.LogError("左值类型错误", node.(parser.BinaryExpression).Left)
			return UnKnownType{}
		}

		// right
		right := node.(parser.BinaryExpression).Right
		var rightValueType AllType
		switch right.Type() {
		case parser.AstTypeBinaryExpression.Name():
			rightValueType = this.visitBinaryExpression(right)
		case parser.AstTypeNumberLiteral.Name():
			rightValueType = this.visitNumberLiteral(right)
		default:
			utils.LogError("右值类型错误", node.(parser.BinaryExpression).Right)
			return UnKnownType{}
		}
		switch leftValueType {
		case NumberType{}:
		default:
			utils.LogError("左值类型错误", node.(parser.BinaryExpression).Left)
			return UnKnownType{}
		}

		switch rightValueType {
		case NumberType{}:
		default:
			utils.LogError("右值类型错误", node.(parser.BinaryExpression).Right)
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
		left := node.(parser.BinaryExpression).Left
		var leftValueType AllType
		switch left.Type() {
		case parser.AstTypeBinaryExpression.Name():
			leftValueType = this.visitBinaryExpression(left)
		case parser.AstTypeNumberLiteral.Name():
			leftValueType = this.visitNumberLiteral(left)
		case parser.AstTypeStringLiteral.Name():
			leftValueType, _ = this.visitStringLiteral(left)
		case parser.AstTypeIdentifier.Name():
			leftValueType, _, _ = this.visitIdentifier(left)
		default:
			utils.LogError("左值类型错误", node.(parser.BinaryExpression).Left)
			return UnKnownType{}
		}

		// right
		right := node.(parser.BinaryExpression).Right
		var rightValueType AllType
		switch right.Type() {
		case parser.AstTypeBinaryExpression.Name():
			rightValueType = this.visitBinaryExpression(right)
		case parser.AstTypeNumberLiteral.Name():
			rightValueType = this.visitNumberLiteral(right)
		case parser.AstTypeStringLiteral.Name():
			rightValueType, _ = this.visitStringLiteral(right)
		case parser.AstTypeIdentifier.Name():
			rightValueType, _, _ = this.visitIdentifier(right)
		default:
			utils.LogError("右值类型错误", node.(parser.BinaryExpression).Right)
			return UnKnownType{}
		}
		utils.LogInfo("visitBinaryExpression", leftValueType, rightValueType)

		if leftValueType != rightValueType {
			utils.LogError("类型不匹配", node.(parser.BinaryExpression).Left, node.(parser.BinaryExpression).Right)
			return UnKnownType{}
		}
		return BooleanType{}
	case "==":
		fallthrough
	case "!=":
		// left
		left := node.(parser.BinaryExpression).Left
		var leftValueType AllType
		switch left.Type() {
		case parser.AstTypeBinaryExpression.Name():
			leftValueType = this.visitBinaryExpression(left)
		case parser.AstTypeNumberLiteral.Name():
			leftValueType = this.visitNumberLiteral(left)
		case parser.AstTypeStringLiteral.Name():
			leftValueType, _ = this.visitStringLiteral(left)
		case parser.AstTypeIdentifier.Name():
			leftValueType, _, _ = this.visitIdentifier(left)
		case parser.AstTypeBooleanLiteral.Name():
			leftValueType = this.visitBooleanLiteral(left)
		case parser.AstTypeNullLiteral.Name():
			leftValueType = this.visitNullLiteral(left)
		default:
			utils.LogError("左值类型错误", node.(parser.BinaryExpression).Left)
			return UnKnownType{}
		}

		// right
		right := node.(parser.BinaryExpression).Right
		var rightValueType AllType
		switch right.Type() {
		case parser.AstTypeBinaryExpression.Name():
			rightValueType = this.visitBinaryExpression(right)
		case parser.AstTypeNumberLiteral.Name():
			rightValueType = this.visitNumberLiteral(right)
		case parser.AstTypeStringLiteral.Name():
			rightValueType, _ = this.visitStringLiteral(right)
		case parser.AstTypeIdentifier.Name():
			rightValueType, _, _ = this.visitIdentifier(right)
		case parser.AstTypeBooleanLiteral.Name():
			rightValueType = this.visitBooleanLiteral(right)
		case parser.AstTypeNullLiteral.Name():
			rightValueType = this.visitNullLiteral(right)
		default:
			utils.LogError("右值类型错误", node.(parser.BinaryExpression).Right)
			return UnKnownType{}
		}
		utils.LogInfo("visitBinaryExpression", leftValueType, rightValueType)

		if leftValueType != rightValueType {
			utils.LogError("类型不匹配", node.(parser.BinaryExpression).Left, node.(parser.BinaryExpression).Right)
			return UnKnownType{}
		}
		return BooleanType{}
	case "and":
		fallthrough
	case "or":
		// left
		left := node.(parser.BinaryExpression).Left
		var leftValueType AllType
		switch left.Type() {
		case parser.AstTypeBinaryExpression.Name():
			leftValueType = this.visitBinaryExpression(left)
		case parser.AstTypeIdentifier.Name():
			leftValueType, _, _ = this.visitIdentifier(left)
		case parser.AstTypeBooleanLiteral.Name():
			leftValueType = this.visitBooleanLiteral(left)
		default:
			utils.LogError("左值类型错误", node.(parser.BinaryExpression).Left)
			return UnKnownType{}
		}

		// right
		right := node.(parser.BinaryExpression).Right
		var rightValueType AllType
		switch right.Type() {
		case parser.AstTypeBinaryExpression.Name():
			rightValueType = this.visitBinaryExpression(right)
		case parser.AstTypeIdentifier.Name():
			rightValueType, _, _ = this.visitIdentifier(right)
		case parser.AstTypeBooleanLiteral.Name():
			rightValueType = this.visitBooleanLiteral(right)
		default:
			utils.LogError("右值类型错误", node.(parser.BinaryExpression).Right)
			return UnKnownType{}
		}
		utils.LogInfo("in and not visitBinaryExpression", leftValueType, rightValueType)

		if (leftValueType.ValueType() != BooleanType{}.ValueType()) && (rightValueType.ValueType() != BooleanType{}.ValueType()) {
			utils.LogError("类型不匹配", node.(parser.BinaryExpression).Left, node.(parser.BinaryExpression).Right)
			return UnKnownType{}
		}
		return BooleanType{}
	default:
		utils.LogError("not support binary expression operator", node.(parser.BinaryExpression).Operator)
		return UnKnownType{}
	}
}

// not
func (this *SemanticAnalysisV2) visitUnaryExpression(node parser.Node) AllType {
	if node.Type() != parser.AstTypeUnaryExpression.Name() {
		return UnKnownType{}
	}
	op := node.(parser.UnaryExpression).Operator
	switch op {
	case "not":
		var vValueType AllType
		// not 后面只能加 identifier / bool
		v := node.(parser.UnaryExpression).Value
		switch v.Type() {
		case parser.AstTypeBinaryExpression.Name():
			vValueType = this.visitBinaryExpression(v)
		case parser.AstTypeIdentifier.Name():
			vValueType, _, _ = this.visitIdentifier(v)
		case parser.AstTypeBooleanLiteral.Name():
			vValueType = this.visitBooleanLiteral(v)
		}
		boolType := BooleanType{}.ValueType()
		if vValueType.ValueType() != boolType {
			utils.LogError("not value type error", vValueType)
			return UnKnownType{}
		}
		return BooleanType{}
	case "-":
		var vValueType AllType
		// - 后面只能加 identifier / number
		v := node.(parser.UnaryExpression).Value
		switch v.Type() {
		case parser.AstTypeBinaryExpression.Name():
			vValueType = this.visitBinaryExpression(v)
		case parser.AstTypeIdentifier.Name():
			vValueType, _, _ = this.visitIdentifier(v)
		case parser.AstTypeNumberLiteral.Name():
			vValueType = this.visitNumberLiteral(v)
		}
		numberType := NumberType{}.ValueType()
		if vValueType.ValueType() != numberType {
			utils.LogError("not value type error", vValueType)
			return UnKnownType{}
		}
		return NumberType{}
	default:
		utils.LogError("not support unary expression operator", node.(parser.UnaryExpression).Operator)
		return UnKnownType{}
	}

}

func (this *SemanticAnalysisV2) visitIdentifier(node parser.Node) (valueType AllType, variableName string, varType string) {
	if node.Type() != parser.AstTypeIdentifier.Name() {
		return VoidType{}, "", varType
	}
	signature, ok := this.CurrentScope.LookupSignature(node.(parser.Identifier).Value)
	if ok {
		return signature.ReturnType, signature.Name, varType
	}
	return UnKnownType{}, node.(parser.Identifier).Value, varType
}

func (this *SemanticAnalysisV2) visitDictKeyLiteral(node parser.Node) (valueType AllType, variableName string, varType string) {
	switch node.Type() {
	case parser.AstTypeStringLiteral.Name():
		valueType, variableName = this.visitStringLiteral(node)
		return valueType, variableName, varType
	case parser.AstTypeNumberLiteral.Name():
		valueType = this.visitNumberLiteral(node)
		return valueType, fmt.Sprintf("%d", int64(node.(parser.NumberLiteral).Value)), varType
	case parser.AstTypeIdentifier.Name():
		return this.visitIdentifier(node)
	default:
		utils.LogError("visitDictKeyLiteral error, unknownType", node)
	}

	return UnKnownType{}, "", ""
}

func (this *SemanticAnalysisV2) visitStringLiteral(node parser.Node) (valueType AllType, value string) {
	if node.Type() != parser.AstTypeStringLiteral.Name() {
		return VoidType{}, ""
	}
	value = node.(parser.StringLiteral).Value
	return StringType{}, value
}

func (this *SemanticAnalysisV2) visitBooleanLiteral(node parser.Node) AllType {
	if node.Type() != parser.AstTypeBooleanLiteral.Name() {
		return VoidType{}
	}
	return BooleanType{}
}

func (this *SemanticAnalysisV2) visitNullLiteral(node parser.Node) AllType {
	if node.Type() != parser.AstTypeNullLiteral.Name() {
		return VoidType{}
	}
	return NullType{}
}

func (this *SemanticAnalysisV2) visitNumberLiteral(node parser.Node) AllType {
	if node.Type() != parser.AstTypeNumberLiteral.Name() {
		return VoidType{}
	}
	return NumberType{}
}

func (this *SemanticAnalysisV2) visitWhileStatement(node parser.Node) AllType {
	if node.Type() != parser.AstTypeWhileStatement.Name() {
		return UnKnownType{}
	}
	condition := node.(parser.WhileStatement).Condition
	var conditionType AllType
	switch condition.Type() {
	case parser.AstTypeBinaryExpression.Name():
		conditionType = this.visitBinaryExpression(condition)
	case parser.AstTypeIdentifier.Name():
		conditionType, _, _ = this.visitIdentifier(condition)
	case parser.AstTypeBooleanLiteral.Name():
		conditionType = this.visitBooleanLiteral(condition)
	default:
		utils.LogError("not support while statement condition type", conditionType)
		return UnKnownType{}
	}

	booleanTypeName := BooleanType{}.ValueType()
	if conditionType.ValueType() != booleanTypeName {
		utils.LogError("while condition type error", conditionType)
		return UnKnownType{}
	}

	body := node.(parser.WhileStatement).Body
	if body.Type() == parser.AstTypeBlockStatement.Name() {
		this.visitBlockStatement(body)
	}

	return VoidType{}
}

func (this *SemanticAnalysisV2) visitForStatement(node parser.Node) AllType {
	if node.Type() != parser.AstTypeForStatement.Name() {
		return UnKnownType{}
	} // 开始新的作用域

	funcScope := NewScopeV2()
	funcScope.SetParent(this.CurrentScope)
	this.CurrentScope = funcScope
	// 这里用 defer 主要是怕提前 return，作用域没回滚
	defer func() {
		this.CurrentScope = this.CurrentScope.Parent
	}()

	init := node.(parser.ForStatement).Init
	var initType AllType
	initType = UnKnownType{}
	switch init.Type() {
	case parser.AstTypeVariableDeclaration.Name():
		this.visitVariableDeclaration(init)
	case parser.AstTypeAssignmentExpression.Name():
		this.visitAssignmentExpression(init)
	case parser.AstTypeBooleanLiteral.Name():
		initType = this.visitBooleanLiteral(init)
	default:
		utils.LogError("not support for statement init type", init)
		return UnKnownType{}
	}
	utils.LogInfo("for init type", initType)

	testExp := node.(parser.ForStatement).Test
	var testType AllType
	testType = UnKnownType{}
	switch testExp.Type() {
	case parser.AstTypeBinaryExpression.Name():
		testType = this.visitBinaryExpression(testExp)
	default:
		utils.LogError("not support for statement test type", testExp)
		return UnKnownType{}
	}
	utils.LogInfo("for testType type", testType)

	update := node.(parser.ForStatement).Update
	var updateType AllType
	updateType = UnKnownType{}
	switch update.Type() {
	case parser.AstTypeBinaryExpression.Name():
		testType = this.visitBinaryExpression(update)
	default:
		utils.LogError("not support for statement update type", update)
		return UnKnownType{}
	}
	utils.LogInfo("for updateType type", updateType)

	body := node.(parser.ForStatement).Body
	if body.Type() == parser.AstTypeBlockStatement.Name() {
		this.visitBlockStatement(body)
	}

	return VoidType{}
}

func (this *SemanticAnalysisV2) visitIfStatement(node parser.Node) AllType {
	if node.Type() == parser.AstTypeIfStatement.Name() {
		return UnKnownType{}
	}
	condition := node.(parser.IfStatement).Condition
	var conditionType AllType
	conditionType = UnKnownType{}
	switch condition.Type() {
	case parser.AstTypeBinaryExpression.Name():
		conditionType = this.visitBinaryExpression(condition)
	case parser.AstTypeIdentifier.Name():
		conditionType, _, _ = this.visitIdentifier(condition)
	case parser.AstTypeBooleanLiteral.Name():
		conditionType = this.visitBooleanLiteral(condition)
	default:
		utils.LogError("not support if statement condition type", conditionType)
		return UnKnownType{}
	}
	utils.LogInfo("if condition type", conditionType)

	bValueType := BooleanType{}.ValueType()
	if conditionType.ValueType() != bValueType {
		utils.LogError("if condition type error", conditionType)
		return UnKnownType{}
	}

	consequent := node.(parser.IfStatement).Consequent
	if consequent.Type() == parser.AstTypeBlockStatement.Name() {
		this.visitBlockStatement(consequent)
	}

	alternate := node.(parser.IfStatement).Alternate
	if alternate != nil {
		if alternate.Type() == parser.AstTypeBlockStatement.Name() {
			this.visitBlockStatement(alternate)
		}
	}

	return VoidType{}
}

// block
func (this *SemanticAnalysisV2) visitBlockStatement(node parser.Node) AllType {
	// 开始新的作用域
	// 进入 block 就是一个新的作用域
	funcScope := NewScopeV2()
	funcScope.SetParent(this.CurrentScope)
	this.CurrentScope = funcScope

	// 防止中途退出作用域没返回
	defer func() {
		utils.LogInfo("visitBlockStatement after current Scope", this.CurrentScope)
		this.CurrentScope = this.CurrentScope.Parent
	}()

	// 默认函数返回都是 null
	var retValueType AllType
	retValueType = UnKnownType{}
	utils.LogInfo("visitBlockStatement visit node before", node.Type())
	for _, item := range node.(parser.BlockStatement).Body {
		utils.LogInfo("visitBlockStatement visit item", item.Type())
		switch item.Type() {
		// 变量定义
		case parser.AstTypeVariableDeclaration.Name():
			this.visitVariableDeclaration(item)
		// 赋值
		case parser.AstTypeAssignmentExpression.Name():
			this.visitAssignmentExpression(item)
		// for 循环 break
		case parser.AstTypeBreakStatement.Name():
			this.visitBreakStatement(item)
		// for 循环 continue
		case parser.AstTypeContinueStatement.Name():
			this.visitContinueStatement(item)
		// return
		case parser.AstTypeReturnStatement.Name():
			retValueType = this.visitReturnStatement(item)
		// 函数调用
		case parser.AstTypeCallExpression.Name():
			this.visitCallExpression(item)
		default:
			utils.LogError("not support block statement type", item)
		}
	}
	utils.LogInfo("visitBlockStatement before current Scope", retValueType)
	this.CurrentScope.LogNowScope()
	return retValueType
}

// break
func (this *SemanticAnalysisV2) visitBreakStatement(node parser.Node) AllType {
	return VoidType{}
}

func (this *SemanticAnalysisV2) visitContinueStatement(node parser.Node) AllType {
	return VoidType{}
}

func (this *SemanticAnalysisV2) visitReturnStatement(node parser.Node) AllType {
	if node.Type() != parser.AstTypeReturnStatement.Name() {
		return UnKnownType{}
	}
	v := node.(parser.ReturnStatement).Value
	if v == nil {
		return VoidType{}
	}

	var rightType AllType
	rightType = VoidType{}
	switch v.Type() {
	case parser.AstTypeBinaryExpression.Name():
		rightType = this.visitBinaryExpression(v)
	case parser.AstTypeNumberLiteral.Name():
		rightType = this.visitNumberLiteral(v)
	case parser.AstTypeNullLiteral.Name():
		rightType = this.visitNullLiteral(v)
	case parser.AstTypeBooleanLiteral.Name():
		rightType = this.visitBooleanLiteral(v)
	case parser.AstTypeStringLiteral.Name():
		rightType, _ = this.visitStringLiteral(v)
	case parser.AstTypeArrayLiteral.Name():
		rightType = this.visitArrayLiteral(v)
	case parser.AstTypeUnaryExpression.Name():
		rightType = this.visitUnaryExpression(v)
	case parser.AstTypeDictLiteral.Name():
		rightType = this.visitDictLiteral(v)
	case parser.AstTypeIdentifier.Name():
		var varName string
		_, varName, _ = this.visitIdentifier(v)
		symbol, ok := this.CurrentScope.LookupSignature(varName)
		if ok {
			rightType = symbol.ReturnType
		} else {
			utils.LogError("undeclared variable", varName)
			return UnKnownType{}
		}
	default:
		utils.LogError("not support return value type", v.Type())
		return UnKnownType{}
	}
	utils.LogInfo("rightType", rightType)
	return rightType
}

// array
func (this *SemanticAnalysisV2) visitArrayLiteral(node parser.Node) AllType {
	if node.Type() != parser.AstTypeArrayLiteral.Name() {
		return VoidType{}
	}
	// 做一下限制，不能多种类型混合在数组里面一起
	if len(node.(parser.ArrayLiteral).Values) > 0 {
		firstElement := node.(parser.ArrayLiteral).Values[0]
		var firstElementType AllType
		checkType := firstElement.Type()
		switch checkType {
		case parser.AstTypeStringLiteral.Name():
			firstElementType, _ = this.visitStringLiteral(firstElement)
		case parser.AstTypeNullLiteral.Name():
			firstElementType = this.visitNumberLiteral(firstElement)
		case parser.AstTypeBooleanLiteral.Name():
			firstElementType = this.visitBooleanLiteral(firstElement)
		case parser.AstTypeNumberLiteral.Name():
			firstElementType = this.visitNumberLiteral(firstElement)
		default:
			utils.LogError("array literal type error", checkType, node.(parser.ArrayLiteral).Values[0].Type())
		}
		// 数组
		for _, item := range node.(parser.ArrayLiteral).Values {
			if item.Type() != checkType {
				utils.LogError("array literal type error", checkType, item.Type())
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
func (this *SemanticAnalysisV2) visitDictLiteral(node parser.Node) AllType {
	if node.Type() != parser.AstTypeDictLiteral.Name() {
		return VoidType{}
	}
	// 做一下限制，不能多种类型混合在dict
	// key 也不能重复
	if len(node.(parser.DictLiteral).Values) <= 0 {
		return DictType{
			KeyType: VoidType{},
			VType:   VoidType{},
		}
	}
	keyNames := []string{}
	firstVType, _ := this.visitPropertyAssignment(node.(parser.DictLiteral).Values[0])
	for _, value := range node.(parser.DictLiteral).Values {
		dictVType, keyName := this.visitPropertyAssignment(value)
		if dictVType != firstVType {
			utils.LogError("dict literal value type not the same error", firstVType, dictVType)
			return VoidType{}
		}
		inSlice := utils.InStringSlice(keyNames, keyName)
		if inSlice {
			utils.LogError("dict literal duplicate key", keyName)
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
func (this *SemanticAnalysisV2) visitPropertyAssignment(node parser.Node) (vType AllType, keyName string) {
	if node.Type() != parser.AstTypePropertyAssignment.Name() {
		return VoidType{}, ""
	}

	kv := node.(parser.PropertyAssignment)
	// key 必须是 1.数字、string、变量也得是 数字 or string
	var keyValueType AllType
	keyValueType, keyName, _ = this.visitDictKeyLiteral(kv.Key)

	switch keyValueType {
	case StringType{}, NumberType{}:
	default:
		utils.LogError("property assignment key type error", keyValueType, kv.Key.Type())
	}

	v := kv.Value
	switch v.Type() {
	case parser.AstTypeFunctionExpression.Name():
		vType = this.visitFunctionExpression(v)
	case parser.AstTypeBinaryExpression.Name():
		vType = this.visitBinaryExpression(v)
	case parser.AstTypeNumberLiteral.Name():
		vType = this.visitNumberLiteral(v)
	case parser.AstTypeNullLiteral.Name():
		vType = this.visitNullLiteral(v)
	case parser.AstTypeStringLiteral.Name():
		vType, _ = this.visitStringLiteral(v)
	case parser.AstTypeArrayLiteral.Name():
		vType = this.visitArrayLiteral(v)
	case parser.AstTypeDictLiteral.Name():
		vType = this.visitDictLiteral(v)
	case parser.AstTypeIdentifier.Name():
		var varName string
		_, varName, _ = this.visitIdentifier(v)
		symbol, ok := this.CurrentScope.LookupSignature(varName)
		if ok {
			vType = symbol.ReturnType
		} else {
			utils.LogError("undeclared variable", varName)
			return
		}
	default:
		utils.LogError("property assignment type error", v.Type(), node.Type())
	}
	return vType, keyName
}

// 调用函数
func (this *SemanticAnalysisV2) visitCallExpression(node parser.Node) (AllType, string) {
	if node.Type() != parser.AstTypeCallExpression.Name() {
		return UnKnownType{}, ""
	}

	n := node.(parser.CallExpression)

	if n.Object.Type() == parser.AstTypeMemberExpression.Name() {
		return this.visitMemberExpression(node.(parser.CallExpression).Object)
	} else if n.Object.Type() == parser.AstTypeIdentifier.Name() {

		valueType, variableName, _ := this.visitIdentifier(n.Object)

		// 如果是函数就是直接调用它返回值
		functionValueType := FunctionType{}
		if valueType.ValueType() == functionValueType.ValueType() {
			fn := valueType.(FunctionType)
			// 这里要检查
			// 1. 参数个数
			// todo 2. 参数类型
			if len(fn.Params) != len(n.Args) {
				utils.LogError("call expression param number error", len(fn.Params), len(n.Args))
			}

			return fn.ReturnType, variableName
		}
		return valueType, variableName
	}
	utils.LogInfo("visitCallExpression", n.Object.Type())
	return UnKnownType{}, ""
}

// Todo array dot 语法 [] .
// 这里写的有点乱，要看看怎么搞比较好
func (this *SemanticAnalysisV2) visitMemberExpression(node parser.Node) (AllType, string) {
	var variableName string
	if node.Type() != parser.AstTypeMemberExpression.Name() {
		return UnKnownType{}, ""
	}
	et := node.(parser.MemberExpression).ElementType
	switch et {
	case "dot":
		//logInfo("visitMemberExpression node.(before)", node.(MemberExpression))
		var memberName string
		var memberType AllType
		memberType, memberName, _ = this.visitIdentifier(node.(parser.MemberExpression).Object)
		// 有可能是 this、cls 这类，也有可能是实例化的名字
		if memberName == lexer.TokenTypeCls.Name() || memberName == lexer.TokenTypeThis.Name() {
			var variableType AllType
			variableType, variableName, _ = this.visitIdentifier(node.(parser.MemberExpression).Property)
			return variableType, variableName
		} else {
			// 1.处理.new()
			var variableValueType AllType
			var propertyName string
			variableValueType, propertyName, _ = this.visitIdentifier(node.(parser.MemberExpression).Property)
			utils.LogInfo("visitMemberExpression variableName", variableValueType, variableName)
			if propertyName == lexer.TokenTypeNew.Name() {
				variableValueType, _, _ = this.visitIdentifier(node.(parser.MemberExpression).Object)
				return InstanceType{
					ClassType: ClassType{
						MemberSignatures: variableValueType.(ClassType).MemberSignatures,
						SuperType:        variableValueType.(ClassType).SuperType,
					},
				}, propertyName
			}
			// 2.处理类方法
			ins := InstanceType{}
			if memberType.ValueType() == ins.ValueType() {
				for _, signature := range memberType.(InstanceType).ClassType.MemberSignatures {
					if signature.Name == propertyName {
						return signature.ReturnType, propertyName
					}
				}
			}
			// 3.暂时没处理
			return UnKnownType{}, ""
		}

	// 这里是调用数组或者调用dict
	case "array_dict":
	}
	return UnKnownType{}, ""
}

func (this *SemanticAnalysisV2) visitClassExpression(node parser.Node) AllType {
	// 1.检查 super 的 class 是否存在
	if node.Type() != parser.AstTypeClassExpression.Name() {
		return UnKnownType{}
	}

	classNameExp := node.(parser.ClassExpression).Name
	var className string
	var classType AllType
	classType, className, _ = this.visitIdentifier(classNameExp)

	// 这里是检测 这个 class 有没有和其他变量重名
	UnKnownClassType := UnKnownType{}
	if classType.ValueType() != UnKnownClassType.ValueType() {
		return UnKnownType{}
	}

	superClass := node.(parser.ClassExpression).SuperClass
	superClassName := ""
	var superClassType AllType
	// 1.检查 super 的 class 是否存在
	if superClass != nil {
		superClassType, superClassName, _ = this.visitIdentifier(superClass)

		// 表示现在找不到这个类
		if superClassType.ValueType() == UnKnownClassType.ValueType() {
			utils.LogError("super class not found", superClassName)
			return UnKnownType{}
		}
	}

	// 开始新的作用域
	funcScope := NewScopeV2()
	funcScope.SetParent(this.CurrentScope)
	this.CurrentScope = funcScope

	// 现在开始处理这个类
	classBody := node.(parser.ClassExpression).Body
	memberSignatures := make([]Signature, 0)
	utils.LogInfo("visitClassExpression classBody", classBody.Type())
	if classBody.Type() == parser.AstTypeClassBodyStatement.Name() {
		memberSignatures = this.visitClassBodyStatement(classBody)
	}
	thisClassType := ClassType{
		MemberSignatures: memberSignatures,
		SuperType:        superClassType,
	}
	// 回滚到上级作用域
	this.CurrentScope = this.CurrentScope.Parent

	// 类不改变, 直接赋值 const
	// 在父级放 class
	this.CurrentScope.AddSignature(className, thisClassType, false, "const")
	return thisClassType
}

func (this *SemanticAnalysisV2) visitClassBodyStatement(node parser.Node) []Signature {
	if node.Type() != parser.AstTypeClassBodyStatement.Name() {
		return nil
	}
	signatures := make([]Signature, 0)
	for _, item := range node.(parser.ClassBodyStatement).Body {
		utils.LogInfo("visitClassBodyStatement visit item", item.Type())
		switch item.Type() {
		case parser.AstTypeClassVariableDeclaration.Name():
			signature := this.visitClassVariableDeclaration(item)
			signatures = append(signatures, signature)
		default:
			utils.LogError("unknown class body statement", item.Type())
		}
	}
	ifHasNewFunc := false
	for _, signature := range signatures {
		if signature.Name == lexer.TokenTypeNew.Name() {
			ifHasNewFunc = true
			break
		}
	}
	if !ifHasNewFunc {
		utils.LogError("class init has not new func", node.Type())
		return nil
	}
	return signatures
}

func (this *SemanticAnalysisV2) visitClassVariableDeclaration(node parser.Node) Signature {
	//type visitClassVariableDeclaration struct {
	//	Kind  string // kind属性
	//	Name  Node   // name属性
	//	Value Node   // value属性
	//}
	left := node.(parser.ClassVariableDeclaration).Name
	var variableName string
	switch left.Type() {
	case parser.AstTypeMemberExpression.Name():
		_, variableName = this.visitMemberExpression(left)
	case parser.AstTypeIdentifier.Name():
		_, variableName, _ = this.visitIdentifier(left)
	default:
		utils.LogError("invalid class variable declaration", left)
		return Signature{}
	}

	var valueType AllType
	right := node.(parser.ClassVariableDeclaration).Value
	switch right.Type() {
	// 访问函数
	case parser.AstTypeFunctionExpression.Name():
		valueType = this.visitFunctionExpression(right)
	case parser.AstTypeBinaryExpression.Name():
		valueType = this.visitBinaryExpression(right)
	case parser.AstTypeNumberLiteral.Name():
		valueType = this.visitNumberLiteral(right)
	case parser.AstTypeNullLiteral.Name():
		valueType = this.visitNullLiteral(right)
	case parser.AstTypeBooleanLiteral.Name():
		valueType = this.visitBooleanLiteral(right)
	case parser.AstTypeStringLiteral.Name():
		valueType, _ = this.visitStringLiteral(right)
	case parser.AstTypeArrayLiteral.Name():
		valueType = this.visitArrayLiteral(right)
	case parser.AstTypeUnaryExpression.Name():
		valueType = this.visitUnaryExpression(right)
	case parser.AstTypeDictLiteral.Name():
		valueType = this.visitDictLiteral(right)
	case parser.AstTypeIdentifier.Name():
		var varName string
		_, varName, _ = this.visitIdentifier(right)
		symbol, ok := this.CurrentScope.LookupSignature(varName)
		if ok {
			valueType = symbol.ReturnType
		} else {
			utils.LogError("undeclared variable", varName)
			return Signature{}
		}
	default:
		utils.LogError("invalid class variable declaration", right)
		return Signature{}
	}
	// 先这么写 false
	this.CurrentScope.AddSignature(variableName, valueType, false, "const")

	fmt.Printf("visitClassVariableDeclaration this.CurrentScope: %+v\n", this.CurrentScope)
	return Signature{
		Name:       variableName,
		ReturnType: valueType,
		IsStatic:   false,
		VarType:    "const",
	}
}

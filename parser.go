package main

import (
	"fmt"
	"strconv"
)

type BaseParser struct {
	lexer    *TokenList
	Position int
	Cache    []Token
}

func NewBaseParser(lexer *TokenList) *BaseParser {
	return &BaseParser{
		lexer:    lexer,
		Position: 0,
		Cache:    []Token{},
	}
}

func (this *BaseParser) Current() Token {
	if len(this.Cache) == this.Position {
		this.Cache = append(this.Cache, this.lexer.nextToken())
	}
	return this.Cache[this.Position]
}

func (this *BaseParser) Next() Token {
	c := this.Current()
	this.Position += 1
	return c
}

func (this *BaseParser) Mark() int {
	return this.Position
}

func (this *BaseParser) Reset(pos int) {
	this.Position = pos
}

func (this *BaseParser) Expect(TokenType TokenType) bool {
	c := this.Current()
	return c.Type == TokenType
}

func (this *BaseParser) Match(TokenType TokenType) Token {
	if this.Expect(TokenType) {
		return this.Next()
	}
	return Token{
		Type: TokenTypeError,
	}
}

type SansLangParser struct {
	BaseParser
}

func NewSansLangParser(lexer *TokenList) *SansLangParser {
	return &SansLangParser{BaseParser: *NewBaseParser(lexer)}
}

func (this *SansLangParser) Parse() Program {
	programAst := this.astParseProgram()
	fmt.Printf("programAst %+v\n", programAst)
	eofToken := this.Match(TokenTypeEof)
	if !eofToken.Nil() {
		fmt.Errorf("Parse error current:%v pos:%v", this.Current(), this.Position)
	}
	return programAst
}

func (this *SansLangParser) astParseProgram() Program {
	mainAst := Program{
		Body: []Node{},
	}
	body := this.astParseStatements()
	mainAst.Body = body
	return mainAst
}

func (this *SansLangParser) astParseStatements() []Node {
	body := []Node{}
	fmt.Printf("astParseStatements %v\n", this.Current())
	for this.Current().Type != TokenTypeEof {
		subAst := this.astParseStatement()
		fmt.Printf("astParseStatements subAst %+v this.Current():%v\n", subAst, this.Current())
		if subAst != nil {
			body = append(body, subAst)
		} else {
			break
		}
	}
	return body
}

func (this *SansLangParser) astParseVariableDeclaration() Node {
	fmt.Printf("astParseVariableDeclaration %v \n", this.Current())
	if this.Expect(TokenTypeVar) || this.Expect(TokenTypeConst) {
		op := this.Next()
		id := this.astParseCallMemberExpression()
		if id != nil {
			assign := this.Match(TokenTypeAssign)
			if !assign.Error() {
				exp := this.astParseExpression()
				if exp != nil {
					return VariableDeclaration{Kind: op.Value, Name: id, Value: exp}
				}
			}
		}
	}
	return nil
}

func (this *SansLangParser) astParseStatement() Node {
	// statement ->
	//              | variableDeclaration
	//              | blockStatement
	//              | expressionStatement
	//              | returnStatement
	//              | ifStatement
	//              | whileStatement
	//              | forStatement
	//              | breakStatement
	//              | continueStatement
	mark := this.Mark()
	ast := this.astParseVariableDeclaration()
	if ast != nil {
		return ast
	}
	this.Reset(mark)

	ast = this.astParseExpressionStatement()
	if ast != nil {
		return ast
	}
	this.Reset(mark)

	ast = this.astParseIfStatement()
	if ast != nil {
		return ast
	}
	this.Reset(mark)

	ast = this.astParseWhileStatement()
	if ast != nil {
		return ast
	}
	this.Reset(mark)

	ast = this.astParseForStatement()
	if ast != nil {
		return ast
	}
	this.Reset(mark)

	ast = this.astParseBreakStatement()
	if ast != nil {
		return ast
	}
	this.Reset(mark)

	ast = this.astParseContinueStatement()
	if ast != nil {
		return ast
	}
	this.Reset(mark)

	ast = this.astParseReturnStatement()
	if ast != nil {
		return ast
	}
	this.Reset(mark)
	return nil
}

func (this *SansLangParser) astParseAssignment() Node {
	key := this.astParseString()
	if key != nil && this.Expect(TokenTypeColon) {
		value := this.Match(TokenTypeColon)
		if !value.Error() {
			exp := this.astParseExpression()
			if exp != nil {
				return PropertyAssignment{Key: key, Value: exp}
			}
		}
	}
	return nil
}

func (this *SansLangParser) astParseClassExpression() Node {
	// classExpression: 'class' identifier super ('(' identifier? ')') classBodyStatement
	mark := this.Mark()
	classToken := this.Match(TokenTypeClass)
	if !classToken.Error() {
		id := this.astParseIdentifier()
		if id != nil {
			var superClass Node
			if this.Expect(TokenTypeSuper) {
				superToken := this.Match(TokenTypeSuper)
				if !superToken.Error() {
					this.Match(TokenTypeLParen)
					superClass = this.astParseIdentifier()
					this.Match(TokenTypeRParen)
				}
			}
			body := this.astParseClassBody()
			return ClassExpression{
				Name:       id,
				SuperClass: superClass,
				Body:       body,
			}
		}
	}
	this.Reset(mark)
	return nil
}

func (this *SansLangParser) astParseClassBody() Node {
	// classBodyStatement: '{' classBodyStatements '}'
	lb := this.Match(TokenTypeLBrace)
	if !lb.Error() {
		body := this.astParseClassBodyStatements()
		this.Match(TokenTypeRBrace)
		return ClassBodyStatement{Body: body}
	}
	return nil
}

func (this *SansLangParser) astParseClassBodyStatements() []Node {
	// 处理不同的函数定义
	body := []Node{}
	fmt.Printf("astParseClassBodyStatements %v\n", this.Current())
	for this.Current().Type != TokenTypeRBrace {
		subAst := this.astParseClassBodyStatement()
		fmt.Printf("astParseClassBodyStatements subAst %+v this.Current():%v\n", subAst, this.Current())
		if subAst != nil {
			body = append(body, subAst)
		} else {
			break
		}
	}
	if body != nil && len(body) != 0 {
		return body
	}
	return nil
}

func (this *SansLangParser) astParseClassBodyStatement() Node {
	mark := this.Mark()
	ast := this.astParseClassVariableDeclaration()
	if ast != nil {
		return ast
	}
	this.Reset(mark)

	ast = this.astParseClassExpressionStatement()
	if ast != nil {
		return ast
	}
	this.Reset(mark)
	return nil
}

func (this *SansLangParser) astParseClassVariableDeclaration() Node {
	fmt.Printf("astParseClassVariableDeclaration %v \n", this.Current())
	// const this.age = 1
	// const cls.age = 1
	// const cls.new = function(){}
	// const new = function() {}
	if this.Expect(TokenTypeConst) || this.Expect(TokenTypeVar) {
		op := this.Next()
		id := this.astParseCallClassMemberExpression()
		if id != nil {
			assign := this.Match(TokenTypeAssign)
			if !assign.Error() {
				exp := this.astParseExpression()
				if exp != nil {
					return ClassVariableDeclaration{Kind: op.Value, Name: id, Value: exp}
				}
			}
		}
	}
	return nil
}

func (this *SansLangParser) astParseClassExpressionStatement() Node {
	fmt.Printf("astParseClassExpressionStatement %v \n", this.Current())
	exp := this.astParseExpression()
	if exp != nil {
		fmt.Printf("astParseExpressionStatement %v\n", exp)
		return exp
	}
	return nil
}

func (this *SansLangParser) astParseBlockStatement() Node {
	lbraceToken := this.Match(TokenTypeLBrace)
	if !lbraceToken.Error() {
		body := []Node{}
		for this.Current().Type != TokenTypeRBrace {
			subAst := this.astParseStatement()
			fmt.Printf("astParseBlockStatement subAst %+v this.Current():%v\n", subAst, this.Current())
			if subAst != nil {
				body = append(body, subAst)
			} else {
				break
			}
		}
		this.Match(TokenTypeRBrace)
		return BlockStatement{Body: body}
	}
	return nil
}

func (this *SansLangParser) astParseContinueStatement() Node {
	continueToken := this.Match(TokenTypeContinue)
	if !continueToken.Error() {
		return ContinueStatement{}
	}
	return nil
}

func (this *SansLangParser) astParseBreakStatement() Node {
	fmt.Printf("astParseBreakStatement %v\n", this.Current())
	breakToken := this.Match(TokenTypeBreak)
	if !breakToken.Error() {
		return BreakStatement{}
	}
	return nil
}

func (this *SansLangParser) astParseReturnStatement() Node {
	fmt.Printf("astParseReturnStatement %v\n", this.Current())
	returnToken := this.Match(TokenTypeReturn)
	if !returnToken.Error() {
		return ReturnStatement{}
	}
	return nil
}

func (this *SansLangParser) astParseIfStatement() Node {
	// ifStatement -> 'if' '(' expression ')' blockStatement ('else' (blockStatement | ifStatement) )?
	if this.Expect(TokenTypeIf) {
		this.Match(TokenTypeIf)
		this.Match(TokenTypeLParen)
		condition := this.astParseExpression()
		if condition != nil {
			rp := this.Match(TokenTypeRParen)
			if !rp.Error() {
				consequent := this.astParseBlockStatement()
				if consequent != nil {
					pos := this.Mark()
					var alternate Node
					if this.Expect(TokenTypeElse) {
						this.Match(TokenTypeElse)
						alternate = this.astParseBlockStatement()
						if alternate == nil {
							alternate = this.astParseIfStatement()
						} else {
							this.Reset(pos)
						}
					}
					return IfStatement{Condition: condition, Consequent: consequent, Alternate: alternate}
				}
			}
		}
	}
	return nil
}

func (this *SansLangParser) astParseWhileStatement() Node {
	// ifStatement -> 'if' '(' expression ')' blockStatement ('else' (blockStatement | ifStatement) )?
	if this.Expect(TokenTypeWhile) {
		this.Match(TokenTypeWhile)
		this.Match(TokenTypeLParen)
		condition := this.astParseExpression()
		if condition != nil {
			rp := this.Match(TokenTypeRParen)
			if !rp.Error() {
				body := this.astParseBlockStatement()
				if body != nil {
					return WhileStatement{Condition: condition, Body: body}
				}
			}
		}
	}
	return nil
}

func (this *SansLangParser) astParseForStatement() Node {
	// forStatement: 'for' '(' (expressionStatement | variableDeclaration)? ';' expression? ';' expression? ')' blockStatement
	if this.Expect(TokenTypeFor) {
		this.Match(TokenTypeFor)
		this.Match(TokenTypeLParen)
		init := this.astParseExpressionStatement()
		if init == nil {
			init = this.astParseVariableDeclaration()
		}
		semi := this.Match(TokenTypeSemi)
		if !semi.Error() {
			test := this.astParseExpression()
			if test != nil {
				semi = this.Match(TokenTypeSemi)
				if !semi.Error() {
					update := this.astParseExpression()
					if update != nil {
						this.Match(TokenTypeRParen)
						body := this.astParseBlockStatement()
						return ForStatement{Init: init, Test: test, Update: update, Body: body}
					}
				}
			}
		}
	}

	return nil
}

func (this *SansLangParser) astParseFunctionExpression() Node {
	// functionExpression -> 'function' '(' formalParameterList ')' blockStatement
	params := []Node{}
	//Cannot call a pointer method on 'this.Match(TokenTypeFunction)'
	funcToken := this.Match(TokenTypeFunction)
	if !funcToken.Error() {
		lparenToken := this.Match(TokenTypeLParen)
		if !lparenToken.Error() {
			params = this.astParseFormalParameterList()
			rparenToken := this.Match(TokenTypeRParen)
			if !rparenToken.Error() {
				body := this.astParseBlockStatement()
				if body != nil {
					return FunctionExpression{Params: params, Body: body}
				}
			}
		}
	}
	return nil
}

func (this *SansLangParser) astParseFormalParameterList() []Node {
	// formalParameterList -> (id ',')*
	// TODO: 默认参数
	params := []Node{}
	for this.Expect(TokenTypeId) {
		id := this.astParseIdentifier()
		params = append(params, id)
		if this.Expect(TokenTypeComma) {
			this.Match(TokenTypeComma)
		} else {
			break
		}
	}
	return params
}

func (this *SansLangParser) astParseArgsWithParen() []Node {
	// todo 支持默认参数
	// ()
	args := []Node{}
	if this.Expect(TokenTypeLParen) {
		this.Match(TokenTypeLParen)
		for this.Expect(TokenTypeNumeric) || this.Expect(TokenTypeId) || this.Expect(TokenTypeString) {
			// 字面量
			arg := this.astParseLiteral()
			if arg != nil {
				args = append(args, arg)
			}
			// ,
			t := this.Match(TokenTypeComma)
			if t.Error() {
				break
			}
		}
		t := this.Match(TokenTypeRParen)
		if !t.Error() {
			return args
		}
	}
	return nil
}

func (this *SansLangParser) astParseCallMemberExpression() Node {
	// factor arguments callMemberExpressionTail
	// factor '.' identifier callMemberExpressionTail;
	// factor '[' expression ']' callMemberExpressionTail
	subAst := this.astParseFactor()
	mark := this.Mark()
	if subAst != nil {
		args := this.astParseArgsWithParen()
		if args != nil {
			n := CallExpression{Object: subAst, Args: args}
			node := this.astParseCallMemberExpressionTail(n)
			return node
		}
		// 点语法
		if this.Expect(TokenTypeDot) {
			this.Match(TokenTypeDot)
			prop := this.astParseIdentifier()
			if prop != nil {
				node := MemberExpression{
					Object:      subAst,
					Property:    prop,
					ElementType: "dot",
				}
				return this.astParseCallMemberExpressionTail(node)
			}
		}
		this.Reset(mark)
		// 数组
		if this.Expect(TokenTypeLBracket) {
			this.Match(TokenTypeLBracket)
			prop := this.astParseExpression()
			if prop != nil {
				this.Match(TokenTypeRBracket)
				node := MemberExpression{
					Object:      subAst,
					Property:    prop,
					ElementType: "array",
				}
				return this.astParseCallMemberExpressionTail(node)
			}
		}
	}
	this.Reset(mark)
	return subAst
}

func (this *SansLangParser) astParseCallClassMemberExpression() Node {
	// factor arguments callMemberExpressionTail
	// factor '.' identifier callMemberExpressionTail;
	// factor '[' expression ']' callMemberExpressionTail
	subAst := this.astParseFactor()
	mark := this.Mark()
	if subAst != nil {
		args := this.astParseArgsWithParen()
		if args != nil {
			n := CallExpression{Object: subAst, Args: args}
			node := this.astParseCallClassMemberExpressionTail(n)
			return node
		}
		// 点语法
		if this.Expect(TokenTypeDot) {
			this.Match(TokenTypeDot)
			prop := this.astParseIdentifier()
			if prop != nil {
				node := MemberExpression{
					Object:      subAst,
					Property:    prop,
					ElementType: "dot",
				}
				return this.astParseCallClassMemberExpressionTail(node)
			}
		}
		this.Reset(mark)
		// 数组
		if this.Expect(TokenTypeLBracket) {
			this.Match(TokenTypeLBracket)
			prop := this.astParseExpression()
			if prop != nil {
				this.Match(TokenTypeRBracket)
				node := MemberExpression{
					Object:      subAst,
					Property:    prop,
					ElementType: "array",
				}
				return this.astParseCallClassMemberExpressionTail(node)
			}
		}
	}
	this.Reset(mark)
	return subAst
}

func (this *SansLangParser) astParseCallClassMemberExpressionTail(node Node) Node {
	// arguments callMemberExpressionTail
	// （'.' identifier callMemberExpressionTail）*
	// '[' expression ']' callMemberExpressionTail *
	// 处理函数调用
	mark := this.Mark()
	args := this.astParseArgsWithParen()
	if args != nil {
		node = CallExpression{Object: node, Args: args}
		return this.astParseCallMemberExpressionTail(node)
	}
	// 处理点语法
	if this.Expect(TokenTypeDot) {
		this.Match(TokenTypeDot)
		prop := this.astParseIdentifier()
		if prop != nil {
			node = MemberExpression{
				Object:      node,
				Property:    prop,
				ElementType: "dot",
			}
			return this.astParseCallMemberExpressionTail(node)
		}
	}
	this.Reset(mark)
	// 处理数组
	if this.Expect(TokenTypeLBracket) {
		this.Match(TokenTypeLBracket)
		prop := this.astParseExpression()
		if prop != nil {
			this.Match(TokenTypeRBracket)
			node = MemberExpression{
				Object:      node,
				Property:    prop,
				ElementType: "array",
			}
			return this.astParseCallMemberExpressionTail(node)
		}
	}
	this.Reset(mark)
	fmt.Printf("astParseCallMemberTail %v\n", this.Current())
	return node
}

func (this *SansLangParser) astParseCallMemberExpressionTail(node Node) Node {
	// arguments callMemberExpressionTail
	// （'.' identifier callMemberExpressionTail）*
	// '[' expression ']' callMemberExpressionTail *
	// 处理函数调用
	mark := this.Mark()
	args := this.astParseArgsWithParen()
	if args != nil {
		node = CallExpression{Object: node, Args: args}
		return this.astParseCallMemberExpressionTail(node)
	}
	// 处理点语法
	if this.Expect(TokenTypeDot) {
		this.Match(TokenTypeDot)
		prop := this.astParseIdentifier()
		if prop != nil {
			node = MemberExpression{
				Object:      node,
				Property:    prop,
				ElementType: "dot",
			}
			return this.astParseCallMemberExpressionTail(node)
		}
	}
	this.Reset(mark)
	// 处理数组
	if this.Expect(TokenTypeLBracket) {
		this.Match(TokenTypeLBracket)
		prop := this.astParseExpression()
		if prop != nil {
			this.Match(TokenTypeRBracket)
			node = MemberExpression{
				Object:      node,
				Property:    prop,
				ElementType: "array",
			}
			return this.astParseCallMemberExpressionTail(node)
		}
	}
	this.Reset(mark)
	fmt.Printf("astParseCallMemberTail %v\n", this.Current())
	return node
}

// not
func (this *SansLangParser) astParseNotExpression() Node {
	if this.Expect(TokenTypeNot) {
		notValue := this.Match(TokenTypeNot)
		if !notValue.Error() {
			rightAst := this.astParseCallMemberExpression()
			if rightAst != nil {
				return UnaryExpression{Value: rightAst, Operator: notValue.Value}
			}
		}
	}
	leftAst := this.astParseCallMemberExpression()
	return leftAst
}

// * /
func (this *SansLangParser) astParseMulDivExpression() Node {
	leftAst := this.astParseNotExpression()
	if leftAst != nil {
		for this.Expect(TokenTypeMul) || this.Expect(TokenTypeDiv) || this.Expect(TokenTypeMod) {
			op := this.Current()
			this.Next()
			rightAst := this.astParseNotExpression()
			if rightAst != nil {
				leftAst = BinaryExpression{Left: leftAst, Operator: op.Value, Right: rightAst}
			} else {
				break
			}
		}
	}
	fmt.Printf("astParseMulDivExpression %v\n", leftAst)
	return leftAst
}

// + -
func (this *SansLangParser) astParseAddSubExpression() Node {
	// 处理左值，如果左值有，才会继续往下走，否则直接返回 null
	// 这样是因为有可能处理到类似 a + 1 而不是 a * 1 的情况，
	// 如果是 a + 1, 就直接返回 a 就好了
	leftAst := this.astParseMulDivExpression()
	if leftAst != nil {
		for this.Expect(TokenTypePlus) || this.Expect(TokenTypeMinus) {
			op := this.Current()
			this.Next()
			rightAst := this.astParseMulDivExpression()
			if rightAst != nil {
				leftAst = BinaryExpression{Left: leftAst, Operator: op.Value, Right: rightAst}
			} else {
				break
			}
		}
	}
	fmt.Printf("astParseAddSubExpression %v\n", leftAst)
	return leftAst
}

// < <= > >=
func (this *SansLangParser) astParseCompareExpression() Node {
	leftAst := this.astParseAddSubExpression()
	if leftAst != nil {
		for this.Expect(TokenTypeLessThan) || this.Expect(TokenTypeGreaterThan) || this.Expect(TokenTypeLessThanEquals) || this.Expect(TokenTypeGreaterThanEquals) {
			op := this.Current()
			this.Next()
			rightAst := this.astParseAddSubExpression()
			if rightAst != nil {
				leftAst = BinaryExpression{Left: leftAst, Operator: op.Value, Right: rightAst}
			} else {
				break
			}
		}
	}
	fmt.Printf("astParseCompareExpression %v\n", leftAst)
	return leftAst
}

func (this *SansLangParser) astParseEqualsAndNotEqualExpression() Node {
	leftAst := this.astParseCompareExpression()
	if leftAst != nil {
		for this.Expect(TokenTypeNotEquals) || this.Expect(TokenTypeEquals) {
			op := this.Current()
			this.Next()
			rightAst := this.astParseCompareExpression()
			if rightAst != nil {
				leftAst = BinaryExpression{Left: leftAst, Operator: op.Value, Right: rightAst}
			} else {
				break
			}
		}
	}
	fmt.Printf("astParseEqualsAndNotEqualExpression %v\n", leftAst)
	return leftAst
}

func (this *SansLangParser) astParseAndOrExpression() Node {
	leftAst := this.astParseEqualsAndNotEqualExpression()
	if leftAst != nil {
		for this.Expect(TokenTypeAnd) || this.Expect(TokenTypeOr) {
			op := this.Current()
			this.Next()
			rightAst := this.astParseEqualsAndNotEqualExpression()
			if rightAst != nil {
				leftAst = BinaryExpression{Left: leftAst, Operator: op.Value, Right: rightAst}
			} else {
				break
			}
		}
	}
	fmt.Printf("astParseAndOrExpression %v\n", leftAst)
	return leftAst
}

func (this *SansLangParser) astParseAssignmentExpression() Node {
	leftAst := this.astParseAndOrExpression()
	if leftAst != nil {
		op := this.Current()
		if this.Expect(TokenTypeAssign) {
			this.Next()
			rightAst := this.astParseAndOrExpression()
			if rightAst != nil {
				return AssignmentExpression{Left: leftAst, Operator: op.Value, Right: rightAst}
			}
		}
	}
	fmt.Printf("astParseAssignmentExpression %v\n", leftAst)
	return leftAst
}

func (this *SansLangParser) astParseAssignExpression() Node {
	leftAst := this.astParseAssignmentExpression()
	if leftAst != nil {
		op := this.Current()
		if this.Expect(TokenTypePlusAssign) || this.Expect(TokenTypeMinusAssign) || this.Expect(TokenTypeMulAssign) || this.Expect(TokenTypeDivAssign) {
			this.Next()
			rightAst := this.astParseAssignmentExpression()
			if rightAst != nil {
				return BinaryExpression{Left: leftAst, Operator: op.Value, Right: rightAst}
			}
		}
	}
	fmt.Printf("astParseAssignExpression %v\n", leftAst)
	return leftAst
}

func (this *SansLangParser) astParseExpressionStatement() Node {
	exp := this.astParseExpression()
	if exp != nil {
		fmt.Printf("astParseExpressionStatement %v\n", exp)
		return exp
	}
	return nil
}

func (this *SansLangParser) astParseExpression() Node {
	// 优先级由高到低
	// expression ->
	// | identifier
	// | literal
	// | expression '[' expression ']'
	// | expression '.' identifier
	// | expression arguments
	// | ( '+' | '-' | '!' ) expression
	// | expression ( '*' | '/' | '%' ) expression
	// | expression ( '+' | '-' ) expression
	// | expression ( '<<' | '>>' ) expression
	// | expression ( '<' | '>' | '<=' | '>=' ) expression
	// | expression ( '==' | '!=' | '===' | '!==' ) expression
	// | expression '&' expression
	// | expression '|' expression
	// | expression '&&' expression
	// | expression '||' expression
	// | assignmentExpression
	// | expression ( '+=' | '-=' | '*=' | '/=' ) expression
	// | '(' expression ')'
	exp := this.astParseAssignExpression()
	if exp != nil {
		return exp
	}
	fmt.Printf("astParseExpression out %v\n", exp)
	return nil
}

func (this *SansLangParser) astParseLiteral() Node {

	number := this.astParseNumber()
	if number != nil {
		return number
	}

	str := this.astParseString()
	if str != nil {
		return str
	}

	boolValue := this.astParseBoolean()
	if boolValue != nil {
		return boolValue
	}
	nullValue := this.astParseNull()
	if nullValue != nil {
		return nullValue
	}
	arrayValue := this.astParseArray()
	if arrayValue != nil {
		return arrayValue
	}
	dictValue := this.astParseDict()
	if dictValue != nil {
		return dictValue
	}

	return nil
}

func (this *SansLangParser) astParseDict() Node {
	kvs := []Node{}
	if this.Expect(TokenTypeLBrace) {
		lb := this.Match(TokenTypeLBrace)
		if !lb.Error() {
			kv := this.astParseAssignment()
			if kv != nil {
				kvs = append(kvs, kv)
			}
			for this.Expect(TokenTypeComma) {
				comma := this.Match(TokenTypeComma)
				if !comma.Error() {
					kv = this.astParseAssignment()
					if kv != nil {
						kvs = append(kvs, kv)
					} else {
						break
					}
				} else {
					break
				}
			}
			rb := this.Match(TokenTypeRBrace)
			if !rb.Error() {
				return DictLiteral{Values: kvs}
			}
		}
	}
	return nil
}

func (this *SansLangParser) astParseFactor() Node {
	//  factor:
	//     '(' expression ')'
	//     | 'this'
	//     | 'class'
	//     | literal
	//     | identifier
	//     | functionExpression
	//     | classExpression
	mark := this.Mark()
	if this.Expect(TokenTypeLParen) {
		lp := this.Match(TokenTypeLParen)
		if !lp.Error() {
			exp := this.astParseExpression()
			if exp != nil {
				this.Match(TokenTypeRParen)
				return exp
			}
		}
	}
	this.Reset(mark)

	value := this.astParseLiteral()
	if value != nil {
		return value
	}

	value = this.astParseIdentifier()
	if value != nil {
		return value
	}

	value = this.astParseFunctionExpression()
	if value != nil {
		return value
	}
	value = this.astParseClassExpression()
	if value != nil {
		return value
	}

	if this.Expect(TokenTypeClass) {
		this.Match(TokenTypeClass)
		return ClassLiteral{}
	}
	return nil
}

func (this *SansLangParser) astParseIdentifier() Node {
	if this.Expect(TokenTypeId) {
		id := this.Match(TokenTypeId)
		return Identifier{Value: id.Value}
	}
	return nil
}

func (this *SansLangParser) astParseString() Node {
	if this.Expect(TokenTypeString) {
		id := this.Match(TokenTypeString)
		return StringLiteral{Value: id.Value}
	}
	return nil
}

func (this *SansLangParser) astParseNumber() Node {
	if this.Expect(TokenTypeNumeric) {
		id := this.Match(TokenTypeNumeric)
		floatValue, err := strconv.ParseFloat(id.Value, 64)
		if err != nil {
			fmt.Printf("Parse number error: %s\n", err)
			return nil
		}
		return NumberLiteral{Value: floatValue}
	}
	return nil
}

func (this *SansLangParser) astParseArray() Node {
	exps := []Node{}
	if this.Expect(TokenTypeLBracket) {
		lb := this.Match(TokenTypeLBracket)
		if !lb.Error() {
			exp := this.astParseExpression()
			exps = append(exps, exp)
			for this.Expect(TokenTypeComma) {
				comma := this.Match(TokenTypeComma)
				if !comma.Error() {
					exp = this.astParseExpression()
					if exp != nil {
						exps = append(exps, exp)
					} else {
						break
					}
				} else {
					break
				}
			}
			rb := this.Match(TokenTypeRBracket)
			if !rb.Error() {
				return ArrayLiteral{Values: exps}
			}
		}
	}
	return nil
}

func (this *SansLangParser) astParseNull() Node {
	if this.Expect(TokenTypeNull) {
		this.Match(TokenTypeNull)
		return NullLiteral{}
	}
	return nil
}

func (this *SansLangParser) astParseBoolean() Node {
	if this.Expect(TokenTypeBoolean) {
		boolValue := this.Match(TokenTypeBoolean)
		var v bool
		if boolValue.Value == "true" {
			v = true
		}
		return BooleanLiteral{
			Value: v,
		}
	}
	return nil
}

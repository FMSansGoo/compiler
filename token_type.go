package main

import (
	"fmt"
)

type TokenType struct {
	name  string
	value int64
}

var (
	validTokenTypes = []TokenType{}

	// KeywordType
	TokenTypeOr       = newTokenType("or", 0)
	TokenTypeAnd      = newTokenType("and", 1)
	TokenTypeNot      = newTokenType("not", 2)
	TokenTypeClass    = newTokenType("class", 3)
	TokenTypeWhile    = newTokenType("while", 4)
	TokenTypeFor      = newTokenType("for", 5)
	TokenTypeBreak    = newTokenType("break", 6)
	TokenTypeContinue = newTokenType("continue", 7)
	TokenTypeVar      = newTokenType("var", 8)
	TokenTypeIf       = newTokenType("if", 9)
	TokenTypeElse     = newTokenType("else", 10)
	TokenTypeConst    = newTokenType("const", 11)
	TokenTypeCon      = newTokenType("con", 12)
	TokenTypeFunction = newTokenType("function", 13)
	TokenTypeReturn   = newTokenType("return", 14)
	TokenTypeNull     = newTokenType("null", 15)
	TokenTypeSuper    = newTokenType("super", 16)

	// SpecialType
	TokenTypeNumeric = newTokenType("numeric", 17)
	TokenTypeBoolean = newTokenType("boolean", 18)
	TokenTypeString  = newTokenType("string", 19)
	TokenTypeComment = newTokenType("comment", 19)
	TokenTypeId      = newTokenType("id", 20)
	TokenTypeEof     = newTokenType("eof", 21)

	// NormalType
	TokenTypePlus              = newTokenType("plus", 22)
	TokenTypeMinus             = newTokenType("minus", 23)
	TokenTypeMul               = newTokenType("mul", 24)
	TokenTypeDiv               = newTokenType("div", 25)
	TokenTypeMod               = newTokenType("mod", 26) // %
	TokenTypePlusAssign        = newTokenType("plusAssign", 27)
	TokenTypeMinusAssign       = newTokenType("minusAssign", 28)
	TokenTypeMulAssign         = newTokenType("mulAssign", 29)
	TokenTypeDivAssign         = newTokenType("divAssign", 30)
	TokenTypeLParen            = newTokenType("lparen", 31)   // (
	TokenTypeRParen            = newTokenType("rparen", 32)   // )
	TokenTypeAssign            = newTokenType("assign", 33)   // =
	TokenTypeDot               = newTokenType("dot", 34)      // .
	TokenTypeSemi              = newTokenType("semi", 35)     //;
	TokenTypeComma             = newTokenType("comma", 36)    // ,
	TokenTypeColon             = newTokenType("colon", 37)    // :
	TokenTypeLBrace            = newTokenType("lbrace", 38)   // {
	TokenTypeRBrace            = newTokenType("rbrace", 39)   // }
	TokenTypeLBracket          = newTokenType("lbracket", 40) // [
	TokenTypeRBracket          = newTokenType("rbracket", 41) // ]
	TokenTypeEquals            = newTokenType("equals", 42)
	TokenTypeNotEquals         = newTokenType("notEquals", 43)
	TokenTypeLessThan          = newTokenType("lessThan", 44)
	TokenTypeGreaterThan       = newTokenType("greaterThan", 45)
	TokenTypeLessThanEquals    = newTokenType("lessThanEquals", 46)
	TokenTypeGreaterThanEquals = newTokenType("greaterThanEquals", 47)
	TokenTypeRightShift        = newTokenType("rightShift", 48)
	TokenTypeLeftShift         = newTokenType("leftShift", 49)
	TokenTypeBitAnd            = newTokenType("bitAnd", 50)
	TokenTypeBitOr             = newTokenType("bitOr", 51)
	TokenTypeBitNot            = newTokenType("bitNot", 52)

	// bool
	TokenTypeTrue  = newTokenType("true", 53)
	TokenTypeFalse = newTokenType("false", 54)

	// error
	TokenTypeError = newTokenType("error", 55)

	// class
	TokenTypeCls = newTokenType("cls", 56)
	// this
	TokenTypeThis = newTokenType("this", 56)
)

func newTokenType(name string, value int64) TokenType {
	o := TokenType{name: name, value: value}
	validTokenTypes = append(validTokenTypes, o)
	return o
}

func TokenTypeAll() []TokenType {
	return validTokenTypes
}

func (t TokenType) Value() int64 {
	if !t.valid() {
		panic(fmt.Errorf("invalid TokenType: (%+v)", t))
	}
	return t.value
}

func (t TokenType) ValuePtr() *int64 {
	v := t.Value()
	return &v
}

func (t TokenType) Name() string {
	if !t.valid() {
		panic(fmt.Errorf("invalid TokenType: (%+v)", t))
	}
	return t.name
}

func (t TokenType) NamePtr() *string {
	n := t.Name()
	return &n
}

func (t TokenType) valid() bool {
	for _, v := range TokenTypeAll() {
		if v == t {
			return true
		}
	}
	return false
}

func NewTokenTypeFromValue(i int64) TokenType {
	for _, v := range TokenTypeAll() {
		if v.value == i {
			return v
		}
	}
	panic(fmt.Errorf("invalid TokenType value: (%+v)", i))
}

func GetTokenTypeFromName(s string) TokenType {
	for _, v := range TokenTypeAll() {
		if v.name == s {
			return v
		}
	}
	panic(fmt.Errorf("invalid TokenType name: (%+v)", s))
}

func ValidateTokenTypeValue(i int64) bool {
	for _, v := range TokenTypeAll() {
		if v.value == i {
			return true
		}
	}
	return false
}

func ValidateTokenTypeName(s string) bool {
	for _, v := range TokenTypeAll() {
		if v.name == s {
			return true
		}
	}
	return false
}

func ValidateTokenTypeInKeyWordType(s string) bool {
	for _, v := range TokenTypeAll() {
		if s == v.name {
			if v.value >= TokenTypeOr.Value() && v.value <= TokenTypeSuper.Value() {
				return true
			}
		}

	}
	return false
}

func ValidateTokenTypeInBool(s string) bool {
	for _, v := range TokenTypeAll() {
		if s == v.name {
			if v.value == TokenTypeTrue.Value() || v.value == TokenTypeFalse.Value() {
				return true
			}
		}
	}
	return false
}

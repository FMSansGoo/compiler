package lexer

import (
	"fmt"
	"strings"
)

var EscapeCharacterDict = map[string]string{
	"'":  "'",
	`"`:  `"`,
	"\\": "\\",
	"b":  "\b",
	"f":  "\f",
	"n":  "\n",
	"r":  "\r",
	"t":  "\t",
	"v":  "\v",
}

type BaseLexer struct {
	Code     string
	Position int
}

func NewBaseLexer(code string) *BaseLexer {
	return &BaseLexer{
		Code:     code,
		Position: 0,
	}
}

func (l *BaseLexer) Current() string {
	//fmt.Printf("Current: %d\n", l.Position)
	return string(l.Code[l.Position])
}

func (l *BaseLexer) Peek(offset int) string {
	// - 1 为默认值
	if offset == -1 {
		offset = 1
	}
	return string(l.Code[l.Position+offset])
}

func (l *BaseLexer) Mark() int {
	return l.Position
}

func (l *BaseLexer) Reset(pos int) {
	l.Position = pos
}

func (l *BaseLexer) Expect(char string, offset int) bool {
	if offset == -1 {
		offset = 0
	}
	return l.Peek(offset) == char
}

func (l *BaseLexer) Match(char string) bool {
	if l.Expect(char, -1) {
		l.Advance(-1)
		return true
	} else {
		return false
	}
}

func (l *BaseLexer) Advance(offset int) (end bool) {
	if offset == -1 {
		offset = 1
	}
	l.Position += offset
	if l.Position >= len(l.Code) {
		l.Position = len(l.Code)
		end = true
	}
	return end
}

func (l *BaseLexer) Error(args ...string) {
	fmt.Printf("error: %s\n", strings.Join(args, " "))
}

type SansLangLexer struct {
	BaseLexer
}

func NewSansLangLexer(code string) *SansLangLexer {
	return &SansLangLexer{
		BaseLexer: *NewBaseLexer(code),
	}
}

func (this *SansLangLexer) isSpace(char string) bool {
	return strings.Contains(" \t\r\n", char)
}

func (this *SansLangLexer) skipSpace() {
	for this.isSpace(this.Current()) {
		end := this.Advance(-1)
		if end {
			return
		}
	}
}

func (this *SansLangLexer) isComment() bool {
	return this.Expect("/", -1) && this.Expect("/", 1)
}

func (this *SansLangLexer) skipComment() {
	for !this.Expect("\n", -1) && (this.Position < len(this.Code)) {
		end := this.Advance(-1)
		if end {
			return
		}
	}
}

func (this *SansLangLexer) isDigit(char string) bool {
	return strings.Contains("0123456789", char)
}

func (this *SansLangLexer) isString(char string) bool {
	return strings.Contains("`", char) || strings.Contains("\"", char) || strings.Contains("'", char)
}

func (this *SansLangLexer) isId(char string) bool {
	return strings.Contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_", char)
}

func (this *SansLangLexer) getId() string {
	mark := this.Mark()
	for this.isId(this.Current()) || this.isDigit(this.Current()) {
		end := this.Advance(-1)
		if end {
			return this.Code[mark:this.Position]
		}
	}
	return this.Code[mark:this.Position]
}

func (this *SansLangLexer) guessType(char string) (tt string) {
	switch {
	case ValidateTokenTypeInBool(char):
		return TokenTypeBoolean.Name()
	case ValidateTokenTypeInKeyWordType(char):
		return GetTokenTypeFromName(char).Name()
	default:
		return TokenTypeId.Name()
	}
}

func (this *SansLangLexer) string() string {
	// string: '"' [^\n]* '"' | '\'' [^\n]* '\'' | '`' .* '`'
	str := ""
	prefix := this.Current()
	this.Advance(-1)
	for !this.Expect(prefix, -1) && this.Position < len(this.Code) {
		if prefix != "`" && this.Expect("\n", -1) {
			this.Error("字符串解析错误")
		}
		if this.Match("\\") {
			v, ok := EscapeCharacterDict[this.Current()]
			if ok {
				str += v
			} else {
				str += this.Current()
			}
		} else {
			str += this.Current()
		}
		this.Advance(-1)
	}
	this.Advance(-1)
	return str
}
func (this *SansLangLexer) number() string {
	//fmt.Printf("in number?")
	mark := this.Mark()
	if this.Current() == "0" {
		end := this.Advance(-1)
		if end {
			return this.Code[mark:this.Position]
		}
	} else if strings.Contains("123456789", this.Current()) {
		end := this.Advance(-1)
		if end {
			return this.Code[mark:this.Position]
		}
		for this.isDigit(this.Current()) {
			end = this.Advance(-1)
			if end {
				return this.Code[mark:this.Position]
			}
		}
	}
	if mark != this.Mark() {
		// 前面有匹配到
		if this.Current() == "." {
			this.Advance(-1)
			for this.isDigit(this.Current()) {
				end := this.Advance(-1)
				if end {
					return this.Code[mark:this.Position]
				}
			}
		}
		return this.Code[mark:this.Mark()]
	}
	this.Reset(mark)
	return ""
}

func (this *SansLangLexer) nextToken() Token {
	for this.Position < len(this.Code) {
		switch {
		case this.isSpace(this.Current()):
			this.skipSpace()
		case this.isComment():
			this.skipComment()
		case this.isDigit(this.Current()):
			return NewToken(TokenTypeNumeric, this.number())
		case this.isString(this.Current()):
			return NewToken(TokenTypeString, this.string())
		case this.Match("!"):
			if this.Match("=") {
				return NewToken(TokenTypeNotEquals, "!=")
			}
		case this.Match("+"):
			if this.Match("=") {
				return NewToken(TokenTypePlusAssign, "+=")
			}
			return NewToken(TokenTypePlus, "+")
		case this.Match("-"):
			if this.Match("=") {
				return NewToken(TokenTypeMinusAssign, "-=")
			}
			return NewToken(TokenTypeMinus, "-")
		case this.Match("*"):
			if this.Match("=") {
				return NewToken(TokenTypeMulAssign, "*=")
			}
			return NewToken(TokenTypeMul, "*")
		case this.Match("%"):
			return NewToken(TokenTypeMod, "%")
		case this.Match("/"):
			if this.Match("=") {
				return NewToken(TokenTypeDivAssign, "/=")
			}
			return NewToken(TokenTypeDiv, "/")
		case this.Match("("):
			return NewToken(TokenTypeLParen, "(")
		case this.Match(")"):
			return NewToken(TokenTypeRParen, ")")
		case this.Match("{"):
			return NewToken(TokenTypeLBrace, "{")
		case this.Match("}"):
			return NewToken(TokenTypeRBrace, "}")
		case this.Match("["):
			return NewToken(TokenTypeLBracket, "[")
		case this.Match("]"):
			return NewToken(TokenTypeRBracket, "]")
		case this.Match("<"):
			if this.Match("=") {
				return NewToken(TokenTypeLessThanEquals, "<=")
			}
			return NewToken(TokenTypeLessThan, "<")
		case this.Match(">"):
			if this.Match("=") {
				return NewToken(TokenTypeGreaterThanEquals, ">=")
			}
			return NewToken(TokenTypeGreaterThan, ">")
		case this.Match("="):
			if this.Match("=") {
				return NewToken(TokenTypeEquals, "==")
			}
			return NewToken(TokenTypeAssign, "=")
		case this.Match("."):
			return NewToken(TokenTypeDot, ".")
		case this.Match(";"):
			return NewToken(TokenTypeSemi, ";")
		case this.Match(":"):
			return NewToken(TokenTypeColon, ":")
		case this.Match(","):
			return NewToken(TokenTypeComma, ",")
		case this.isId(this.Current()):
			id := this.getId()
			idType := this.guessType(id)
			fmt.Printf("id:%v idType:%s\n", id, idType)
			if idType == TokenTypeBoolean.Name() {
				return NewToken(TokenTypeBoolean, id)
			} else {
				return NewToken(GetTokenTypeFromName(idType), id)
			}
		default:
			// 如果都不是，则证明遇到未知的字符，需要报错
			this.Error(fmt.Sprintf("无法识别的字符 char：%s, pos:%d", this.Current(), this.Position))
			return NewToken(TokenTypeEof, TokenTypeEof.name)
		}
	}
	return NewToken(TokenTypeEof, TokenTypeEof.name)
}

func (this *SansLangLexer) TokenList() []Token {
	token := this.nextToken()
	ret := []Token{}
	for token.Type != TokenTypeEof {
		ret = append(ret, token)
		token = this.nextToken()
	}
	ret = append(ret, token)
	return ret

}

type TokenList struct {
	Tokens []Token
	index  int
}

func (tl *TokenList) NextToken() Token {
	//fmt.Printf("Current: %d\n", l.Position)
	if tl.index == len(tl.Tokens) {
		return Token{
			Type: TokenTypeError,
		}
	}
	var t = tl.Tokens[tl.index]
	tl.index += 1
	return t
}

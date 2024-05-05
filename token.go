package main

type Token struct {
	Type  TokenType
	Value string
}

func (t *Token) String() string {
	return t.Type.Name()
}

func (t *Token) Nil() bool {
	return t.Type.name == "" || t.Type.value < 0
}

func (t *Token) Error() bool {
	return t.Type == TokenTypeError
}

func NewToken(TokenType TokenType, value string) Token {
	o := Token{Type: TokenType, Value: value}
	return o
}

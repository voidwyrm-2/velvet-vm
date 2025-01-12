package tokens

import "fmt"

type TokenType uint8

const (
	None TokenType = iota
	Number
	Address
	String
	Bool
	Ident
	Directive
	Label
	OpenBracket
	CloseBracket
)

func (tt TokenType) Str() string {
	return []string{
		"None",
		"Number",
		"Address",
		"String",
		"Bool",
		"Ident",
		"Directive",
		"Label",
		"OpenBracket",
		"CloseBracket",
	}[tt]
}

type Token struct {
	kind      TokenType
	lit       string
	start, ln int
}

func New(kind TokenType, lit string, start, ln int) Token {
	return Token{kind: kind, lit: lit, start: start, ln: ln}
}

func Empty() Token {
	return New(None, "", -1, -1)
}

func (t Token) GetKind() TokenType {
	return t.kind
}

func (t Token) GetLn() int {
	return t.ln
}

func (t Token) IsKind(kind TokenType) bool {
	return t.kind == kind
}

func (t Token) IsLit(lit string) bool {
	return t.lit == lit
}

func (t Token) Str() string {
	return fmt.Sprintf("{%s, '%s', %d, %d}", t.kind.Str(), t.lit, t.start, t.ln)
}

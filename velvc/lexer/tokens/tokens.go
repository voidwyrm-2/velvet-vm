package tokens

type TokenType uint8

const (
	None TokenType = iota
	Number
	String
	Ident
	Label
	Directive
)

type Token struct {
	kind           TokenType
	lit            string
	start, end, ln int
}

func New(kind TokenType, lit string, start, end, ln int) Token {
	return Token{kind: kind, lit: lit, start: start, end: end, ln: ln}
}

func Empty() Token {
	return New(None, "", -1, -1, -1)
}

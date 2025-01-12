package lexer

import (
	"fmt"
	"strings"

	"github.com/voidwyrm-2/velvet-vm/velvc/lexer/tokens"
)

func isNum(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isIdent(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '.' || ch == '-' || ch == '/'
}

type Lexer struct {
	text         string
	idx, col, ln int
	ch           rune
}

func New(text string) Lexer {
	return Lexer{text: text, idx: -1, col: 0, ln: 1, ch: -1}
}

func (l Lexer) errfPos(ln, col int, format string, a ...any) error {
	l.ln = ln
	l.col = col
	return l.errf(format, a...)
}

func (l Lexer) errf(format string, a ...any) error {
	return fmt.Errorf(fmt.Sprintf("error on line %d, col %d: ", l.ln, l.col)+format, a...)
}

func (l Lexer) charTok(kind tokens.TokenType) tokens.Token {
	return tokens.New(kind, string(l.ch), l.col, l.ln)
}

func (l *Lexer) advance() {
	l.idx++
	l.col++
	if l.idx < len(l.text) {
		l.ch = rune(l.text[l.idx])
	} else {
		l.ch = -1
	}

	if l.ch == '\n' {
		l.ln++
		l.col = 1
	}
}

func (l Lexer) peek() rune {
	if l.idx+1 < len(l.text) {
		return rune(l.text[l.idx+1])
	}
	return -1
}

func (l Lexer) isNum() bool {
	return isNum(l.ch)
}

func (l Lexer) isIdent() bool {
	return l.isNum() || isIdent(l.ch)
}

/*
0: normal;
1: negitive;
2: address;
else: panic
*/
func (l *Lexer) collectNumber(kind uint8) (tokens.Token, error) {
	start := l.col
	startln := l.ln
	s := ""

	if kind == 2 {
		l.advance()
	}

	for l.ch != -1 && (l.isNum() || l.ch == '_') {
		s += string(l.ch)
		l.advance()
	}

	tkind := tokens.Number
	switch kind {
	case 0:
	case 1:
		s = "-" + s
	case 2:
		tkind = tokens.Address
	default:
		panic(fmt.Sprintf("invalid kind %d", kind))
	}

	if s[0] == '_' {
		if kind == 2 {
			return tokens.Token{}, l.errfPos(start-1, startln, "number literals cannot start with underscores")
		}
		return tokens.Token{}, l.errfPos(start, startln, "number literals cannot start with underscores")
	} else if s[len(s)-1] == '_' {
		return tokens.Token{}, l.errfPos(l.idx-start, startln, "number literals cannot end with underscores")
	}

	return tokens.New(tkind, strings.ReplaceAll(s, "_", ""), start, startln), nil
}

/*
0: normal;
1: directive;
2: label;
else: panic
*/
func (l *Lexer) collectIdent(kind uint8) tokens.Token {
	start := l.col
	startln := l.ln
	s := ""

	if kind == 1 || kind == 2 {
		l.advance()
	}

	for l.ch != -1 && ((l.ch >= 'a' && l.ch <= 'z') || (l.ch >= 'A' && l.ch <= 'Z') || l.ch == '_' || l.ch == '.' || l.ch == '-') {
		s += string(l.ch)
		l.advance()
	}

	tkind := tokens.Ident
	switch kind {
	case 0:
	case 1:
		tkind = tokens.Directive
	case 2:
		tkind = tokens.Label
	default:
		panic(fmt.Sprintf("invalid kind %d", kind))
	}

	if s == "true" || s == "false" {
		tkind = tokens.Bool
	}

	return tokens.New(tkind, s, start, startln)
}

func (l *Lexer) collectString() (tokens.Token, error) {
	start := l.col
	startln := l.ln
	s := ""
	escaped := false

	l.advance()

	for l.ch != -1 && l.ch != '\n' && l.ch != '"' {
		if escaped {
			switch l.ch {
			case '\\', '"', '\'':
				s += string(l.ch)
			case 'n':
				s += "\n"
			case 'r':
				s += "\r"
			case 't':
				s += "\t"
			case '0':
				{
					s += " "
					b := []byte(s)
					b[len(b)-1] = 0
					s = string(b)
				}
			default:
				return tokens.Token{}, l.errfPos(startln, start, "invalid escape character '%c'", l.ch)
			}
			escaped = false
		} else if l.ch == '\\' {
			escaped = true
		} else {
			s += string(l.ch)
		}
		l.advance()
	}

	if l.ch != '"' {
		return tokens.Token{}, l.errfPos(startln, start, "unterminated string literal")
	}

	l.advance()

	return tokens.New(tokens.String, s, start, startln), nil
}

func (l *Lexer) Lex() ([]tokens.Token, error) {
	toks := []tokens.Token{}

	if l.idx == -1 {
		l.advance()
	}

	for l.ch != -1 {
		switch l.ch {
		case ' ', '\n', '\r', '\t':
			l.advance()
		case '[':
			toks = append(toks, l.charTok(tokens.OpenBracket))
			l.advance()
		case ']':
			toks = append(toks, l.charTok(tokens.CloseBracket))
			l.advance()
		case '"':
			if tok, err := l.collectString(); err != nil {
				return []tokens.Token{}, err
			} else {
				toks = append(toks, tok)
			}
		default:
			if l.ch == '/' && l.peek() == '/' {
				for l.ch != -1 && l.ch != '\n' {
					l.advance()
				}
			} else if l.ch == '/' && l.peek() == '*' {
				start, startln := l.col, l.ln
				ended := false

				for l.ch != -1 {
					if l.ch == '*' && l.peek() == '/' {
						ended = true
						break
					}
					l.advance()
				}

				if !ended {
					return []tokens.Token{}, l.errfPos(startln, start, "unterminated multiline comment")
				}
			} else if l.isNum() {
				tok, err := l.collectNumber(0)
				if err != nil {
					return []tokens.Token{}, err
				}
				toks = append(toks, tok)
			} else if l.ch == '&' && isNum(l.peek()) {
				tok, err := l.collectNumber(2)
				if err != nil {
					return []tokens.Token{}, err
				}
				toks = append(toks, tok)
			} else if l.ch == '-' && isNum(l.peek()) {
				tok, err := l.collectNumber(1)
				if err != nil {
					return []tokens.Token{}, err
				}
				toks = append(toks, tok)
			} else if l.ch == '@' && isIdent(l.peek()) {
				toks = append(toks, l.collectIdent(1))
			} else if l.ch == '.' && isIdent(l.peek()) {
				toks = append(toks, l.collectIdent(2))
			} else if l.isIdent() {
				toks = append(toks, l.collectIdent(0))
			} else {
				return []tokens.Token{}, l.errf("illegal character '%s'", string(l.ch))
			}
		}
	}

	return toks, nil
}

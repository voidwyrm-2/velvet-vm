package parser

import (
	"slices"

	"github.com/voidwyrm-2/velvet-vm/velvc/lexer/tokens"
	"github.com/voidwyrm-2/velvet-vm/velvc/parser/nodes"
	"github.com/voidwyrm-2/velvet-vm/velvc/parser/nodes/directive"
	"github.com/voidwyrm-2/velvet-vm/velvc/parser/nodes/halt"
	"github.com/voidwyrm-2/velvet-vm/velvc/parser/nodes/jump"
	"github.com/voidwyrm-2/velvet-vm/velvc/parser/nodes/label"
	"github.com/voidwyrm-2/velvet-vm/velvc/parser/nodes/otherinstruction"
	"github.com/voidwyrm-2/velvet-vm/velvc/parser/nodes/pushcall"
	"github.com/voidwyrm-2/velvet-vm/velvc/parser/nodes/setget"
)

func expect(tokens []tokens.Token, expected ...tokens.TokenType) error {
	if len(tokens) > len(expected) {
		t := tokens[len(expected)]
		return t.Err("expected EOL, but found '%s' instead", t.GetLit())
	} else if len(tokens) < len(expected) {
		t := tokens[len(tokens)-1]
		return t.Err("expected '%s', but found EOL instead", expected[len(tokens)].Str())
	}

	for i, t := range tokens {
		if !t.IsKind(expected[i]) {
			return t.Err("expected '%s', but found '%s' instead", expected[i].Str(), t.GetLit())
		}
	}
	return nil
}

type Parser struct {
	tokenL [][]tokens.Token
}

func New(tokenL [][]tokens.Token) Parser {
	filtered := [][]tokens.Token{}
	for _, tl := range tokenL {
		if len(tl) > 0 {
			filtered = append(filtered, tl)
		}
	}

	return Parser{tokenL: tokenL}
}

func (p Parser) Parse() ([]nodes.Node, error) {
	ns := []nodes.Node{}

	for _, l := range p.tokenL {
		head := l[0]
		l = l[1:]
		switch head.GetKind() {
		case tokens.Label:
			if err := expect(l); err != nil {
				return []nodes.Node{}, err
			}
			ns = append(ns, label.New(head))
		case tokens.Directive:
			switch head.GetLit() {
			case "vars":
				if err := expect(l, tokens.Number); err != nil {
					return []nodes.Node{}, err
				}
				ns = append(ns, directive.New(head, l))
			default:
				return []nodes.Node{}, head.Err("unknown directive '%s'", head.GetLit())
			}
		case tokens.Ident:
			switch head.GetLit() {
			case "halt":
				if err := expect(l, tokens.Number); err != nil {
					return []nodes.Node{}, err
				}
				ns = append(ns, halt.New(head, l[0]))
				continue
			case "push", "pusherr":
				if err := expect(l, tokens.Number); false {
				} else if err2 := expect(l, tokens.String); false {
				} else if err3 := expect(l, tokens.Number); false {
				} else if err4 := expect(l, tokens.Bool); false {
				} else if err5A := expect(l, tokens.Address); false {
				} else if err5B := expect(l, tokens.Address); false {
				} else if err6A := expect(l, tokens.Bool); false {
				} else {
					newL := make([]tokens.Token, len(l))
					copy(newL, l)
					slices.Reverse(newL)
					if err6B := expect(newL, tokens.Bool); false {
					} else if err == nil || err2 == nil || err3 == nil || err4 == nil || (err5A == nil && err5B == nil) || (err6A == nil && err6B == nil) {
					} else if err5A == nil && err5B != nil {
						return []nodes.Node{}, err5B
					} else if err6A == nil && err6B != nil {
						return []nodes.Node{}, err6B
					} else if err6A != nil && err6B == nil {
						return []nodes.Node{}, err6A
					} else {
						return []nodes.Node{}, err
					}
				}
				ns = append(ns, pushcall.New(head, l))
				continue
			case "error",
				"reset",
				"eq",
				"neq",
				"not",
				"lt",
				"gt",
				"lte",
				"gte",
				"add",
				"sub",
				"mul",
				"div",
				"pow",
				"log",
				"neg",
				"and",
				"or",
				"xor":
				if err := expect(l); err != nil {
					return []nodes.Node{}, err
				}
				ns = append(ns, pushcall.New(tokens.NewLit(tokens.Ident, "call"), []tokens.Token{head}))
				continue
			case "call":
				if err := expect(l, tokens.Ident); false {
				} else if err2A := expect(l, tokens.Address); false {
				} else if err2B := expect(l, tokens.Address, tokens.Address); err == nil || (err2A == nil && err2B == nil) {
					ns = append(ns, pushcall.New(head, l))
					continue
				} else if err2A == nil && err2B != nil {
					return []nodes.Node{}, err2B
				} else {
					return []nodes.Node{}, err
				}
			case "set", "get":
				if err := expect(l, tokens.Number); err != nil {
					return []nodes.Node{}, err
				}
				ns = append(ns, setget.New(head, l[0]))
				continue
			case "j", "jt", "jf", "je", "jne", "br", "brt", "brf", "bre", "brne":
				if err := expect(l, tokens.Ident); err != nil {
					return []nodes.Node{}, err
				}
				ns = append(ns, jump.New(head, l[0]))
				continue
			case "nop", "ret", "pop", "dup", "swap", "rot":
				if err := expect(l); err != nil {
					return []nodes.Node{}, err
				}
				ns = append(ns, otherinstruction.New(head))
				continue
			}
			fallthrough
		default:
			return []nodes.Node{}, head.Err("unexpected token '%s'", head.GetLit())
		}
	}

	return ns, nil
}

package directive

import (
	"fmt"
	"strings"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
	"github.com/voidwyrm-2/velvet-vm/velvc/lexer/tokens"
)

type DirectiveNode struct {
	name tokens.Token
	args []tokens.Token
}

func New(name tokens.Token, args []tokens.Token) DirectiveNode {
	return DirectiveNode{name: name, args: args}
}

func (dn DirectiveNode) Generate(ve *emitter.VelvEmitter) error {
	a := []any{}
	for _, arg := range dn.args {
		if v, err := arg.Convert(); err != nil {
			return err
		} else if ok := ve.HasLabel(arg.GetLit()); arg.IsKind(tokens.Ident) {
			if !ok {
				return arg.Err("label '%s' does not exist", arg.GetLit())
			}
			a = append(a, int(ve.GetLabel(arg.GetLit())))
		} else {
			a = append(a, v)
		}
	}
	ve.DoDirective(dn.name.GetLit(), a...)
	return nil
}

func (dn DirectiveNode) Str() string {
	formatted := []string{}
	for _, t := range dn.args {
		formatted = append(formatted, t.Str())
	}
	return fmt.Sprintf("{ins: %s, args: %s}", dn.name.Str(), strings.Join(formatted, ", "))
}

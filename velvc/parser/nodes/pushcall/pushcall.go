package pushcall

import (
	"fmt"
	"strings"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
	"github.com/voidwyrm-2/velvet-vm/velvc/lexer/tokens"
)

func assert[T any](v T, _ error) T {
	return v
}

type PushCallNode struct {
	instruction tokens.Token
	args        []tokens.Token
	ins         int
}

func New(instruction tokens.Token, args []tokens.Token) PushCallNode {
	return PushCallNode{instruction: instruction, args: args, ins: 0}
}

func (pcn PushCallNode) GenerateList() ([]any, error) {
	ls := []any{}
	for pcn.ins < len(pcn.args) {
		switch pcn.args[pcn.ins].GetKind() {
		case tokens.Number:
			ls = append(ls, assert(pcn.args[pcn.ins].Convert()).(int))
			/*case tokens.Address:
			if anyX, err := pcn.args[pcn.ins].Convert(); err != nil {
				return []any{},err
			} else if anyY, err := pcn.args[pcn.ins].Convert(); err != nil {
				return []any{},err
			} else {
				ve.Emit(emitter.Push, 0, uint16(anyX.(int)), uint16(anyY.(int)))
			}*/
		case tokens.Bool:
			ls = append(ls, assert(pcn.args[pcn.ins].Convert()).(bool))
		case tokens.String:
			ls = append(ls, assert(pcn.args[pcn.ins].Convert()).(string))
		case tokens.OpenBracket:
			if subls, err := pcn.GenerateList(); err != nil {
				return []any{}, err
			} else {
				ls = append(ls, subls)
			}
		case tokens.CloseBracket:
			return ls, nil
		}
		pcn.ins += 1
	}
	return ls, nil
}

func (pcn PushCallNode) Generate(ve *emitter.VelvEmitter) error {
	if pcn.instruction.IsLit("call") {
		if len(pcn.args) == 0 {
			ve.EmitNA(emitter.Call, 1)
		} else if pcn.args[0].IsKind(tokens.Address) {
			ve.Emit(emitter.Call, 0, uint16(assert(pcn.args[pcn.ins].Convert()).(int)), uint16(assert(pcn.args[pcn.ins+1].Convert()).(int)))
		} else {
			ve.EmitString(emitter.Call, 0, pcn.args[0].GetLit())
		}
	} else {
		switch pcn.args[pcn.ins].GetKind() {
		case tokens.Number:
			ve.Emit32(emitter.Push, 0, uint32(assert(pcn.args[pcn.ins].Convert()).(int)))
		case tokens.Address:
			ve.Emit(emitter.Push, 2, uint16(assert(pcn.args[pcn.ins].Convert()).(int)), uint16(assert(pcn.args[pcn.ins+1].Convert()).(int)))
		case tokens.Bool:
			ve.EmitString(emitter.Push, 1, assert(pcn.args[pcn.ins].Convert()).(string))
		case tokens.String:
			ve.EmitString(emitter.Push, 2, assert(pcn.args[pcn.ins].Convert()).(string))
		case tokens.Ident:
			ve.EmitString(emitter.Push, 4, assert(pcn.args[pcn.ins].Convert()).(string))

		case tokens.OpenBracket:
			if ls, err := pcn.GenerateList(); err != nil {
				return err
			} else {
				ve.EmitList(emitter.Push, 3, ls...)
			}
		}
	}
	return nil
}

func (pcn PushCallNode) Str() string {
	formatted := []string{}
	for _, t := range pcn.args {
		formatted = append(formatted, t.Str())
	}
	return fmt.Sprintf("{ins: %s, args: %s}", pcn.instruction.Str(), strings.Join(formatted, ", "))
}

package otherinstruction

import (
	"fmt"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
	"github.com/voidwyrm-2/velvet-vm/velvc/lexer/tokens"
)

type OtherInstructionNode struct {
	instruction tokens.Token
}

func New(instruction tokens.Token) OtherInstructionNode {
	return OtherInstructionNode{instruction: instruction}
}

func (oin OtherInstructionNode) Generate(ve *emitter.VelvEmitter) error {
	ve.EmitBasic(map[string]emitter.Opcode{
		"nop":  0,
		"ret":  1,
		"pop":  5,
		"dup":  6,
		"swap": 7,
		"rot":  8,
	}[oin.instruction.GetLit()])
	return nil
}

func (oin OtherInstructionNode) Str() string {
	return fmt.Sprintf("{ins: %s}", oin.instruction.Str())
}

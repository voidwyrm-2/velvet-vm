package jump

import (
	"fmt"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
	"github.com/voidwyrm-2/velvet-vm/velvc/lexer/tokens"
)

type JumpNode struct {
	instruction, label tokens.Token
}

func New(instruction, label tokens.Token) JumpNode {
	return JumpNode{instruction: instruction, label: label}
}

func (jn JumpNode) Generate(ve *emitter.VelvEmitter) error {
	ve.Emit32(emitter.Jump, map[string]uint8{"j": 0, "jt": 1, "jf": 2, "je": 3, "jne": 4}[jn.instruction.GetLit()], ve.GetLabel(jn.label.GetLit()))
	return nil
}

func (jn JumpNode) Str() string {
	return fmt.Sprintf("{ins: %s, label: %s}", jn.instruction.Str(), jn.label.Str())
}

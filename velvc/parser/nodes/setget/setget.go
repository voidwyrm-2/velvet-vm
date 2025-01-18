package setget

import (
	"fmt"
	"strconv"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
	"github.com/voidwyrm-2/velvet-vm/velvc/lexer/tokens"
)

func assert[T any](v T, _ error) T {
	return v
}

type SetgetNode struct {
	instruction, varIndex tokens.Token
}

func New(instruction, varIndex tokens.Token) SetgetNode {
	return SetgetNode{instruction: instruction, varIndex: varIndex}
}

func (sn SetgetNode) Generate(ve *emitter.VelvEmitter) error {
	var f uint8 = 0
	if sn.instruction.GetLit() == "get" {
		f = 1
	}

	ve.Emit32(emitter.Set, f, uint32(assert(strconv.Atoi(sn.varIndex.GetLit()))))

	return nil
}

func (sn SetgetNode) Str() string {
	return fmt.Sprintf("{ins: %s, varIndex: %s}", sn.instruction.Str(), sn.varIndex.Str())
}

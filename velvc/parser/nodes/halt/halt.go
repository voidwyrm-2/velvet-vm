package halt

import (
	"fmt"
	"strconv"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
	"github.com/voidwyrm-2/velvet-vm/velvc/lexer/tokens"
)

func assert[T any](v T, _ error) T {
	return v
}

type HaltNode struct {
	instruction, exitCode tokens.Token
}

func New(instruction, exitCode tokens.Token) HaltNode {
	return HaltNode{instruction: instruction, exitCode: exitCode}
}

func (hn HaltNode) Generate(ve *emitter.VelvEmitter) error {
	ve.EmitNF32(emitter.Halt, uint32(assert(strconv.Atoi(hn.exitCode.GetLit()))))
	return nil
}

func (hn HaltNode) Str() string {
	return fmt.Sprintf("{ins: %s, exitCode: %s}", hn.instruction.Str(), hn.exitCode.Str())
}

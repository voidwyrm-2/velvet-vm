package jump

import (
	"fmt"
	"strings"

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
	jumpType := jn.instruction.GetLit()
	isBranch := 0
	if strings.HasPrefix(jumpType, "br") {
		jumpType = jumpType[2:]
		isBranch = 1
	} else {
		jumpType = jumpType[1:]
	}

	if isBranch != 0 && isBranch != 1 {
		panic(fmt.Sprintf("isBranch is %d instead of 0 or 1", isBranch))
	}

	ve.Emit32(emitter.Jump, map[string]uint8{"": 0, "t": 1, "f": 2, "e": 3, "ne": 4}[jumpType]<<(3*isBranch), ve.GetLabel(jn.label.GetLit()))
	return nil
}

func (jn JumpNode) Str() string {
	return fmt.Sprintf("{ins: %s, label: %s}", jn.instruction.Str(), jn.label.Str())
}

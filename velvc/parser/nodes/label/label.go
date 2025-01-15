package label

import (
	"fmt"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
	"github.com/voidwyrm-2/velvet-vm/velvc/lexer/tokens"
)

type LabelNode struct {
	name tokens.Token
}

func New(name tokens.Token) LabelNode {
	return LabelNode{name: name}
}

func (ln LabelNode) Generate(ve *emitter.VelvetAsm) error {
	ve.CreateLabel(ln.name.GetLit())
	return nil
}

func (ln LabelNode) Str() string {
	return fmt.Sprintf("{name: %s}", ln.name.Str())
}

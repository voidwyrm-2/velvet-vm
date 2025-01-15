package generation

import (
	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
	"github.com/voidwyrm-2/velvet-vm/velvc/parser/nodes"
)

type Generator struct {
	nodes []nodes.Node
	ve    emitter.VelvetAsm
}

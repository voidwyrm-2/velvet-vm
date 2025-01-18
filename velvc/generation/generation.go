package generation

import (
	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
	"github.com/voidwyrm-2/velvet-vm/velvc/parser/nodes"
)

type Generator struct {
	nodes []nodes.Node
	ve    emitter.VelvEmitter
}

func New(nodes []nodes.Node, vars uint16) Generator {
	return Generator{nodes: nodes, ve: emitter.New(vars)}
}

func (g Generator) Generate() error {
	for _, n := range g.nodes {
		if err := n.Generate(&g.ve); err != nil {
			return err
		}
	}
	return nil
}

func (g Generator) Write(filename string) error {
	return g.ve.Write(filename)
}

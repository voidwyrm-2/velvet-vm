package nodes

import "github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"

type Node interface {
	Generate(ve *emitter.VelvEmitter) error
	Str() string
}

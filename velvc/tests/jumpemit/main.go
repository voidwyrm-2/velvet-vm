package main

import (
	"fmt"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
)

func main() {
	ve := emitter.New(0)

	ve.EmitNF32(emitter.Push, 10)
	ve.CreateLabel("loop")
	ve.EmitString(emitter.Push, 2, "Hello there.")
	ve.EmitString(emitter.Call, 0, "println")
	ve.Emit32(emitter.Jump, 0, ve.GetLabel("loop"))

	ve.Halt(0)

	if err := ve.Write("jumpemit-test.cvelv"); err != nil {
		fmt.Println(err.Error())
	}
}

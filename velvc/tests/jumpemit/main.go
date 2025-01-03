package main

import (
	"fmt"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
)

func main() {
	ve := emitter.New(0)

	// infinite loop that prints "Hello there."
	/*
		ve.CreateLabel("loop")
		ve.EmitString(emitter.Push, 2, "Hello there.")
		ve.EmitString(emitter.Call, 0, "println")
		ve.Emit32(emitter.Jump, 0, ve.GetLabel("loop"))
	*/

	// loop that prints "Hello there." ten times
	ve.EmitNF32(emitter.Push, 0)
	ve.CreateLabel("loop")
	ve.EmitString(emitter.Push, 2, "Hello there.")
	ve.EmitString(emitter.Call, 0, "println")
	ve.EmitBasic(emitter.Dup)
	ve.EmitNF32(emitter.Push, 1)
	ve.EmitString(emitter.Call, 0, "add")
	ve.EmitBasic(emitter.Dup)
	ve.EmitNF32(emitter.Push, 10)
	ve.EmitString(emitter.Call, 0, "lt")
	ve.Emit32(emitter.Jump, 1, ve.GetLabel("loop"))

	ve.Halt(0)

	if err := ve.Write("jumpemit-test.cvelv"); err != nil {
		fmt.Println(err.Error())
	}
}

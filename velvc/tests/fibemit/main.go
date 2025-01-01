package main

import (
	"fmt"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
)

func main() {
	ve := emitter.New(1)

	ve.CreateLabel("numloop")
	ve.EmitString(emitter.Push, 2, "please input a number: ")
	ve.EmitString(emitter.Call, 0, "print")
	ve.Emit32(emitter.Jump, 3, ve.GetLabel("numloop"))

	ve.EmitString(emitter.Call, 0, "println")
	ve.Halt(0)

	if err := ve.Write("fibemit-test.cvelv"); err != nil {
		fmt.Println(err.Error())
	}
}

package main

import (
	"fmt"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
)

func main() {
	ve := emitter.New(0)

	ve.EmitString(emitter.Push, 2, "hello there.")
	ve.Emit32(emitter.Push, 0, 4)
	ve.EmitString(emitter.Call, 0, "index")
	ve.EmitBasic(emitter.Dup)
	ve.EmitString(emitter.Call, 0, "println")
	ve.EmitString(emitter.Call, 0, "putcln")

	ve.EmitList(emitter.Push, 3, "LITTLE", "TIMMY", "FELL", "IN", "THE", "WELL")
	ve.Emit32(emitter.Push, 0, 1)
	ve.EmitString(emitter.Call, 0, "index")
	ve.EmitString(emitter.Call, 0, "println")

	ve.Halt(0)

	if err := ve.Write("index-test.cvelv"); err != nil {
		fmt.Println(err.Error())
	}
}

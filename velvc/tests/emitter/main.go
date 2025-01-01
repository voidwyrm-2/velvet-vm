package main

import (
	"fmt"

	"github.com/voidwyrm-2/velvet-vm/velvc/generation/emitter"
)

func main() {
	ve := emitter.New(0)
	ve.EmitString(emitter.Push, 2, "Hello there.")
	ve.EmitString(emitter.Call, 0, "println")
	ve.Halt(0)

	if err := ve.Write("emit-test.cvelv"); err != nil {
		fmt.Println(err.Error())
	}
}

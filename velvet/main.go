package main

import (
	"fmt"
	"io"
	"os"

	"github.com/voidwyrm-2/velvet-vm/velvet/vm"
)

func readFile(fileName string) ([]uint8, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return []uint8{}, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'velvet <file>'")
		return
	}

	content, err := readFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	virmac := vm.New()
	if err = virmac.Run(content); err != nil {
		fmt.Println(err.Error())
		return
	}
}

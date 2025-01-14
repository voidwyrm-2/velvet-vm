package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	_ "embed"

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

//go:embed version.txt
var version string

func main() {
	version = strings.TrimSpace(version)

	dumpStackAfterEachInstruction := flag.Bool("show", false, "Print the stack after each instruction")
	dumpAtEnd := flag.Bool("show-end", false, "Dump the stack at the end of the program")
	showVersion := flag.Bool("v", false, "Show the current Velvet version")

	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("expected 'velvet <file>'")
		return
	}

	content, err := readFile(args[0])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	virmac := vm.New()
	if err = virmac.Run(content, *dumpStackAfterEachInstruction); err != nil {
		fmt.Println(err.Error())
		return
	} else if *dumpAtEnd && !*dumpStackAfterEachInstruction {
		fmt.Println(virmac.DumpStack())
	}
}

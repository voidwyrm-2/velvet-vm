# Velvet

Velvet (or VelvetVM) is a bytecode virtual machine made to be used for real applications (at some point)

The list of instructions can be found in [instructions.md](<./instructions.md>)

The bytecode specification can be found in [bytecode.md](<./bytecode.md>)

Examples of the bytecode assembly can be found in [examples](<./examples>)

An example of what the bytecode itself looks like can be found in [bytecode-test.bc](<./examples/bytecode-test.bc>) and [bytecode-test.bin](<./examples/bytecode-test.bin>)

## Installation

You can either get an executable from the releases or compile it yourself

**Compiling it yourself**

Install Go if you don't have it, either from the [website](<https://go.dev>) or via a package manager(`brew install go`/`sudo port install go` on Mac, `choco install go` on Windows, you're on your own on Linux)

Next follow these instructions in your terminal:
1. Cd to a folder where you want the VM source folder to be
1. `git clone https://github.com/voidwyrm-2/velvet-vm`
3. `cd velvet-vm/velvet`
4. `go build -o velvet .`
5. `velvet ../examples/bytecode-test.cvelv`

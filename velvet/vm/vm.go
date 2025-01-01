package vm

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/voidwyrm-2/velvet-vm/velvet/vm/stack"
)

const InstructionSize int = 7

func getInstruction(bytes []uint8, pc int) (uint16, struct {
	flags [8]bool
	num   uint8
}, struct {
	one, two uint16
	both     uint32
},
) {
	instruction := bytes[pc : pc+7]
	flags := [8]bool{}

	for i := range 8 {
		if (instruction[2]>>i)&1 == 1 {
			flags[i] = true
		} else {
			flags[i] = false
		}
	}

	args := struct {
		one  uint16
		two  uint16
		both uint32
	}{
		one: uint16(instruction[3])<<8 + uint16(instruction[4]),
		two: uint16(instruction[5])<<8 + uint16(instruction[6]),
	}
	args.both = uint32(args.one)<<16 + uint32(args.two)

	return (uint16(instruction[0]) << 8) + uint16(instruction[1]), struct {
		flags [8]bool
		num   uint8
	}{
		flags: flags, num: instruction[2],
	}, args
}

func initDataGetter(bytes []uint8, dataAddr int) (func(addr, length uint16) (string, error), error) {
	getter := func(addr, length uint16) (string, error) {
		if dataAddr+int(addr)+int(length) > len(bytes) {
			return "", fmt.Errorf("data section address '%d' is not valid", dataAddr)
		}
		return string(bytes[dataAddr+int(addr) : dataAddr+int(addr)+int(length)]), nil
	}

	if int(dataAddr) >= len(bytes) {
		return getter, fmt.Errorf("data section address '%d' is not valid", dataAddr)
	}
	return getter, nil
}

type VelvetVM struct {
	stack     stack.Stack
	callables map[string]func(stack stack.Stack) bool
}

func New() VelvetVM {
	return VelvetVM{
		stack:     stack.New(),
		callables: stdfn,
	}
}

func (vm VelvetVM) VerifyBytecode(bytes []uint8) (int, int, bool) {
	if len(bytes) < 20 {
		return 0, 0, false
	} else if !strings.HasPrefix(string(bytes), "Velvet Scarlatina") {
		return 0, 0, false
	}

	return int(bytes[17]), (int(bytes[18]) << 8) + int(bytes[19]), true
}

func (vm VelvetVM) Run(bytes []uint8) error {
	var (
		vars     []stack.StackValue
		dataAddr int
		errFlag  bool
	)

	if _vars, _dataAddr, ok := vm.VerifyBytecode(bytes); !ok {
		return errors.New("malformed bytecode format")
	} else {
		vars, dataAddr = make([]stack.StackValue, _vars), _dataAddr
	}

	getData, err := initDataGetter(bytes, dataAddr)
	if err != nil {
		return err
	}

	_ = vars
	_ = dataAddr
	_ = errFlag

	callstack := []int{}

	pc := 20
	for {
		if pc+7 >= len(bytes) {
			return errors.New("end of bytes reached")
		}

		opcode, fb, args := getInstruction(bytes, pc)

		_ = fb

		switch opcode {
		case 0: // nop
			pc += InstructionSize
		case 1: // ret
			if len(callstack) > 0 {
				addr := callstack[len(callstack)-1]
				callstack = callstack[:len(callstack)-1]
				pc = addr
			} else {
				pc += InstructionSize
			}
		case 2: // halt
			os.Exit(int(args.one))
			pc += InstructionSize
		case 3: // call
			if fnName, err := getData(args.one, args.two); err != nil {
				return err
			} else if fn, ok := vm.callables[fnName]; !ok {
				return fmt.Errorf("function '%s' does not exist", fnName)
			} else {
				errFlag = fn(vm.stack)
			}
			pc += InstructionSize
		case 4: // push
			switch fb.num {
			case 1:
				vm.stack.Push(stack.NewBoolValue(args.one == 1))
			case 2:
				if str, err := getData(args.one, args.two); err != nil {
					return err
				} else {
					vm.stack.Push(stack.NewStringValue(str))
				}
			default:
				vm.stack.Push(stack.NewNumberValue(float32(args.both)))
			}
			pc += InstructionSize
		case 5: // pop
			vm.stack.Expect(stack.Any)
			vm.stack.Pop()
			pc += InstructionSize
		case 6: // dup
			vm.stack.Expect(stack.Any)
			item := vm.stack.Pop()
			vm.stack.Push(item)
			vm.stack.Push(item)
			pc += InstructionSize
		case 7: // swap
			vm.stack.Expect(stack.Any, stack.Any)
			x, y := vm.stack.Pop(), vm.stack.Pop()
			vm.stack.Push(y)
			vm.stack.Push(x)
			pc += InstructionSize
		case 8: // rot
			vm.stack.Expect(stack.Any, stack.Any, stack.Any)
			x, y, z := vm.stack.Pop(), vm.stack.Pop(), vm.stack.Pop()
			vm.stack.Push(z)
			vm.stack.Push(y)
			vm.stack.Push(x)
			pc += InstructionSize
		case 9: // set/get
			if int(args.one) >= len(vars) {
				return fmt.Errorf("%d is not a valid variable index", args.one)
			} else if fb.flags[7] {
				vm.stack.Push(vars[int(args.one)])
			} else {
				vm.stack.Expect(stack.Any)
				vars[int(args.one)] = vm.stack.Pop()
			}
			pc += InstructionSize
		case 10: // j/jt/jf/je/jne
			cond := true

			switch fb.num {
			case 1:
				vm.stack.Expect(stack.Bool)
				cond = vm.stack.Pop().GetBool()
			case 2:
				vm.stack.Expect(stack.Bool)
				cond = !vm.stack.Pop().GetBool()
			case 3:
				cond = errFlag
			case 4:
				cond = !errFlag
			}

			if cond {
				pc = int(args.both)
			} else {
				pc += InstructionSize
			}
		default:
			return fmt.Errorf("invalid opcode '%d'", opcode)
		}
	}

	panic("unreachable")
}

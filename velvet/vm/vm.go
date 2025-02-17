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

func initDataGetter(bytes []byte, dataAddr int) (func(addr uint16, length uint) ([]byte, error), error) {
	bytesGetter := func(addr uint16, length uint) ([]byte, error) {
		if dataAddr+int(addr)+int(length) > len(bytes) {
			return []byte{}, fmt.Errorf("data section address '%d' is not valid", dataAddr)
		}
		return bytes[dataAddr+int(addr) : dataAddr+int(addr)+int(length)], nil
	}

	/*
		stringGetter := func(addr uint16, length uint) (string, error) {
			b, e := bytesGetter(addr, length)
			return string(b), e
		}
	*/

	if int(dataAddr) >= len(bytes) {
		return bytesGetter, fmt.Errorf("data section address '%d' is not valid", dataAddr)
	}
	return bytesGetter, nil
}

func makeListFromBytes(lb []byte, getBytes func(addr uint16, length uint) ([]byte, error)) ([]stack.StackValue, error) {
	if len(lb) == 0 {
		return []stack.StackValue{}, nil
	}

	itemBytes := []struct {
		kind         uint8
		addr, length uint16
	}{}

	for i := 0; i < len(lb); i += 5 {
		itemBytes = append(itemBytes, struct {
			kind   uint8
			addr   uint16
			length uint16
		}{
			kind:   lb[i],
			addr:   (uint16(lb[i+1]) << 8) + uint16(lb[i+2]),
			length: (uint16(lb[i+3]) << 8) + uint16(lb[i+4]),
		})
	}

	items := []stack.StackValue{}

	for _, it := range itemBytes {
		switch it.kind {
		case stack.String: // string
			if str, err := getBytes(it.addr, uint(it.length)); err != nil {
				return []stack.StackValue{}, err
			} else {
				items = append(items, stack.NewStringValue(string(str)))
			}
		case stack.Bool: // bool
			if b, err := getBytes(it.addr, 1); err != nil {
				return []stack.StackValue{}, err
			} else {
				items = append(items, stack.NewBoolValue(b[0] == 1))
			}
		case stack.List: // list
			if sublsb, err := getBytes(it.addr, uint(it.length)*5); err != nil {
				return []stack.StackValue{}, err
			} else if subls, err := makeListFromBytes(sublsb, getBytes); err != nil {
				return []stack.StackValue{}, err
			} else {
				items = append(items, stack.NewListValue(subls...))
			}
		/*case stack.Function:
		if fnName, err := getBytes(args.one, uint(args.two)); err != nil {
			return []stack.StackValue, err
		} else if fn, ok := vm.callables[string(fnName)]; !ok {
			return []stack.StackValue,fmt.Errorf("function '%s' does not exist", string(fnName))
		} else {
			vm.stack.Push(stack.NewFuncValue(fn))
		}*/
		default: // number
			if b, err := getBytes(it.addr, 4); err != nil {
				return []stack.StackValue{}, err
			} else {
				items = append(items, stack.NewNumberValue(float32(int((uint(b[0])<<24)+(uint(b[1])<<16)+(uint(b[2])<<8)+uint(b[3])))))
			}
		}
	}

	return items, nil
}

type VelvetVM struct {
	stack     stack.Stack
	callables map[string]func(st *stack.Stack) error
}

func New() VelvetVM {
	return VelvetVM{
		stack:     stack.New(),
		callables: stdfn,
	}
}

func (vm VelvetVM) DumpStack() string {
	return vm.stack.Dump()
}

func (vm VelvetVM) VerifyBytecode(bytes []byte) (struct {
	isLibrary, f2, f3, f4, f5, f6, f7, f8 bool
}, int, int, int, bool,
) {
	if len(bytes) < 32 {
		return struct {
			isLibrary bool
			f2        bool
			f3        bool
			f4        bool
			f5        bool
			f6        bool
			f7        bool
			f8        bool
		}{}, 0, 0, 0, false
	} else if !strings.HasPrefix(string(bytes), "Velvet Scarlatina") {
		return struct {
			isLibrary bool
			f2        bool
			f3        bool
			f4        bool
			f5        bool
			f6        bool
			f7        bool
			f8        bool
		}{}, 0, 0, 0, false
	}

	return struct {
		isLibrary bool
		f2        bool
		f3        bool
		f4        bool
		f5        bool
		f6        bool
		f7        bool
		f8        bool
	}{isLibrary: (bytes[17] >> 7) == 1}, (int(bytes[18]) << 8) + int(bytes[19]), (int(bytes[20]) << 24) + (int(bytes[21]) << 16) + (int(bytes[22]) << 8) + int(bytes[23]), (int(bytes[24]) << 24) + (int(bytes[25]) << 16) + (int(bytes[26]) << 8) + int(bytes[27]), true
}

func (vm VelvetVM) Run(bytes []byte, dumpStackAfterEachInstruction, dumpVarsAfterEachInstruction bool) error {
	var (
		flags struct {
			isLibrary, f2, f3, f4, f5, f6, f7, f8 bool
		}
		vars                  []stack.StackValue
		dataAddr, entryOffset int
		errFlag               bool
		errReg                string
	)

	setErr := func(e error) {
		if e != nil {
			errFlag = true
			errReg = e.Error()
		}
	}

	if _flags, _vars, _dataAddr, _entryOffset, ok := vm.VerifyBytecode(bytes); !ok {
		return errors.New("malformed bytecode format")
	} else {
		flags, vars, dataAddr, entryOffset = _flags, make([]stack.StackValue, _vars), _dataAddr, _entryOffset
	}

	if entryOffset > dataAddr {
		return fmt.Errorf("the start of execution is inside the data section")
	}

	getBytes, err := initDataGetter(bytes, dataAddr)
	if err != nil {
		return err
	}

	if flags.isLibrary {
		return errors.New("this Velvet bytecode executable has been declared as a library meaning it cannot be directly run, it must be imported by a non-library Velvet executable")
	}

	callstack := []int{}

	pc := 32 + (entryOffset * 7)
	for {
		if pc+7 >= len(bytes) {
			return errors.New("end of bytes reached")
		}

		opcode, fb, args := getInstruction(bytes, pc)

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
			if fb.flags[0] {
				vm.stack.Expect(stack.Function)
				setErr(vm.stack.Pop().GetFunc()(&vm.stack))
			} else {
				if fnName, err := getBytes(args.one, uint(args.two)); err != nil {
					return err
				} else if string(fnName) == "getErr" {
					vm.stack.Push(stack.NewStringValue(errReg))
				} else if fn, ok := vm.callables[string(fnName)]; !ok {
					return fmt.Errorf("function '%s' does not exist", string(fnName))
				} else {
					setErr(fn(&vm.stack))
				}
			}
			pc += InstructionSize
		case 4: // push
			switch fb.num {
			case 1: // bool
				vm.stack.Push(stack.NewBoolValue(args.one != 0))
			case 2: // string
				if str, err := getBytes(args.one, uint(args.two)); err != nil {
					return err
				} else {
					vm.stack.Push(stack.NewStringValue(string(str)))
				}
			case 3: // list
				if lb, err := getBytes(args.one, uint(args.two)*5); err != nil {
					return err
				} else if ls, err := makeListFromBytes(lb, getBytes); err != nil {
					return err
				} else {
					vm.stack.Push(stack.NewListValue(ls...))
				}
			case 4: // function
				if fnName, err := getBytes(args.one, uint(args.two)); err != nil {
					return err
				} else if fn, ok := vm.callables[string(fnName)]; !ok {
					return fmt.Errorf("function '%s' does not exist", string(fnName))
				} else {
					vm.stack.Push(stack.NewFuncValue(fn))
				}
			case 5: // error register
				vm.stack.Push(stack.NewStringValue(errReg))
			default:
				vm.stack.Push(stack.NewNumberValue(float32(int(args.both))))
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
			} else if fb.num&1 == 1 {
				vm.stack.Push(vars[int(args.one)])
			} else {
				vm.stack.Expect(stack.Any)
				vars[int(args.one)] = vm.stack.Pop()
			}
			pc += InstructionSize
		case 10: // j/jt/jf/je/jne or br/brt/brf/bre/brne
			cond := true

			jumpType, isBranch := exactIsBranch(fb.num)

			switch jumpType {
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
				if isBranch {
					callstack = append(callstack, pc)
				}
				pc = int(args.both)
			} else {
				pc += InstructionSize
			}
		default:
			return fmt.Errorf("invalid opcode '%d'", opcode)
		}

		if dumpStackAfterEachInstruction {
			fmt.Println(vm.stack.Dump())
		}

		if dumpVarsAfterEachInstruction {
			if dumpStackAfterEachInstruction {
				fmt.Println("")
			}

			fmtVars := []string{}
			for _, v := range vars {
				fmtVars = append(fmtVars, v.Dump())
			}
		}
	}

	panic("unreachable")
}

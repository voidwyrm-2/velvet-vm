package emitter

import (
	"fmt"
	"os"
	"reflect"
)

type Opcode uint16

const (
	Nop Opcode = iota
	Ret
	Halt
	Call
	Push
	Pop
	Dup
	Swap
	Rot
	Set
	Jump
)

func (o Opcode) Name() string {
	return []string{
		"Nop",
		"Ret",
		"Halt",
		"Call",
		"Push",
		"Pop",
		"Dup",
		"Swap",
		"Rot",
		"Set",
		"Jump",
	}[o]
}

func writeFile(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

type VelvetAsm struct {
	vars  uint16
	flags struct {
		isLibrary, f2, f3, f4, f5, f6, f7, f8 bool
	}
	programEntry uint32
	instructions [][7]byte
	labels       map[string]uint32
	staticCache  map[any][2]uint16
	data         []byte
}

func New(vars uint16) VelvetAsm {
	return VelvetAsm{vars: vars, instructions: [][7]byte{}, labels: map[string]uint32{}, staticCache: map[any][2]uint16{}, data: []byte{}}
}

func (va VelvetAsm) flagsToUint8() uint8 {
	return (boolToUint8(va.flags.isLibrary) << 7) + (boolToUint8(va.flags.f2) << 6) + (boolToUint8(va.flags.f3) << 5) + (boolToUint8(va.flags.f4) << 4) + (boolToUint8(va.flags.f5) << 3) + (boolToUint8(va.flags.f6) << 2) + (boolToUint8(va.flags.f7) << 1) + boolToUint8(va.flags.f8)
}

func (va VelvetAsm) Write(filename string) error {
	output := []byte("Velvet Scarlatina")

	output = append(output, va.flagsToUint8())

	output = append(output, uint8(va.vars>>8), uint8(va.vars))

	dataAddr := uint32(32 + len(va.instructions)*7)
	output = append(output, uint8(dataAddr>>24), uint8(dataAddr>>16), uint8(dataAddr>>8), uint8(dataAddr))

	output = append(output, uint8(va.programEntry>>24), uint8(va.programEntry>>16), uint8(va.programEntry>>8), uint8(va.programEntry))

	output = append(output, 0, 0, 0, 0)

	for _, ins := range va.instructions {
		output = append(output, ins[0:]...)
	}
	output = append(output, va.data...)

	return writeFile(filename, output)
}

func (va VelvetAsm) Data() []byte {
	return va.data
}

func (va VelvetAsm) DataString() string {
	return string(va.data)
}

func (va *VelvetAsm) AddNumber(value uint32) (uint16, uint16) {
	if pos, ok := va.staticCache[value]; ok {
		return pos[0], pos[1]
	}
	addr := uint16(len(va.data))
	va.data = append(va.data, byte(value>>24), byte((value>>16)&0xff), byte((value>>8)&0xff), byte(value))
	va.staticCache[value] = [2]uint16{addr, 4}
	return addr, 4
}

func (va *VelvetAsm) AddBool(value bool) (uint16, uint16) {
	if pos, ok := va.staticCache[value]; ok {
		return pos[0], pos[1]
	}

	addr := uint16(len(va.data))
	if value {
		va.data = append(va.data, 1)
	} else {
		va.data = append(va.data, 0)
	}

	va.staticCache[value] = [2]uint16{addr, 1}
	return addr, 1
}

func (va *VelvetAsm) AddString(value string) (uint16, uint16) {
	// fmt.Println(value, va.staticCache)
	// fmt.Println(va.DataString())
	if pos, ok := va.staticCache[value]; ok {
		// fmt.Println("cache hit with: " + value)
		return pos[0], pos[1]
	}

	//	fmt.Println("cache miss with: " + value)

	addr := uint16(len(va.data))
	va.data = append(va.data, []byte(value)...)

	va.staticCache[value] = [2]uint16{addr, uint16(len(value))}
	return addr, uint16(len(value))
}

func (va *VelvetAsm) AddList(values ...any) (uint16, uint16) {
	// fmt.Println(fmt.Sprintf("%v", values))
	if pos, ok := va.staticCache[fmt.Sprintf("%v", values)]; ok {
		return pos[0], pos[1]
	}

	spl16 := func(n uint16) []byte {
		return []byte{byte(n >> 8), byte(n)}
	}

	addedBytes := []byte{}

	for _, val := range values {
		switch v := val.(type) {
		case int:
			{
				valAddr, valLen := va.AddNumber(uint32(v))
				addedBytes = append(addedBytes, 0b0)
				addedBytes = append(addedBytes, spl16(valAddr)...)
				addedBytes = append(addedBytes, spl16(valLen)...)
			}
		case string:
			{
				valAddr, valLen := va.AddString(v)
				addedBytes = append(addedBytes, 0b1)
				addedBytes = append(addedBytes, spl16(valAddr)...)
				addedBytes = append(addedBytes, spl16(valLen)...)
			}
		case bool:
			{
				valAddr, valLen := va.AddBool(v)
				addedBytes = append(addedBytes, 0b10)
				addedBytes = append(addedBytes, spl16(valAddr)...)
				addedBytes = append(addedBytes, spl16(valLen)...)
			}
		case []any:
			{
				valAddr, valLen := va.AddList(v)
				addedBytes = append(addedBytes, 0b100)
				addedBytes = append(addedBytes, spl16(valAddr)...)
				addedBytes = append(addedBytes, spl16(valLen)...)
			}
		default:
			panic(fmt.Sprintf("'%s' is not a valid type", reflect.TypeOf(val).Name()))
		}
	}

	addr := uint16(len(va.data))

	va.data = append(va.data, addedBytes...)

	va.staticCache[fmt.Sprintf("%v", values)] = [2]uint16{addr, uint16(len(values))}

	return addr, uint16(len(values))
}

func (va *VelvetAsm) CreateLabel(name string) {
	if _, ok := va.labels[name]; ok {
		panic(fmt.Sprintf("label '%s' already exists", name))
	}
	va.labels[name] = uint32(32 + len(va.instructions)*7)
}

func (va VelvetAsm) HasLabel(name string) bool {
	_, ok := va.labels[name]
	return ok
}

func (va VelvetAsm) GetLabel(name string) uint32 {
	if addr, ok := va.labels[name]; !ok {
		panic(fmt.Sprintf("label '%s' does not exist", name))
	} else {
		return addr
	}
}

func (va *VelvetAsm) Emit(op Opcode, flag uint8, one, two uint16) {
	va.instructions = append(va.instructions, [7]byte{byte(op >> 8), byte(op), flag, uint8(one >> 8), uint8(one), uint8(two >> 8), uint8(two)})
}

func (va *VelvetAsm) Emit32(op Opcode, flag uint8, both uint32) {
	va.Emit(op, flag, uint16(both>>16), uint16(both))
}

func (va *VelvetAsm) EmitString(op Opcode, flag uint8, str string) {
	addr, length := va.AddString(str)
	va.Emit(op, flag, addr, length)
}

func (va *VelvetAsm) EmitList(op Opcode, flag uint8, values ...any) {
	addr, length := va.AddList(values...)
	va.Emit(op, flag, addr, length)
}

func (va *VelvetAsm) EmitNF(op Opcode, one, two uint16) {
	va.Emit(op, 0, one, two)
}

func (va *VelvetAsm) EmitNF32(op Opcode, both uint32) {
	va.Emit32(op, 0, both)
}

func (va *VelvetAsm) EmitNA(op Opcode, flag uint8) {
	va.Emit(op, flag, 0, 0)
}

func (va *VelvetAsm) EmitBasic(op Opcode) {
	va.Emit(op, 0, 0, 0)
}

func (va *VelvetAsm) Halt(code int8) {
	va.Emit(Halt, 0, uint16(code), 0)
}

func (va *VelvetAsm) DoDirective(name string, args ...any) {
	switch name {
	case "vars":
		va.vars = uint16(args[0].(int))
	case "entry":
		va.SetEntry(uint32(args[0].(int)))
	default:
		panic("unknown directive '" + name + "'")
	}
}

func (va *VelvetAsm) SetEntry(entryOffset uint32) {
	va.programEntry = entryOffset
}

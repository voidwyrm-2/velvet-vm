package emitter

import (
	"os"
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

type VelvetAsm struct {
	vars         uint8
	instructions [][7]byte
	data         []byte
}

func New(vars uint8) VelvetAsm {
	return VelvetAsm{vars: vars, instructions: [][7]byte{}, data: []byte{}}
}

func (va VelvetAsm) Write(filename string) error {
	output := []byte("Velvet Scarlatina")
	output = append(output, va.vars)

	dataAddr := uint16(20 + len(va.instructions)*7)
	output = append(output, uint8(dataAddr>>8), uint8(dataAddr))

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

func (va *VelvetAsm) Emit(op Opcode, flag uint8, one, two uint16) {
	va.instructions = append(va.instructions, [7]byte{byte(op >> 8), byte(op), flag, uint8(one >> 8), uint8(one), uint8(two >> 8), uint8(two)})
}

func (va *VelvetAsm) Emit32(op Opcode, flag uint8, both uint32) {
	va.Emit(op, flag, uint16(both>>16), uint16(both))
}

func (va *VelvetAsm) EmitString(op Opcode, flag uint8, str string) {
	addr, length := uint16(len(va.data)), uint16(len(str))
	va.data = append(va.data, []byte(str)...)
	va.Emit(op, flag, addr, length)
}

func (va *VelvetAsm) Halt(code int8) {
	va.Emit(Halt, 0, uint16(code), 0)
}

package z80

type Z80 struct {
	A  byte
	F  byte
	B  byte
	C  byte
	D  byte
	E  byte
	H  byte
	L  byte
	SP uint16
	PC uint16

	IX uint16
	IY uint16
	I  byte
	R  byte
	Ap byte // A'
	Fp byte // F'
	Bp byte // B'
	Cp byte // C'
	Dp byte // D'
	Ep byte // E'
	Hp byte // H'
	Lp byte // L'

	memory MemoryAccess
}

func New(memory MemoryAccess) *Z80 {
	cpu := &Z80{
		memory: memory,
		PC:     0x100,
	}
	return cpu
}

func (cpu *Z80) Step() (int, error) {
	opcode := cpu.memory.ReadByte(cpu.PC, false)
	cpu.PC += 1
	println(opcode)
	// switch opcode {
	// case
	// }
	return 0, nil
}

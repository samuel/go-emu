package cpu6502

import (
	"fmt"
	"os"
)

const (
	IV_NMI   = 0xFFFA
	IV_RESET = 0xFFFC
	IV_IRQ   = 0xFFFE

	FLAG_C = 0x01
	FLAG_Z = 0x02
	FLAG_I = 0x04
	FLAG_D = 0x08
	FLAG_B = 0x10
	FLAG_V = 0x40
	FLAG_S = 0x80
)

type CPU6502 struct {
	A, X, Y byte

	// P byte // processor status
	// 0 - C - carry flag
	// 1 - Z - zero flag
	// 2 - I - interrupt enable/disable (0=enabled)
	// 3 - D - decimal mode status flag (1=BCD)
	// 4 - B - set when a software interrupt (BRK) is executed
	// 5 - not used (1 all the time)
	// 6 - V - Overflow
	// 7 - S/N - Sign/Negative
	CarryFlag              bool
	ZeroFlag               bool
	InterruptsDisabledFlag bool
	DecimalFlag            bool // true=BCD
	SoftwareInterruptFlag  bool
	OverflowFlag           bool
	SignFlag               bool

	PC uint16 // program counter (PCL=low, PCH=high)
	SP uint8  // stack pointer

	Cycles uint64

	NMICounter int // steps till NMI

	memory MemoryAccess
}

func NewCPU6502(memory MemoryAccess) *CPU6502 {
	cpu := &CPU6502{
		memory: memory,
		SP:     0xFD,
		// P: 0x34, 24?
		InterruptsDisabledFlag: true}
	// SoftwareInterruptFlag: true}
	return cpu
}

func (cpu *CPU6502) GetP() byte {
	var flags byte = 1 << 5
	if cpu.CarryFlag {
		flags |= FLAG_C
	}
	if cpu.ZeroFlag {
		flags |= FLAG_Z
	}
	if cpu.InterruptsDisabledFlag {
		flags |= FLAG_I
	}
	if cpu.DecimalFlag {
		flags |= FLAG_D
	}
	if cpu.SoftwareInterruptFlag {
		flags |= FLAG_B
	}
	if cpu.OverflowFlag {
		flags |= FLAG_V
	}
	if cpu.SignFlag {
		flags |= FLAG_S
	}
	return flags
}

func (cpu *CPU6502) SetP(flags byte) {
	cpu.CarryFlag = flags&FLAG_C != 0
	cpu.ZeroFlag = flags&FLAG_Z != 0
	cpu.InterruptsDisabledFlag = flags&FLAG_I != 0
	cpu.DecimalFlag = flags&FLAG_D != 0
	cpu.SoftwareInterruptFlag = flags&FLAG_B != 0
	cpu.OverflowFlag = flags&FLAG_V != 0
	cpu.SignFlag = flags&FLAG_S != 0
}

func (cpu *CPU6502) FlagString() string {
	flags := ""
	not_set := "_"
	if cpu.CarryFlag {
		flags += "C"
	} else {
		flags += not_set
	}
	if cpu.ZeroFlag {
		flags += "Z"
	} else {
		flags += not_set
	}
	if cpu.InterruptsDisabledFlag {
		flags += "I"
	} else {
		flags += not_set
	}
	if cpu.DecimalFlag {
		flags += "D"
	} else {
		flags += not_set
	}
	if cpu.SoftwareInterruptFlag {
		flags += "B"
	} else {
		flags += not_set
	}
	if cpu.OverflowFlag {
		flags += "V"
	} else {
		flags += not_set
	}
	if cpu.SignFlag {
		flags += "S"
	} else {
		flags += not_set
	}
	return flags
}

func (cpu *CPU6502) String() string {
	return fmt.Sprintf("{PC:%04x SP:%02x A:%02x X:%02x Y:%02x P:%02x:%s}",
		cpu.PC, cpu.SP, cpu.A, cpu.X, cpu.Y, cpu.GetP(), cpu.FlagString())
}

func (cpu *CPU6502) ReadByte(address uint16, peek bool) byte {
	return cpu.memory.ReadByte(address, peek)
}

// Read a 16-bit unsigned int dealing with page wrap
func (cpu *CPU6502) ReadUI16(address uint16, peek bool) uint16 {
	// There's a bug in the 6502 where Indirect addressing doesn't advance pages
	// 02ff -> bytes 02ff & 0200 rather than 02ff 0300
	return uint16(cpu.memory.ReadByte(address, peek)) |
		(uint16(cpu.memory.ReadByte((address+1)&0xff+(address&0xff00), peek)) << 8)
}

func (cpu *CPU6502) ReadOpcode() (OpcodeSpec, uint16) {
	return ReadOpcode(cpu.memory, cpu.PC)
}

func (cpu *CPU6502) PushByte(value byte) {
	cpu.memory.WriteByte(0x100+uint16(cpu.SP), value)
	cpu.SP--
}

func (cpu *CPU6502) PopByte() byte {
	cpu.SP++
	return cpu.memory.ReadByte(0x100+uint16(cpu.SP), false)
}

func (cpu *CPU6502) PushAddress(addr uint16) {
	cpu.PushByte(byte(addr >> 8))
	cpu.PushByte(byte(addr & 0xff))
}

func (cpu *CPU6502) PopAddress() uint16 {
	return uint16(cpu.PopByte()) | (uint16(cpu.PopByte()) << 8)
}

func (cpu *CPU6502) Step() (int, error) {
	if cpu.NMICounter > 0 {
		cpu.NMICounter--
		if cpu.NMICounter == 0 {
			cpu.PushAddress(cpu.PC)
			cpu.SoftwareInterruptFlag = false
			cpu.PushByte(cpu.GetP())
			cpu.PC = cpu.ReadUI16(IV_NMI, false)
			cpu.Cycles += uint64(5)
			return 5, nil
		}
	}

	opcode, opval := cpu.ReadOpcode()
	cpu.PC += uint16(opcode.Size)

	var addr uint16
	var value byte
	var cycles int = opcode.Cycles
	if cycles < 0 {
		cycles = -cycles
	}
	cycles2 := opcode.Size
	if cycles2 < 2 {
		cycles2 = 2
	}
	switch opcode.AddressingMode {
	default:
		panic(fmt.Sprintf("Unhandled addresing mode %d", opcode.AddressingMode))
	case AMImplied:
		// do nothing
	case AMAccumulator:
		// do nothing
		value = cpu.A
	case AMImmediate:
		// do nothing
		value = byte(opval)
	case AMAbsolute, AMZeroPage:
		addr = opval
		if opcode.Instruction.Read {
			cycles2++
			value = cpu.memory.ReadByte(addr, false)
			if opcode.Instruction.Write {
				// For read-modify-write instructions the value gets written twice
				cpu.memory.WriteByte(addr, value)
				cycles2 += 2
			}
		} else if opcode.Instruction.Write {
			cycles2++
		}
	case AMZeroPageX, AMZeroPageY:
		cpu.memory.ReadByte(opval, false) // dummy read
		if opcode.AddressingMode == AMZeroPageX {
			addr = (opval + uint16(cpu.X)) & 0x00ff
		} else {
			addr = (opval + uint16(cpu.Y)) & 0x00ff
		}
		cycles2 += 2
		if opcode.Instruction.Read {
			value = cpu.memory.ReadByte(addr, false)
			if opcode.Instruction.Write {
				// For read-modify-write instructions the value gets written twice
				cpu.memory.WriteByte(addr, value)
				cycles2 += 2
			}
		}
	case AMAbsoluteX, AMAbsoluteY, AMIndirectY:
		addrt := opval
		if opcode.AddressingMode == AMAbsoluteX {
			addr = addrt + uint16(cpu.X)
		} else if opcode.AddressingMode == AMAbsoluteY {
			addr = addrt + uint16(cpu.Y)
		} else {
			addrt = cpu.ReadUI16(opval, false)
			addr = addrt + uint16(cpu.Y)
			cycles2 += 2
		}
		value = cpu.memory.ReadByte(addrt&0xff00|addr&0x00ff, false)
		cycles2++
		if opcode.Instruction.Read {
			if opcode.Instruction.Write {
				// For read-modify-write instructions the value gets written twice
				value = cpu.memory.ReadByte(addr, false)
				cpu.memory.WriteByte(addr, value)
				cycles2 += 3
			} else if addr&0xff00 != addrt&0xff00 {
				value = cpu.memory.ReadByte(addr, false)
				cycles++
				cycles2++
			}
		} else if opcode.Instruction.Write {
			cycles2++
		}
	case AMRelative:
		addr = cpu.PC + uint16(int8(opval))
	case AMIndirectX:
		cpu.memory.ReadByte(opval, false) // dummy read
		addr = uint16(byte(opval) + cpu.X)
		addr = cpu.ReadUI16(addr, false)
		cycles2 += 4
		if opcode.Instruction.Read {
			value = cpu.memory.ReadByte(addr, false)
			if opcode.Instruction.Write {
				// For read-modify-write instructions the value gets written twice
				cpu.memory.WriteByte(addr, value)
				cycles2 += 2
			}
		}
	case AMIndirect:
		addr = cpu.ReadUI16(opval, false)
		cycles2 += 2
	}

	jump := false // jump to 'addr' and account for clock

	switch opcode.Instruction.Num {
	default:
		panic("Unhandled opcode " + opcode.Instruction.Name)
	case I_AAX.Num: // undocumented
		value = cpu.X & cpu.A
		cpu.memory.WriteByte(addr, value)
	case I_ADC.Num:
		res := uint16(value) + uint16(cpu.A)
		if cpu.CarryFlag {
			res++
		}
		cpu.ZeroFlag = (res & 0xff) == 0
		if cpu.DecimalFlag {
			if (cpu.A^value^byte(res))&0x10 == 0x10 {
				res += 6
			}
			if res&0xf0 > 0x90 {
				res += 0x60
			}
		}
		cpu.SignFlag = res&0x80 != 0
		cpu.OverflowFlag = !((cpu.A^value)&0x80 != 0) && ((uint16(cpu.A)^res)&0x80 != 0)
		cpu.CarryFlag = res&0x100 != 0
		cpu.A = byte(res & 0xff)
	case I_AND.Num:
		cpu.A &= value
		cpu.SignFlag = cpu.A&0x80 > 0
		cpu.ZeroFlag = cpu.A == 0
	case I_ASL.Num:
		cpu.CarryFlag = value&0x80 != 0
		value = value << 1
		cpu.SignFlag = value&0x80 != 0
		cpu.ZeroFlag = value == 0
		if opcode.AddressingMode == AMAccumulator {
			cpu.A = value
		} else {
			cpu.memory.WriteByte(addr, value)
		}
	case I_BCC.Num:
		if !cpu.CarryFlag {
			jump = true
		}
	case I_BCS.Num:
		if cpu.CarryFlag {
			jump = true
		}
	case I_BEQ.Num:
		if cpu.ZeroFlag {
			jump = true
		}
	case I_BIT.Num:
		cpu.SignFlag = value&0x80 != 0
		cpu.OverflowFlag = value&0x40 != 0
		cpu.ZeroFlag = value&cpu.A == 0
	case I_BMI.Num:
		if cpu.SignFlag {
			jump = true
		}
	case I_BPL.Num:
		if !cpu.SignFlag {
			jump = true
		}
	case I_BNE.Num:
		if !cpu.ZeroFlag {
			jump = true
		}
	case I_BRK.Num:
		cpu.PushAddress(cpu.PC + 1)
		cpu.SoftwareInterruptFlag = true
		cpu.PushByte(cpu.GetP())
		cpu.InterruptsDisabledFlag = true
		cpu.PC = cpu.ReadUI16(IV_IRQ, false)
		cycles2 += 5
	case I_BVC.Num:
		if !cpu.OverflowFlag {
			jump = true
		}
	case I_BVS.Num:
		if cpu.OverflowFlag {
			jump = true
		}
	case I_CLC.Num:
		cpu.CarryFlag = false
	case I_CLD.Num:
		cpu.DecimalFlag = false
	case I_CLI.Num:
		cpu.InterruptsDisabledFlag = false
	case I_CLV.Num:
		cpu.OverflowFlag = false
	case I_CMP.Num:
		res := cpu.A - value
		cpu.CarryFlag = cpu.A >= value
		cpu.SignFlag = res&0x80 != 0
		cpu.ZeroFlag = res == 0
	case I_CPX.Num:
		res := cpu.X - value
		cpu.CarryFlag = cpu.X >= value
		cpu.SignFlag = res&0x80 != 0
		cpu.ZeroFlag = res == 0
	case I_CPY.Num:
		res := cpu.Y - value
		cpu.CarryFlag = cpu.Y >= value
		cpu.SignFlag = res&0x80 != 0
		cpu.ZeroFlag = res == 0
	case I_DCP.Num: // undocumented - equivalent to DEC, CMP
		value--
		cpu.memory.WriteByte(addr, value)
		res := cpu.A - value
		cpu.CarryFlag = cpu.A >= value
		cpu.SignFlag = res&0x80 != 0
		cpu.ZeroFlag = res == 0
	case I_DEC.Num:
		value--
		cpu.SignFlag = value&0x80 != 0
		cpu.ZeroFlag = value == 0
		cpu.memory.WriteByte(addr, value)
	case I_DEX.Num:
		cpu.X -= 1
		cpu.SignFlag = cpu.X&0x80 != 0
		cpu.ZeroFlag = cpu.X == 0
	case I_DEY.Num:
		cpu.Y -= 1
		cpu.SignFlag = cpu.Y&0x80 != 0
		cpu.ZeroFlag = cpu.Y == 0
	case I_EOR.Num:
		cpu.A ^= value
		cpu.SignFlag = cpu.A&0x80 != 0
		cpu.ZeroFlag = cpu.A == 0
	case I_INC.Num:
		value++
		cpu.SignFlag = value&0x80 != 0
		cpu.ZeroFlag = value == 0
		cpu.memory.WriteByte(addr, value)
	case I_INX.Num:
		cpu.X += 1
		cpu.SignFlag = cpu.X&0x80 != 0
		cpu.ZeroFlag = cpu.X == 0
	case I_INY.Num:
		cpu.Y += 1
		cpu.SignFlag = cpu.Y&0x80 != 0
		cpu.ZeroFlag = cpu.Y == 0
	case I_ISC.Num: // undocumented - equivalent to INC, SBC
		value++
		cpu.memory.WriteByte(addr, value)
		temp := uint16(cpu.A) - uint16(value)
		if !cpu.CarryFlag {
			temp--
		}
		cpu.SignFlag = temp&0x80 != 0
		cpu.ZeroFlag = temp&0xff == 0
		cpu.OverflowFlag = ((cpu.A^byte(temp))&0x80 != 0) && ((cpu.A^value)&0x80 != 0)
		if cpu.DecimalFlag {
			var carry byte = 1
			if cpu.CarryFlag {
				carry = 0
			}
			if ((cpu.A & 0x0f) - carry) < (value & 0x0f) {
				/* EP */
				temp -= 6
			}
			if temp > 0x99 {
				temp -= 0x60
			}
		}
		cpu.CarryFlag = temp < 0x100
		cpu.A = byte(temp & 0xff)
	case I_JMP.Num:
		cpu.PC = addr
	case I_JSR.Num:
		cpu.PushAddress(cpu.PC - 1)
		cpu.PC = addr
		cycles2 += 3
	case I_LAX.Num: // undocumented
		cpu.A = value
		cpu.X = value
		cpu.SignFlag = cpu.A&0x80 != 0
		cpu.ZeroFlag = cpu.A == 0
	case I_LDA.Num:
		cpu.A = value
		cpu.SignFlag = cpu.A&0x80 != 0
		cpu.ZeroFlag = cpu.A == 0
	case I_LDX.Num:
		cpu.X = value
		cpu.SignFlag = cpu.X&0x80 != 0
		cpu.ZeroFlag = cpu.X == 0
	case I_LDY.Num:
		cpu.Y = value
		cpu.SignFlag = cpu.Y&0x80 != 0
		cpu.ZeroFlag = cpu.Y == 0
	case I_LSR.Num:
		cpu.CarryFlag = value&0x01 > 0
		value >>= 1
		cpu.SignFlag = value&0x80 != 0
		cpu.ZeroFlag = value == 0
		if opcode.AddressingMode == AMAccumulator {
			cpu.A = value
		} else {
			cpu.memory.WriteByte(addr, value)
		}
	case I_NOP.Num, I_DOP.Num, I_TOP.Num, I_NP2.Num:
		// no-op
	case I_ORA.Num:
		cpu.A |= value
		cpu.SignFlag = cpu.A&0x80 != 0
		cpu.ZeroFlag = cpu.A == 0
	case I_PHA.Num:
		cpu.PushByte(cpu.A)
		cycles2++
	case I_PHP.Num:
		cpu.PushByte(cpu.GetP() | FLAG_B) // B flag always pushed as 1
		cycles2++
	case I_PLA.Num:
		cpu.A = cpu.PopByte()
		cpu.SignFlag = cpu.A&0x80 != 0
		cpu.ZeroFlag = cpu.A == 0
		cycles2 += 2
	case I_PLP.Num:
		cpu.SetP(cpu.PopByte() & ^byte(FLAG_B)) // B flag discarded
		cycles2 += 2
	case I_RLA.Num: // undocumented - equivalent to ROL, AND
		var carry byte = 0
		if cpu.CarryFlag {
			carry = 1
		}
		cpu.CarryFlag = value&0x80 > 0
		value = (value << 1) | carry
		if opcode.AddressingMode == AMAccumulator {
			cpu.A = value
		} else {
			cpu.memory.WriteByte(addr, value)
		}
		cpu.A &= value
		cpu.SignFlag = cpu.A&0x80 > 0
		cpu.ZeroFlag = cpu.A == 0
	case I_RTI.Num:
		cpu.SetP(cpu.PopByte())
		cpu.PC = cpu.PopAddress()
		cycles2 += 4
	case I_ROL.Num:
		var carry byte = 0
		if cpu.CarryFlag {
			carry = 1
		}
		cpu.CarryFlag = value&0x80 > 0
		value = (value << 1) | carry
		cpu.SignFlag = value&0x80 > 0
		cpu.ZeroFlag = value == 0
		if opcode.AddressingMode == AMAccumulator {
			cpu.A = value
		} else {
			cpu.memory.WriteByte(addr, value)
		}
	case I_ROR.Num:
		var carry byte = 0
		if cpu.CarryFlag {
			carry = 0x80
		}
		cpu.CarryFlag = value&0x01 > 0
		value = (value >> 1) | carry
		cpu.SignFlag = value&0x80 > 0
		cpu.ZeroFlag = value == 0
		if opcode.AddressingMode == AMAccumulator {
			cpu.A = value
		} else {
			cpu.memory.WriteByte(addr, value)
		}
	case I_RRA.Num: // undocumented - equivalent to ROR, ADC
		var carry byte = 0
		if cpu.CarryFlag {
			carry = 0x80
		}
		cpu.CarryFlag = value&0x01 > 0
		value = (value >> 1) | carry
		if opcode.AddressingMode == AMAccumulator {
			cpu.A = value
		} else {
			cpu.memory.WriteByte(addr, value)
		}
		res := uint16(value) + uint16(cpu.A)
		if cpu.CarryFlag {
			res++
		}
		cpu.ZeroFlag = (res & 0xff) == 0
		if cpu.DecimalFlag {
			if (cpu.A^value^byte(res))&0x10 == 0x10 {
				res += 6
			}
			if res&0xf0 > 0x90 {
				res += 0x60
			}
		}
		cpu.SignFlag = res&0x80 != 0
		cpu.OverflowFlag = !((cpu.A^value)&0x80 != 0) && ((uint16(cpu.A)^res)&0x80 != 0)
		cpu.CarryFlag = res&0x100 != 0
		cpu.A = byte(res & 0xff)
	case I_RTS.Num:
		cpu.PC = cpu.PopAddress() + 1
		cycles2 += 4
	case I_SBC.Num:
		temp := uint16(cpu.A) - uint16(value)
		if !cpu.CarryFlag {
			temp--
		}
		cpu.SignFlag = temp&0x80 != 0
		cpu.ZeroFlag = temp&0xff == 0
		cpu.OverflowFlag = ((cpu.A^byte(temp))&0x80 != 0) && ((cpu.A^value)&0x80 != 0)
		if cpu.DecimalFlag {
			var carry byte = 1
			if cpu.CarryFlag {
				carry = 0
			}
			if ((cpu.A & 0x0f) - carry) < (value & 0x0f) {
				/* EP */
				temp -= 6
			}
			if temp > 0x99 {
				temp -= 0x60
			}
		}
		cpu.CarryFlag = temp < 0x100
		cpu.A = byte(temp & 0xff)
	case I_SEC.Num:
		cpu.CarryFlag = true
	case I_SED.Num:
		cpu.DecimalFlag = true
	case I_SEI.Num:
		cpu.InterruptsDisabledFlag = true
	case I_SLO.Num: // undocumented - equivalent to ASL, ORA
		cpu.CarryFlag = value&0x80 != 0
		value = value << 1
		if opcode.AddressingMode == AMAccumulator {
			cpu.A = value
		} else {
			cpu.memory.WriteByte(addr, value)
		}

		cpu.A |= value
		cpu.SignFlag = cpu.A&0x80 != 0
		cpu.ZeroFlag = cpu.A == 0
	case I_SRE.Num: // undocumented - equivalent to LSR, EOR
		cpu.CarryFlag = value&0x01 > 0
		value >>= 1
		if opcode.AddressingMode == AMAccumulator {
			cpu.A = value
		} else {
			cpu.memory.WriteByte(addr, value)
		}
		cpu.A ^= value
		cpu.SignFlag = cpu.A&0x80 != 0
		cpu.ZeroFlag = cpu.A == 0
	case I_STA.Num:
		cpu.memory.WriteByte(addr, cpu.A)
	case I_STX.Num:
		cpu.memory.WriteByte(addr, cpu.X)
	case I_STY.Num:
		cpu.memory.WriteByte(addr, cpu.Y)
	case I_TAX.Num:
		cpu.X = cpu.A
		cpu.SignFlag = cpu.X&0x80 != 0
		cpu.ZeroFlag = cpu.X == 0
	case I_TAY.Num:
		cpu.Y = cpu.A
		cpu.SignFlag = cpu.Y&0x80 != 0
		cpu.ZeroFlag = cpu.Y == 0
	case I_TSX.Num:
		cpu.X = cpu.SP
		cpu.SignFlag = cpu.X&0x80 != 0
		cpu.ZeroFlag = cpu.X == 0
	case I_TXA.Num:
		cpu.A = cpu.X
		cpu.SignFlag = cpu.A&0x80 != 0
		cpu.ZeroFlag = cpu.A == 0
	case I_TXS.Num:
		cpu.SP = cpu.X
	case I_TYA.Num:
		cpu.A = cpu.Y
		cpu.SignFlag = cpu.A&0x80 != 0
		cpu.ZeroFlag = cpu.A == 0
	}

	if cycles2 != cycles {
		fmt.Printf("%s am:%d expected:%d was:%d\n", opcode.Instruction.Name, opcode.AddressingMode, cycles, cycles2)
		os.Exit(1)
	}

	if jump {
		if cpu.PC&0xff00 != addr&0xff00 {
			cycles += 2
		} else {
			cycles += 1
		}
		cpu.PC = addr
	}

	cpu.Cycles += uint64(cycles)

	return cycles, nil
}

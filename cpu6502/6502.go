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
    CarryFlag bool
    ZeroFlag bool
    InterruptsDisabledFlag bool
    DecimalFlag bool            // true=BCD
    SoftwareInterruptFlag bool
    OverflowFlag bool
    SignFlag bool

    PC uint16 // program counter (PCL=low, PCH=high)
    SP uint8 // stack pointer

    Cycles uint64
    PPUCycle int
    ScanLine int

    memory MemoryAccess
}

func NewCPU6502(memory MemoryAccess) *CPU6502 {
    cpu := &CPU6502{
            memory: memory,
            SP: 0xFD,
            // P: 0x34, 24?
            InterruptsDisabledFlag: true,
            // SoftwareInterruptFlag: true,
            ScanLine: 241}
    return cpu
}

func (cpu *CPU6502) GetP() byte {
    var flags byte = 1 << 5
    if cpu.CarryFlag { flags |= FLAG_C }
    if cpu.ZeroFlag { flags |= FLAG_Z }
    if cpu.InterruptsDisabledFlag { flags |= FLAG_I }
    if cpu.DecimalFlag { flags |= FLAG_D }
    if cpu.SoftwareInterruptFlag { flags |= FLAG_B }
    if cpu.OverflowFlag { flags |= FLAG_V }
    if cpu.SignFlag { flags |= FLAG_S }
    return flags
}

func (cpu *CPU6502) SetP(flags byte) {
    cpu.CarryFlag = flags & FLAG_C != 0
    cpu.ZeroFlag = flags & FLAG_Z != 0
    cpu.InterruptsDisabledFlag = flags & FLAG_I != 0
    cpu.DecimalFlag = flags & FLAG_D != 0
    cpu.SoftwareInterruptFlag = flags & FLAG_B != 0
    cpu.OverflowFlag = flags & FLAG_V != 0
    cpu.SignFlag = flags & FLAG_S != 0
}

func (cpu *CPU6502) FlagString() string {
    flags := ""
    not_set := "_"
    if cpu.CarryFlag { flags += "C" } else { flags += not_set }
    if cpu.ZeroFlag { flags += "Z" } else { flags += not_set }
    if cpu.InterruptsDisabledFlag { flags += "I" } else { flags += not_set }
    if cpu.DecimalFlag { flags += "D" } else { flags += not_set }
    if cpu.SoftwareInterruptFlag { flags += "B" } else { flags += not_set }
    if cpu.OverflowFlag { flags += "V" } else { flags += not_set }
    if cpu.SignFlag { flags += "S" } else { flags += not_set }
    return flags
}

func (cpu *CPU6502) String() string {
    return fmt.Sprintf("CPU6502{PC:%04x SP:%02x A:%02x X:%02x Y:%02x P:%02x:%s CYC:%d SL:%d}",
        cpu.PC, cpu.SP, cpu.A, cpu.X, cpu.Y, cpu.GetP(), cpu.FlagString(), cpu.PPUCycle, cpu.ScanLine)
}

func (cpu *CPU6502) ReadByte(address uint16) byte {
    return cpu.memory.ReadByte(address)
}

func (cpu *CPU6502) ReadOpcode() (OpcodeSpec, uint16) {
    return ReadOpcode(cpu.memory, cpu.PC)
}

func (cpu *CPU6502) PushByte(value byte) {
    cpu.memory.WriteByte(0x100 + uint16(cpu.SP), value)
    cpu.SP--
}

func (cpu *CPU6502) PopByte() byte {
    cpu.SP++
    return cpu.memory.ReadByte(0x100 + uint16(cpu.SP))
}

func (cpu *CPU6502) PushAddress(addr uint16) {
    cpu.PushByte(byte(addr >> 8))
    cpu.PushByte(byte(addr & 0xff))
}

func (cpu *CPU6502) PopAddress() uint16 {
    return uint16(cpu.PopByte()) | (uint16(cpu.PopByte()) << 8)
}

func (cpu *CPU6502) Step() os.Error {
    opcode, opval := cpu.ReadOpcode()
    cpu.PC += uint16(opcode.Size)

    var addr uint16
    var value byte
    var cycles int = opcode.Cycles
    page_cycles := false
    if cycles < 0 {
        cycles = -cycles
        page_cycles = true
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
        value = cpu.memory.ReadByte(addr)
    case AMZeroPageX:
        addr = (opval + uint16(cpu.X)) & 0x00ff
        value = cpu.memory.ReadByte(addr)
    case AMAbsoluteX:
        addr = opval + uint16(cpu.X)
        if page_cycles && addr & 0xff00 != opval & 0xff00 {
            cycles += 1
        }
        value = cpu.memory.ReadByte(addr)
    case AMZeroPageY:
        addr = (opval + uint16(cpu.Y)) & 0x00ff
        value = cpu.memory.ReadByte(addr)
    case AMAbsoluteY:
        addr = opval + uint16(cpu.Y)
        if page_cycles && addr & 0xff00 != opval & 0xff00 {
            cycles += 1
        }
        value = cpu.memory.ReadByte(addr)
    case AMRelative:
        addr = cpu.PC + uint16(int8(opval))
    case AMIndirectX:
        addr = uint16(byte(opval) + cpu.X)
        addr = (uint16(cpu.memory.ReadByte(addr)) |
                (uint16(cpu.memory.ReadByte((addr + 1) & 0xff + (addr & 0xff00))) << 8))
        value = cpu.memory.ReadByte(addr)
    case AMIndirectY:
        addrt := (uint16(cpu.memory.ReadByte(opval)) |
                  (uint16(cpu.memory.ReadByte((opval + 1) & 0xff + (opval & 0xff00))) << 8))
        addr = addrt + uint16(cpu.Y)
        if page_cycles && addrt & 0xff00 != addr & 0xff00 {
            // page crossing
            cycles += 1
        }
        value = cpu.memory.ReadByte(addr)
    case AMIndirect:
        // There's a bug in the 6502 where Indirect addressing doesn't advance pages
        // 02ff -> bytes 02ff & 0200 rather than 02ff 0300
        addr = (uint16(cpu.memory.ReadByte(opval)) |
                (uint16(cpu.memory.ReadByte((opval + 1) & 0xff + (opval & 0xff00))) << 8))
    }
    
    jump := false // jump to 'addr' and account for clock

    switch opcode.Instruction {
    default:
        panic("Unhandled opcode "+opcode.InstructionName)
    case I_AAX: // undocumented
        value = cpu.X & cpu.A
        cpu.memory.WriteByte(addr, value)
    case I_ADC:
        res := uint16(value) + uint16(cpu.A)
        if cpu.CarryFlag { res++ }
        cpu.ZeroFlag = (res & 0xff) == 0
        if cpu.DecimalFlag {
            if (cpu.A ^ value ^ byte(res)) & 0x10 == 0x10 {
                res += 6
            }
            if res & 0xf0 > 0x90 {
                res += 0x60
            }
        }
        cpu.SignFlag = res & 0x80 != 0
        cpu.OverflowFlag = !((cpu.A ^ value) & 0x80 != 0) && ((uint16(cpu.A) ^ res) & 0x80 != 0)
        cpu.CarryFlag = res & 0x100 != 0
        cpu.A = byte(res & 0xff)
    case I_AND:
        cpu.A &= value
        cpu.SignFlag = cpu.A & 0x80 > 0
        cpu.ZeroFlag = cpu.A == 0
    case I_ASL:
        cpu.CarryFlag = value & 0x80 != 0
        value = value << 1
        cpu.SignFlag = value & 0x80 != 0
        cpu.ZeroFlag = value == 0
        if opcode.AddressingMode == AMAccumulator {
            cpu.A = value
        } else {
            cpu.memory.WriteByte(addr, value)
        }
    case I_BCC:
        if !cpu.CarryFlag { jump = true }
    case I_BCS:
        if cpu.CarryFlag { jump = true }
    case I_BEQ:
        if cpu.ZeroFlag { jump = true }
    case I_BIT:
        cpu.SignFlag = value & 0x80 != 0
        cpu.OverflowFlag = value & 0x40 != 0
        cpu.ZeroFlag = value & cpu.A == 0
    case I_BMI:
        if cpu.SignFlag { jump = true }
    case I_BPL:
        if !cpu.SignFlag { jump = true }
    case I_BNE:
        if !cpu.ZeroFlag { jump = true }
    case I_BVC:
        if !cpu.OverflowFlag { jump = true }
    case I_BVS:
        if cpu.OverflowFlag { jump = true }
    case I_CLC:
        cpu.CarryFlag = false
    case I_CLD:
        cpu.DecimalFlag = false
    case I_CLV:
        cpu.OverflowFlag = false
    case I_CMP:
        res := cpu.A - value
        cpu.CarryFlag = cpu.A >= value
        cpu.SignFlag = res & 0x80 != 0
        cpu.ZeroFlag = res == 0
    case I_CPX:
        res := cpu.X - value
        cpu.CarryFlag = cpu.X >= value
        cpu.SignFlag = res & 0x80 != 0
        cpu.ZeroFlag = res == 0
    case I_CPY:
        res := cpu.Y - value
        cpu.CarryFlag = cpu.Y >= value
        cpu.SignFlag = res & 0x80 != 0
        cpu.ZeroFlag = res == 0
    case I_DCP: // undocumented - equivalent to DEC, CMP
        value--
        cpu.memory.WriteByte(addr, value)
        res := cpu.A - value
        cpu.CarryFlag = cpu.A >= value
        cpu.SignFlag = res & 0x80 != 0
        cpu.ZeroFlag = res == 0
    case I_DEC:
        value--
        cpu.SignFlag = value & 0x80 != 0
        cpu.ZeroFlag = value == 0
        cpu.memory.WriteByte(addr, value)
    case I_DEX:
        cpu.X -= 1
        cpu.SignFlag = cpu.X & 0x80 != 0
        cpu.ZeroFlag = cpu.X == 0
    case I_DEY:
        cpu.Y -= 1
        cpu.SignFlag = cpu.Y & 0x80 != 0
        cpu.ZeroFlag = cpu.Y == 0
    case I_EOR:
        cpu.A ^= value
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_INC:
        value++
        cpu.SignFlag = value & 0x80 != 0
        cpu.ZeroFlag = value == 0
        cpu.memory.WriteByte(addr, value)
    case I_INX:
        cpu.X += 1
        cpu.SignFlag = cpu.X & 0x80 != 0
        cpu.ZeroFlag = cpu.X == 0
    case I_INY:
        cpu.Y += 1
        cpu.SignFlag = cpu.Y & 0x80 != 0
        cpu.ZeroFlag = cpu.Y == 0
    case I_ISC: // undocumented - equivalent to INC, SBC
        value++
        cpu.memory.WriteByte(addr, value)
        temp := uint16(cpu.A) - uint16(value)
        if !cpu.CarryFlag {
            temp--
        }
        cpu.SignFlag = temp & 0x80 != 0
        cpu.ZeroFlag = temp & 0xff == 0
        cpu.OverflowFlag = ((cpu.A ^ byte(temp)) & 0x80 != 0) && ((cpu.A ^ value) & 0x80 != 0)
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
    case I_JMP:
        cpu.PC = addr
    case I_JSR:
        cpu.PushAddress(cpu.PC-1)
        cpu.PC = addr
    case I_LAX: // undocumented
        cpu.A = value
        cpu.X = value
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_LDA:
        cpu.A = value
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_LDX:
        cpu.X = value
        cpu.SignFlag = cpu.X & 0x80 != 0
        cpu.ZeroFlag = cpu.X == 0
    case I_LDY:
        cpu.Y = value
        cpu.SignFlag = cpu.Y & 0x80 != 0
        cpu.ZeroFlag = cpu.Y == 0
    case I_LSR:
        cpu.CarryFlag = value & 0x01 > 0
        value >>= 1
        cpu.SignFlag = value & 0x80 != 0
        cpu.ZeroFlag = value == 0
        if opcode.AddressingMode == AMAccumulator {
            cpu.A = value
        } else {
            cpu.memory.WriteByte(addr, value)
        }
    case I_NOP, I_DOP, I_TOP, I_NP2:
        // no-op
    case I_ORA:
        cpu.A |= value
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_PHA:
        cpu.PushByte(cpu.A)
    case I_PHP:
        cpu.PushByte(cpu.GetP() | FLAG_B) // B flag always pushed as 1
    case I_PLA:
        cpu.A = cpu.PopByte()
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_PLP:
        cpu.SetP(cpu.PopByte() & ^byte(FLAG_B)) // B flag discarded
    case I_RLA: // undocumented - equivalent to ROL, AND
        var carry byte = 0
        if cpu.CarryFlag {
            carry = 1
        }
        cpu.CarryFlag = value & 0x80 > 0
        value = (value << 1) | carry
        if opcode.AddressingMode == AMAccumulator {
            cpu.A = value
        } else {
            cpu.memory.WriteByte(addr, value)
        }
        cpu.A &= value
        cpu.SignFlag = cpu.A & 0x80 > 0
        cpu.ZeroFlag = cpu.A == 0
    case I_RTI:
        cpu.SetP(cpu.PopByte())
        cpu.PC = cpu.PopAddress()
    case I_ROL:
        var carry byte = 0
        if cpu.CarryFlag {
            carry = 1
        }
        cpu.CarryFlag = value & 0x80 > 0
        value = (value << 1) | carry
        cpu.SignFlag = value & 0x80 > 0
        cpu.ZeroFlag = value == 0
        if opcode.AddressingMode == AMAccumulator {
            cpu.A = value
        } else {
            cpu.memory.WriteByte(addr, value)
        }
    case I_ROR:
        var carry byte = 0
        if cpu.CarryFlag {
            carry = 0x80
        }
        cpu.CarryFlag = value & 0x01 > 0
        value = (value >> 1) | carry
        cpu.SignFlag = value & 0x80 > 0
        cpu.ZeroFlag = value == 0
        if opcode.AddressingMode == AMAccumulator {
            cpu.A = value
        } else {
            cpu.memory.WriteByte(addr, value)
        }
    case I_RRA: // undocumented - equivalent to ROR, ADC
        var carry byte = 0
        if cpu.CarryFlag {
            carry = 0x80
        }
        cpu.CarryFlag = value & 0x01 > 0
        value = (value >> 1) | carry
        if opcode.AddressingMode == AMAccumulator {
            cpu.A = value
        } else {
            cpu.memory.WriteByte(addr, value)
        }
        res := uint16(value) + uint16(cpu.A)
        if cpu.CarryFlag { res++ }
        cpu.ZeroFlag = (res & 0xff) == 0
        if cpu.DecimalFlag {
            if (cpu.A ^ value ^ byte(res)) & 0x10 == 0x10 {
                res += 6
            }
            if res & 0xf0 > 0x90 {
                res += 0x60
            }
        }
        cpu.SignFlag = res & 0x80 != 0
        cpu.OverflowFlag = !((cpu.A ^ value) & 0x80 != 0) && ((uint16(cpu.A) ^ res) & 0x80 != 0)
        cpu.CarryFlag = res & 0x100 != 0
        cpu.A = byte(res & 0xff)
    case I_RTS:
        cpu.PC = cpu.PopAddress() + 1
    case I_SBC:
        temp := uint16(cpu.A) - uint16(value)
        if !cpu.CarryFlag {
            temp--
        }
        cpu.SignFlag = temp & 0x80 != 0
        cpu.ZeroFlag = temp & 0xff == 0
        cpu.OverflowFlag = ((cpu.A ^ byte(temp)) & 0x80 != 0) && ((cpu.A ^ value) & 0x80 != 0)
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
    case I_SEC:
        cpu.CarryFlag = true
    case I_SED:
        cpu.DecimalFlag = true
    case I_SEI:
        cpu.InterruptsDisabledFlag = true
    case I_SLO: // undocumented - equivalent to ASL, ORA
        cpu.CarryFlag = value & 0x80 != 0
        value = value << 1
        if opcode.AddressingMode == AMAccumulator {
            cpu.A = value
        } else {
            cpu.memory.WriteByte(addr, value)
        }

        cpu.A |= value
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_SRE: // undocumented - equivalent to LSR, EOR
        cpu.CarryFlag = value & 0x01 > 0
        value >>= 1
        if opcode.AddressingMode == AMAccumulator {
            cpu.A = value
        } else {
            cpu.memory.WriteByte(addr, value)
        }
        cpu.A ^= value
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_STA:
        cpu.memory.WriteByte(addr, cpu.A)
    case I_STX:
        cpu.memory.WriteByte(addr, cpu.X)
    case I_STY:
        cpu.memory.WriteByte(addr, cpu.Y)
    case I_TAX:
        cpu.X = cpu.A
        cpu.SignFlag = cpu.X & 0x80 != 0
        cpu.ZeroFlag = cpu.X == 0
    case I_TAY:
        cpu.Y = cpu.A
        cpu.SignFlag = cpu.Y & 0x80 != 0
        cpu.ZeroFlag = cpu.Y == 0
    case I_TSX:
        cpu.X = cpu.SP
        cpu.SignFlag = cpu.X & 0x80 != 0
        cpu.ZeroFlag = cpu.X == 0
    case I_TXA:
        cpu.A = cpu.X
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_TXS:
        cpu.SP = cpu.X
    case I_TYA:
        cpu.A = cpu.Y
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    }

    if jump {
        if cpu.PC & 0xff00 != addr & 0xff00 {
            cycles += 2
        } else {
            cycles += 1
        }
        cpu.PC = addr
    }

    cpu.Cycles += uint64(cycles)
    cpu.PPUCycle += 3*cycles
    if cpu.PPUCycle >= 341 {
        cpu.PPUCycle -= 341
        cpu.ScanLine++
        if cpu.ScanLine > 260 {
            cpu.ScanLine -= 262
        }
    }

    return nil
}
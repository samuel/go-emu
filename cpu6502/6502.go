package cpu6502

import (
    "fmt"
    "os"
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

    instructionCounter uint64

    memory MemoryAccess
}

func NewCPU6502(memory MemoryAccess) *CPU6502 {
    cpu := &CPU6502{
            memory: memory,
            SP: 0xFD,
            // P: 0x34,
            InterruptsDisabledFlag: true,
            SoftwareInterruptFlag: true}
    return cpu
}

func (cpu *CPU6502) String() string {
    return fmt.Sprintf("CPU6502{PC:%04x SP:%04x A:%02x X:%02x Y:%02x}", cpu.PC, cpu.SP, cpu.A, cpu.X, cpu.Y)
}

func (cpu *CPU6502) ReadOpcode() (DecodedOpcode, os.Error) {
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
    cpu.PushByte(byte(addr & 0xff))
    cpu.PushByte(byte(addr >> 8))
}

func (cpu *CPU6502) PopAddress() uint16 {
    return (uint16(cpu.PopByte()) << 8) | uint16(cpu.PopByte())
}

func (cpu *CPU6502) Step() os.Error {
    opcode, err := cpu.ReadOpcode()
    if err != nil {
        return err
    }
    cpu.PC += uint16(opcode.Size)

    var addr int = -1
    var value int = -1
    switch opcode.Spec.AddressingMode {
    default:
        panic(fmt.Sprintf("Unhandled addresing mode %d", opcode.Spec.AddressingMode))
    case AMImplied:
        // do nothing
    case AMAccumulator:
        // do nothing
        value = int(cpu.A)
    case AMImmediate:
        // do nothing
        value = opcode.Value
    case AMAbsolute, AMZeroPage:
        addr = opcode.Value
        value = int(cpu.memory.ReadByte(uint16(addr)))
    case AMRelative:
        addr = int(cpu.PC) + int(int8(opcode.Value))
    }

    switch opcode.Spec.Instruction {
    default:
        panic("Unhandled opcode "+opcode.Spec.InstructionName)
    case I_AND:
        cpu.A &= byte(value)
        cpu.SignFlag = cpu.A & 0x80 > 0
        cpu.ZeroFlag = cpu.A == 0
    case I_ASL:
        cpu.CarryFlag = value & 0x80 != 0
        value = (value << 1) & 0xff
        cpu.SignFlag = value & 0x80 != 0
        cpu.ZeroFlag = value == 0
        if addr < 0 {
            cpu.A = byte(value)
        } else {
            cpu.memory.WriteByte(uint16(addr), byte(value))
        }
    case I_BEQ:
        if cpu.ZeroFlag {
            cpu.PC = uint16(addr)
        }
    case I_BPL:
        if !cpu.SignFlag {
            cpu.PC = uint16(addr)
        }
    case I_CLD:
        cpu.DecimalFlag = false
    case I_JMP:
        cpu.PC = uint16(addr)
    case I_JSR:
        cpu.PushAddress(cpu.PC-1)
        cpu.PC = uint16(addr)
    case I_LDA:
        cpu.A = byte(value)
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_LDX:
        cpu.X = byte(value)
        cpu.SignFlag = cpu.X & 0x80 != 0
        cpu.ZeroFlag = cpu.X == 0
    case I_LDY:
        cpu.X = byte(value)
        cpu.SignFlag = cpu.Y & 0x80 != 0
        cpu.ZeroFlag = cpu.Y == 0
    case I_ORA:
        cpu.A |= byte(value)
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_PHA:
        cpu.PushByte(cpu.A)
    case I_PLA:
        cpu.A = cpu.PopByte()
        cpu.SignFlag = cpu.A & 0x80 != 0
        cpu.ZeroFlag = cpu.A == 0
    case I_RTS:
        cpu.PC = cpu.PopAddress() + 1
    case I_SEI:
        cpu.InterruptsDisabledFlag = true
    case I_STA:
        cpu.memory.WriteByte(uint16(addr), cpu.A)
    case I_STX:
        cpu.memory.WriteByte(uint16(addr), cpu.X)
    case I_STY:
        cpu.memory.WriteByte(uint16(addr), cpu.Y)
    case I_TXS:
        cpu.SP = cpu.X
    }

    return nil
}
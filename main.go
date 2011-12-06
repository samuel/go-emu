package main

import (
    "flag"
    "fmt"
    "log"

    "cpu6502"
    "nes"
)

var (
    f_trace = flag.Bool("t", false, "print trace while running")
    f_rom   = flag.String("r", "", "ROM file")
)

func parseFlags() {
    flag.Parse()
    if *f_rom == "" {
        log.Fatal("ROM is required (-r)")
    }
}

func main() {
    parseFlags()
    cart, err := nes.LoadCartFile(*f_rom)
    if err != nil {
        panic(err.String())
    }

    fmt.Println(cart)

    state, err := nes.NewNESState(cart)
    if err != nil {
        panic(err.String())
    }

    fmt.Println(state)
    // state.CPU.PC = 0xc000

    // for i := 0; i < 10000000; i++ {
    for {
        if (*f_trace) {
            opcode, val := state.CPU.ReadOpcode()
            fmt.Printf("%.4X  %02X ", state.CPU.PC, opcode.Opcode)
            if opcode.Size == 2 {
                fmt.Printf("%02X    ", val & 0xff)
            } else if opcode.Size == 3 {
                fmt.Printf("%02X %02X ", val & 0xff, val >> 8)
            } else {
                fmt.Printf("      ")
            }

            if opcode.Instruction.Num >= 99 || opcode.Opcode == 0x1a || opcode.Opcode == 0x3a || opcode.Opcode == 0x5a || opcode.Opcode == 0xeb {
                fmt.Printf("*")
            } else {
                fmt.Printf(" ")
            }

            mem := ""
            switch opcode.AddressingMode {
            case cpu6502.AMZeroPage, cpu6502.AMAbsolute:
                switch opcode.Instruction.Num {
                default:
                    mem = fmt.Sprintf("= %02X", state.CPU.ReadByte(val, true))
                case cpu6502.I_JMP.Num, cpu6502.I_JSR.Num:
                }
            case cpu6502.AMIndirect:
                mem = fmt.Sprintf("= %02X%02X", state.CPU.ReadByte(val+1, true), state.CPU.ReadByte(val, true))
            case cpu6502.AMIndirectX:
                addr1 := uint16(byte(val) + state.CPU.X)
                addr2 := (uint16(state.CPU.ReadByte(addr1, true)) |
                          (uint16(state.CPU.ReadByte((addr1 + 1) & 0xff + (addr1 & 0xff00), true)) << 8))
                value := state.CPU.ReadByte(addr2, true)
                mem = fmt.Sprintf("@ %02X = %04X = %02X", addr1, addr2, value)
            case cpu6502.AMIndirectY:
                addr1 := (uint16(state.CPU.ReadByte(val, true)) |
                          (uint16(state.CPU.ReadByte((val + 1) & 0xff + (val & 0xff00), true)) << 8))
                addr2 := uint16(addr1 + uint16(state.CPU.Y))
                value := state.CPU.ReadByte(addr2, true)
                mem = fmt.Sprintf("= %04X @ %04X = %02X", addr1, addr2, value)
            case cpu6502.AMAbsoluteX:
                addr := val + uint16(state.CPU.X)
                value := state.CPU.ReadByte(addr, true)
                mem = fmt.Sprintf("@ %04X = %02X", addr, value)
            case cpu6502.AMAbsoluteY:
                addr := val + uint16(state.CPU.Y)
                value := state.CPU.ReadByte(addr, true)
                mem = fmt.Sprintf("@ %04X = %02X", addr, value)
            case cpu6502.AMZeroPageX:
                addr := (val + uint16(state.CPU.X)) & 0x00ff
                value := state.CPU.ReadByte(addr, true)
                mem = fmt.Sprintf("@ %02X = %02X", addr, value)
            case cpu6502.AMZeroPageY:
                addr := (val + uint16(state.CPU.Y)) & 0x00ff
                value := state.CPU.ReadByte(addr, true)
                mem = fmt.Sprintf("@ %02X = %02X", addr, value)
            }

            fmt.Printf("%s %-27s A:%02X X:%02X Y:%02X P:%02X SP:%02X CYC:%4d SL:%d %s\n",
                opcode.Instruction.Name, opcode.FormatArguments(val, state.CPU.PC+2) + " " + mem,
                state.CPU.A, state.CPU.X, state.CPU.Y, state.CPU.GetP(), state.CPU.SP,
                state.PPUCycle, state.Scanline, state.CPU.FlagString())
        }

        state.Step()
    }

    // cpu6502.Disassemble(cart.PRGPages[len(cart.PRGPages)-1][pc-0xc000:])
}

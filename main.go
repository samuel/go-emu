package main

import (
    "fmt"

    "cpu6502"
    "nes"
)

func main() {
    cart, err := nes.LoadCartFile("nestest.nes")
    if err != nil {
        panic(err.String())
    }

    fmt.Println(cart)

    state, err := nes.NewNESState(cart)
    if err != nil {
        panic(err.String())
    }

    fmt.Println(state)
    state.CPU.PC = 0xc000

    for i := 0; i < 4096; i++ {
        opcode, val := state.CPU.ReadOpcode()
        if opcode.Opcode == cpu6502.I_XXX {
            panic("Unknown opcode")
        }
        fmt.Printf("%.4x %s %-20s %s\n", state.CPU.PC, opcode.InstructionName, opcode.FormatArguments(val), state.CPU)

        state.CPU.Step()
    }

    // cpu6502.Disassemble(cart.PRGPages[len(cart.PRGPages)-1][pc-0xc000:])
}
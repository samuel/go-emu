package main

import (
    "fmt"

    // "cpu6502"
    "nes"
)

func main() {
    cart, err := nes.LoadCartFile("SMB2.nes")
    if err != nil {
        panic(err.String())
    }

    state, err := nes.NewNESState(cart)
    if err != nil {
        panic(err.String())
    }

    fmt.Println(state)

    for i := 0; i < 50; i++ {
        opcode, err := state.CPU.ReadOpcode()
        if err != nil {
            panic(err.String())
        }
        fmt.Printf("%.4x %-20s %s\n", state.CPU.PC, opcode, state.CPU)

        state.CPU.Step()
    }
 
    // cpu6502.Disassemble(cart.PRGPages[len(cart.PRGPages)-1][pc-0xc000:])
}
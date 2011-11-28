package nes

import (
    "fmt"
    "os"
    "cpu6502"
)

const CPU_RESET_VECTOR = 0xFFFC

// CPU Memory Map (16bit buswidth, 0-FFFFh)
//   0000h-07FFh   Internal 2K Work RAM (mirrored to 800h-1FFFh)
//   2000h-2007h   Internal PPU Registers (mirrored to 2008h-3FFFh)
//   4000h-4017h   Internal APU Registers
//   4018h-5FFFh   Cartridge Expansion Area almost 8K
//   6000h-7FFFh   Cartridge SRAM Area 8K
//   8000h-FFFFh   Cartridge PRG-ROM Area 32K
// CPU Reset vector located at [FFFC], even smaller carts must have memory at that
// location. Larger carts may use whatever external mappers to access more than
// the usual 32K.

type NESState struct {
    workingRam [2048]byte   // 0000h-07FFh   Internal 2K Work RAM (mirrored to 800h-1FFFh)
    cartSRAM [8192]byte     // 6000h-7FFFh   Cartridge SRAM Area 8K
    ppuRegisters [8]byte    // 2000h-2007h (mirrored to 2008h-3fffh)
    mapper Mapper
    CPU *cpu6502.CPU6502
}

func NewNESState(cart *Cart) (*NESState, os.Error) {
    mapper, err := NewMapper(cart)
    if err != nil {
        return nil, err
    }

    state := &NESState{mapper:mapper}

    // TODO: Set workingRam to 0xFF except 0x0008=0xf7, 0x0009=0xef, 0x000a=0xdf, 0x000f=0xbf

    state.CPU = cpu6502.NewCPU6502(state)
    state.CPU.PC = uint16(state.ReadByte(CPU_RESET_VECTOR)) | (uint16(state.ReadByte(CPU_RESET_VECTOR+1)) << 8)

    return state, nil
}

func (nes *NESState) ReadByte(address uint16) byte {
    if address >= 0x0000 && address <= 0x07ff {
        return nes.workingRam[address]
    }
    if address == 0x2002 { // PPU Status Register
        // 7 = VBlank flag (reset on read and end of vblank)
        // 6 = sprite 0 hit (1=background-to-Sprite0 collision)
        // 5 = lost sprites (1=more than 8 sprites in 1 scanline)
        // 4-0 = unused
        return 0x80 // VBlank
    }
    if address >= 0x2000 && address <= 0x3fff {
        return nes.ppuRegisters[(address - 0x2000) & 7]
    }
    if address >= 0x8000 && address <= 0xffff {
        return nes.mapper.ReadByte(address)
    }
    panic("unknown address")
}

func (nes *NESState) WriteByte(address uint16, value byte) {
    if address >= 0x8000 && address <= 0xffff {
        nes.mapper.WriteByte(address, value)
    } else if address >= 0x0000 && address <= 0x07ff {
        nes.workingRam[address] = value
    } else if address >= 0x2000 && address <= 0x3fff {
        nes.ppuRegisters[(address - 0x2000) & 7] = value
    } else {
        panic("unknown address")
    }
}

func (nes *NESState) String() string {
    return fmt.Sprintf("NESState{%s}", nes.CPU)
}
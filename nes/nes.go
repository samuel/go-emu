package nes

import (
    "fmt"
    "os"
    "cpu6502"
)

const (
    SCANLINES = 262
    SCANLINE_VBLANK = 240 //243
    PIXELS_PER_SCANLINE = 1364 //341
    CPU_CYCLES_PER_VIDEO_CYCLE = 12

    // Register 2002h
    BIT_VBLANK = 0x80
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
    apuRegisters [32]byte   // 4000h-4017h
    PPUCycle int
    Scanline int
    VBlank bool
    mapper Mapper
    CPU *cpu6502.CPU6502

    testing bool
}

func NewNESState(cart *Cart) (*NESState, os.Error) {
    mapper, err := NewMapper(cart)
    if err != nil {
        return nil, err
    }

    state := &NESState{mapper:mapper} //, Scanline:241}

    // TODO: Set workingRam to 0xFF except 0x0008=0xf7, 0x0009=0xef, 0x000a=0xdf, 0x000f=0xbf

    state.CPU = cpu6502.NewCPU6502(state)
    state.CPU.PC = uint16(state.ReadByte(CPU_RESET_VECTOR, false)) | (uint16(state.ReadByte(CPU_RESET_VECTOR+1, false)) << 8)

    // TODO
    state.testing = true

    return state, nil
}

// func (nes *NESState) GetScanline() int {
//     return (nes.CPU.Cycles * 3) / 341 - 21
// }

func (nes *NESState) Step() {
    cycles, _ := nes.CPU.Step()
    nes.PPUCycle += CPU_CYCLES_PER_VIDEO_CYCLE*cycles
    if nes.PPUCycle >= PIXELS_PER_SCANLINE {
        nes.PPUCycle -= PIXELS_PER_SCANLINE
        nes.Scanline++
        if nes.Scanline >= SCANLINES {
            nes.Scanline -= SCANLINES
            nes.VBlank = false
        } else if nes.Scanline >= SCANLINE_VBLANK {
            nes.VBlank = true
        }
    }
}

func (nes *NESState) ReadByte(address uint16, peek bool) byte {
    if address >= 0x0000 && address <= 0x07ff {
        return nes.workingRam[address]
    }
    if address >= 0x2000 && address <= 0x3fff {
        trans := (address - 0x2000) & 7
        if trans == 2 { // PPU Status Register
            // 7 = VBlank flag (reset on read and end of vblank)
            // 6 = sprite 0 hit (1=background-to-Sprite0 collision)
            // 5 = lost sprites (1=more than 8 sprites in 1 scanline)
            // 4-0 = unused
            var val byte = 0
            if nes.VBlank { val |= BIT_VBLANK }
            if !peek {
                nes.VBlank = false
            }
            return val // VBlank
        }
        return nes.ppuRegisters[trans]
    }
    if address >= 0x4000 && address <= 0x4017 {
        return nes.apuRegisters[address - 0x4000]
    }
    if address >= 0x6000 && address <= 0x7fff {
        return nes.cartSRAM[address - 0x6000]
    }
    if address >= 0x8000 && address <= 0xffff {
        return nes.mapper.ReadByte(address, peek)
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
    } else if address >= 0x4000 && address <= 0x4017 {
        nes.apuRegisters[address - 0x4000] = value
    } else if address >= 0x6000 && address <= 0x7fff {
        if nes.testing {
            if address == 0x6000 {
                fmt.Printf("%.2x\n", value)
                if value < 0x80 {
                    os.Exit(1)
                }
            } else if address >= 0x6004 {
                fmt.Printf("%c", value)
            }
        }
        nes.cartSRAM[address - 0x6000] = value
    } else {
        panic("unknown address")
    }
}

func (nes *NESState) String() string {
    return fmt.Sprintf("{CPU:%s Mapper:%s}", nes.CPU, nes.mapper)
}

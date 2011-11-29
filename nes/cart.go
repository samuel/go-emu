package nes

import (
    "bytes"
    "fmt"
    "io"
    "os"
)

var (
    ErrInvalidCartFormat = os.NewError("invalid cart format")
    MapperNames = []string{"",
        "MMC1", "UNROM", "CNROM", "MMC3", "MMC5",                                     // 1-5
        "FFE F4xxx", "AOROM", "FFE F3xxx", "MMC2", "MMC4",                            // 6-10
        "Color Dreams", "", "", "", "100-in-1",                                       // 11-15
        "Bandai", "FFE F8xxx", "Jaleco SS8806", "Namcot 106", "",                     // 16-20
        "Konami VRC4", "Konami VRC2 type A", "Konami VRC2 type B", "Konami VRC6", "", // 21-25
        "", "", "", "", "",                                                           // 26-30
        "", "Irem G-101", "Taito TC0190", "Nina-1", ""}                               // 31-35
        // 64: Tengen RAMBO-1
        // 65: Irem H-3001
        // 66: GNROM
        // 68: Sunsoft Mapper #4
        // 69: Sunsoft FME-7
        // 71: Camerica
        // 78: Irem 74HC161/32
        // 91: HK-SF3
        // 72-in-1 (no number assigned)
        // 110-in-1 (no number assigned)
)

type Cart struct {
    Num8KRAMPages byte "Usualy 00h=None-or-not-specified"
    Mapper byte
    FourScreenVRAMLayout bool
    BatteryBacked bool
    VerticalMirroring bool "otherwise horizontal mirroring"
    PC10 bool "PC18 game (arcade machine with additional 8K Z80-ROM)"
    VSUnisystem bool "VS Unisystem game (arcade machine with different palette)"

    Trainer []byte // 256 bytes
    PRGPages []byte // 16KB pages
    CHRPages [][]byte // 8KB pages
    PlayChoiceROM []byte // 8KB
    
    Title string
}

func LoadCart(r io.Reader) (*Cart, os.Error) {
    var b [1024]byte

    // Read header
    if _, err := io.ReadFull(r, b[:16]); err != nil {
        return nil, err
    }

    if !bytes.Equal(b[:4], []byte("NES\x1a")) {
        return nil, ErrInvalidCartFormat
    }

    cart := &Cart{
        Mapper: b[6] >> 4,
        FourScreenVRAMLayout: b[6] & 8 != 0,
        BatteryBacked: b[6] & 2 != 0,
        VerticalMirroring: b[6] & 1 != 0}
    if b[0xf] == 0 {
        cart.Mapper |= b[7] & 0xf0
        cart.PC10 = b[7] & 2 != 0
        cart.VSUnisystem = b[7] & 1 != 0
        cart.Num8KRAMPages = b[8]
    }
    if b[6] & 4 != 0 {
        cart.Trainer = make([]byte, 256)
        if _, err := io.ReadFull(r, cart.Trainer); err != nil {
            return cart, err
        }
    }
    // 16K PRG-ROM pages + 8K CHR-ROM pages
    prg_count := int(b[4])
    chr_count := int(b[5])
    pages := make([]byte, 16384*prg_count+8192*chr_count)
    if _, err := io.ReadFull(r, pages); err != nil {
        return cart, err
    }
    // o := 0
    // for i := 0; i < prg_count; i, o = i+1, o+16384 {
    //  cart.PRGPages = append(cart.PRGPages, pages[o:o+16384])
    // }
    o := 16384*prg_count
    cart.PRGPages = pages[:o]
    for i := 0; i < chr_count; i, o = i+1, o+8192 {
        cart.CHRPages = append(cart.CHRPages, pages[o:o+8192])
    }
    // 8K Play Choice 10 Z80-ROM
    if cart.PC10 {
        cart.PlayChoiceROM = make([]byte, 8192)
        if _, err := io.ReadFull(r, cart.PlayChoiceROM); err != nil {
            return cart, err
        }
    }
    // 128 zero-padded title
    _, err := io.ReadFull(r, b[:128])
    if err == nil {
        i := bytes.IndexByte(b[:128], 0)
        cart.Title = string(b[:i])
    } else if err.String() != "EOF" {
        return cart, err
    }

    return cart, nil
}

func LoadCartFile(filename string) (*Cart, os.Error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    return LoadCart(file)
}

func (cart *Cart) String() string {
    var extra string
    if (cart.VerticalMirroring) {
        extra = "Vertical"
    } else {
        extra = "Horizontal"
    }
    if cart.PC10 {
        extra += ", PC10"
    }
    if cart.VSUnisystem {
        extra += ", VSUnisystem"
    }
    if cart.Title != "" {
        extra += ", Title=" + cart.Title
    }
    return fmt.Sprintf("Cart{Mapper=%d Num8KRAMPages=%d, %dx16k PRG, %dx8k CHR, %s}",
        cart.Mapper, cart.Num8KRAMPages, len(cart.PRGPages)/0x4000, len(cart.CHRPages), extra)
}
package nes

//
// // SPR-RAM Memory Map (8bit buswidth, 0-FFh)
// //   00-FF         Sprite Attributes (256 bytes, for 64 sprites / 4 bytes each)
// // Sprite RAM is directly built-in in the PPU chip. SPR-RAM is not connected to
// // CPU or PPU bus, and can be accessed via I/O Ports only.
//
// // I/O Map
// //   2000h - PPU Control Register 1 (W)
//   // Bit7  Execute NMI on VBlank             (0=Disabled, 1=Enabled)
//   // Bit6  PPU Master/Slave Selection        (0=Master, 1=Slave) (Not used in NES)
//   // Bit5  Sprite Size                       (0=8x8, 1=8x16)
//   // Bit4  Pattern Table Address Background  (0=VRAM 0000h, 1=VRAM 1000h)
//   // Bit3  Pattern Table Address 8x8 Sprites (0=VRAM 0000h, 1=VRAM 1000h)
//   // Bit2  Port 2007h VRAM Address Increment (0=Increment by 1, 1=Increment by 32)
//   // Bit1-0 Name Table Scroll Address        (0-3=VRAM 2000h,2400h,2800h,2C00h)
//   // (That is, Bit0=Horizontal Scroll by 256, Bit1=Vertical Scroll by 240)
// //   2001h - PPU Control Register 2 (W)
// //   2002h - PPU Status Register (R)
// //   2003h - SPR-RAM Address Register (W)
// //   2004h - SPR-RAM Data Register (RW)
// //   2005h - PPU Background Scrolling Offset (W2)
// //   2006h - VRAM Address Register (W2)
// //   2007h - VRAM Read/Write Data Register (RW)
// //   4000h - APU Channel 1 (Rectangle) Volume/Decay
// //   4001h - APU Channel 1 (Rectangle) Sweep
// //   4002h - APU Channel 1 (Rectangle) Frequency
// //   4003h - APU Channel 1 (Rectangle) Length
// //   4004h - APU Channel 2 (Rectangle) Volume/Decay
// //   4005h - APU Channel 2 (Rectangle) Sweep
// //   4006h - APU Channel 2 (Rectangle) Frequency
// //   4007h - APU Channel 2 (Rectangle) Length
// //   4008h - APU Channel 3 (Triangle) Linear Counter
// //   4009h - APU Channel 3 (Triangle) N/A
// //   400Ah - APU Channel 3 (Triangle) Frequency
// //   400Bh - APU Channel 3 (Triangle) Length
// //   400Ch - APU Channel 4 (Noise) Volume/Decay
// //   400Dh - APU Channel 4 (Noise) N/A
// //   400Eh - APU Channel 4 (Noise) Frequency
// //   400Fh - APU Channel 4 (Noise) Length
// //   4010h - APU Channel 5 (DMC) Play mode and DMA frequency
// //   4011h - APU Channel 5 (DMC) Delta counter load register
// //   4012h - APU Channel 5 (DMC) Address load register
// //   4013h - APU Channel 5 (DMC) Length register
// //   4014h - SPR-RAM DMA Register
// //   4015h - DMC/IRQ/length counter status/channel enable register (RW)
// //   4016h - Joypad #1 (RW)
// //   4017h - Joypad #2/APU SOFTCLK (RW)
// // Additionally, external hardware may contain further ports:
// //   4020h - VS Unisystem Coin Acknowlege
// //   4020h-40FFh - Famicom Disk System (FDS)
// //   4100h-FFFFh - Various addresses used by various cartridge mappers
//
// type Ports struct {
//     ppuports [8]byte
//     apuports
// }
//
// func NewPorts() *Ports {
//     p := Ports{}
//     // p.ports = make([]byte, 8, 8)
//     return &p
// }
//
// func (p *Ports) ReadByte(address uint16) byte {
//     if address >= 0x2000 && address <= 0x2007 {
//         return 0
//     }
//     if address >= 0x4000 && address <= 0x4017 {
//         return 1
//     }
//     panic("unknown port memory address")
// }
//
// func (p *Ports) WriteByte(address uint16, value byte) {
// }
//
// func (p *Ports) Slice(address uint16) []byte {
//     panic("Can't slice Ports")
// }

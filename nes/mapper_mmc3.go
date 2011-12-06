package nes

/*
+--------------------------------------------------------------------+
¦ A great majority of newer NES games (early 90's) use this mapper,  ¦
¦ both U.S. and Japanese. Among the better-known MMC3 titles are     ¦
¦ Super Mario Bros. 2 and 3, MegaMan 3, 4, 5, and 6, and Crystalis.  ¦
+--------------------------------------------------------------------+

+-------+   +------------------------------------------------------+
¦ $8000 +---¦ CPxxxNNN                                             ¦
+-------+   ¦ ¦¦   +-+                                             ¦
            ¦ ¦¦    +--- Command Number                            ¦
            ¦ ¦¦          0 - Select 2 1K VROM pages at PPU $0000  ¦
            ¦ ¦¦          1 - Select 2 1K VROM pages at PPU $0800  ¦
            ¦ ¦¦          2 - Select 1K VROM page at PPU $1000     ¦
            ¦ ¦¦          3 - Select 1K VROM page at PPU $1400     ¦
            ¦ ¦¦          4 - Select 1K VROM page at PPU $1800     ¦
            ¦ ¦¦          5 - Select 1K VROM page at PPU $1C00     ¦
            ¦ ¦¦          6 - Select first switchable ROM page     ¦
            ¦ ¦¦          7 - Select second switchable ROM page    ¦
            ¦ ¦¦                                                   ¦
            ¦ ¦+-------- PRG Address Select                        ¦
            ¦ ¦           0 - Enable swapping for $8000 and $A000  ¦
            ¦ ¦           1 - Enable swapping for $A000 and $C000  ¦
            ¦ ¦                                                    ¦
            ¦ +--------- CHR Address Select                        ¦
            ¦             0 - Use normal address for commands 0-5  ¦
            ¦             1 - XOR command 0-5 address with $1000   ¦
            +------------------------------------------------------+

+-------+   +----------------------------------------------+
¦ $8001 +---¦ PPPPPPPP                                     ¦
+-------+   ¦ +------+                                     ¦
            ¦    ¦                                         ¦
            ¦    ¦                                         ¦
            ¦    +------- Page Number for Command          ¦
            ¦              Activates the command number    ¦
            ¦              written to bits 0-2 of $8000    ¦
            +----------------------------------------------+

+-------+   +----------------------------------------------+
¦ $A000 +---¦ xxxxxxxM                                     ¦
+-------+   ¦        ¦                                     ¦
            ¦        ¦                                     ¦
            ¦        ¦                                     ¦
            ¦        +--- Mirroring Select                 ¦
            ¦              0 - Horizontal mirroring        ¦
            ¦              1 - Vertical mirroring          ¦
            ¦ NOTE: I don't have any confidence in the     ¦
            ¦       accuracy of this information.          ¦
            +----------------------------------------------+

+-------+   +----------------------------------------------+
¦ $A001 +---¦ Sxxxxxxx                                     ¦
+-------+   ¦ ¦                                            ¦
            ¦ ¦                                            ¦
            ¦ ¦                                            ¦
            ¦ +---------- SaveRAM Toggle                   ¦
            ¦              0 - Disable $6000-$7FFF         ¦
            ¦              1 - Enable $6000-$7FFF          ¦
            +----------------------------------------------+

+-------+   +----------------------------------------------+
¦ $C000 +---¦ IIIIIIII                                     ¦
+-------+   ¦ +------+                                     ¦
            ¦    ¦                                         ¦
            ¦    ¦                                         ¦
            ¦    +------- IRQ Counter Register             ¦
            ¦              The IRQ countdown value is      ¦
            ¦              stored here.                    ¦
            +----------------------------------------------+

+-------+   +----------------------------------------------+
¦ $C001 +---¦ IIIIIIII                                     ¦
+-------+   ¦ +------+                                     ¦
            ¦    ¦                                         ¦
            ¦    ¦                                         ¦
            ¦    +------- IRQ Latch Register               ¦
            ¦              A temporary value is stored     ¦
            ¦              here.                           ¦
            +----------------------------------------------+

+-------+   +----------------------------------------------+
¦ $E000 +---¦ xxxxxxxx                                     ¦
+-------+   ¦ +------+                                     ¦
            ¦    ¦                                         ¦
            ¦    ¦                                         ¦
            ¦    +------- IRQ Control Register 0           ¦
            ¦              Any value written here will     ¦
            ¦              disable IRQ's and copy the      ¦
            ¦              latch register to the actual    ¦
            ¦              IRQ counter register.           ¦
            +----------------------------------------------+

+-------+   +----------------------------------------------+
¦ $E001 +---¦ xxxxxxxx                                     ¦
+-------+   ¦ +------+                                     ¦
            ¦    ¦                                         ¦
            ¦    ¦                                         ¦
            ¦    +------- IRQ Control Register 1           ¦
            ¦              Any value written here will     ¦
            ¦              enable IRQ's.                   ¦
            +----------------------------------------------+

Notes: - Two of the 8K ROM banks in the PRG area are switchable.
          The other two are "hard-wired" to the last two banks in
          the cart. The default setting is switchable banks at
          $8000 and $A000, with banks 0 and 1 being swapped in
          at reset. Through bit 6 of $8000, the hard-wiring can
          be made to affect $8000 and $E000 instead of $C000 and
          $E000. The switchable banks, whatever their addresses,
          can be swapped through commands 6 and 7.
       - A cart will first write the command and base select number
          to $8000, then the value to be used to $8001.
       - On carts with VROM, the first 8K of VROM is swapped into
          PPU $0000 on reset. On carts without VROM, as always, there
          is 8K of VRAM at PPU $0000.
*/

type MapperMMC3 struct {
    cart *Cart
    prg_banks []int
}

func (m *MapperMMC3) String() string {
    return "{Type:MMC3}"
}

func NewMapperMMC3(cart *Cart) *MapperMMC3 {
    num_pages := len(cart.PRGPages) / 8192
    return &MapperMMC3{cart:cart, prg_banks:[]int{0, 1, num_pages-2, num_pages-1}}
    // for o := 0, 0; o < len(cart.PRGPages); o += 8192 {
    //     mapper.prg_banks = append(mapper.prg_banks, cart.PRGPages[o:o+8192])
    // }
}

func (m *MapperMMC3) ReadByte(address uint16, peek bool) byte {
    addr := m.translateAddress(address)
    return m.cart.PRGPages[addr]
}

func (m *MapperMMC3) WriteByte(address uint16, value byte) {
    addr := m.translateAddress(address)
    m.cart.PRGPages[addr] = value
}

func (m *MapperMMC3) translateAddress(address uint16) int {
    if address < 0x8000 {
        panic("address out of range")
    }
    return m.prg_banks[(address - 0x8000) / 0x2000] * 8192 + int(address - 0x8000) & 0x1fff
}

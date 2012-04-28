package nes

// +----------------+
// ¦ Mapper 0: NROM ¦
// +----------------+

type MapperNROM struct {
	cart *Cart
}

func (m *MapperNROM) String() string {
	return "{Type:NROM}"
}

func NewMapperNROM(cart *Cart) *MapperNROM {
	return &MapperNROM{cart: cart}
}

func (m *MapperNROM) ReadByte(address uint16, peek bool) byte {
	addr := m.translateAddress(address)
	return m.cart.PRGPages[addr]
}

func (m *MapperNROM) WriteByte(address uint16, value byte) {
	addr := m.translateAddress(address)
	m.cart.PRGPages[addr] = value
}

func (m *MapperNROM) translateAddress(address uint16) uint16 {
	if address < 0x8000 {
		panic("address out of range")
	}
	// Mirror first 16K if that's all there is
	if len(m.cart.PRGPages) <= 0x4000 && address >= 0xc000 {
		return address - 0xc000
	}
	return address - 0x8000
}

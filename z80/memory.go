package z80

type MemoryAccess interface {
	ReadByte(address uint16, peek bool) byte
	WriteByte(address uint16, value byte)
}

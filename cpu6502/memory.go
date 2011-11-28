package cpu6502

type MemoryAccess interface {
    ReadByte(address uint16) byte
    WriteByte(address uint16, value byte)
}

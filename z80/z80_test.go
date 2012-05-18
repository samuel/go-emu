package z80

import (
	"testing"
)

type TestMemory struct {
	bytes []byte
}

func NewTestMemory(bytes []byte) *TestMemory {
	mem := TestMemory{}
	mem.bytes = make([]byte, 0x200, 0x200)
	for i := 0; i < len(bytes); i++ {
		mem.bytes[i] = bytes[i]
	}
	return &mem
}

func (m *TestMemory) ReadByte(addr uint16, peek bool) byte {
	return m.bytes[addr]
}

func (m *TestMemory) WriteByte(addr uint16, value byte) {
	m.bytes[addr] = value
}

func TestStack(t *testing.T) {
	memory := NewTestMemory([]byte{0})
	cpu := New(memory)
	_ = cpu
}

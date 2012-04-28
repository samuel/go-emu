package cpu6502

import (
    "testing"
)

type InstructionTest struct {
    bytes []byte
}

// []InstructionTest{
//     {0x29, 0xf0}
// }

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
    cpu := NewCPU6502(memory)
    sp := cpu.SP

    cpu.PushByte(1)
    if cpu.SP != sp - 1 { t.Errorf("Pushing to stack did not decrement SP") }
    cpu.PushByte(2)
    if cpu.PopByte() != 2 { t.Errorf("Stack returned wrong value") }
    if cpu.PopByte() != 1 { t.Errorf("Stack returned wrong value") }
    if cpu.SP != sp { t.Errorf("Popping from stack didn't increment SP") }

    cpu.PushAddress(0x1234)
    if cpu.SP != sp - 2 { t.Errorf("Pushing address to stack did not decrement SP by 2") }
    if cpu.PopAddress() != 0x1234 { t.Errorf("Popping address from stack returned wrong value") }
    if cpu.SP != sp { t.Errorf("Popping address from stack didn't increment SP by 2") }
}

func TestAND(t *testing.T) {
    memory := NewTestMemory([]byte{
        0x29, 0xf0,  // AND #$f0
        0x29, 0x0f}) // AND #$0f
    cpu := NewCPU6502(memory)
    cpu.A = 0xff
    cpu.Step()
    if cpu.A != 0xf0 { t.Errorf("AND/Immediate didn't produce correct result")  }
    if cpu.ZeroFlag { t.Errorf("AND/Immediate set zero flag when it shouldn't have") }
    if !cpu.SignFlag { t.Errorf("AND/Immediate failed to set sign flag") }
    cpu.Step()
    if cpu.A != 0x00 { t.Errorf("AND/Immediate didn't produce correct result")  }
    if !cpu.ZeroFlag { t.Errorf("AND/Immediate failed to set zero flag") }
    if cpu.SignFlag { t.Errorf("AND/Immediate set sign flag when it shouldn't have") }
}

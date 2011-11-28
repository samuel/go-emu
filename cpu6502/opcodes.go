package cpu6502

import (
    "fmt"
    "os"
)

const (
    // Addressing Modes
    AMImmediate   = 1
    AMZeroPage    = 2
    AMZeroPageX   = 3
    AMZeroPageY   = 4
    AMAbsolute    = 5
    AMAbsoluteX   = 6
    AMAbsoluteY   = 7
    AMIndirect    = 8
    AMIndirectX   = 9
    AMIndirectY   = 10
    AMAccumulator = 11
    AMRelative    = 12
    AMImplied     = 13

    I_XXX = 0
    I_ADC = 1
    I_AND = 2
    I_ASL = 3
    I_BCC = 4
    I_BCS = 5
    I_BEQ = 6
    I_BIT = 7
    I_BMI = 8
    I_BNE = 9
    I_BPL = 10
    I_BRK = 11
    I_BVC = 12
    I_BVS = 13
    I_CLC = 14
    I_CLD = 15
    I_CLI = 16
    I_CLV = 17
    I_CMP = 18
    I_CPX = 19
    I_CPY = 20
    I_DEC = 21
    I_DEX = 22
    I_DEY = 23
    I_EOR = 24
    I_INC = 25
    I_INX = 26
    I_INY = 27
    I_JMP = 28
    I_JSR = 29
    I_LDA = 30
    I_LDX = 31
    I_LDY = 32
    I_LSR = 33
    I_NOP = 34
    I_ORA = 35
    I_PHA = 36
    I_PHP = 37
    I_PLA = 38
    I_PLP = 39
    I_ROL = 40
    I_ROR = 41
    I_RTI = 42
    I_RTS = 43
    I_SBC = 44
    I_SEC = 45
    I_SED = 46
    I_SEI = 47
    I_STA = 48
    I_STX = 49
    I_STY = 50
    I_TAX = 51
    I_TAY = 52
    I_TSX = 53
    I_TXA = 54
    I_TXS = 55
    I_TYA = 56
)

type OpcodeSpec struct {
    Opcode int
    Instruction int
    InstructionName string
    AddressingMode int
    // Size int
}

type DecodedOpcode struct {
    Spec OpcodeSpec
    Size int
    Arguments string
    Value int
}

var (
    opcodes = [256]OpcodeSpec{
        {0x00, I_BRK, "BRK", AMImplied},     {0x01, I_ORA, "ORA", AMIndirectX},
        {0x02, I_XXX, "", 0},                {0x03, I_XXX, "", 0},
        {0x04, I_XXX, "", 0},                {0x05, I_ORA, "ORA", AMZeroPage},
        {0x06, I_ASL, "ASL", AMZeroPage},    {0x07, I_XXX, "", 0},
        {0x08, I_PHP, "PHP", AMImplied},     {0x09, I_ORA, "ORA", AMImmediate},
        {0x0a, I_ASL, "ASL", AMAccumulator}, {0x0b, I_XXX, "", 0},
        {0x0c, I_XXX, "", 0},                {0x0d, I_ORA, "ORA", AMAbsolute},
        {0x0e, I_ASL, "ASL", AMAbsolute},    {0x0f, I_XXX, "", 0},
        {0x10, I_BPL, "BPL", AMRelative},    {0x11, I_ORA, "ORA", AMIndirectY},
        {0x12, I_XXX, "", 0},                {0x13, I_XXX, "", 0},
        {0x14, I_XXX, "", 0},                {0x15, I_ORA, "ORA", AMZeroPageX},
        {0x16, I_ASL, "ASL", AMZeroPageX},   {0x17, I_XXX, "", 0},
        {0x18, I_CLC, "CLC", AMImplied},     {0x19, I_ORA, "ORA", AMAbsoluteY},
        {0x1a, I_XXX, "", 0},                {0x1b, I_XXX, "", 0},
        {0x1c, I_XXX, "", 0},                {0x1d, I_ORA, "ORA", AMAbsoluteX},
        {0x1e, I_ASL, "ASL", AMAbsoluteX},   {0x1f, I_XXX, "", 0},
        {0x20, I_JSR, "JSR", AMAbsolute},    {0x21, I_AND, "AND", AMIndirectX},
        {0x22, I_XXX, "", 0},                {0x23, I_XXX, "", 0},
        {0x24, I_BIT, "BIT", AMZeroPage},    {0x25, I_AND, "AND", AMZeroPage},
        {0x26, I_ROL, "ROL", AMZeroPage},    {0x27, I_XXX, "", 0},
        {0x28, I_PLP, "PLP", AMImplied},     {0x29, I_AND, "AND", AMImmediate},
        {0x2a, I_ROL, "ROL", AMAccumulator}, {0x2b, I_XXX, "", 0},
        {0x2c, I_BIT, "BIT", AMAbsolute},    {0x2d, I_AND, "AND", AMAbsolute},
        {0x2e, I_ROL, "ROL", AMAbsolute},    {0x2f, I_XXX, "", 0},
        {0x30, I_BMI, "BMI", AMRelative},    {0x31, I_AND, "AND", AMIndirectY},
        {0x32, I_XXX, "", 0},                {0x33, I_XXX, "", 0},
        {0x34, I_XXX, "", 0},                {0x35, I_AND, "AND", AMZeroPageX},
        {0x36, I_ROL, "ROL", AMZeroPageX},   {0x37, I_XXX, "", 0},
        {0x38, I_SEC, "SEC", AMImplied},     {0x39, I_AND, "AND", AMAbsoluteY},
        {0x3a, I_XXX, "", 0},                {0x3b, I_XXX, "", 0},
        {0x3c, I_XXX, "", 0},                {0x3d, I_AND, "AND", AMAbsoluteX},
        {0x3e, I_ROL, "ROL", AMAbsoluteX},   {0x3f, I_XXX, "", 0},
        {0x40, I_RTI, "RTI", AMImplied},     {0x41, I_EOR, "EOR", AMIndirectX},
        {0x42, I_XXX, "", 0},                {0x43, I_XXX, "", 0},
        {0x44, I_XXX, "", 0},                {0x45, I_EOR, "EOR", AMZeroPage},
        {0x46, I_LSR, "LSR", AMZeroPage},    {0x47, I_XXX, "", 0},
        {0x48, I_PHA, "PHA", AMImplied},     {0x49, I_EOR, "EOR", AMImmediate},
        {0x4a, I_LSR, "LSR", AMAccumulator}, {0x4b, I_XXX, "", 0},
        {0x4c, I_JMP, "JMP", AMAbsolute},    {0x4d, I_EOR, "EOR", AMAbsolute},
        {0x4e, I_LSR, "LSR", AMAbsolute},    {0x4f, I_XXX, "", 0},
        {0x50, I_BVC, "BVC", AMRelative},    {0x51, I_EOR, "EOR", AMIndirectY},
        {0x52, I_XXX, "", 0},                {0x53, I_XXX, "", 0},
        {0x54, I_XXX, "", 0},                {0x55, I_EOR, "EOR", AMZeroPageX},
        {0x56, I_LSR, "LSR", AMZeroPageX},   {0x57, I_XXX, "", 0},
        {0x58, I_CLI, "CLI", AMImplied},     {0x59, I_EOR, "EOR", AMAbsoluteY},
        {0x5a, I_XXX, "", 0},                {0x5b, I_XXX, "", 0},
        {0x5c, I_XXX, "", 0},                {0x5d, I_EOR, "EOR", AMAbsoluteX},
        {0x5e, I_LSR, "LSR", AMAbsoluteX},   {0x5f, I_XXX, "", 0},
        {0x60, I_RTS, "RTS", AMImplied},     {0x61, I_ADC, "ADC", AMIndirectX},
        {0x62, I_XXX, "", 0},                {0x63, I_XXX, "", 0},
        {0x64, I_XXX, "", 0},                {0x65, I_ADC, "ADC", AMZeroPage},
        {0x66, I_ROR, "ROR", AMZeroPage},    {0x67, I_XXX, "", 0},
        {0x68, I_PLA, "PLA", AMImplied},     {0x69, I_ADC, "ADC", AMImmediate},
        {0x6a, I_ROR, "ROR", AMAccumulator}, {0x6b, I_XXX, "", 0},
        {0x6c, I_JMP, "JMP", AMIndirect},    {0x6d, I_ADC, "ADC", AMAbsolute},
        {0x6e, I_ROR, "ROR", AMAbsolute},    {0x6f, I_XXX, "", 0},
        {0x70, I_BVS, "BVS", AMRelative},    {0x71, I_ADC, "ADC", AMIndirectY},
        {0x72, I_XXX, "", 0},                {0x73, I_XXX, "", 0},
        {0x74, I_XXX, "", 0},                {0x75, I_ADC, "ADC", AMZeroPageX},
        {0x76, I_ROR, "ROR", AMZeroPageX},   {0x77, I_XXX, "", 0},
        {0x78, I_SEI, "SEI", AMImplied},     {0x79, I_ADC, "ADC", AMAbsoluteY},
        {0x7a, I_XXX, "", 0},                {0x7b, I_XXX, "", 0},
        {0x7c, I_XXX, "", 0},                {0x7d, I_ADC, "ADC", AMAbsoluteX},
        {0x7e, I_ROR, "ROR", AMAbsoluteX},   {0x7f, I_XXX, "", 0},
        {0x80, I_XXX, "", 0},                {0x81, I_STA, "STA", AMIndirectX},
        {0x82, I_XXX, "", 0},                {0x83, I_XXX, "", 0},
        {0x84, I_STY, "STY", AMZeroPage},    {0x85, I_STA, "STA", AMZeroPage},
        {0x86, I_STX, "STX", AMZeroPage},    {0x87, I_XXX, "", 0},
        {0x88, I_DEY, "DEY", AMImplied},     {0x89, I_XXX, "", 0},
        {0x8a, I_TXA, "TXA", AMImplied},     {0x8b, I_XXX, "", 0},
        {0x8c, I_STY, "STY", AMAbsolute},    {0x8d, I_STA, "STA", AMAbsolute},
        {0x8e, I_STX, "STX", AMAbsolute},    {0x8f, I_XXX, "", 0},
        {0x90, I_BCC, "BCC", AMRelative},    {0x91, I_STA, "STA", AMIndirectY},
        {0x92, I_XXX, "", 0},                {0x93, I_XXX, "", 0},
        {0x94, I_STY, "STY", AMZeroPageX},   {0x95, I_STA, "STA", AMZeroPageX},
        {0x96, I_STX, "STX", AMZeroPageY},   {0x97, I_XXX, "", 0},
        {0x98, I_TYA, "TYA", AMImplied},     {0x99, I_STA, "STA", AMAbsoluteY},
        {0x9a, I_TXS, "TXS", AMImplied},     {0x9b, I_XXX, "", 0},
        {0x9c, I_XXX, "", 0},                {0x9d, I_STA, "STA", AMAbsoluteX},
        {0x9e, I_XXX, "", 0},                {0x9f, I_XXX, "", 0},
        {0xa0, I_LDY, "LDY", AMImmediate},   {0xa1, I_LDA, "LDA", AMIndirectX},
        {0xa2, I_LDX, "LDX", AMImmediate},   {0xa3, I_XXX, "", 0},
        {0xa4, I_LDY, "LDY", AMZeroPage},    {0xa5, I_LDA, "LDA", AMZeroPage},
        {0xa6, I_LDX, "LDX", AMZeroPage},    {0xa7, I_XXX, "", 0},
        {0xa8, I_TAY, "TAY", AMImplied},     {0xa9, I_LDA, "LDA", AMImmediate},
        {0xaa, I_TAX, "TAX", AMImplied},     {0xab, I_XXX, "", 0},
        {0xac, I_LDY, "LDY", AMAbsolute},    {0xad, I_LDA, "LDA", AMAbsolute},
        {0xae, I_LDX, "LDX", AMAbsolute},    {0xaf, I_XXX, "", 0},
        {0xb0, I_BCS, "BCS", AMRelative},    {0xb1, I_LDA, "LDA", AMIndirectY},
        {0xb2, I_XXX, "", 0},                {0xb3, I_XXX, "", 0},
        {0xb4, I_LDY, "LDY", AMZeroPageX},   {0xb5, I_LDA, "LDA", AMZeroPageX},
        {0xb6, I_LDX, "LDX", AMZeroPageY},   {0xb7, I_XXX, "", 0},
        {0xb8, I_CLV, "CLV", AMImplied},     {0xb9, I_LDA, "LDA", AMAbsoluteY},
        {0xba, I_TSX, "TSX", AMImplied},     {0xbb, I_XXX, "", 0},
        {0xbc, I_LDY, "LDY", AMAbsoluteX},   {0xbd, I_LDA, "LDA", AMAbsoluteX},
        {0xbe, I_LDX, "LDX", AMAbsoluteY},   {0xbf, I_XXX, "", 0},
        {0xc0, I_CPY, "CPY", AMImmediate},   {0xc1, I_CMP, "CMP", AMIndirectX},
        {0xc2, I_XXX, "", 0},                {0xc3, I_XXX, "", 0},
        {0xc4, I_CPY, "CPY", AMZeroPage},    {0xc5, I_CMP, "CMP", AMZeroPage},
        {0xc6, I_DEC, "DEC", AMZeroPage},    {0xc7, I_XXX, "", 0},
        {0xc8, I_INY, "INY", AMImplied},     {0xc9, I_CMP, "CMP", AMImmediate},
        {0xca, I_DEC, "DEC", AMImplied},     {0xcb, I_XXX, "", 0},
        {0xcc, I_CPY, "CPY", AMAbsolute},    {0xcd, I_CMP, "CMP", AMAbsolute},
        {0xce, I_DEC, "DEC", AMAbsolute},    {0xcf, I_XXX, "", 0},
        {0xd0, I_BNE, "BNE", AMRelative},    {0xd1, I_CMP, "CMP", AMIndirectY},
        {0xd2, I_XXX, "", 0},                {0xd3, I_XXX, "", 0},
        {0xd4, I_XXX, "", 0},                {0xd5, I_CMP, "CMP", AMZeroPageX},
        {0xd6, I_DEC, "DEC", AMZeroPageX},   {0xd7, I_XXX, "", 0},
        {0xd8, I_CLD, "CLD", AMImplied},     {0xd9, I_CMP, "CMP", AMAbsoluteY},
        {0xda, I_XXX, "", 0},                {0xdb, I_XXX, "", 0},
        {0xdc, I_XXX, "", 0},                {0xdd, I_CMP, "CMP", AMAbsoluteX},
        {0xde, I_DEC, "DEC", AMAbsoluteX},   {0xdf, I_XXX, "", 0},
        {0xe0, I_CPX, "CPX", AMImmediate},   {0xe1, I_SBC, "SBC", AMIndirectX},
        {0xe2, I_XXX, "", 0},                {0xe3, I_XXX, "", 0},
        {0xe4, I_CPX, "CPX", AMZeroPage},    {0xe5, I_SBC, "SBC", AMZeroPage},
        {0xe6, I_INC, "INC", AMZeroPage},    {0xe7, I_XXX, "", 0},
        {0xe8, I_INX, "INX", AMImplied},     {0xe9, I_SBC, "SBC", AMImmediate},
        {0xea, I_NOP, "NOP", AMImplied},     {0xeb, I_XXX, "", 0},
        {0xec, I_CPX, "CPX", AMAbsolute},    {0xed, I_SBC, "SBC", AMAbsolute},
        {0xee, I_INC, "INC", AMAbsolute},    {0xef, I_XXX, "", 0},
        {0xf0, I_BEQ, "BEQ", AMRelative},    {0xf1, I_SBC, "SBC", AMIndirectY},
        {0xf2, I_XXX, "", 0},                {0xf3, I_XXX, "", 0},
        {0xf4, I_XXX, "", 0},                {0xf5, I_SBC, "SBC", AMZeroPageX},
        {0xf6, I_INC, "INC", AMZeroPageX},   {0xf7, I_XXX, "", 0},
        {0xf8, I_SED, "SED", AMImplied},     {0xf9, I_SBC, "SBC", AMAbsoluteY},
        {0xfa, I_XXX, "", 0},                {0xfb, I_XXX, "", 0},
        {0xfc, I_XXX, "", 0},                {0xfd, I_SBC, "SBC", AMAbsoluteX},
        {0xfe, I_INC, "INC", AMAbsoluteX},   {0xff, I_XXX, "", 0}}
)

func ReadOpcode(memory MemoryAccess, address uint16) (DecodedOpcode, os.Error) {
    switch spec := opcodes[memory.ReadByte(address)]; spec.AddressingMode {
    case AMImplied, AMAccumulator:
        return DecodedOpcode{spec, 1, "", 0}, nil
    case AMImmediate, AMIndirect, AMIndirectX, AMIndirectY, AMRelative, AMZeroPage, AMZeroPageX, AMZeroPageY:
        return DecodedOpcode{spec, 2, "", int(memory.ReadByte(address+1))}, nil
    case AMAbsolute, AMAbsoluteX, AMAbsoluteY:
        low := memory.ReadByte(address+1)
        high := memory.ReadByte(address+2)
        value := int(high) << 8 | int(low)
        return DecodedOpcode{spec, 3, "", value}, nil
    }
    return DecodedOpcode{}, os.NewError("unknown opcode")
}

func (dec DecodedOpcode) GetArguments() string {
    switch dec.Spec.AddressingMode {
    case AMAccumulator:
        dec.Arguments = "A"
    case AMImmediate:
        dec.Arguments = fmt.Sprintf("#$%.2x", dec.Value)
    case AMIndirect:
        dec.Arguments = fmt.Sprintf("($%.2x)", dec.Value)
    case AMIndirectX:
        dec.Arguments = fmt.Sprintf("($%.2x,X)", dec.Value)
    case AMIndirectY:
        dec.Arguments = fmt.Sprintf("($%.2x),Y", dec.Value)
    case AMRelative, AMZeroPage:
        dec.Arguments = fmt.Sprintf("$%.2x", dec.Value)
    case AMZeroPageX:
        dec.Arguments = fmt.Sprintf("$%.2x,X", dec.Value)
    case AMZeroPageY:
        dec.Arguments = fmt.Sprintf("$%.2x,Y", dec.Value)
    case AMAbsolute:
        dec.Arguments = fmt.Sprintf("$%.4x", dec.Value)
    case AMAbsoluteX:
        dec.Arguments = fmt.Sprintf("$%.4x,X", dec.Value)
    case AMAbsoluteY:
        dec.Arguments = fmt.Sprintf("$%.4x,Y", dec.Value)
    }
    return dec.Arguments
}

func (dec DecodedOpcode) String() string {
    return fmt.Sprintf("%s %s", dec.Spec.InstructionName, dec.GetArguments())
}
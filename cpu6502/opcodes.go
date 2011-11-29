package cpu6502

import (
    "fmt"
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
    // Undocumented/invalid opcodes
    I_UNDOCUMENTED = 57
    I_KIL = 57 // Stop program counter (processor lock up)
    I_DOP = 58 // double NOP
    I_SLO = 59 // Shift left one bit in memory, then OR accumulator with memory. Status flags: N, Z, C
    I_AAC = 60 // (ANC) AND byte with accumulator. If result is negative then carry is set. Status flags: N, Z, C
    I_TOP = 61 // triple NOP
    I_NP2 = 62 // NOP - undocumented
    I_RLA = 63 // rotate one bit left in memory, then AND accumulator with memory. Status flags: N, Z, C
    I_SRE = 64 // Shift right one bit in memory, then EOR accumulator with memory. Status flags: N,Z,C
    I_ASR = 65 // [ARR] AND byte with accumulator, then shift right one bit in accumulator. Status flags: N,Z,C
    I_RRA = 66 // Rotate one bit right in memory, then add memory to accumulator (with carry). Status flags: N,V,Z,C
    I_AAR = 67 // AND byte with accumulator, then rotate one bit right in accumulator and check bit 5 and 6:
               // If both bits are 1: set C, clear V.
               // If both bits are 0: clear C and V.
               // If only bit 5 is 1: set V, clear C.
               // If only bit 6 is 1: set C and V.
               // Status flags: N,V,Z,C
    I_AAX = 68 // (SAX) [AXS] AND X register with accumulator and store result in memory. Status flags: N,Z
    I_XAA = 69 // (ANE) Exact operation unknown. Read the referenced documents for more information and observations.
    I_AXA = 70 // (SHA) AND X register with accumulator then AND result with 7 and store in memory. Status flags: -
    I_XAS = 71 // (SHS) [TAS] AND X register with accumulator and store result in stack
               // pointer, then AND stack pointer with the high byte of the
               // target address of the argument + 1. Store result in memory.
               // S = X AND A, M = S AND HIGH(arg) + 1
               // Status flags: -
    I_SYA = 72 // (SHY) [SAY] AND Y register with the high byte of the target address of the argument + 1. Store the result in memory.
               // M = Y AND HIGH(arg) + 1
               // Status flags: -
    I_SXA = 73 // (SHX) [XAS] AND X register with the high byte of the target address of the argument + 1. Store the result in memory.
               // M = X AND HIGH(arg) + 1
               // Status flags: -
    I_LAX = 74 // Load accumulator and X register with memory. Status flags: N,Z
    I_ATX = 75 // (LXA) [OAL] AND byte with accumulator, then transfer accumulator to X register. Status flags: N,Z
    I_LAR = 76 // (LAE) [LAS] AND memory with stack pointer, transfer result to accumulator, X register and stack pointer. Status flags: N,Z
    I_DCP = 77 // [DCM] Subtract 1 from memory (without borrow). Status flags: C
    I_AXS = 78 // (SBX) [SAX] AND X register with accumulator and store result in X register, then subtract byte from X register (without borrow). Status flags: N,Z,C
    I_ISC = 79 // (ISB) [INS] Increase memory by one, then subtract memory from accumulator (with borrow). Status flags: N,V,Z,C
    I_SB2 = I_SBC // Same as legal opcode $E9 (SBC #byte)
)

type OpcodeSpec struct {
    Opcode int
    Instruction int
    InstructionName string
    Size int
    AddressingMode int
    Cycles int
}

var (
    opcodes = [256]OpcodeSpec{
        {0x00, I_BRK, "BRK", 1, AMImplied, 7},     {0x01, I_ORA, "ORA", 2, AMIndirectX, 6},
        {0x02, I_KIL, "KIL", 1, AMImplied, 0},     {0x03, I_SLO, "SLO", 2, AMIndirectX, 8},
        {0x04, I_DOP, "DOP", 2, AMZeroPage, 3},    {0x05, I_ORA, "ORA", 2, AMZeroPage, 3},
        {0x06, I_ASL, "ASL", 2, AMZeroPage, 5},    {0x07, I_SLO, "SLO", 2, AMZeroPage, 5},
        {0x08, I_PHP, "PHP", 1, AMImplied, 3},     {0x09, I_ORA, "ORA", 2, AMImmediate, 2},
        {0x0a, I_ASL, "ASL", 1, AMAccumulator, 2}, {0x0b, I_AAC, "AAC", 2, AMImmediate, 2},
        {0x0c, I_TOP, "TOP", 3, AMAbsolute, 4},    {0x0d, I_ORA, "ORA", 3, AMAbsolute, 4},
        {0x0e, I_ASL, "ASL", 3, AMAbsolute, 6},    {0x0f, I_SLO, "SLO", 3, AMAbsolute, 6},
        {0x10, I_BPL, "BPL", 2, AMRelative, -2},   {0x11, I_ORA, "ORA", 2, AMIndirectY, 5},
        {0x12, I_KIL, "KIL", 1, AMImplied, 0},     {0x13, I_SLO, "SLO", 2, AMIndirectY, 8},
        {0x14, I_DOP, "DOP", 2, AMZeroPageX, 4},   {0x15, I_ORA, "ORA", 2, AMZeroPageX, 4},
        {0x16, I_ASL, "ASL", 2, AMZeroPageX, 6},   {0x17, I_SLO, "SLO", 2, AMZeroPageX, 6},
        {0x18, I_CLC, "CLC", 1, AMImplied, 2},     {0x19, I_ORA, "ORA", 3, AMAbsoluteY, -4},
        {0x1a, I_NP2, "NOP", 1, AMImplied, 2},     {0x1b, I_SLO, "SLO", 3, AMAbsoluteY, 7},
        {0x1c, I_TOP, "TOP", 3, AMAbsoluteX, -4},  {0x1d, I_ORA, "ORA", 3, AMAbsoluteX, -4},
        {0x1e, I_ASL, "ASL", 3, AMAbsoluteX, 7},   {0x1f, I_SLO, "SLO", 3, AMAbsoluteX, 7},
        {0x20, I_JSR, "JSR", 3, AMAbsolute, 6},    {0x21, I_AND, "AND", 2, AMIndirectX, 6},
        {0x22, I_KIL, "KIL", 1, AMImplied, 0},     {0x23, I_RLA, "RLA", 2, AMIndirectX, 8},
        {0x24, I_BIT, "BIT", 2, AMZeroPage, 3},    {0x25, I_AND, "AND", 2, AMZeroPage, 3},
        {0x26, I_ROL, "ROL", 2, AMZeroPage, 5},    {0x27, I_RLA, "RLA", 2, AMZeroPage, 5},
        {0x28, I_PLP, "PLP", 1, AMImplied, 4},     {0x29, I_AND, "AND", 2, AMImmediate, 2},
        {0x2a, I_ROL, "ROL", 1, AMAccumulator, 2}, {0x2b, I_AAC, "AAC", 2, AMImmediate, 2},
        {0x2c, I_BIT, "BIT", 3, AMAbsolute, 4},    {0x2d, I_AND, "AND", 3, AMAbsolute, 4},
        {0x2e, I_ROL, "ROL", 3, AMAbsolute, 6},    {0x2f, I_RLA, "RLA", 3, AMAbsolute, 6},
        {0x30, I_BMI, "BMI", 2, AMRelative, -2},   {0x31, I_AND, "AND", 2, AMIndirectY, 5},
        {0x32, I_KIL, "KIL", 1, AMImplied, 0},     {0x33, I_RLA, "RLA", 2, AMIndirectY, 8},
        {0x34, I_DOP, "DOP", 2, AMZeroPageX, 4},   {0x35, I_AND, "AND", 2, AMZeroPageX, 4},
        {0x36, I_ROL, "ROL", 2, AMZeroPageX, 6},   {0x37, I_RLA, "RLA", 2, AMZeroPageX, 6},
        {0x38, I_SEC, "SEC", 1, AMImplied, 2},     {0x39, I_AND, "AND", 3, AMAbsoluteY, -4},
        {0x3a, I_NP2, "NOP", 1, AMImplied, 2},     {0x3b, I_RLA, "RLA", 3, AMAbsoluteY, 7},
        {0x3c, I_TOP, "TOP", 3, AMAbsoluteX, -4},  {0x3d, I_AND, "AND", 3, AMAbsoluteX, -4},
        {0x3e, I_ROL, "ROL", 3, AMAbsoluteX, 7},   {0x3f, I_RLA, "RLA", 3, AMAbsoluteX, 7},
        {0x40, I_RTI, "RTI", 1, AMImplied, 6},     {0x41, I_EOR, "EOR", 2, AMIndirectX, 6},
        {0x42, I_KIL, "KIL", 1, AMImplied, 0},     {0x43, I_SRE, "SRE", 2, AMIndirectX, 8},
        {0x44, I_DOP, "DOP", 2, AMZeroPage, 3},    {0x45, I_EOR, "EOR", 2, AMZeroPage, 3},
        {0x46, I_LSR, "LSR", 2, AMZeroPage, 5},    {0x47, I_SRE, "SRE", 2, AMZeroPage, 5},
        {0x48, I_PHA, "PHA", 1, AMImplied, 3},     {0x49, I_EOR, "EOR", 2, AMImmediate, 2},
        {0x4a, I_LSR, "LSR", 1, AMAccumulator, 2}, {0x4b, I_ASR, "ASR", 2, AMImmediate, 2},
        {0x4c, I_JMP, "JMP", 3, AMAbsolute, 3},    {0x4d, I_EOR, "EOR", 3, AMAbsolute, 4},
        {0x4e, I_LSR, "LSR", 3, AMAbsolute, 6},    {0x4f, I_SRE, "SRE", 3, AMAbsolute, 6},
        {0x50, I_BVC, "BVC", 2, AMRelative, -2},   {0x51, I_EOR, "EOR", 2, AMIndirectY, -5},
        {0x52, I_KIL, "KIL", 1, AMImplied, 0},     {0x53, I_SRE, "SRE", 2, AMIndirectY, 8},
        {0x54, I_DOP, "DOP", 2, AMZeroPageX, 4},   {0x55, I_EOR, "EOR", 2, AMZeroPageX, 4},
        {0x56, I_LSR, "LSR", 2, AMZeroPageX, 6},   {0x57, I_SRE, "SRE", 2, AMZeroPageX, 6},
        {0x58, I_CLI, "CLI", 1, AMImplied, 2},     {0x59, I_EOR, "EOR", 3, AMAbsoluteY, -4},
        {0x5a, I_NP2, "NOP", 1, AMImplied, 2},     {0x5b, I_SRE, "SRE", 3, AMAbsoluteY, 7},
        {0x5c, I_TOP, "TOP", 3, AMAbsoluteX, -4},  {0x5d, I_EOR, "EOR", 3, AMAbsoluteX, -4},
        {0x5e, I_LSR, "LSR", 3, AMAbsoluteX, 7},   {0x5f, I_SRE, "SRE", 3, AMAbsoluteX, 7},
        {0x60, I_RTS, "RTS", 1, AMImplied, 6},     {0x61, I_ADC, "ADC", 2, AMIndirectX, 6},
        {0x62, I_KIL, "KIL", 1, AMImplied, 0},     {0x63, I_RRA, "RRA", 2, AMIndirectX, 8},
        {0x64, I_DOP, "DOP", 2, AMZeroPage, 3},    {0x65, I_ADC, "ADC", 2, AMZeroPage, 3},
        {0x66, I_ROR, "ROR", 2, AMZeroPage, 5},    {0x67, I_RRA, "RRA", 2, AMZeroPage, 5},
        {0x68, I_PLA, "PLA", 1, AMImplied, 4},     {0x69, I_ADC, "ADC", 2, AMImmediate, 2},
        {0x6a, I_ROR, "ROR", 1, AMAccumulator, 2}, {0x6b, I_AAR, "AAR", 2, AMImmediate, 2},
        {0x6c, I_JMP, "JMP", 3, AMIndirect, 5},    {0x6d, I_ADC, "ADC", 3, AMAbsolute, 4},
        {0x6e, I_ROR, "ROR", 3, AMAbsolute, 6},    {0x6f, I_RRA, "RRA", 3, AMAbsolute, 6},
        {0x70, I_BVS, "BVS", 2, AMRelative, -2},   {0x71, I_ADC, "ADC", 2, AMIndirectY, -5},
        {0x72, I_KIL, "KIL", 1, AMImplied, 0},     {0x73, I_RRA, "RRA", 2, AMIndirectY, 8},
        {0x74, I_DOP, "DOP", 2, AMZeroPageX, 4},   {0x75, I_ADC, "ADC", 2, AMZeroPageX, 4},
        {0x76, I_ROR, "ROR", 2, AMZeroPageX, 6},   {0x77, I_RRA, "RRA", 2, AMZeroPageX, 6},
        {0x78, I_SEI, "SEI", 1, AMImplied, 2},     {0x79, I_ADC, "ADC", 3, AMAbsoluteY, -4},
        {0x7a, I_NP2, "NOP", 1, AMImplied, 2},     {0x7b, I_RRA, "RRA", 3, AMAbsoluteY, 7},
        {0x7c, I_TOP, "TOP", 3, AMAbsoluteX, -4},  {0x7d, I_ADC, "ADC", 3, AMAbsoluteX, -4},
        {0x7e, I_ROR, "ROR", 3, AMAbsoluteX, 7},   {0x7f, I_RRA, "RRA", 3, AMAbsoluteX, 7},
        {0x80, I_DOP, "DOP", 2, AMImmediate, 2},   {0x81, I_STA, "STA", 2, AMIndirectX, 6},
        {0x82, I_DOP, "DOP", 2, AMImmediate, 2},   {0x83, I_AAX, "AAX", 2, AMIndirectX, 6},
        {0x84, I_STY, "STY", 2, AMZeroPage, 3},    {0x85, I_STA, "STA", 2, AMZeroPage, 3},
        {0x86, I_STX, "STX", 2, AMZeroPage, 3},    {0x87, I_AAX, "AAX", 2, AMZeroPage, 3},
        {0x88, I_DEY, "DEY", 1, AMImplied, 2},     {0x89, I_DOP, "DOP", 2, AMImmediate, 2},
        {0x8a, I_TXA, "TXA", 1, AMImplied, 2},     {0x8b, I_XAA, "XAA", 2, AMImmediate, 2},
        {0x8c, I_STY, "STY", 3, AMAbsolute, 4},    {0x8d, I_STA, "STA", 3, AMAbsolute, 4},
        {0x8e, I_STX, "STX", 3, AMAbsolute, 4},    {0x8f, I_AAX, "AAX", 3, AMAbsolute, 4},
        {0x90, I_BCC, "BCC", 2, AMRelative, -2},   {0x91, I_STA, "STA", 2, AMIndirectY, 6},
        {0x92, I_KIL, "KIL", 1, AMImplied, 0},     {0x93, I_AXA, "AXA", 2, AMIndirectY, 6},
        {0x94, I_STY, "STY", 2, AMZeroPageX, 4},   {0x95, I_STA, "STA", 2, AMZeroPageX, 4},
        {0x96, I_STX, "STX", 2, AMZeroPageY, 4},   {0x97, I_AAX, "AAX", 2, AMZeroPageY, 4},
        {0x98, I_TYA, "TYA", 1, AMImplied, 2},     {0x99, I_STA, "STA", 3, AMAbsoluteY, 5},
        {0x9a, I_TXS, "TXS", 1, AMImplied, 2},     {0x9b, I_XAS, "XAS", 3, AMAbsoluteY, 5},
        {0x9c, I_SYA, "SYA", 3, AMAbsoluteX, 5},   {0x9d, I_STA, "STA", 3, AMAbsoluteX, 5},
        {0x9e, I_SXA, "SXA", 3, AMAbsoluteY, 5},   {0x9f, I_AXA, "AXA", 3, AMAbsoluteY, 5},
        {0xa0, I_LDY, "LDY", 2, AMImmediate, 2},   {0xa1, I_LDA, "LDA", 2, AMIndirectX, 6},
        {0xa2, I_LDX, "LDX", 2, AMImmediate, 2},   {0xa3, I_LAX, "LAX", 2, AMIndirectX, 6},
        {0xa4, I_LDY, "LDY", 2, AMZeroPage, 3},    {0xa5, I_LDA, "LDA", 2, AMZeroPage, 3},
        {0xa6, I_LDX, "LDX", 2, AMZeroPage, 3},    {0xa7, I_LAX, "LAX", 2, AMZeroPage, 3},
        {0xa8, I_TAY, "TAY", 1, AMImplied, 2},     {0xa9, I_LDA, "LDA", 2, AMImmediate, 2},
        {0xaa, I_TAX, "TAX", 1, AMImplied, 2},     {0xab, I_ATX, "ATX", 2, AMImmediate, 2},
        {0xac, I_LDY, "LDY", 3, AMAbsolute, 4},    {0xad, I_LDA, "LDA", 3, AMAbsolute, 4},
        {0xae, I_LDX, "LDX", 3, AMAbsolute, 4},    {0xaf, I_LAX, "LAX", 3, AMAbsolute, 4},
        {0xb0, I_BCS, "BCS", 2, AMRelative, -2},   {0xb1, I_LDA, "LDA", 2, AMIndirectY, -5},
        {0xb2, I_KIL, "KIL", 1, AMImplied, 0},     {0xb3, I_LAX, "LAX", 2, AMIndirectY, -5},
        {0xb4, I_LDY, "LDY", 2, AMZeroPageX, 4},   {0xb5, I_LDA, "LDA", 2, AMZeroPageX, 4},
        {0xb6, I_LDX, "LDX", 2, AMZeroPageY, 4},   {0xb7, I_LAX, "LAX", 2, AMZeroPageY, 4},
        {0xb8, I_CLV, "CLV", 1, AMImplied, 2},     {0xb9, I_LDA, "LDA", 3, AMAbsoluteY, -4},
        {0xba, I_TSX, "TSX", 1, AMImplied, 2},     {0xbb, I_LAR, "LAR", 3, AMAbsoluteY, -4},
        {0xbc, I_LDY, "LDY", 3, AMAbsoluteX, -4},  {0xbd, I_LDA, "LDA", 3, AMAbsoluteX, -4},
        {0xbe, I_LDX, "LDX", 3, AMAbsoluteY, -4},  {0xbf, I_LAX, "LAX", 3, AMAbsoluteY, -4},
        {0xc0, I_CPY, "CPY", 2, AMImmediate, 2},   {0xc1, I_CMP, "CMP", 2, AMIndirectX, 6},
        {0xc2, I_DOP, "DOP", 2, AMImmediate, 2},   {0xc3, I_DCP, "DCP", 2, AMIndirectX, 8},
        {0xc4, I_CPY, "CPY", 2, AMZeroPage, 3},    {0xc5, I_CMP, "CMP", 2, AMZeroPage, 3},
        {0xc6, I_DEC, "DEC", 2, AMZeroPage, 5},    {0xc7, I_DCP, "DCP", 2, AMZeroPage, 5},
        {0xc8, I_INY, "INY", 1, AMImplied, 2},     {0xc9, I_CMP, "CMP", 2, AMImmediate, 2},
        {0xca, I_DEX, "DEX", 1, AMImplied, 2},     {0xcb, I_AXS, "AXS", 2, AMImmediate, 2},
        {0xcc, I_CPY, "CPY", 3, AMAbsolute, 4},    {0xcd, I_CMP, "CMP", 3, AMAbsolute, 4},
        {0xce, I_DEC, "DEC", 3, AMAbsolute, 6},    {0xcf, I_DCP, "DCP", 3, AMAbsolute, 6},
        {0xd0, I_BNE, "BNE", 2, AMRelative, -2},   {0xd1, I_CMP, "CMP", 2, AMIndirectY, -5},
        {0xd2, I_KIL, "KIL", 1, AMImplied, 0},     {0xd3, I_DCP, "DCP", 2, AMIndirectY, 8},
        {0xd4, I_DOP, "DOP", 2, AMZeroPageX, 4},   {0xd5, I_CMP, "CMP", 2, AMZeroPageX, 4},
        {0xd6, I_DEC, "DEC", 2, AMZeroPageX, 6},   {0xd7, I_DCP, "DCP", 2, AMZeroPageX, 6},
        {0xd8, I_CLD, "CLD", 1, AMImplied, 2},     {0xd9, I_CMP, "CMP", 3, AMAbsoluteY, -4},
        {0xda, I_NP2, "NOP", 1, AMImplied, 2},     {0xdb, I_DCP, "DCP", 3, AMAbsoluteY, 7},
        {0xdc, I_TOP, "TOP", 3, AMAbsoluteX, -4},  {0xdd, I_CMP, "CMP", 3, AMAbsoluteX, -4},
        {0xde, I_DEC, "DEC", 3, AMAbsoluteX, 7},   {0xdf, I_DCP, "DCP", 3, AMAbsoluteX, 7},
        {0xe0, I_CPX, "CPX", 2, AMImmediate, 2},   {0xe1, I_SBC, "SBC", 2, AMIndirectX, 6},
        {0xe2, I_DOP, "DOP", 2, AMImmediate, 2},   {0xe3, I_ISC, "ISC", 2, AMIndirectX, 8},
        {0xe4, I_CPX, "CPX", 2, AMZeroPage, 3},    {0xe5, I_SBC, "SBC", 2, AMZeroPage, 3},
        {0xe6, I_INC, "INC", 2, AMZeroPage, 5},    {0xe7, I_ISC, "ISC", 2, AMZeroPage, 5},
        {0xe8, I_INX, "INX", 1, AMImplied, 2},     {0xe9, I_SBC, "SBC", 2, AMImmediate, 2},
        {0xea, I_NOP, "NOP", 1, AMImplied, 2},     {0xeb, I_SB2, "SBC", 2, AMImmediate, 2},
        {0xec, I_CPX, "CPX", 3, AMAbsolute, 4},    {0xed, I_SBC, "SBC", 3, AMAbsolute, 4},
        {0xee, I_INC, "INC", 3, AMAbsolute, 6},    {0xef, I_ISC, "ISC", 3, AMAbsolute, 6},
        {0xf0, I_BEQ, "BEQ", 2, AMRelative, -2},   {0xf1, I_SBC, "SBC", 2, AMIndirectY, -5},
        {0xf2, I_KIL, "KIL", 1, AMImplied, 0},     {0xf3, I_ISC, "ISC", 2, AMIndirectY, 8},
        {0xf4, I_DOP, "DOP", 2, AMZeroPageX, 4},   {0xf5, I_SBC, "SBC", 2, AMZeroPageX, 4},
        {0xf6, I_INC, "INC", 2, AMZeroPageX, 6},   {0xf7, I_ISC, "ISC", 2, AMZeroPageX, 6},
        {0xf8, I_SED, "SED", 1, AMImplied, 2},     {0xf9, I_SBC, "SBC", 3, AMAbsoluteY, -4},
        {0xfa, I_NP2, "NOP", 1, AMImplied, 2},     {0xfb, I_ISC, "ISC", 3, AMAbsoluteY, 7},
        {0xfc, I_TOP, "TOP", 3, AMAbsoluteX, -4},  {0xfd, I_SBC, "SBC", 3, AMAbsoluteX, -4},
        {0xfe, I_INC, "INC", 3, AMAbsoluteX, 7},   {0xff, I_ISC, "ISC", 3, AMAbsoluteX, 7}}
)

func ReadOpcode(memory MemoryAccess, address uint16) (OpcodeSpec, uint16) {
    op := opcodes[memory.ReadByte(address)]

    var value uint16 = 0
    if op.Size == 2 {
        value = uint16(memory.ReadByte(address+1))
    } else if op.Size == 3 {
        value = uint16(memory.ReadByte(address+1)) | (uint16(memory.ReadByte(address+2)) << 8)
    }
    // switch op.AddressingMode {
    // case AMRelative:
    //     value = address + uint16(op.Size) + uint16(int8(value))
    // }
    return op, value
}

func (op OpcodeSpec) FormatArguments(value uint16, address uint16) string {
    var arguments string = ""
    switch op.AddressingMode {
    case AMAccumulator:
        arguments = "A"
    case AMImmediate:
        arguments = fmt.Sprintf("#$%02X", value)
    case AMIndirect:
        arguments = fmt.Sprintf("($%04X)", value)
    case AMIndirectX:
        arguments = fmt.Sprintf("($%02X,X)", value)
    case AMIndirectY:
        arguments = fmt.Sprintf("($%02X),Y", value)
    case AMRelative:
        arguments = fmt.Sprintf("$%02X", address + uint16(int8(value)))
    case AMZeroPage:
        arguments = fmt.Sprintf("$%02X", value)
    case AMZeroPageX:
        arguments = fmt.Sprintf("$%02X,X", value)
    case AMZeroPageY:
        arguments = fmt.Sprintf("$%02X,Y", value)
    case AMAbsolute:
        arguments = fmt.Sprintf("$%04X", value)
    case AMAbsoluteX:
        arguments = fmt.Sprintf("$%04X,X", value)
    case AMAbsoluteY:
        arguments = fmt.Sprintf("$%04X,Y", value)
    }
    return arguments
}

func (op OpcodeSpec) String() string {
    return fmt.Sprintf("%s %d", op.InstructionName, op.AddressingMode)
}
package cpu6502

import (
    "fmt"
)

const (
    // Addressing Modes
    AMImmediate int  = iota
    AMZeroPage
    AMZeroPageX
    AMZeroPageY
    AMAbsolute
    AMAbsoluteX
    AMAbsoluteY
    AMIndirect
    AMIndirectX
    AMIndirectY
    AMAccumulator
    AMRelative
    AMImplied
)

type InstructionSpec struct {
    Num int
    Name string
    Read bool
    Write bool
}

type OpcodeSpec struct {
    Opcode int
    Instruction InstructionSpec
    Size int
    AddressingMode int
    Cycles int
}

var (
    // Instructions
    I_ADC = InstructionSpec{1, "ADC", true, false}
    I_AND = InstructionSpec{2, "AND", true, false}
    I_ASL = InstructionSpec{3, "ASL", true, true}
    I_BCC = InstructionSpec{4, "BCC", false, false}
    I_BCS = InstructionSpec{5, "BCS", false, false}
    I_BEQ = InstructionSpec{6, "BEQ", false, false}
    I_BIT = InstructionSpec{7, "BIT", true, false}
    I_BMI = InstructionSpec{8, "BMI", false, false}
    I_BNE = InstructionSpec{9, "BNE", false, false}
    I_BPL = InstructionSpec{10, "BPL", false, false}
    I_BRK = InstructionSpec{11, "BRK", false, false}
    I_BVC = InstructionSpec{12, "BVC", false, false}
    I_BVS = InstructionSpec{13, "BVS", false, false}
    I_CLC = InstructionSpec{14, "CLC", false, false}
    I_CLD = InstructionSpec{15, "CLD", false, false}
    I_CLI = InstructionSpec{16, "CLI", false, false}
    I_CLV = InstructionSpec{17, "CLV", false, false}
    I_CMP = InstructionSpec{18, "CMP", true, false}
    I_CPX = InstructionSpec{19, "CPX", true, false}
    I_CPY = InstructionSpec{20, "CPY", true, false}
    I_DEC = InstructionSpec{21, "DEC", true, true}
    I_DEX = InstructionSpec{22, "DEX", false, false}
    I_DEY = InstructionSpec{23, "DEY", false, false}
    I_EOR = InstructionSpec{24, "EOR", true, false}
    I_INC = InstructionSpec{25, "INC", true, true}
    I_INX = InstructionSpec{26, "INX", false, false}
    I_INY = InstructionSpec{27, "INY", false, false}
    I_JMP = InstructionSpec{28, "JMP", false, false}
    I_JSR = InstructionSpec{29, "JSR", false, false}
    I_LDA = InstructionSpec{30, "LDA", true, false}
    I_LDX = InstructionSpec{31, "LDX", true, false}
    I_LDY = InstructionSpec{32, "LDY", true, false}
    I_LSR = InstructionSpec{33, "LSR", true, true}
    I_NOP = InstructionSpec{34, "NOP", false, false}
    I_ORA = InstructionSpec{35, "ORA", true, false}
    I_PHA = InstructionSpec{36, "PHA", false, false}
    I_PHP = InstructionSpec{37, "PHP", false, false}
    I_PLA = InstructionSpec{38, "PLA", false, false}
    I_PLP = InstructionSpec{39, "PLP", false, false}
    I_ROL = InstructionSpec{40, "ROL", true, true}
    I_ROR = InstructionSpec{41, "ROR", true, true}
    I_RTI = InstructionSpec{42, "RTI", false, false}
    I_RTS = InstructionSpec{43, "RTS", false, false}
    I_SBC = InstructionSpec{44, "SBC", true, false}
    I_SEC = InstructionSpec{45, "SEC", false, false}
    I_SED = InstructionSpec{46, "SED", false, false}
    I_SEI = InstructionSpec{47, "SEI", false, false}
    I_STA = InstructionSpec{48, "STA", false, true}
    I_STX = InstructionSpec{49, "STX", false, true}
    I_STY = InstructionSpec{50, "STY", false, true}
    I_TAX = InstructionSpec{51, "TAX", false, false}
    I_TAY = InstructionSpec{52, "TAY", false, false}
    I_TSX = InstructionSpec{53, "TSX", false, false}
    I_TXA = InstructionSpec{54, "TXA", false, false}
    I_TXS = InstructionSpec{55, "TXS", false, false}
    I_TYA = InstructionSpec{56, "TYA", false, false}
    // Undocumented/invalid opcodes
    I_KIL = InstructionSpec{57, "KIL", false, false} // Stop program counter (processor lock up)
    I_DOP = InstructionSpec{58, "DOP", false, false} // double NOP
    I_SLO = InstructionSpec{59, "SLO", true, true} // Shift left one bit in memory, then OR accumulator with memory. Status flags: N, Z, C
    I_AAC = InstructionSpec{60, "AAC", false, false} // (ANC) AND byte with accumulator. If result is negative then carry is set. Status flags: N, Z, C
    I_TOP = InstructionSpec{61, "TOP", false, false} // triple NOP
    I_NP2 = InstructionSpec{62, "NP2", false, false} // NOP - undocumented
    I_RLA = InstructionSpec{63, "RLA", true, true} // rotate one bit left in memory, then AND accumulator with memory. Status flags: N, Z, C
    I_SRE = InstructionSpec{64, "SRE", true, true} // Shift right one bit in memory, then EOR accumulator with memory. Status flags: N,Z,C
    I_ASR = InstructionSpec{65, "ASR", false, false} // [ALR] AND byte with accumulator, then shift right one bit in accumulator. Status flags: N,Z,C
    I_RRA = InstructionSpec{66, "RRA", true, true} // Rotate one bit right in memory, then add memory to accumulator (with carry). Status flags: N,V,Z,C
    I_ARR = InstructionSpec{67, "ARR", false, false} // AND byte with accumulator, then rotate one bit right in accumulator and check bit 5 and 6:
               // If both bits are 1: set C, clear V.
               // If both bits are 0: clear C and V.
               // If only bit 5 is 1: set V, clear C.
               // If only bit 6 is 1: set C and V.
               // Status flags: N,V,Z,C
    I_AAX = InstructionSpec{68, "AAX", false, true} // (SAX) [AXS] AND X register with accumulator and store result in memory. Status flags: N,Z
    I_XAA = InstructionSpec{69, "XAA", false, false} // (ANE) Exact operation unknown. Read the referenced documents for more information and observations.
    I_AXA = InstructionSpec{70, "AXA", false, true} // (SHA) AND X register with accumulator then AND result with 7 and store in memory. Status flags: -
    I_XAS = InstructionSpec{71, "XAS", true, false} // (SHS) [TAS] AND X register with accumulator and store result in stack
               // pointer, then AND stack pointer with the high byte of the
               // target address of the argument + 1. Store result in memory.
               // S = X AND A, M = S AND HIGH(arg) + 1
               // Status flags: -
    I_SYA = InstructionSpec{72, "SYA", false, true} // (SHY) [SAY] AND Y register with the high byte of the target address of the argument + 1. Store the result in memory.
               // M = Y AND HIGH(arg) + 1
               // Status flags: -
    I_SXA = InstructionSpec{73, "SXA", false, true} // (SHX) [XAS] AND X register with the high byte of the target address of the argument + 1. Store the result in memory.
               // M = X AND HIGH(arg) + 1
               // Status flags: -
    I_LAX = InstructionSpec{74, "LAX", true, false} // Load accumulator and X register with memory. Status flags: N,Z
    I_ATX = InstructionSpec{75, "ATX", false, false} // (LXA) [OAL] AND byte with accumulator, then transfer accumulator to X register. Status flags: N,Z
    I_LAR = InstructionSpec{76, "LAR", true, false} // (LAE) [LAS] AND memory with stack pointer, transfer result to accumulator, X register and stack pointer. Status flags: N,Z
    I_DCP = InstructionSpec{77, "DCP", true, true} // [DCM] Subtract 1 from memory (without borrow). Status flags: C
    I_AXS = InstructionSpec{78, "AXS", false, false} // (SBX) [SAX] AND X register with accumulator and store result in X register, then subtract byte from X register (without borrow). Status flags: N,Z,C
    I_ISC = InstructionSpec{79, "ISC", true, true} // (ISB) [INS] Increase memory by one, then subtract memory from accumulator (with borrow). Status flags: N,V,Z,C
    I_SB2 = InstructionSpec{I_SBC.Num, "SB2", true, false} // Same as legal opcode $E9 (SBC #byte)
)

var (
    opcodes = [256]OpcodeSpec{
        // BRK is actually a 2 byte instruction
        {0x00, I_BRK, 1, AMImplied, 7},     {0x01, I_ORA, 2, AMIndirectX, 6},
        {0x02, I_KIL, 1, AMImplied, 0},     {0x03, I_SLO, 2, AMIndirectX, 8},
        {0x04, I_DOP, 2, AMZeroPage, 3},    {0x05, I_ORA, 2, AMZeroPage, 3},
        {0x06, I_ASL, 2, AMZeroPage, 5},    {0x07, I_SLO, 2, AMZeroPage, 5},
        {0x08, I_PHP, 1, AMImplied, 3},     {0x09, I_ORA, 2, AMImmediate, 2},
        {0x0a, I_ASL, 1, AMAccumulator, 2}, {0x0b, I_AAC, 2, AMImmediate, 2},
        {0x0c, I_TOP, 3, AMAbsolute, 4},    {0x0d, I_ORA, 3, AMAbsolute, 4},
        {0x0e, I_ASL, 3, AMAbsolute, 6},    {0x0f, I_SLO, 3, AMAbsolute, 6},
        {0x10, I_BPL, 2, AMRelative, -2},   {0x11, I_ORA, 2, AMIndirectY, 5},
        {0x12, I_KIL, 1, AMImplied, 0},     {0x13, I_SLO, 2, AMIndirectY, 8},
        {0x14, I_DOP, 2, AMZeroPageX, 4},   {0x15, I_ORA, 2, AMZeroPageX, 4},
        {0x16, I_ASL, 2, AMZeroPageX, 6},   {0x17, I_SLO, 2, AMZeroPageX, 6},
        {0x18, I_CLC, 1, AMImplied, 2},     {0x19, I_ORA, 3, AMAbsoluteY, -4},
        {0x1a, I_NP2, 1, AMImplied, 2},     {0x1b, I_SLO, 3, AMAbsoluteY, 7},
        {0x1c, I_TOP, 3, AMAbsoluteX, -4},  {0x1d, I_ORA, 3, AMAbsoluteX, -4},
        {0x1e, I_ASL, 3, AMAbsoluteX, 7},   {0x1f, I_SLO, 3, AMAbsoluteX, 7},
        {0x20, I_JSR, 3, AMAbsolute, 6},    {0x21, I_AND, 2, AMIndirectX, 6},
        {0x22, I_KIL, 1, AMImplied, 0},     {0x23, I_RLA, 2, AMIndirectX, 8},
        {0x24, I_BIT, 2, AMZeroPage, 3},    {0x25, I_AND, 2, AMZeroPage, 3},
        {0x26, I_ROL, 2, AMZeroPage, 5},    {0x27, I_RLA, 2, AMZeroPage, 5},
        {0x28, I_PLP, 1, AMImplied, 4},     {0x29, I_AND, 2, AMImmediate, 2},
        {0x2a, I_ROL, 1, AMAccumulator, 2}, {0x2b, I_AAC, 2, AMImmediate, 2},
        {0x2c, I_BIT, 3, AMAbsolute, 4},    {0x2d, I_AND, 3, AMAbsolute, 4},
        {0x2e, I_ROL, 3, AMAbsolute, 6},    {0x2f, I_RLA, 3, AMAbsolute, 6},
        {0x30, I_BMI, 2, AMRelative, -2},   {0x31, I_AND, 2, AMIndirectY, 5},
        {0x32, I_KIL, 1, AMImplied, 0},     {0x33, I_RLA, 2, AMIndirectY, 8},
        {0x34, I_DOP, 2, AMZeroPageX, 4},   {0x35, I_AND, 2, AMZeroPageX, 4},
        {0x36, I_ROL, 2, AMZeroPageX, 6},   {0x37, I_RLA, 2, AMZeroPageX, 6},
        {0x38, I_SEC, 1, AMImplied, 2},     {0x39, I_AND, 3, AMAbsoluteY, -4},
        {0x3a, I_NP2, 1, AMImplied, 2},     {0x3b, I_RLA, 3, AMAbsoluteY, 7},
        {0x3c, I_TOP, 3, AMAbsoluteX, -4},  {0x3d, I_AND, 3, AMAbsoluteX, -4},
        {0x3e, I_ROL, 3, AMAbsoluteX, 7},   {0x3f, I_RLA, 3, AMAbsoluteX, 7},
        {0x40, I_RTI, 1, AMImplied, 6},     {0x41, I_EOR, 2, AMIndirectX, 6},
        {0x42, I_KIL, 1, AMImplied, 0},     {0x43, I_SRE, 2, AMIndirectX, 8},
        {0x44, I_DOP, 2, AMZeroPage, 3},    {0x45, I_EOR, 2, AMZeroPage, 3},
        {0x46, I_LSR, 2, AMZeroPage, 5},    {0x47, I_SRE, 2, AMZeroPage, 5},
        {0x48, I_PHA, 1, AMImplied, 3},     {0x49, I_EOR, 2, AMImmediate, 2},
        {0x4a, I_LSR, 1, AMAccumulator, 2}, {0x4b, I_ASR, 2, AMImmediate, 2},
        {0x4c, I_JMP, 3, AMAbsolute, 3},    {0x4d, I_EOR, 3, AMAbsolute, 4},
        {0x4e, I_LSR, 3, AMAbsolute, 6},    {0x4f, I_SRE, 3, AMAbsolute, 6},
        {0x50, I_BVC, 2, AMRelative, -2},   {0x51, I_EOR, 2, AMIndirectY, -5},
        {0x52, I_KIL, 1, AMImplied, 0},     {0x53, I_SRE, 2, AMIndirectY, 8},
        {0x54, I_DOP, 2, AMZeroPageX, 4},   {0x55, I_EOR, 2, AMZeroPageX, 4},
        {0x56, I_LSR, 2, AMZeroPageX, 6},   {0x57, I_SRE, 2, AMZeroPageX, 6},
        {0x58, I_CLI, 1, AMImplied, 2},     {0x59, I_EOR, 3, AMAbsoluteY, -4},
        {0x5a, I_NP2, 1, AMImplied, 2},     {0x5b, I_SRE, 3, AMAbsoluteY, 7},
        {0x5c, I_TOP, 3, AMAbsoluteX, -4},  {0x5d, I_EOR, 3, AMAbsoluteX, -4},
        {0x5e, I_LSR, 3, AMAbsoluteX, 7},   {0x5f, I_SRE, 3, AMAbsoluteX, 7},
        {0x60, I_RTS, 1, AMImplied, 6},     {0x61, I_ADC, 2, AMIndirectX, 6},
        {0x62, I_KIL, 1, AMImplied, 0},     {0x63, I_RRA, 2, AMIndirectX, 8},
        {0x64, I_DOP, 2, AMZeroPage, 3},    {0x65, I_ADC, 2, AMZeroPage, 3},
        {0x66, I_ROR, 2, AMZeroPage, 5},    {0x67, I_RRA, 2, AMZeroPage, 5},
        {0x68, I_PLA, 1, AMImplied, 4},     {0x69, I_ADC, 2, AMImmediate, 2},
        {0x6a, I_ROR, 1, AMAccumulator, 2}, {0x6b, I_ARR, 2, AMImmediate, 2},
        {0x6c, I_JMP, 3, AMIndirect, 5},    {0x6d, I_ADC, 3, AMAbsolute, 4},
        {0x6e, I_ROR, 3, AMAbsolute, 6},    {0x6f, I_RRA, 3, AMAbsolute, 6},
        {0x70, I_BVS, 2, AMRelative, -2},   {0x71, I_ADC, 2, AMIndirectY, -5},
        {0x72, I_KIL, 1, AMImplied, 0},     {0x73, I_RRA, 2, AMIndirectY, 8},
        {0x74, I_DOP, 2, AMZeroPageX, 4},   {0x75, I_ADC, 2, AMZeroPageX, 4},
        {0x76, I_ROR, 2, AMZeroPageX, 6},   {0x77, I_RRA, 2, AMZeroPageX, 6},
        {0x78, I_SEI, 1, AMImplied, 2},     {0x79, I_ADC, 3, AMAbsoluteY, -4},
        {0x7a, I_NP2, 1, AMImplied, 2},     {0x7b, I_RRA, 3, AMAbsoluteY, 7},
        {0x7c, I_TOP, 3, AMAbsoluteX, -4},  {0x7d, I_ADC, 3, AMAbsoluteX, -4},
        {0x7e, I_ROR, 3, AMAbsoluteX, 7},   {0x7f, I_RRA, 3, AMAbsoluteX, 7},
        {0x80, I_DOP, 2, AMImmediate, 2},   {0x81, I_STA, 2, AMIndirectX, 6},
        {0x82, I_DOP, 2, AMImmediate, 2},   {0x83, I_AAX, 2, AMIndirectX, 6},
        {0x84, I_STY, 2, AMZeroPage, 3},    {0x85, I_STA, 2, AMZeroPage, 3},
        {0x86, I_STX, 2, AMZeroPage, 3},    {0x87, I_AAX, 2, AMZeroPage, 3},
        {0x88, I_DEY, 1, AMImplied, 2},     {0x89, I_DOP, 2, AMImmediate, 2},
        {0x8a, I_TXA, 1, AMImplied, 2},     {0x8b, I_XAA, 2, AMImmediate, 2},
        {0x8c, I_STY, 3, AMAbsolute, 4},    {0x8d, I_STA, 3, AMAbsolute, 4},
        {0x8e, I_STX, 3, AMAbsolute, 4},    {0x8f, I_AAX, 3, AMAbsolute, 4},
        {0x90, I_BCC, 2, AMRelative, -2},   {0x91, I_STA, 2, AMIndirectY, 6},
        {0x92, I_KIL, 1, AMImplied, 0},     {0x93, I_AXA, 2, AMIndirectY, 6},
        {0x94, I_STY, 2, AMZeroPageX, 4},   {0x95, I_STA, 2, AMZeroPageX, 4},
        {0x96, I_STX, 2, AMZeroPageY, 4},   {0x97, I_AAX, 2, AMZeroPageY, 4},
        {0x98, I_TYA, 1, AMImplied, 2},     {0x99, I_STA, 3, AMAbsoluteY, 5},
        {0x9a, I_TXS, 1, AMImplied, 2},     {0x9b, I_XAS, 3, AMAbsoluteY, 5},
        {0x9c, I_SYA, 3, AMAbsoluteX, 5},   {0x9d, I_STA, 3, AMAbsoluteX, 5},
        {0x9e, I_SXA, 3, AMAbsoluteY, 5},   {0x9f, I_AXA, 3, AMAbsoluteY, 5},
        {0xa0, I_LDY, 2, AMImmediate, 2},   {0xa1, I_LDA, 2, AMIndirectX, 6},
        {0xa2, I_LDX, 2, AMImmediate, 2},   {0xa3, I_LAX, 2, AMIndirectX, 6},
        {0xa4, I_LDY, 2, AMZeroPage, 3},    {0xa5, I_LDA, 2, AMZeroPage, 3},
        {0xa6, I_LDX, 2, AMZeroPage, 3},    {0xa7, I_LAX, 2, AMZeroPage, 3},
        {0xa8, I_TAY, 1, AMImplied, 2},     {0xa9, I_LDA, 2, AMImmediate, 2},
        {0xaa, I_TAX, 1, AMImplied, 2},     {0xab, I_ATX, 2, AMImmediate, 2},
        {0xac, I_LDY, 3, AMAbsolute, 4},    {0xad, I_LDA, 3, AMAbsolute, 4},
        {0xae, I_LDX, 3, AMAbsolute, 4},    {0xaf, I_LAX, 3, AMAbsolute, 4},
        {0xb0, I_BCS, 2, AMRelative, -2},   {0xb1, I_LDA, 2, AMIndirectY, -5},
        {0xb2, I_KIL, 1, AMImplied, 0},     {0xb3, I_LAX, 2, AMIndirectY, -5},
        {0xb4, I_LDY, 2, AMZeroPageX, 4},   {0xb5, I_LDA, 2, AMZeroPageX, 4},
        {0xb6, I_LDX, 2, AMZeroPageY, 4},   {0xb7, I_LAX, 2, AMZeroPageY, 4},
        {0xb8, I_CLV, 1, AMImplied, 2},     {0xb9, I_LDA, 3, AMAbsoluteY, -4},
        {0xba, I_TSX, 1, AMImplied, 2},     {0xbb, I_LAR, 3, AMAbsoluteY, -4},
        {0xbc, I_LDY, 3, AMAbsoluteX, -4},  {0xbd, I_LDA, 3, AMAbsoluteX, -4},
        {0xbe, I_LDX, 3, AMAbsoluteY, -4},  {0xbf, I_LAX, 3, AMAbsoluteY, -4},
        {0xc0, I_CPY, 2, AMImmediate, 2},   {0xc1, I_CMP, 2, AMIndirectX, 6},
        {0xc2, I_DOP, 2, AMImmediate, 2},   {0xc3, I_DCP, 2, AMIndirectX, 8},
        {0xc4, I_CPY, 2, AMZeroPage, 3},    {0xc5, I_CMP, 2, AMZeroPage, 3},
        {0xc6, I_DEC, 2, AMZeroPage, 5},    {0xc7, I_DCP, 2, AMZeroPage, 5},
        {0xc8, I_INY, 1, AMImplied, 2},     {0xc9, I_CMP, 2, AMImmediate, 2},
        {0xca, I_DEX, 1, AMImplied, 2},     {0xcb, I_AXS, 2, AMImmediate, 2},
        {0xcc, I_CPY, 3, AMAbsolute, 4},    {0xcd, I_CMP, 3, AMAbsolute, 4},
        {0xce, I_DEC, 3, AMAbsolute, 6},    {0xcf, I_DCP, 3, AMAbsolute, 6},
        {0xd0, I_BNE, 2, AMRelative, -2},   {0xd1, I_CMP, 2, AMIndirectY, -5},
        {0xd2, I_KIL, 1, AMImplied, 0},     {0xd3, I_DCP, 2, AMIndirectY, 8},
        {0xd4, I_DOP, 2, AMZeroPageX, 4},   {0xd5, I_CMP, 2, AMZeroPageX, 4},
        {0xd6, I_DEC, 2, AMZeroPageX, 6},   {0xd7, I_DCP, 2, AMZeroPageX, 6},
        {0xd8, I_CLD, 1, AMImplied, 2},     {0xd9, I_CMP, 3, AMAbsoluteY, -4},
        {0xda, I_NP2, 1, AMImplied, 2},     {0xdb, I_DCP, 3, AMAbsoluteY, 7},
        {0xdc, I_TOP, 3, AMAbsoluteX, -4},  {0xdd, I_CMP, 3, AMAbsoluteX, -4},
        {0xde, I_DEC, 3, AMAbsoluteX, 7},   {0xdf, I_DCP, 3, AMAbsoluteX, 7},
        {0xe0, I_CPX, 2, AMImmediate, 2},   {0xe1, I_SBC, 2, AMIndirectX, 6},
        {0xe2, I_DOP, 2, AMImmediate, 2},   {0xe3, I_ISC, 2, AMIndirectX, 8},
        {0xe4, I_CPX, 2, AMZeroPage, 3},    {0xe5, I_SBC, 2, AMZeroPage, 3},
        {0xe6, I_INC, 2, AMZeroPage, 5},    {0xe7, I_ISC, 2, AMZeroPage, 5},
        {0xe8, I_INX, 1, AMImplied, 2},     {0xe9, I_SBC, 2, AMImmediate, 2},
        {0xea, I_NOP, 1, AMImplied, 2},     {0xeb, I_SB2, 2, AMImmediate, 2},
        {0xec, I_CPX, 3, AMAbsolute, 4},    {0xed, I_SBC, 3, AMAbsolute, 4},
        {0xee, I_INC, 3, AMAbsolute, 6},    {0xef, I_ISC, 3, AMAbsolute, 6},
        {0xf0, I_BEQ, 2, AMRelative, -2},   {0xf1, I_SBC, 2, AMIndirectY, -5},
        {0xf2, I_KIL, 1, AMImplied, 0},     {0xf3, I_ISC, 2, AMIndirectY, 8},
        {0xf4, I_DOP, 2, AMZeroPageX, 4},   {0xf5, I_SBC, 2, AMZeroPageX, 4},
        {0xf6, I_INC, 2, AMZeroPageX, 6},   {0xf7, I_ISC, 2, AMZeroPageX, 6},
        {0xf8, I_SED, 1, AMImplied, 2},     {0xf9, I_SBC, 3, AMAbsoluteY, -4},
        {0xfa, I_NP2, 1, AMImplied, 2},     {0xfb, I_ISC, 3, AMAbsoluteY, 7},
        {0xfc, I_TOP, 3, AMAbsoluteX, -4},  {0xfd, I_SBC, 3, AMAbsoluteX, -4},
        {0xfe, I_INC, 3, AMAbsoluteX, 7},   {0xff, I_ISC, 3, AMAbsoluteX, 7}}
)

func ReadOpcode(memory MemoryAccess, address uint16) (OpcodeSpec, uint16) {
    op := opcodes[memory.ReadByte(address, false)]

    var value uint16 = 0
    if op.Size == 2 {
        value = uint16(memory.ReadByte(address+1, false))
    } else if op.Size == 3 {
        value = uint16(memory.ReadByte(address+1, false)) | (uint16(memory.ReadByte(address+2, false)) << 8)
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
    return fmt.Sprintf("%s %d", op.Instruction.Name, op.AddressingMode)
}

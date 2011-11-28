package cpu6502

import (
    "fmt"
    "os"
)

func DisassembleOne(memory MemoryAccess, address uint16) (DecodedOpcode, os.Error) {
    dec, err := ReadOpcode(memory, address)
    if err != nil {
        return dec, err
    }
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
    return dec, err
}

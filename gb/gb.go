package gb

import (
	"github.com/samuel/go-emu/z80"
)

type GBState struct {
	CPU *z80.Z80
	cart *Cart
}

func New(cart *Cart) (*GBState, error) {
	state := &GBState{cart: cart}

	state.CPU = z80.New(state)

	return state, nil
}

func (gb *GBState) Step() {
	gb.CPU.Step()
}

func (gb *GBState) ReadByte(address uint16, peek bool) byte {
	return gb.cart.memory[address]
}

func (gb *GBState) WriteByte(address uint16, value byte) {
	gb.cart.memory[address] = value
}

package nes

import (
	"fmt"
)

const (
	APU_FRAME_CLOCK_DIVIDER = 89490 // frame clock = cpu clock / clock divider = ~240Hz NTSC

	// Register 4017
	BIT_APU_FRAME_RATE        = 0x80 // Frame Rate Select  (0=NTSC=60Hz=240Hz/4, 1=PAL=48Hz=240Hz/5)
	BIT_APU_FRAME_IRQ_DISABLE = 0x40 // Frame IRQ Disable  (0=Enable Frame IRQ, 1=Disable Frame IRQ)
)

type APUState struct {
	FrameIRQEnabled bool
	FrameRate       int // NTSC=4, PAL=5
}

func NewAPUState() (*APUState, error) {
	state := &APUState{
		FrameIRQEnabled: true,
		FrameRate:       4}
	return state, nil
}

func (apu *APUState) Pulse(cycles int) {
}

func (apu *APUState) ReadByte(address uint16, peek bool) byte {
	if address < 0x4000 && address > 0x4017 {
		panic("Invalid APU address")
	}
	return 0
}

func (apu *APUState) WriteByte(address uint16, value byte) {
	if address < 0x4000 && address > 0x4017 {
		panic("Invalid APU address")
	}
	if address == 0x4017 {
		apu.FrameIRQEnabled = value&BIT_APU_FRAME_IRQ_DISABLE == 0
		if value&BIT_APU_FRAME_RATE == 0 {
			apu.FrameRate = 4
		} else {
			apu.FrameRate = 5
		}
	}
}

func (apu *APUState) String() string {
	return fmt.Sprintf("{FrameIRQEnabled:%s FrameRate:%d}", apu.FrameIRQEnabled, apu.FrameRate)
}

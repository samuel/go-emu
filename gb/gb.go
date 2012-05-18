package gb

import (
	"github.com/samuel/go-nes/z80"
)

type GBState struct {
	CPU *z80.Z80
}
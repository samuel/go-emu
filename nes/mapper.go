package nes

import "errors"

const (
	MAPPER_NROM = 0
	MAPPER_MMC1 = 1
	MAPPER_MMC3 = 4
)

type Mapper interface {
	ReadByte(address uint16, peek bool) byte
	WriteByte(address uint16, value byte)
}

func NewMapper(cart *Cart) (Mapper, error) {
	switch cart.Mapper {
	case MAPPER_NROM:
		return NewMapperNROM(cart), nil
	case MAPPER_MMC1:
		return NewMapperMMC1(cart), nil
	case MAPPER_MMC3:
		return NewMapperMMC3(cart), nil
	}
	return nil, errors.New("mapper not implemented")
}

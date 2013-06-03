package gb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

// http://gbdev.gg8.se/wiki/articles/The_Cartridge_Header
type Cart struct {
	Title   string
	CGB     bool
	CGBOnly bool

	memory []byte
}

func (cart *Cart) String() string {
	cgb := "no"
	if cart.CGBOnly {
		cgb = "only"
	} else if cart.CGB {
		cgb = "yes"
	}
	return fmt.Sprintf("{Title:%s CGB:%s Manufacturer:%.4x MemoryLen:%d}",
		cart.Title, cgb, cart.ManufacturerCode(), len(cart.memory))
}

func (cart *Cart) Logo() []byte {
	return cart.memory[0x104:0x134]
}

func (cart *Cart) ManufacturerCode() []byte {
	return cart.memory[0x13f:0x143]
}

func LoadCart(r io.Reader) (*Cart, error) {
	cart := &Cart{
		memory: make([]byte, 0, 32*1024),
	}

	// var b [32*1024]byte
	b := make([]byte, 32*1024)

	for {
		n, err := r.Read(b[:32*1024])
		cart.memory = append(cart.memory, b[:n]...)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	checksum := byte(0)
	for i := 0x134; i <= 0x14c; i++ {
		checksum = checksum - cart.memory[i] - 1
	}
	if checksum != cart.memory[0x14d] {
		return nil, errors.New("header checksum failed")
	}

	i := bytes.IndexByte(cart.memory[0x134:0x144], 0)
	cart.Title = string(cart.memory[0x134 : 0x134+i])
	cart.CGB = cart.memory[0x143] == 0x80 || cart.memory[0x143] == 0xc0
	cart.CGBOnly = cart.memory[0x143] == 0xc0

	return cart, nil
}

func LoadCartFile(filename string) (*Cart, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadCart(file)
}

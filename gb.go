package main

import (
	"flag"
	"fmt"
	"log"

	// "github.com/samuel/go-emu/z80"
	"github.com/samuel/go-emu/gb"
)

var (
	f_trace = flag.Bool("t", false, "print trace while running")
	f_rom   = flag.String("r", "", "ROM file")
)

func parseFlags() {
	flag.Parse()
	if *f_rom == "" {
		log.Fatal("ROM is required (-r)")
	}
}

func main() {
	parseFlags()
	cart, err := gb.LoadCartFile(*f_rom)
	if err != nil {
		panic(err)
	}

	state, err := gb.New(cart)
	if err != nil {
		panic(err)
	}

	fmt.Println(cart)
	fmt.Printf("%+v\n", state)

	for i := 0; i < 10000; i++ {
		state.Step()
	}
}

package main

import (
	"log"
	"os"
)

const startAddress = 0x200
const startFontAddress = 0x050
const maxMemorySize = 4096

type Chip8 struct {
	memory     [maxMemorySize]uint8
	v          [16]uint8
	i          uint16
	delayTimer uint8
	soundTimer uint8
	pc         uint16
	sp         uint8
	stack      [16]uint16
	display    [displayWidth * displayHeight]uint8
}

func NewChip8() *Chip8 {
	return &Chip8{
		pc: startAddress,
	}
}

func (c *Chip8) loadRom(romName string) {
	bytes, err := os.ReadFile(romName)
	if err != nil {
		log.Fatalln(err)
	}

	if len(bytes) > maxMemorySize-startAddress {
		log.Fatalln("ROM size exceeds memory limit")
	}

	copy(c.memory[startAddress:], bytes)
}

func (c *Chip8) fetch() uint16 {
	return uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
}

func (c *Chip8) execute(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := opcode & 0x000F
	nn := opcode & 0x00FF
	nnn := opcode & 0x0FFF

	switch opcode & 0xF000 {
	case 0x0000:
		switch nn {
		case 0xE0:
			c.display = [displayWidth * displayHeight]uint8{}
		case 0xEE:
			if c.sp == 0 {
				log.Fatalln("stack underflow")
			}
			c.sp -= 1
			c.pc = c.stack[c.sp]
		default:
			log.Fatalf("Unknown opcode: 0x%04X\n", opcode)
		}
	case 0x1000:
		c.pc = nnn
	case 0x2000:
		if c.sp >= uint8(len(c.stack)) {
			log.Fatalln("stack overflow")
		}
		c.stack[c.sp] = c.pc
		c.sp += 1
		c.pc = nnn
	case 0x3000:
	case 0x4000:
	case 0x5000:
	case 0x6000:
	case 0x7000:
	case 0x8000:
	case 0x9000:
	case 0xA000:
	case 0xB000:
	case 0xC000:
	case 0xD000:
	case 0xE000:
	case 0xF000:
	default:
		log.Fatalf("Unknown opcode: 0x%04X\n", opcode)
	}
}

func (c *Chip8) setTimer() {
	if c.delayTimer > 0 {
		c.delayTimer -= 1
	}
	if c.soundTimer > 0 {
		c.soundTimer -= 1
	}
}

func (c *Chip8) cycle() {
	// chip8: 500Hz
	for i := 0; i < 8; i += 1 {
		opcode := c.fetch()
		c.pc += 2
		c.execute(opcode)
	}
}

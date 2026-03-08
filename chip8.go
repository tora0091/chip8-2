package main

import (
	"log"
	"math/rand"
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
		if c.v[x] == uint8(nn) {
			c.pc += 2
		}
	case 0x4000:
		if c.v[x] != uint8(nn) {
			c.pc += 2
		}
	case 0x5000:
		if n == 0 {
			if c.v[x] == c.v[y] {
				c.pc += 2
			}
		} else {
			log.Fatalf("Unknown opcode: 0x%04X\n", opcode)
		}
	case 0x6000:
		c.v[x] = uint8(nn)
	case 0x7000:
		c.v[x] += uint8(nn)
	case 0x8000:
		switch n {
		case 0x0:
			c.v[x] = c.v[y]
		case 0x1:
			c.v[x] |= c.v[y]
		case 0x2:
			c.v[x] &= c.v[y]
		case 0x3:
			c.v[x] ^= c.v[y]
		case 0x4:
			ans := uint16(c.v[x]) + uint16(c.v[y])
			c.v[x] = uint8(ans)
			if ans > 0xFF {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
		case 0x5:
			if c.v[x] >= c.v[y] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] -= c.v[y]
		case 0x6:
			ans := c.v[x] & 0x1
			if ans == 1 {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] >>= 1
		case 0x7:
			if c.v[y] >= c.v[x] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 1
			}
			c.v[x] = c.v[y] - c.v[x]
		case 0xE:
			ans := c.v[x] >> 7
			if ans == 1 {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] <<= 1
		default:
			log.Fatalf("Unknown opcode: 0x%04X\n", opcode)
		}
	case 0x9000:
		if c.v[x] != c.v[y] {
			c.pc += 2
		}
	case 0xA000:
		c.i = nnn
	case 0xB000:
		c.pc = nnn + uint16(c.v[0])
	case 0xC000:
		c.v[x] = uint8(rand.Intn(256)) & uint8(nn)
	case 0xD000:
	case 0xE000:
	case 0xF000:
		switch nn {
		case 0x07:
			c.v[x] = c.delayTimer
		case 0x0A:

		case 0x15:
			c.delayTimer = c.v[x]
		case 0x18:
			c.soundTimer = c.v[x]
		case 0x1E:
			c.i += uint16(c.v[x])
		case 0x29:

		case 0x33:
			value := c.v[x]
			c.memory[c.i] = value / 100
			c.memory[c.i+1] = (value / 10) % 10
			c.memory[c.i+2] = value % 10
		case 0x55:
			for idx := uint16(0); idx < 16; idx += 1 {
				c.memory[c.i+idx] = c.v[idx]
			}
		case 0x65:
			for idx := uint16(0); idx < 16; idx += 1 {
				c.v[idx] = c.memory[c.i+idx]
			}
		default:
			log.Fatalf("Unknown opcode: 0x%04X\n", opcode)
		}
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
	opcode := c.fetch()
	c.pc += 2
	c.execute(opcode)
}

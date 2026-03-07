package main

import (
	"log"
	"os"
	"time"
)

const displayWidth = 64
const displayHeight = 32
const displayScale = 15

func main() {
	romName := getRomName()

	chip8 := NewChip8()
	chip8.loadRom(romName)

	// 60Hz
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()
	for range ticker.C {
		chip8.cycle()
		// timer: 60Hz
		chip8.setTimer()
	}
}

func getRomName() string {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s [ROM name]\n", os.Args[0])
	}

	path := os.Args[1]
	_, err := os.Stat(path)
	if err != nil {
		log.Fatalln(err)
	}

	return path
}

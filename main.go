package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const displayWidth = 64
const displayHeight = 32
const displayScale = 15

type Game struct {
	chip8 *Chip8
}

func (g *Game) Update() error {
	// chip8: 500Hz
	for i := 0; i < 8; i += 1 {
		g.chip8.cycle()
	}

	// timer: 60Hz
	g.chip8.setTimer()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, world")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return displayWidth * displayScale, displayHeight * displayScale
}

func main() {
	romName := getRomName()

	chip8 := NewChip8()
	chip8.loadRom(romName)

	ebiten.SetWindowSize(displayWidth*displayScale, displayHeight*displayScale)
	ebiten.SetWindowTitle("Chip-8 emulator!!!")

	if err := ebiten.RunGame(&Game{chip8: chip8}); err != nil {
		log.Fatalln(err)
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

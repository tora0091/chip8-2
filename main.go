package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const displayWidth = 64
const displayHeight = 32
const displayScale = 15

type Game struct {
	chip8 *Chip8
}

func (g *Game) Update() error {
	g.updateKeypad()

	if g.chip8.waitingForKey {
		for i := 0; i < 16; i += 1 {
			if g.chip8.keypad[i] {
				g.chip8.v[g.chip8.waitingRegister] = uint8(i)
				break
			}
		}
	}

	// chip8: 500Hz
	for i := 0; i < 8; i += 1 {
		g.chip8.cycle()
	}

	// timer: 60Hz
	g.chip8.updateTimer()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for idx, val := range g.chip8.display {
		if val == 0 {
			continue
		}

		x := float32(idx%displayWidth) * displayScale
		y := float32(idx/displayWidth) * displayScale

		vector.FillRect(
			screen,
			x,
			y,
			displayScale,
			displayScale,
			color.White,
			false,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return displayWidth * displayScale, displayHeight * displayScale
}

func (g *Game) updateKeypad() {
	g.chip8.keypad[0x0] = ebiten.IsKeyPressed(ebiten.KeyX)
	g.chip8.keypad[0x1] = ebiten.IsKeyPressed(ebiten.Key1)
	g.chip8.keypad[0x2] = ebiten.IsKeyPressed(ebiten.Key2)
	g.chip8.keypad[0x3] = ebiten.IsKeyPressed(ebiten.Key3)
	g.chip8.keypad[0x4] = ebiten.IsKeyPressed(ebiten.KeyQ)
	g.chip8.keypad[0x5] = ebiten.IsKeyPressed(ebiten.KeyW)
	g.chip8.keypad[0x6] = ebiten.IsKeyPressed(ebiten.KeyE)
	g.chip8.keypad[0x7] = ebiten.IsKeyPressed(ebiten.KeyA)
	g.chip8.keypad[0x8] = ebiten.IsKeyPressed(ebiten.KeyS)
	g.chip8.keypad[0x9] = ebiten.IsKeyPressed(ebiten.KeyD)
	g.chip8.keypad[0xA] = ebiten.IsKeyPressed(ebiten.KeyZ)
	g.chip8.keypad[0xB] = ebiten.IsKeyPressed(ebiten.KeyC)
	g.chip8.keypad[0xC] = ebiten.IsKeyPressed(ebiten.Key4)
	g.chip8.keypad[0xD] = ebiten.IsKeyPressed(ebiten.KeyR)
	g.chip8.keypad[0xE] = ebiten.IsKeyPressed(ebiten.KeyF)
	g.chip8.keypad[0xF] = ebiten.IsKeyPressed(ebiten.KeyV)
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

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"snakeGame/game/core"
	"snakeGame/game/ui"
)

func main() {
	// Define the game board dimensions in grid cells.
	gridWidth, gridHeight := 20, 20       // 20x20 grid cells
	cellSize := 16                        // each cell is 16 pixels square
	screenWidth := gridWidth * cellSize   // internal pixel width
	screenHeight := gridHeight * cellSize // internal pixel height
	start := image.Pt(gridWidth/2, gridHeight/2)

	// Set up the window for desktop. We double the pixel dimensions for a retro-scaled look
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Game")

	// Initialize the Game instance (from our game package).
	g := core.NewGame(start, gridWidth, gridHeight, cellSize)

	if err := ui.LoadTheme(); err != nil {
		log.Fatal(err)
	}

	// Start the game loop. Ebiten will call g.Update, g.Draw, g.Layout appropriately.
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

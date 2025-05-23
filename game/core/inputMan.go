package core

import (
	"github.com/hajimehoshi/ebiten/v2" // Ebiten game engine
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

// handleInput updates the snake direction based on keyboard input.
func (g *Game) handleInput() {
	// Prevent reverse movement
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) && g.Snake.Dir.Y != 1 {
		g.Snake.PendingDir = image.Pt(0, -1)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) && g.Snake.Dir.Y != -1 {
		g.Snake.PendingDir = image.Pt(0, 1)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) && g.Snake.Dir.X != 1 {
		g.Snake.PendingDir = image.Pt(-1, 0)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) && g.Snake.Dir.X != -1 {
		g.Snake.PendingDir = image.Pt(1, 0)
	}
}

func (g *Game) handlePauseToggle() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.State.TogglePause()
	}
}

func (g *Game) handleRetryInput() {
	// Keyboard restart
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.resetGame()
		return
	}

	// Mouse click restart
	mx, my := ebiten.CursorPosition()
	centerX := g.screenWidth / 2
	centerY := g.screenHeight / 2
	x := centerX - 80
	y := centerY
	w, h := 160, 30

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if mx >= x && mx <= x+w && my >= y && my <= y+h {
			g.resetGame()
		}
	}
}

package render

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// UIManager handles overlay rendering like score, pause, and game over prompts.
type UIManager struct {}

// NewUIManager creates a new instance of UIManager.
func NewUIManager() *UIManager {
	return &UIManager{}
}

// DrawStatus draws the in-game score and level.
func (ui *UIManager) DrawStatus(screen *ebiten.Image, score, level int) {
	text := fmt.Sprintf("Score: %d  Level: %d", score, level)
	ebitenutil.DebugPrintAt(screen, text, 5, 5)
}

// DrawPauseOverlay draws a "PAUSED" message in the center of the screen.
func (ui *UIManager) DrawPauseOverlay(screen *ebiten.Image, screenWidth, screenHeight int) {
	text := "PAUSED"
	ebitenutil.DebugPrintAt(screen, text, screenWidth/2-30, screenHeight/2 - 20)
}

// DrawGameOverOverlay draws the game over screen with retry prompt.
func (ui *UIManager) DrawGameOverOverlay(screen *ebiten.Image, screenWidth, screenHeight int) {
	centerX := screenWidth / 2
	centerY := screenHeight / 2

	ebitenutil.DebugPrintAt(screen, "GAME OVER", centerX-40, centerY-40)

	x := centerX - 80
	y := centerY
	w := 160
	h := 30

	// Retry button background
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), color.RGBA{100, 100, 100, 255}, false)

	// Retry text
	ebitenutil.DebugPrintAt(screen, "\u27F3 Play Again (Press Enter)", x+20, y+8)
}


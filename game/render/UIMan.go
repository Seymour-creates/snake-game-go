package render

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// UIManager handles overlay rendering like score, pause, and game over prompts.
type UIManager struct{}

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
	ebitenutil.DebugPrintAt(screen, text, screenWidth/2-30, screenHeight/2-20)
}

// DrawGameOverOverlay draws the game over screen with retry prompt.
func (ui *UIManager) DrawGameOverOverlay(screen *ebiten.Image, screenWidth, screenHeight int) {
	centerX := screenWidth / 2
	centerY := screenHeight / 2

	// Game Over
	ebitenutil.DebugPrintAt(screen, "GAME OVER", centerX-40, centerY-40)

	// Retry Text
	retryMsg := "\u27F3 Play Again (Press Enter)"
	charWidth := 7 // estimated width per char in debug font
	textWidth := len(retryMsg) * charWidth
	x := centerX - textWidth/2
	y := centerY + 8

	ebitenutil.DebugPrintAt(screen, retryMsg, x, y)
}

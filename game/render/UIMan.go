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

// DrawGameOverOverlay draws the game over screen with selectable options.
func (ui *UIManager) DrawGameOverOverlay(screen *ebiten.Image, screenWidth, screenHeight int, selected int) {
	centerX := screenWidth / 2
	centerY := screenHeight / 2

	// Game Over
	ebitenutil.DebugPrintAt(screen, "GAME OVER", centerX-40, centerY-40)

	// Options
	options := []string{"Play Again", "Main Menu", "Exit Game"}
	menuStartY := centerY + 8
	for i, item := range options {
		prefix := "  "
		if i == selected {
			prefix = "> "
		}
		text := prefix + item
		charWidth := 7 // estimated width per char in debug font
		textWidth := len(text) * charWidth
		x := centerX - textWidth/2
		y := menuStartY + i*30
		ebitenutil.DebugPrintAt(screen, text, x, y)
	}
}

// DrawTitleScreen draws the title screen with selectable menu tiles.
func (ui *UIManager) DrawTitleScreen(screen *ebiten.Image, screenWidth, screenHeight int, selected int) {
	title := "SNAKE GAME"
	menu := []string{"Play Game", "Settings"}

	// Draw title centered at top
	titleX := screenWidth/2 - len(title)*7/2
	titleY := screenHeight / 4
	ebitenutil.DebugPrintAt(screen, title, titleX, titleY)

	// Draw menu options
	menuStartY := titleY + 40
	for i, item := range menu {
		prefix := "  "
		if i == selected {
			prefix = "> "
		}
		text := prefix + item
		textX := screenWidth/2 - len(text)*7/2
		textY := menuStartY + i*30
		ebitenutil.DebugPrintAt(screen, text, textX, textY)
	}
}

// DrawSettingsScreen draws a placeholder for the settings screen.
func (ui *UIManager) DrawSettingsScreen(screen *ebiten.Image, screenWidth, screenHeight int) {
	msg := "Settings - Coming Soon"
	msgX := screenWidth/2 - len(msg)*7/2
	msgY := screenHeight / 2
	ebitenutil.DebugPrintAt(screen, msg, msgX, msgY)
}

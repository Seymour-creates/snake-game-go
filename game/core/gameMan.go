package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"snakeGame/game/entities"
	"snakeGame/game/render"
	"snakeGame/game/ui"
)

// Game implements the ebiten.Game interface and holds the game state.
// It contains all entities and game data (Snake, food, score, level, etc.).
// Ebiten requires Game to implement Update, Draw, and Layout methods.
// This struct is designed for easy extension (e.g., adding more levels, skins, music).
type Game struct {
	State                     *StateManager
	Speed                     *SpeedManager
	UI                        *render.UIManager
	Renderer                  *render.Renderer
	SpriteManager             *render.SpriteManager
	screenWidth, screenHeight int                       // screen size in pixels for rendering
	gridWidth, gridHeight     int                       // grid size in cells (play field dimensions)
	cellSize                  int                       // pixel size of one grid cell
	Snake                     *entities.SnakeController // the player-controlled Snake
	Food                      *entities.Food            // the current food item on the board
	gameOver                  bool                      // whether the game is over
	frameCount                int                       // Tracks number of frames since last move
	frameDelay                int                       // Delay between moves (in frames)
	showRetry                 bool                      // used to toggle retry prompt visibility
}

// NewGame initializes a new game state with a Snake and an initial food.
func NewGame(startPosition image.Point, gridWidth, gridHeight, cellSize int) *Game {
	snake := entities.NewSnakeController(startPosition, gridWidth, gridHeight)
	log.Println("Attempting to load sprite sheet...")

	sprites := render.NewSpriteManager(cellSize)

	log.Println("SpriteManager initialized")
	g := &Game{
		State:         NewStateManager(),
		Speed:         NewSpeedManager(),
		UI:            render.NewUIManager(),
		Renderer:      render.NewRenderer(sprites),
		SpriteManager: sprites,
		Snake:         snake,
		Food:          entities.NewFood(snake, gridWidth, gridHeight),
		gridWidth:     gridWidth,
		gridHeight:    gridHeight,
		cellSize:      cellSize,
		screenWidth:   gridWidth * cellSize,
		screenHeight:  gridHeight * cellSize,
		gameOver:      false,
		showRetry:     false,
		frameDelay:    20,
	}

	return g
}

// Update advances the game state by one frame (called ~60 times per second by Ebiten).
func (g *Game) Update() error {
	// Pause/resume toggle should still work while game is running
	g.handlePauseToggle()

	// Allow reset via Enter or mouse click even if game is over
	if g.State.GameOver {
		g.handleRetryInput()
		return nil
	}

	// Exit early if paused or game is over
	if !g.State.IsRunning() {
		return nil
	}

	// Handle user input for Snake direction.
	g.handleInput()

	// Frame throttling — only update logic every N frames
	if !g.Speed.ShouldUpdate() {
		return nil
	}

	// Calculate the Snake's new head position based on current direction.
	g.Snake.ApplyPendingDirection(g.gridWidth, g.gridHeight)
	newHead := g.Snake.NextHeadPosition()

	// Collision check: Wall boundaries || Snake runs into itself.
	if g.checkCollision(newHead) {
		g.State.SetGameOver()
		return nil
	}

	// Check if food is eaten.
	if g.checkFoodEaten(newHead) {
		g.handleFoodEaten()
	} else {
		g.Snake.MoveForward()
	}

	return nil // No error, continue game.
}

// Draw renders the game state to the screen (called every frame after Update).
func (g *Game) Draw(screen *ebiten.Image) {

	screenWidth := g.gridWidth * g.cellSize
	screenHeight := g.gridHeight * g.cellSize

	// Draw theme
	ui.DrawBackground(screen, screenWidth, screenHeight, g.cellSize)
	ui.DrawBorder(screen, screenWidth, screenHeight, g.cellSize)

	// Draw the Snake.
	g.Renderer.DrawSnake(screen, g.Snake)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.Food.Pos.X*g.cellSize), float64(g.Food.Pos.Y*g.cellSize))

	// Draw the food.
	ui.DrawFood(screen, g.Food.Pos, g.cellSize)

	if g.State.Paused {
		g.UI.DrawPauseOverlay(screen, g.screenWidth, g.screenHeight)
	}

	// Overlay game status text.
	if g.State.GameOver {
		g.UI.DrawGameOverOverlay(screen, g.screenWidth, g.screenHeight)
	} else {
		g.UI.DrawStatus(screen, g.State.Score, g.State.Level)
		if g.State.Paused {
			g.UI.DrawPauseOverlay(screen, g.screenWidth, g.screenHeight)
		}
	}
}

func (g *Game) resetGame() {
	start := image.Pt(g.gridWidth/2, g.gridHeight/2)
	g.Snake = entities.NewSnakeController(start, g.gridWidth, g.gridHeight)
	g.Food = entities.NewFood(g.Snake, g.gridWidth, g.gridHeight)
	g.State.GameOver = false
	g.frameCount = 0
	g.showRetry = false
}

// Layout returns the internal screen size (logical resolution) for the game.
func (g *Game) Layout(_, _ int) (int, int) {
	// We use a fixed logical size; Ebiten will scale the actual window accordingly.
	return g.screenWidth, g.screenHeight
}

func (g *Game) checkFoodEaten(newHead image.Point) bool {
	return newHead == g.Food.Pos
}

func (g *Game) handleFoodEaten() {
	g.Snake.Grow()
	g.State.IncreaseScore()
	g.Speed.AdjustDelayByLevel(g.State.Level)
	g.Food = entities.NewFood(g.Snake, g.gridWidth, g.gridHeight)
}

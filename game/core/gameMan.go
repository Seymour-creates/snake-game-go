package core

import (
	"image"
	"log"
	"os"
	"snakeGame/game/audio"
	"snakeGame/game/entities"
	"snakeGame/game/render"
	"snakeGame/game/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// GameScreen represents the current screen/state of the game.
type GameScreen int

const (
	ScreenTitle GameScreen = iota
	ScreenPlaying
	ScreenSettings
	ScreenGameOver
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
	SoundMan                  *audio.SoundManager
	CurrentScreen             GameScreen
	menuSelected              int // 0 = Play Game, 1 = Settings
	gameOverSelected          int // 0 = Play Again, 1 = Main Menu, 2 = Exit Game
}

// NewGame initializes a new game state with a Snake and an initial food.
func NewGame(startPosition image.Point, gridWidth, gridHeight, cellSize int) *Game {
	snake := entities.NewSnakeController(startPosition, gridWidth, gridHeight)
	log.Println("Attempting to load sprite sheet...")

	sprites := render.NewSpriteManager(cellSize)

	log.Println("SpriteManager initialized")
	g := &Game{
		State:            NewStateManager(),
		Speed:            NewSpeedManager(),
		UI:               render.NewUIManager(),
		Renderer:         render.NewRenderer(sprites),
		SpriteManager:    sprites,
		Snake:            snake,
		Food:             entities.NewFood(snake, gridWidth, gridHeight),
		gridWidth:        gridWidth,
		gridHeight:       gridHeight,
		cellSize:         cellSize,
		screenWidth:      gridWidth * cellSize,
		screenHeight:     gridHeight * cellSize,
		gameOver:         false,
		showRetry:        false,
		frameDelay:       20,
		SoundMan:         audio.NewSoundManager(),
		CurrentScreen:    ScreenTitle,
		menuSelected:     0,
		gameOverSelected: 0,
	}

	g.loadSounds()
	return g
}

// Update advances the game state by one frame (called ~60 times per second by Ebiten).
func (g *Game) Update() error {
	if g.CurrentScreen == ScreenTitle {
		// Handle menu navigation
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.menuSelected = (g.menuSelected + 1) % 2 // wrap
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.menuSelected = (g.menuSelected + 1) % 2 // wrap
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if g.menuSelected == 0 {
				g.CurrentScreen = ScreenPlaying
			} else if g.menuSelected == 1 {
				g.CurrentScreen = ScreenSettings
			}
		}
		return nil
	}
	if g.CurrentScreen == ScreenGameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.gameOverSelected = (g.gameOverSelected + 2) % 3 // wrap up
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.gameOverSelected = (g.gameOverSelected + 1) % 3 // wrap down
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			switch g.gameOverSelected {
			case 0: // Play Again
				g.resetGame()
				g.CurrentScreen = ScreenPlaying
			case 1: // Main Menu
				g.resetGame()
				g.CurrentScreen = ScreenTitle
			case 2: // Exit Game
				os.Exit(0)
			}
		}
		return nil
	}
	// Pause/resume toggle should still work while game is running
	g.handlePauseToggle()

	// Allow reset via Enter or mouse click even if game is over
	if g.State.GameOver {
		g.CurrentScreen = ScreenGameOver
		g.gameOverSelected = 0
		return nil
	}

	// Handle background music play/pause
	if g.State.Paused {
		g.SoundMan.PauseLoopingSound("bgm")
	} else {
		g.SoundMan.PlayLoopingSound("bgm")
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
		// Play apple bite sound
		g.SoundMan.PlaySound("bite")
	} else {
		g.Snake.MoveForward()
	}

	return nil // No error, continue game.
}

// Draw renders the game state to the screen (called every frame after Update).
func (g *Game) Draw(screen *ebiten.Image) {
	if g.CurrentScreen == ScreenTitle {
		g.UI.DrawTitleScreen(screen, g.screenWidth, g.screenHeight, g.menuSelected)
		return
	}
	if g.CurrentScreen == ScreenGameOver {
		g.SoundMan.PauseLoopingSound("bgm")
		g.UI.DrawGameOverOverlay(screen, g.screenWidth, g.screenHeight, g.gameOverSelected)
		return
	}
	if g.CurrentScreen == ScreenSettings {
		g.UI.DrawSettingsScreen(screen, g.screenWidth, g.screenHeight)
		return
	}

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
		g.SoundMan.PauseLoopingSound("bgm")
		g.UI.DrawGameOverOverlay(screen, g.screenWidth, g.screenHeight, g.gameOverSelected)
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
	g.State.Reset()
	g.frameCount = 0
	g.showRetry = false
	g.loadSounds()
	// Restart background music from the beginning
	if g.SoundMan != nil {
		if player, ok := g.SoundMan.LoopingPlayer("bgm"); ok {
			player.Rewind()
			player.Play()
		}
	}
	g.CurrentScreen = ScreenTitle
	g.menuSelected = 0
	g.gameOverSelected = 0
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

// Helper to load all sounds
func (g *Game) loadSounds() {
	if g.SoundMan == nil {
		return
	}
	g.SoundMan.ClearSounds()
	// Load background music (looping)
	bgData, err := os.ReadFile("game/audio/assets/sound_snake_movement_fixed.wav")
	if err == nil {
		if err := g.SoundMan.LoadLoopingSound("bgm", bgData); err != nil {
			log.Printf("Failed to load background music: %v", err)
		} else {
			log.Printf("Background music loaded successfully")
		}
	} else {
		log.Printf("Failed to read background music file: %v", err)
	}
	// Load apple bite sound
	biteData, err := os.ReadFile("game/audio/assets/sound_apple_bite_fixed.wav")
	if err == nil {
		if err := g.SoundMan.LoadSound("bite", biteData); err != nil {
			log.Printf("Failed to load apple bite sound: %v", err)
		} else {
			log.Printf("Apple bite sound loaded successfully")
		}
	} else {
		log.Printf("Failed to read apple bite sound file: %v", err)
	}
}

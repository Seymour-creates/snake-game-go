package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"math"
	"os"
)

var (
	grassTile   *ebiten.Image
	vineSide    *ebiten.Image
	vineCorner  *ebiten.Image
	appleSprite *ebiten.Image
)

// LoadTheme loads all UI-related sprites into memory.
// Place the image files in game/ui/assets/ with correct names.
func LoadTheme() error {
	grassTile = loadImage("game/ui/assets/tile_pebble.png")
	vineSide = loadImage("game/ui/assets/border_side_vine.png")
	vineCorner = loadImage("game/ui/assets/border_corner_vine.png")
	appleSprite = loadImage("game/ui/assets/apple_sprite.png")
	return nil
}

func loadImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		panic("Failed to load image: " + path)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic("Failed to decode image: " + path)
	}

	return ebiten.NewImageFromImage(img)
}

// DrawBackground tiles the grass tile to fill the entire screen.
func DrawBackground(screen *ebiten.Image, screenWidth, screenHeight int, cellSize int) {
	tileW, tileH := grassTile.Bounds().Dx(), grassTile.Bounds().Dy()
	scaleX := float64(cellSize) / float64(tileW)
	scaleY := float64(cellSize) / float64(tileH)

	for y := 0; y < screenHeight; y += cellSize {
		for x := 0; x < screenWidth; x += cellSize {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scaleX, scaleY)
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(grassTile, op)
		}
	}
}

// DrawBorder draws a decorative vine border around the game screen.
func DrawBorder(screen *ebiten.Image, screenWidth, screenHeight int, cellSize int) {
	tileH := vineSide.Bounds().Dy()
	scale := float64(cellSize) / float64(tileH)

	tilesX := screenWidth / cellSize
	tilesY := screenHeight / cellSize

	// Top & Bottom sides
	for i := 0; i < tilesX; i++ {
		x := float64(i * cellSize)

		// Top (rotate 90°)
		topOp := &ebiten.DrawImageOptions{}
		topOp.GeoM.Scale(scale, scale)
		topOp.GeoM.Rotate(math.Pi / 2)
		topOp.GeoM.Translate(x+float64(cellSize), 0)
		screen.DrawImage(vineSide, topOp)

		// Bottom (rotate -90°)
		bottomOp := &ebiten.DrawImageOptions{}
		bottomOp.GeoM.Scale(scale, scale)
		bottomOp.GeoM.Rotate(-math.Pi / 2)
		bottomOp.GeoM.Translate(x, float64(screenHeight))
		screen.DrawImage(vineSide, bottomOp)
	}

	// Left & Right sides
	for i := 0; i < tilesY; i++ {
		y := float64(i * cellSize)

		// Left (no rotation)
		leftOp := &ebiten.DrawImageOptions{}
		leftOp.GeoM.Scale(scale, scale)
		leftOp.GeoM.Translate(0, y)
		screen.DrawImage(vineSide, leftOp)

		// Right (rotate 180°)
		rightOp := &ebiten.DrawImageOptions{}
		rightOp.GeoM.Scale(scale, scale)
		rightOp.GeoM.Rotate(math.Pi)
		rightOp.GeoM.Translate(float64(screenWidth), y+float64(cellSize))
		screen.DrawImage(vineSide, rightOp)
	}

	// Four corners
	// Top-left (no rotation)
	tlOp := &ebiten.DrawImageOptions{}
	tlOp.GeoM.Scale(scale, scale)
	tlOp.GeoM.Translate(0, 0)
	screen.DrawImage(vineCorner, tlOp)

	// Top-right (rotate 90°)
	trOp := &ebiten.DrawImageOptions{}
	trOp.GeoM.Scale(scale, scale)
	trOp.GeoM.Rotate(math.Pi / 2)
	trOp.GeoM.Translate(float64(screenWidth), 0)
	screen.DrawImage(vineCorner, trOp)

	// Bottom-right (rotate 180°)
	brOp := &ebiten.DrawImageOptions{}
	brOp.GeoM.Scale(scale, scale)
	brOp.GeoM.Rotate(math.Pi)
	brOp.GeoM.Translate(float64(screenWidth), float64(screenHeight))
	screen.DrawImage(vineCorner, brOp)

	// Bottom-left (rotate -90°)
	blOp := &ebiten.DrawImageOptions{}
	blOp.GeoM.Scale(scale, scale)
	blOp.GeoM.Rotate(-math.Pi / 2)
	blOp.GeoM.Translate(0, float64(screenHeight))
	screen.DrawImage(vineCorner, blOp)
}

func DrawFood(screen *ebiten.Image, pos image.Point, cellSize int) {
	tileW := appleSprite.Bounds().Dx()
	tileH := appleSprite.Bounds().Dy()
	scaleX := float64(cellSize) / float64(tileW)
	scaleY := float64(cellSize) / float64(tileH)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(pos.X*cellSize), float64(pos.Y*cellSize))
	screen.DrawImage(appleSprite, op)
}

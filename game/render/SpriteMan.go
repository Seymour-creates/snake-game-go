package render

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// SnakePart represents the specific sprite type.
type SnakePart int

const (
	Head SnakePart = iota
	Tail
	BodyVertical
	BodyHorizontal
	Bend
)

// SpriteManager manages the sprite atlas.
type SpriteManager struct {
	Parts    map[SnakePart]*ebiten.Image
	CellSize int
}

func loadImageScaled(path string, width, height int) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load image: %s, error: %v", path, err)
	}
	return scaleSprite(img, width, height)
}

func scaleSprite(src *ebiten.Image, targetWidth, targetHeight int) *ebiten.Image {
	dst := ebiten.NewImage(targetWidth, targetHeight)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(targetWidth)/float64(src.Bounds().Dx()), float64(targetHeight)/float64(src.Bounds().Dy()))
	dst.DrawImage(src, op)
	return dst
}

func NewSpriteManager(cellSize int) *SpriteManager {
	sprites := map[SnakePart]*ebiten.Image{
		Head:           loadImageScaled("assets/snake_head_upward.png", cellSize, cellSize),
		Tail:           loadImageScaled("assets/snake_tail.png", cellSize, cellSize),
		BodyVertical:   loadImageScaled("assets/snake_body.png", cellSize, cellSize),
		BodyHorizontal: loadImageScaled("assets/snake_body.png", cellSize, cellSize),
		Bend:           loadImageScaled("assets/snake_body_bend.png", cellSize, cellSize),
	}

	return &SpriteManager{
		Parts:    sprites,
		CellSize: cellSize,
	}
}

func (s *SpriteManager) ResolveSegmentSprite(tileType SnakePart) *ebiten.Image {
	return s.Parts[tileType]
}

func (s *SpriteManager) DrawSegment(screen *ebiten.Image, spriteType SnakePart, pos image.Point, rotation float64) {
	img := s.ResolveSegmentSprite(spriteType)
	if img == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	if rotation != 0 {
		op.GeoM.Translate(-float64(s.CellSize)/2, -float64(s.CellSize)/2)
		op.GeoM.Rotate(rotation)
		op.GeoM.Translate(float64(s.CellSize)/2, float64(s.CellSize)/2)
	}

	op.GeoM.Translate(float64(pos.X*s.CellSize), float64(pos.Y*s.CellSize))
	screen.DrawImage(img, op)
}

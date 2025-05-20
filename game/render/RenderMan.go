package render

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"snakeGame/game/entities"
)

type Renderer struct {
	SpriteManager *SpriteManager
}

func NewRenderer(sm *SpriteManager) *Renderer {
	return &Renderer{SpriteManager: sm}
}

// DrawSnake renders the entire snake segment-by-segment
func (r *Renderer) DrawSnake(screen *ebiten.Image, sc *entities.SnakeController) {
	for seg := sc.Head; seg != nil; seg = seg.Next {
		var spriteType SnakePart
		var angle float64

		switch seg.Tile {
		case entities.TileHead:
			spriteType = Head
			angle = seg.Rotation
		case entities.TileTail:
			spriteType = Tail
			angle = seg.Rotation + 180
		case entities.TileBend:
			spriteType = Bend
			angle = seg.Rotation
		case entities.TileBody:
			// Determine if body is horizontal or vertical for sprite (optional if using single sprite + rotation)
			if seg.Rotation == 90.0 || seg.Rotation == 270.0 {
				spriteType = BodyHorizontal
				angle = seg.Rotation
			} else {
				spriteType = BodyVertical
				angle = seg.Rotation
			}
		default:
			spriteType = BodyVertical // fallback
			angle = seg.Rotation
		}

		r.SpriteManager.DrawSegment(screen, spriteType, seg.Pos, angle*math.Pi/180)
	}
}

// DrawFood renders the food block
func (r *Renderer) DrawFood(screen *ebiten.Image, food *entities.Food, cellSize int) {
	fx := food.Pos.X * cellSize
	fy := food.Pos.Y * cellSize
	vector.DrawFilledRect(screen, float32(fx), float32(fy), float32(cellSize), float32(cellSize), color.RGBA{R: 255, A: 255}, false)
}

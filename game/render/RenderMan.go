package render

import (
	"github.com/hajimehoshi/ebiten/v2"
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

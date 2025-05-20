package entities

import (
	"image"
	"log"
)

// Directions
var (
	Up    = image.Pt(0, -1)
	Down  = image.Pt(0, 1)
	Left  = image.Pt(-1, 0)
	Right = image.Pt(1, 0)
)

// TileType Tile types
type TileType uint8

const (
	TileHead TileType = iota
	TileBody
	TileBend
	TileTail
)

// SnakeSegment Segment in linked list
type SnakeSegment struct {
	Pos      image.Point
	Tile     TileType
	Rotation float64
	Prev     *SnakeSegment
	Next     *SnakeSegment
}

// SnakeController manages the snake state and logic
type SnakeController struct {
	Head             *SnakeSegment
	Tail             *SnakeSegment
	Dir              image.Point
	PendingDir       image.Point
	GridWidth        int
	GridHeight       int
	Growing          bool
}

// NewSnakeController sets up a snake with head, 2 body segments, and a tail
func NewSnakeController(start image.Point, gridWidth, gridHeight int) *SnakeController {
	tail := &SnakeSegment{
		Pos:      image.Pt(start.X-3, start.Y),
		Tile:     TileTail,
		Rotation: 270.0,
	}
	body2 := &SnakeSegment{
		Pos:      image.Pt(start.X-2, start.Y),
		Tile:     TileBody,
		Rotation: 90.0,
		Next:     tail,
	}
	tail.Prev = body2

	body1 := &SnakeSegment{
		Pos:      image.Pt(start.X-1, start.Y),
		Tile:     TileBody,
		Rotation: 90.0,
		Next:     body2,
	}
	body2.Prev = body1

	head := &SnakeSegment{
		Pos:      start,
		Tile:     TileHead,
		Rotation: 90.0,
		Next:     body1,
	}
	body1.Prev = head

	return &SnakeController{
		Head:       head,
		Tail:       tail,
		Dir:        Right,
		PendingDir: Right,
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
	}
}

// MoveForward shifts the snake forward
func (sc *SnakeController) MoveForward() {
	sc.Dir = sc.PendingDir
	newHeadPos := sc.Head.Pos.Add(sc.Dir)

	newHead := &SnakeSegment{
		Pos:      newHeadPos,
		Tile:     TileHead,
		Rotation: directionToAngle(sc.Dir),
		Next:     sc.Head,
	}
	sc.Head.Prev = newHead
	sc.Head = newHead

	// Remove tail if not growing
	if !sc.Growing {
		oldTail := sc.Tail
		sc.Tail = oldTail.Prev
		if sc.Tail != nil {
			sc.Tail.Next = nil
			sc.Tail.Tile = TileTail

			// should face AWAY from previous segment
			if sc.Tail.Prev != nil {
				dir := sc.Tail.Prev.Pos.Sub(sc.Tail.Pos)
				sc.Tail.Rotation = directionToAngle(dir)
			}

		}
	} else {
		sc.Growing = false
	}

	sc.assignBends()
}

// assignBends updates types and rotation of interior segments
func (sc *SnakeController) assignBends() {
	if sc.Head == nil {
		return
	}

	sc.Head.Tile = TileHead
	sc.Head.Rotation = directionToAngle(sc.Dir)

	for seg := sc.Head.Next; seg != nil && seg.Next != nil; seg = seg.Next {
		prev := seg.Prev.Pos
		curr := seg.Pos
		next := seg.Next.Pos

		dir1 := prev.Sub(curr)
		dir2 := next.Sub(curr)

		if dir1.X == dir2.X || dir1.Y == dir2.Y {
			// Straight
			seg.Tile = TileBody
			if dir1.X != 0 {
				seg.Rotation = 0.0 // horizontal
			} else {
				seg.Rotation = 90.0 // vertical
			}
		} else {
			// Bend
			seg.Tile = TileBend
			seg.Rotation = bendAngleForDirections(dir1, dir2)
		}
	}
}

// Grow triggers growth on next move
func (sc *SnakeController) Grow() {
	sc.Growing = true
}

// HeadPos returns current head position
func (sc *SnakeController) HeadPos() image.Point {
	return sc.Head.Pos
}

// NextHeadPosition returns next head position
func (sc *SnakeController) NextHeadPosition() image.Point {
	return sc.Head.Pos.Add(sc.Dir)
}

// ApplyPendingDirection sets new direction if valid
func (sc *SnakeController) ApplyPendingDirection(gridWidth, gridHeight int) {
	if sc.CanMoveTo(sc.PendingDir, gridWidth, gridHeight) {
		sc.Dir = sc.PendingDir
	}
}

// CanMoveTo checks bounds
func (sc *SnakeController) CanMoveTo(dir image.Point, gridWidth, gridHeight int) bool {
	next := sc.Head.Pos.Add(dir)
	return next.X >= 0 && next.X < gridWidth && next.Y >= 0 && next.Y < gridHeight
}

// Occupies returns true if a point is occupied
func (sc *SnakeController) Occupies(pt image.Point) bool {
	seg := sc.Head
	for seg != nil {
		if seg.Pos == pt {
			return true
		}
		seg = seg.Next
	}
	return false
}

// Converts direction to degrees
func directionToAngle(dir image.Point) float64 {
	switch dir {
	case Up:
		return 0.0
	case Right:
		return 90.0
	case Down:
		return 180.0
	case Left:
		return 270.0
	default:
		log.Printf("Unknown dir: %+v\n", dir)
		return 0.0
	}
}

// Returns angle for bend segment
func bendAngleForDirections(d1, d2 image.Point) float64 {
	if (d1 == Up && d2 == Right) || (d1 == Right && d2 == Up) {
		return 0.0
	}
	if (d1 == Right && d2 == Down) || (d1 == Down && d2 == Right) {
		return 90.0
	}
	if (d1 == Down && d2 == Left) || (d1 == Left && d2 == Down) {
		return 180.0
	}
	if (d1 == Left && d2 == Up) || (d1 == Up && d2 == Left) {
		return 270.0
	}
	return 0.0
}

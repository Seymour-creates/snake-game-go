package core

import (
	"image"
)

func (g *Game) checkCollision(newHead image.Point) bool {
	// Check out-of-bounds
	outOfBounds := newHead.X < 0 || newHead.X >= g.gridWidth || newHead.Y < 0 || newHead.Y >= g.gridHeight

	// Determine if the snake is growing (e.g., head touches food)
	growing := newHead == g.Food.Pos

	// Get current tail
	tail := g.Snake.Tail

	// Check self-collision
	selfHit := (!growing && newHead != tail.Pos && g.Snake.Occupies(newHead)) ||
		(growing && g.Snake.Occupies(newHead))

	return outOfBounds || selfHit
}


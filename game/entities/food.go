package entities

import (
	"image"
	"math/rand"
)

// Food represents a food item that the snake can eat.
type Food struct {
	Pos image.Point // grid position of the food
}

// NewFood generates a Food item at a random position not occupied by the snake.
func NewFood(snake *SnakeController, gridWidth, gridHeight int) *Food {
	for {
		// Pick a random position within the grid.
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		candidate := image.Point{X: x, Y: y}
		// Ensure the random position is not on the snake.
		if !snake.Occupies(candidate) {
			return &Food{Pos: candidate}
		}
		// If it collides with the snake, loop and try another position.
	}
}


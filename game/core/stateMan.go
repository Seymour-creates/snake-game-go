package core

import "fmt"

// StateManager manages global game state such as score, level, pause, and game over flags.
type StateManager struct {
	Score    int
	Level    int
	Paused   bool
	GameOver bool
}

// NewStateManager initializes a new game state.
func NewStateManager() *StateManager {
	return &StateManager{
		Score:    0,
		Level:    1,
		Paused:   false,
		GameOver: false,
	}
}

// IncreaseScore increments the score and adjusts level as needed.
func (s *StateManager) IncreaseScore() {
	s.Score++
	if s.Score%5 == 0 {
		s.Level++
	}
}

// Reset restores the game state to its initial values.
func (s *StateManager) Reset() {
	s.Score = 0
	s.Level = 1
	s.Paused = false
	s.GameOver = false
}

// TogglePause switches the pause state.
func (s *StateManager) TogglePause() {
	fmt.Println("Pause toggled")
	s.Paused = !s.Paused
}

// SetGameOver marks the game as over.
func (s *StateManager) SetGameOver() {
	s.GameOver = true
}

// IsRunning returns true if the game is not paused or over.
func (s *StateManager) IsRunning() bool {
	return !s.Paused && !s.GameOver
}


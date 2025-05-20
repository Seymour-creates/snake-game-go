package core

import "math"

// SpeedManager controls how often the game updates the snake's position based on level.
type SpeedManager struct {
	FrameCount int
	FrameDelay int
}

// NewSpeedManager initializes the manager with default delay.
func NewSpeedManager() *SpeedManager {
	return &SpeedManager{
		FrameCount: 0,
		FrameDelay: 20, // default speed (3 moves/sec)
	}
}

// ShouldUpdate returns true if it's time to update the game logic (snake move).
func (s *SpeedManager) ShouldUpdate() bool {
	s.FrameCount++
	if s.FrameCount >= s.FrameDelay {
		s.FrameCount = 0
		return true
	}
	return false
}

// AdjustDelayByLevel updates delay based on level to increase game speed.
func (s *SpeedManager) AdjustDelayByLevel(level int) {
	// Simple formula: base delay - (level * 2), but not lower than 5 frames
	s.FrameDelay = int(math.Max(5, float64(20 - level*2)))
}

// ResetFrameCount clears the frame counter, useful on restart.
func (s *SpeedManager) ResetFrameCount() {
	s.FrameCount = 0
}


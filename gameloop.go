// Package gameloop provides functions for creating game loops.
package gameloop

// Config describes a configuration for a game loop.
type Config struct {
	// TargetFPS is used to calculate the seconds per update (1 / TargetFPS).
	TargetFPS uint16

	// IdleThreshold prevents updating the game if the time
	// elapsed since the previous frame exceeds this number (in seconds).
	IdleThreshold float64

	// CurrentTimeFunc is a callback that should return the current time in
	// seconds. It is used by the game loop to calculate the time elapsed
	// between frames.
	CurrentTimeFunc func() float64

	// ProcessInputFunc is a callback that is called within the game loop. It is
	// called before the UpdateFunc. It should process input logic
	// (ie. keyboard, mouse, gamepad, etc.) and return a flag to signal the
	// game loop to quit.
	//
	// This function will not be called if the elapsed time since the previous
	// frame has exceeded IdleThreshold (ie. if window is minimized, etc.).
	ProcessInputFunc func() (quit bool)

	// UpdateFunc is a callback that is called within the game loop. It is
	// called after ProcessInputFunc and it should contain logic that updates
	// the game's state. This function will be called based on a fixed interval
	// of 1 / TargetFPS (ie. 1 sec / 60 FPS = 0.01667 secs) and it is passed as
	// a parameter (dt).
	//
	// This function will not be called if the elapsed time since the previous
	// frame has exceeded IdleThreshold (ie. if window is minimized, etc.).
	UpdateFunc func(dt float64)

	// RenderFunc is a callback that is called within the game loop. It is
	// called after UpdateFunc and it should contain rendering logic.
	//
	// This function will not be called if the elapsed time since the previous
	// frame has exceeded IdleThreshold (ie. if window is minimized, etc.).
	RenderFunc func()
}

// Create will create a game loop based on a given configuration.
func Create(c Config) func() {
	secsPerUpdate := 1 / float64(c.TargetFPS)

	previous := c.CurrentTimeFunc()
	lag := 0.0

	var current, elapsed float64

	return func() {
		for quit := false; !quit; {
			current = c.CurrentTimeFunc()
			elapsed = current - previous
			previous = current

			if elapsed > c.IdleThreshold {
				continue
			}

			lag += elapsed

			quit = c.ProcessInputFunc()

			for lag >= secsPerUpdate {
				c.UpdateFunc(secsPerUpdate)
				lag -= secsPerUpdate
			}

			c.RenderFunc()
		}
	}
}

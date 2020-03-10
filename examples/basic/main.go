package main

import (
	"log"
	"time"

	"github.com/atamocius/gameloop"
)

func main() {
	// Create a game loop config.
	config := gameloop.Config{
		// TargetFPS is used to calculate the seconds per update
		// (1 / TargetFPS).
		TargetFPS: 60,

		// IdleThreshold prevents updating the game if the time
		// elapsed since the previous frame exceeds this number (in seconds).
		IdleThreshold: 1,

		// CurrentTimeFunc accepts a function that returns the current time in
		// seconds. The gameloop library only provides a scaffold, it is up to
		// the user to provide an implementation. In this case, time's UnixNano
		// method was used but had to be multiplied by 0.000000001 to convert
		// to seconds.
		CurrentTimeFunc: func() float64 {
			return float64(time.Now().UnixNano()) * 1e-9
		},

		// ProcessInputFunc accepts a function that processes input logic
		// (ie. keyboard, mouse, gamepad, etc.) and returns a flag to signal the
		// game loop to quit.
		ProcessInputFunc: func() bool {
			time.Sleep(5 * time.Millisecond) // Simulating work
			log.Println("process input")
			return false
		},

		// UpdateFunc accepts a function that updates the game's state.
		// This function will be called based on a fixed interval
		// of 1 / TargetFPS (ie. 1 sec / 60 FPS = 0.01667 secs) and it is passed
		// as a parameter (dt).
		UpdateFunc: func(dt float64) {
			time.Sleep(5 * time.Millisecond) // Simulating work
			log.Printf("updating, dt: %v\n", dt)
		},

		// RenderFunc accepts a function that contains rendering logic.
		RenderFunc: func() {
			time.Sleep(5 * time.Millisecond) // Simulating work
			log.Println("rendering")
		},
	}

	// Call the gameloop.Create() function and pass the config to create
	// a game loop.
	runLoop := gameloop.Create(config)

	// Run the created game loop.
	runLoop()
}

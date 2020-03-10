# gameloop

A [_fixed time step_](http://gameprogrammingpatterns.com/game-loop.html#play-catch-up) game loop with variable rate rendering. Suitable for [_time-based animations_](http://blog.sklambert.com/using-time-based-animation-implement/#time-based-animation).

## Installation

```bash
go get github.com/atamocius/gameloop
```

## Examples

### Basic Usage

```go
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
```

### Using With SDL

```go
package main

import (
	"fmt"

	"github.com/atamocius/gameloop"
	"github.com/veandco/go-sdl2/sdl"
)

const windowWidth, windowHeight = 350, 350

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_TIMER); err != nil {
		fmt.Printf("error initializing SDL: %v", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Game",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		windowWidth, windowHeight, 0)
	if err != nil {
		fmt.Printf("error creating window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("error creating renderer: %v", err)
	}
	defer renderer.Destroy()

	sq := square{
		ScrollSpeed: 60,
		Rect: sdl.FRect{
			X: 100,
			Y: 100,
			W: 20,
			H: 20,
		},
		Delta: sdl.FPoint{
			X: 4,
			Y: 2,
		},
	}

	render := func(r *sdl.Renderer) {
		r.SetDrawColor(0, 0, 0, 255)
		r.Clear()
		sq.Draw(r)
		r.Present()
	}

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
		// the user to provide an implementation. In this case, SDL's GetTicks
		// was used.
		CurrentTimeFunc: func() float64 {
			return float64(sdl.GetTicks()) * 0.001
		},

		// ProcessInputFunc accepts a function that processes input logic
		// (ie. keyboard, mouse, gamepad, etc.) and returns a flag to signal the
		// game loop to quit.
		ProcessInputFunc: func() (quit bool) {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				if event.GetType() == sdl.QUIT {
					return true
				}
			}
			return false
		},

		// UpdateFunc accepts a function that updates the game's state.
		// This function will be called based on a fixed interval
		// of 1 / TargetFPS (ie. 1 sec / 60 FPS = 0.01667 secs) and it is passed
		// as a parameter (dt).
		UpdateFunc: func(dt float64) {
			sq.Move(float32(dt))
		},

		// RenderFunc accepts a function that contains rendering logic.
		// Since this function does not support parameters, we need to take
		// advantage of closures to create one that wraps our use of
		// SDL's Renderer.
		RenderFunc: func() {
			render(renderer)
		},
	}

	// Call the gameloop.Create() function and pass the config to create
	// a game loop.
	runLoop := gameloop.Create(config)

	// Run the created game loop.
	runLoop()
}

type square struct {
	ScrollSpeed float32
	Rect        sdl.FRect
	Delta       sdl.FPoint
}

func (s *square) Move(dt float32) {
	s.Rect.X += s.Delta.X * dt * s.ScrollSpeed
	s.Rect.Y += s.Delta.Y * dt * s.ScrollSpeed

	if s.Rect.X <= 0 || s.Rect.X >= windowWidth-s.Rect.W {
		s.Delta.X = -s.Delta.X
	}
	if s.Rect.Y <= 0 || s.Rect.Y >= windowHeight-s.Rect.H {
		s.Delta.Y = -s.Delta.Y
	}
}

func (s *square) Draw(r *sdl.Renderer) {
	r.SetDrawColor(255, 0, 0, 255)
	r.FillRectF(&s.Rect)
}
```

## Sources

- [Game Programming Patterns - Game Loop](http://gameprogrammingpatterns.com/game-loop.html)
- [Fix Your Timestep!](https://gafferongames.com/post/fix_your_timestep/)
- [HTML5 Gamer - Why You Should be Using Time-based Animation and How to Implement it](http://blog.sklambert.com/using-time-based-animation-implement/)
- [Kontra - GameLoop](https://straker.github.io/kontra/api/gameLoop)
- [Kontra - GameLoop (source)](https://github.com/straker/kontra/blob/master/src/gameLoop.js)
- [bell0bytes - The Game Loop](https://bell0bytes.eu/the-game-loop/)
- [lwjglgamedev - The Game Loop](https://ahbejarano.gitbook.io/lwjglgamedev/chapter2)
- [G-Engine #3: Game Loop](http://clarkkromenaker.com/post/gengine-03-game-loop/)
- [Lesson 08 - Timing: Frame Rate, Physics, Animation](https://thenumbat.github.io/cpp-course/sdl2/08/08.html)
- [The Game Loop and Frame Rate Management](http://www.brandonfoltz.com/downloads/tutorials/The_Game_Loop_and_Frame_Rate_Management.pdf)

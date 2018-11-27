package main;

import (
  "github.com/veandco/go-sdl2/sdl"
)

func keyevents(action chan Action) {
  a := Action{dir: North, quit: false}
  for {
    for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
      switch e := event.(type) {
        case *sdl.KeyboardEvent:
          a = handle_KeyboardEvent(a, e);
      }
    }
    action <- a
  }
}

func handle_KeyboardEvent(action Action, e *sdl.KeyboardEvent) Action {
  // Check for move
  if (e.Keysym.Sym == sdl.K_UP) && (e.State == sdl.PRESSED) {
    action.dir = North
  }
  if (e.Keysym.Sym == sdl.K_DOWN) && (e.State == sdl.PRESSED) {
    action.dir = South
  }
  if (e.Keysym.Sym == sdl.K_LEFT) && (e.State == sdl.PRESSED) {
    action.dir = East
  }
  if (e.Keysym.Sym == sdl.K_RIGHT) && (e.State == sdl.PRESSED) {
    action.dir = West
  }

  // Check for quit
  if (e.Keysym.Sym == sdl.K_q) {
    action.quit = true;
  }

  return action;
}


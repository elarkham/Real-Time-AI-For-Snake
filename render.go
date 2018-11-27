package main;

import (
  "fmt"
  "os"

  "github.com/veandco/go-sdl2/sdl"
)

type Color struct {
  r uint8
  g uint8
  b uint8
  a uint8
}

var Win_Title string = "Snake"

var Win_W int32 = 500
var Win_H int32 = 500
var Pix_D int32 = 11

// Colors
//var C_BG Color = Color{167,167,167,255}
var C_BG Color = Color{88,88,90,255}

var C_Snake Color = Color{28,172,120,255}
var C_Apple Color = Color{238,32,77,255}

var C_Path Color = Color{254,199,22,255}
var C_Tree Color = Color{31,117,254,255}

func run(state chan State) int {

  // Initilize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to initilize SDL: %s\n", err)
    return 1
	}
	defer sdl.Quit()

  // Create Window
	window, err := sdl.CreateWindow(
    Win_Title,
    sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
    Win_W, Win_H,
    sdl.WINDOW_SHOWN | sdl.WINDOW_RESIZABLE,
  );
	if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
    return 2
	}
	defer window.Destroy()

  // Create Renderer
  ren, err := sdl.CreateRenderer(
    window, -1,
    sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC,
  );
	if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
    return 3
	}
  defer ren.Destroy()

  ren.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

  //-- Draw Recieved State
  for s := range state {

    // set grid size
    ren.SetLogicalSize(s.grid_w * Pix_D, s.grid_h * Pix_D);

    // check for window resize
    sdl.PollEvent()

    // set background color
    ren.SetDrawColor(C_BG.r, C_BG.g, C_BG.b, C_BG.a);
    ren.Clear();

    // draw apple
    draw_point(ren, C_Apple, s.apple, 1)

    // draw snake
    for _, p := range s.snake {
      draw_point(ren, C_Snake, p, 1);
    }

    // get max score
    var max_score float64 = 0;
    for _, n := range s.tree {
      h_score := float64(n.fScore - n.gScore)
      if h_score > max_score {
        max_score = h_score;
      }
    }

    // draw tree
    for _, n := range s.tree {
      score := n.fScore - n.gScore
      var shade float64;
      if score <= 0 {
        shade = 0;
      } else {
        shade = 1 / (max_score / float64(score + 1));
      }
      color := C_Tree
      shaded := Color{
        r: uint8(float64(color.r) + float64(255 - color.r) * shade),
        g: uint8(float64(color.g) + float64(255 - color.g) * shade),
        b: uint8(float64(color.b) + float64(255 - color.b) * shade),
        a: color.a,
      }
      draw_point(ren, shaded, n.p, 3);
    }

    // draw path
    for _, n := range s.path {
      score := n.fScore - n.gScore
      var shade float64;
      if score <= 0 {
        shade = 0;
      } else {
        shade = 1 / (max_score / float64(score + 1));
      }
      color := C_Path
      shaded := Color{
        r: uint8(float64(color.r) + float64(255 - color.r) * shade),
        g: uint8(float64(color.g) + float64(255 - color.g) * shade),
        b: uint8(float64(color.b) + float64(255 - color.b) * shade),
        a: color.a,
      }
      draw_point(ren, shaded, n.p, 3);
    }

	  // render
    ren.Present();
	}

  return 0;
}

func draw_point(ren *sdl.Renderer, c Color, p Point, delta int32) {
  x := p.x * Pix_D + delta
  y := p.y * Pix_D + delta
  d := Pix_D - delta * 2
  rect := sdl.Rect{x, y, d, d}
  ren.SetDrawColor(c.r, c.g, c.b, c.a);
  ren.FillRect(&rect);
}

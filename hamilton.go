package main;

import (
//  "math"
//  "fmt"
//  "os"
)

func hamilton(state_recv <-chan State, action chan<- Action) {
  state := <- state_recv
  height := state.grid_h
  width  := state.grid_w

  circuit := make(map[Point]Point)

  //-- Create Circuit

  // create inner circuit
  var cur Point
  for y := int32(0); y < height; y++ {
    parity := y % 2;
    for x := int32(0); x < width; x++ {
      cur = Point{x, y};
      if (x == (width - 1)) && (parity == 0) {
        circuit[cur] = Point{x, y + 1}
        continue;
      }
      if (x == 1) && (parity != 0) {
        circuit[cur] = Point{x, y + 1}
        continue;
      }
      if parity == 0 {
        circuit[cur] = Point{x + 1, y}
      } else {
        circuit[cur] = Point{x - 1, y}
      }
    }
  }

  // connect end->start
  cur = Point{1, height-1};
  for y := height-1; y >= 0; y-- {
    circuit[cur] = Point{0, y}
    cur = Point{0, y};
  }

  //-- Debug Circuit Output
  /*
  for y := int32(0); y < height; y++ {
    for x := int32(0); x < width; x++ {
      k := Point{x, y};
      v := circuit[k];
      fmt.Printf("%+v -> %+v\n", k, v)
    }
  }
  */

  // Validate Circuit
  /*
  line := 0;
  start := Point{0, 0}
  k := start
  fmt.Println()
  for {
    v := circuit[k];
    fmt.Printf("%+v ->", k)
    k  = v;
    if k == start {
      fmt.Printf("%+v\n", k)
      fmt.Println("Cycled!")
      break;
    }
    line += 1;
    if line == 4 {
      fmt.Println();
      line = 0;
    }
  }
  */

  //-- Determine snake's location on circuit and lock it on
  head := state.snake[len(state.snake) - 1];
  next := circuit[head];
  dir  := get_dir(head, next);
  action <- Action{dir: dir}
  for state := range state_recv {
    head = state.snake[len(state.snake) - 1];
    next = circuit[head];
    dir  = get_dir(head, next);
    //fmt.Printf("%+v -> %+v : %#v\n", head, next, dir)
    action <- Action{dir: dir}
  }
}

func get_dir(p1, p2 Point) Direction {
  dx := p1.x - p2.x;
  dy := p1.y - p2.y;
  p  := Point{dx, dy}

  north_p := Point{ 0,  1};
  south_p := Point{ 0, -1};
  east_p  := Point{ 1,  0};
  west_p  := Point{-1,  0}

  if p == north_p {
    return North;
  }
  if p == south_p {
    return South;
  }
  if p == east_p {
    return East;
  }
  if p == west_p {
    return West;
  }
  return -1;
}

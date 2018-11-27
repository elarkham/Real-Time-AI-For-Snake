package main;

import (
  "os"
  "flag"
  "fmt"
  "time"
  "math/rand"

  "github.com/veandco/go-sdl2/sdl"
)
type Direction int

const (
  North Direction = iota
  South
  East
  West
)

type Point struct {
  x int32
  y int32
}

type State struct {
  dir     Direction
  apple   Point
  snake []Point
  path  []*Node
  tree  []*Node

  // Settings
  player string
  debug  bool
  json   bool

  grid_h int32
  grid_w int32

  seed int

  // Game Stats
  total_actions  int32
  total_timeouts int32

  start_time time.Time
  end_time   time.Time
}

type Action struct {
  dir    Direction
  path []*Node
  tree []*Node
  quit   bool
}

func main() {
  //-- Read Flags
  var grid_h, grid_w, seed int
  flag.IntVar(&grid_h, "h", 50, "Grid Height")
  flag.IntVar(&grid_w, "w", 50, "Grid Width")
  flag.IntVar(&seed, "seed", 1, "Game Seed")

  var player string
  flag.StringVar(&player, "player", "keyboard", "Player of the game")

  delay := time.Duration(50) * time.Millisecond;
  flag.DurationVar(&delay, "delay", delay, "Game speed")

  var debug, json bool
  flag.BoolVar(&debug, "debug", false, "Enables debug mode")
  flag.BoolVar(&json,  "json",  false, "Dumps json stats")

  flag.Parse();

  //-- Initialize State
  s := State{
    // Settings
    grid_h: int32(grid_h),
    grid_w: int32(grid_w),

    player: player,
    seed:   seed,

    debug:  debug,
    json:   json,

    // States
    start_time: time.Now(),
  }

  //-- Setup Render Thread
  draw := make(chan State)
  go run(draw)

  //-- Setup Action Thread
  action  := make(chan Action)
  publish := make(chan State)
  if player == "keyboard" {
    go keyevents(action)
  } else if player == "rta" {
    go rta(publish, action);
  } else if player == "astar" {
    go astar(publish, action, h_dist);
  } else if player == "astar-adv" {
    go astar(publish, action, h_tail);
  } else if player == "hamilton" {
    go hamilton(publish, action);
  } else {
    fmt.Fprintln(os.Stderr, "Invalid Player")
    os.Exit(1);
  }

  //-- Initialize Game Logic
  rand.Seed(int64(seed))

  // Start Point
  y := int32(grid_w/2 + 3)
  x := int32(grid_w/2)
  s.dir = North

  // Create Snake
  s.snake = make([]Point, 0)
  for i := 0; i < 3; i++ {
    point := Point{x, y - int32(i)}
    s.snake = append(s.snake, point)
  }

  // Generate Apple
  s.apple = rand_apple(s)

  // Submit Initial
  draw <- s
  if player != "keyboard" {
    publish <- s
  }
  sdl.Delay(uint32(delay.Seconds() * 1000));

  // Main Loop
  run := true
  for run {
    //-- Action
    var new_dir Direction;
    var path []*Node;
    var tree []*Node;
    select {
    case a := <-action:
      new_dir = a.dir
      path    = a.path
      tree    = a.tree
      run     = !a.quit
      s.total_actions += 1;
      break;
    case <-time.After(time.Millisecond * 1):
      if debug {
        fmt.Fprintln(os.Stderr, "Timeout");
      }
      new_dir = s.dir
      path    = s.path
      s.total_timeouts += 1;
    }

    // prevent snake from moving in opposite direction
    if new_dir != flip_dir(s.dir) {
      s.dir = new_dir;
    }

    // analysis of algorithm work
    if debug {
      s.path = path
      s.tree = tree
    }

    //-- Move
    x, y = next_coord(s.dir, s.snake[len(s.snake) - 1]);
    head := Point{x, y}
    s.snake = append(s.snake, head)

    // check if game won
    score := int32(len(s.snake))
    if score == (s.grid_h * s.grid_w) {
      draw <- s
      break;
    }

    // check if apple was consumed
    if (x == s.apple.x) && (y == s.apple.y) {
      s.apple = rand_apple(s);
      fmt.Println("Score: ", score)
    } else {
      s.snake = s.snake[1:]
    }

    // check if state is valid
    if invalid_move(s) {
      if debug {
        fmt.Fprintln(os.Stderr, "Invalid Move To", head)
      }
      run = false;
    }

	  //-- Render
    draw <- s
    if player != "keyboard" {
      select {
      case publish <- s:
        break;
      case <-time.After(time.Millisecond * 1):
        if debug {
          fmt.Fprintln(os.Stderr, "Publish Timeout")
        }
      }
    }

    // Delay
    sdl.Delay(uint32(delay.Seconds() * 1000));
	}
  s.end_time = time.Now();
  print_results(s);
}

func print_results(s State) {
  stats := get_stats(s)
  print_stats(stats);
  if s.json {
    dump_json(stats);
  }
  sdl.Delay(1000 * 3);
}

func invalid_move(s State) bool {
  h := s.snake[len(s.snake) - 1];

  // Check if inside grid bounds
  if (h.x < 0) || (h.x >= s.grid_w) {
    return true;
  }
  if (h.y < 0) || (h.y >= s.grid_h) {
    return true;
  }

  // Check for tail collision
  for i, p := range s.snake {
    if i == (len(s.snake) - 1) {
      continue;
    }
    if (p.x == h.x) && (p.y == h.y) {
      return true;
    }
  }

  return false;
}

func flip_dir(dir Direction) Direction {
  if (dir == North) {
    return South;
  }
  if (dir == South) {
    return North;
  }
  if (dir == East) {
    return West;
  }
  if (dir == West) {
    return East;
  }
  return -1;
}

func rand_apple(s State) Point {
  x := int32(rand.Intn(int(s.grid_w)));
  y := int32(rand.Intn(int(s.grid_h)));

  // Prevent apple from spawning inside snake
  for _, p := range s.snake {
    if (p.x == x) && (p.y == y) {
      return rand_apple(s);
    }
  }
  return Point{x, y}
}

func next_coord(dir Direction, p Point) (int32, int32) {
  switch dir {
    case North:
      p.y -= 1;
    case South:
      p.y += 1;
    case East:
      p.x -= 1;
    case West:
      p.x += 1;
  }
  return p.x, p.y;
}

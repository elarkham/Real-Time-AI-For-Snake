package main;

import (
  "fmt"
//  "math/rand"
)

func astar( state_recv <-chan  State,
            action     chan<-  Action,
            h func(State, *Node) (int32) ) {
Outer:
  for state := range state_recv {
    // Start A* search
    start := new(Node)
    start.p   = state.snake[len(state.snake)-1]
    start.dir = state.dir

    path, tree := astar_search(state, start, state.apple, h)
    if path == nil {
      // Stall for time
      fmt.Println("Stalling")
      cur := new(Node)
      cur.p   = state.snake[len(state.snake)-1]
      cur.dir = state.dir

      children := children(cur, state);
      if len(children) <= 0 {
        fmt.Println("Search Failed, Game Over");
        action <- Action{dir: North, quit: true}
        return;
      }
      _, best := best_node(children);
      action <- Action{dir: best.dir, quit: false}
      continue;
    }
    // Iterate Over Discovered Path
    for i, node := range path[1:] {
      action <- Action{
          dir: node.dir,
          path: path,
          tree: tree,
          quit: false,
        }
      if i == (len(path[1:]) - 1) {
        continue Outer;
      }
      state = <- state_recv
    }
  }
}


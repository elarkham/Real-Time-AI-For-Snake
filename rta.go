package main;

import (
//  "math"
  "fmt"
//  "os"
)

func rta(state_recv <-chan State, action chan<- Action) {
  for state := range state_recv {
    visted := make(map[Point]bool)

    cur := new(Node)
    cur.p   = state.snake[len(state.snake) - 1]
    cur.dir = state.dir

    visted[cur.p] = true;
    cur.gScore = 0;
    cur.fScore = mdist(cur.p, state.apple) + h_tail(state, cur);

    children := children(cur, state);
    for _, child := range children {
      if visted[child.p] {
        continue;
      }
      visted[child.p] = true;
      child.parent = cur;
      child.gScore = cur.gScore + mdist(cur.p, child.p);
      child.fScore = child.gScore + h_tail(state, child);
    }
    if len(children) <= 0 {
      fmt.Println("Search Failed")
      action <- Action{dir: East, quit: true}
      continue;
    }
    _, best := best_node(children);
    action <- Action{dir: best.dir, quit: false, path: children}
  }
}



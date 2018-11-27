package main;

import (
  "math"
  "fmt"
  "os"
)

type Node struct {
  p   Point
  dir Direction
  fScore int32
  gScore int32
  parent *Node
}

func (n *Node) String() string {
  return fmt.Sprintf("%+v", *n);
}

const MAX_INT32 = int32((^uint32(0)) >> 1)

func astar_search( state State,
                   start *Node,
                   goal Point,
                   h func(State, *Node) (int32),
                 ) ([]*Node, []*Node) {

  visited := make([]*Node, 1)
  open   := make([]*Node, 0)
  closed := make(map[Point]bool)
  gScore := make(map[Point]int32)

  start.gScore = 0;
  start.fScore = mdist(start.p, goal) + h(state, start);
  visited[0] = start;

  open = append(open, start);

  for len(open) > 0 {
    i, cur := best_node(open);
    if cur.p == goal {
      path := make([]*Node, 0);
      path = append(path, cur)
      path = reverse(get_path(path));
      return path, visited;
    }
    open = append(open[:i], open[i+1:]...);
    closed[cur.p] = true;

    children := children(cur, state);
    for _, child := range children {
      if closed[child.p] {
        continue;
      }
      if !contains(open, child) {
        open = append(open, child);
      }
      tmp_gScore := cur.gScore + mdist(cur.p, child.p);
      if curScore, ok := gScore[child.p]; ok {
        if tmp_gScore >= curScore {
          continue;
        }
      }
      visited = append(visited, child);

      child.parent = cur
      child.gScore = tmp_gScore;
      child.fScore = child.gScore + h(state, child);
    }
  }

  return nil, nil
}

func h_dist(state State, node *Node) int32 {
  return mdist(node.p, state.apple);
}

func h_tail(state State, node *Node) int32 {
  tail := state.snake[0];
  h1 := mdist(node.p, state.apple) * 4
  h2 := mdist(node.p, tail) * 3

  h3 := mdist(node.p, Point{node.p.x, state.grid_h - 1}) / 7
  h4 := mdist(node.p, Point{node.p.x, -1}) / 7
  h5 := mdist(node.p, Point{state.grid_w - 1, node.p.y}) / 7
  h6 := mdist(node.p, Point{-1, node.p.y}) / 7

  return 0 +h1 +h2 -h3 -h4 -h5 -h6;
}

func mdist(start, goal Point) int32 {
  x := math.Abs(float64(start.x - goal.x))
  y := math.Abs(float64(start.y - goal.y))
  return int32(x + y)
}

func get_path(path []*Node) []*Node {
  child := path[len(path)-1];
  if child.parent != nil {
    return get_path(append(path, child.parent));
  }
  return path;
}

func children(node *Node, s State) []*Node {
  children  := make([]*Node, 0);
  child_dir := child_dir(node.dir)

  var p Point;
  for _, dir := range child_dir {
    c := node.p
    switch dir {
      case North:
        p.y = c.y - 1;
        p.x = c.x;
        break;
      case South:
        p.y = c.y + 1;
        p.x = c.x;
        break;
      case East:
        p.y = c.y;
        p.x = c.x - 1;
        break;
      case West:
        p.y = c.y;
        p.x = c.x + 1;
        break;
      default:
        continue;
    }
    if !bad_move(p, s) {
      child := new(Node)
      child.p = p
      child.dir = dir
      child.fScore = MAX_INT32

      children = append(children, child);
    }
  }
  return children
}

func bad_move(h Point, s State) bool {
  // Check if inside grid bounds
  if (h.x < 0) || (h.x >= s.grid_w) {
    return true;
  }
  if (h.y < 0) || (h.y >= s.grid_h) {
    return true;
  }

  // Check for tail collision
  for _, p := range s.snake {
    if (p.x == h.x) && (p.y == h.y) {
      return true;
    }
  }

  return false
}

func child_dir(dir Direction) []Direction {
  switch dir {
    case North:
      return []Direction{North, East, West};
    case South:
      return []Direction{South, East, West};
    case East:
      return []Direction{East, North, South};
    case West:
      return []Direction{West, North, South};
    default:
      os.Exit(1);
  }
  return []Direction{};
}

func best_node(open []*Node) (int, *Node) {
  best   := open[0];
  best_i := 0
  for i, node := range open {
    if node.fScore < best.fScore {
      best   = node;
      best_i = i;
    }
  }
  return best_i, best;
}

func is_goal(node *Node, state State) bool {
  p1, p2 := node.p, state.apple
  return (p1.x == p2.x) && (p1.y == p2.y)
}

func contains(nodes []*Node, node *Node) bool {
  for _, n := range nodes {
    if n.p == node.p {
      return true;
    }
  }
  return false;
}

func reverse(a []*Node) []*Node {
  for l, r := 0, len(a)-1; l < r; l, r = l+1, r-1 {
	  a[l], a[r] = a[r], a[l]
  }
  return a;
}

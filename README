Dependencies:
 This project requires this sdl bindings library for go:
 https://github.com/veandco/go-sdl2

 It can be install via:
 `go get -v github.com/veandco/go-sdl2/sdl`

Building:
  All you should have to do is run:
     1. "go build"
  OR 2. "go build snake"

Running:
  Program provides a usage summary if the --help flag is passed

Notes:
  - The delay flag controls how long a tick is and how difficult
    the game will be as a result

  - Game always defaults to seed 1 if you dont manually give it one

  - The debug flag will visualize the pathing for A* and RTA* but
    doesn't display anything for hamilton since it wouldn't be
    useful

  - Json flag generates json output for easy automated testing,
    I intended to do more with this but didn't have the time.

  - By default the game is user controlled, if you want to run
    a specific algorithm you must pass the `-player` flag
    followed by one of the following algoritm names:
    [ astar, astar-adv, hamilton, rta ]

  - astar-adv uses a complicated heuristic that makes paths close
    to the snakes tail more favorable, astar on the otherhand
    just uses manhattan distance as the heuristic. In the case of
    rta, it just uses the advanced heuristic in all cases since
    it always performed better.

package main;

import (
  "encoding/json"
  "io/ioutil"
  "time"
  "fmt"
)

type Stats struct {
  Score float64 `json:"score"`

  Grid_h int32  `json:"grid_h"`
  Grid_w int32  `json:"grid_w"`

  Board_area  float64 `json:"board_area"`
  Board_ratio float64 `json:"board_ratio"`

  Avg_actions_per_apple  float64 `json:"avg_actions_per_apple"`
  Avg_timeouts_per_apple float64 `json:"avg_timeouts_per_apple"`

  Game_duration  time.Duration `json:"game_duration"`
  Avg_apple_time time.Duration `json:"avg_apple_time"`

  Player string `json:"player"`
  Seed   int    `json:"seed"`
  Win    bool   `json:"win"`
}

func get_stats(s State) Stats {
  stats := Stats{}

  stats.Score = float64(len(s.snake))

  stats.Grid_h = s.grid_h;
  stats.Grid_w = s.grid_w;

  stats.Board_area  = float64(s.grid_h * s.grid_w)
  stats.Board_ratio = stats.Score / stats.Board_area

  stats.Avg_actions_per_apple  = float64(s.total_actions) / stats.Score
  stats.Avg_timeouts_per_apple = float64(s.total_timeouts) / stats.Score

  stats.Game_duration  = s.end_time.Sub(s.start_time)

  avg_apple_time_int  := stats.Game_duration.Seconds() / stats.Score
  stats.Avg_apple_time = time.Duration(avg_apple_time_int * 1000) * time.Millisecond

  stats.Player = s.player
  stats.Seed = s.seed
  stats.Win = (stats.Board_ratio == 1.0)

  return stats
}

func dump_json(s Stats) {
  filename := fmt.Sprintf("results-%v-%v.json", s.Player, s.Seed)
  json, _ := json.MarshalIndent(s, "  ", "");
  ioutil.WriteFile(filename, json, 0644)
}

func print_stats(s Stats) {
  fmt.Println("---")
  if s.Win {
    fmt.Println("Game Completed")
  } else {
    fmt.Println("Game Over")
  }
  fmt.Println("Score:", s.Score);
  fmt.Printf("Board Ratio: %v%%\n", s.Board_ratio * 100);
  fmt.Println("Player:", s.Player)
  fmt.Println("Seed:", s.Seed)
  fmt.Println("Duration:", s.Game_duration)
  fmt.Println("---")
  fmt.Println("Average Actions Per Apple: ", s.Avg_actions_per_apple)
  fmt.Println("Average Timeouts Per Apple:", s.Avg_timeouts_per_apple)
  fmt.Println("Average Apple Time:", s.Avg_apple_time)
}



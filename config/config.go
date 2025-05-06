package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
    Laps        int
    LapLen      int
    PenaltyLen  int
    FiringLines int
    Start       string
    StartDelta  time.Duration
}

func Load(filename string) (Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return Config{}, fmt.Errorf("failed to read config: %w", err)
    }
    var raw struct {
        Laps        int    `json:"laps"`
        LapLen      int    `json:"lapLen"`
        PenaltyLen  int    `json:"penaltyLen"`
        FiringLines int    `json:"firingLines"`
        Start       string `json:"start"`
        StartDelta  string `json:"startDelta"`
    }
    if err := json.Unmarshal(data, &raw); err != nil {
        return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    parts := strings.Split(raw.StartDelta, ":")
    if len(parts) != 3 {
        return Config{}, fmt.Errorf("invalid StartDelta format: %s", raw.StartDelta)
    }
    hours, err := strconv.Atoi(parts[0])
    if err != nil {
        return Config{}, fmt.Errorf("invalid hours in StartDelta: %w", err)
    }
    minutes, err := strconv.Atoi(parts[1])
    if err != nil {
        return Config{}, fmt.Errorf("invalid minutes in StartDelta: %w", err)
    }
    seconds, err := strconv.Atoi(parts[2])
    if err != nil {
        return Config{}, fmt.Errorf("invalid seconds in StartDelta: %w", err)
    }
    startDelta := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second

    return Config{
        Laps:        raw.Laps,
        LapLen:      raw.LapLen,
        PenaltyLen:  raw.PenaltyLen,
        FiringLines: raw.FiringLines,
        Start:       raw.Start,
        StartDelta:  startDelta,
    }, nil
}
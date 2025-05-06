package entity

import "time"

const (
	Registered = iota
	Started
	Firing
	Penalty
	Finished
	Disqualified
	CannotContinue
)

type Competitor struct {
	ID           string
	Status       int
	Reason       string
	StartPlanned time.Time
	StartActual  time.Time
	LapTimes     []time.Duration
	PenaltyTime  time.Duration
	Shots        int
	Hits         int
	LastLapStart time.Time
	PenaltyStart time.Time
	FiringRange  string
}

package entity

import "time"

const ()

type Competitor struct {
    ID              string
    Registered      bool
    Started         bool
    Finished        bool
    Disqualified    bool
    CannotContinue  bool
    Reason          string
    StartPlanned    time.Time
    StartActual     time.Time
    LapTimes        []time.Duration
    PenaltyTime     time.Duration
    Shots           int
    Hits            int
    LastLapStart    time.Time
    PenaltyStart    time.Time 
    FiringRange     string 
}
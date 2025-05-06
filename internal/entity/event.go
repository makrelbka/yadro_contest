package entity

import "time"

type Event struct {
    Time         time.Time
    ID           int
    CompetitorID string
    Extra        []string
}
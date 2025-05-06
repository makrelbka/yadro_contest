package library

import (
    "fmt"
    "time"
    "yadro/internal/entity"
    "yadro/internal/usecase/repository"
)

type CompetitorService struct {
    repo *inmemory.InMemoryRepository
}

func NewCompetitorService(repo *inmemory.InMemoryRepository) *CompetitorService {
    return &CompetitorService{repo: repo}
}

func (s *CompetitorService) GetCompetitor(id string) (*entity.Competitor, bool) {
    return s.repo.GetCompetitor(id)
}

func (s *CompetitorService) RegisterCompetitor(c *entity.Competitor) {
    c.Registered = true
    s.repo.CreateCompetitor(c)
}

func (s *CompetitorService) SetStartTime(c *entity.Competitor, t time.Time) {
    c.StartPlanned = t
    s.repo.UpdateCompetitor(c)
}

func (s *CompetitorService) StartCompetitor(c *entity.Competitor, actualTime time.Time) error {
    c.StartActual = actualTime
    if s.repo.Cfg.StartDelta == 0 {
        c.Disqualified = true
        return fmt.Errorf("invalid start delta for competitor(%s)", c.ID)
    }
    if c.StartActual.Sub(c.StartPlanned) > s.repo.Cfg.StartDelta {
        c.Disqualified = true
    } else {
        c.Started = true
        c.LastLapStart = actualTime
    }
    s.repo.UpdateCompetitor(c)
    return nil
}

func (s *CompetitorService) EnterFiringRange(c *entity.Competitor, rangeID string) {
    c.FiringRange = rangeID
    c.Shots = 5
    s.repo.UpdateCompetitor(c)
}

func (s *CompetitorService) HitTarget(c *entity.Competitor) {
    c.Hits++
    s.repo.UpdateCompetitor(c)
}

func (s *CompetitorService) LeaveFiringRange(c *entity.Competitor) {
    c.FiringRange = ""
    s.repo.UpdateCompetitor(c)
}

func (s *CompetitorService) EnterPenaltyLaps(c *entity.Competitor, startTime time.Time) {
    c.PenaltyStart = startTime
    s.repo.UpdateCompetitor(c)
}

func (s *CompetitorService) LeavePenaltyLaps(c *entity.Competitor, endTime time.Time) {
    c.PenaltyTime += endTime.Sub(c.PenaltyStart)
    s.repo.UpdateCompetitor(c)
}

func (s *CompetitorService) EndLap(c *entity.Competitor, lapEndTime time.Time) {
    c.LapTimes = append(c.LapTimes, lapEndTime.Sub(c.LastLapStart))
    c.LastLapStart = lapEndTime
    if len(c.LapTimes) == s.repo.Cfg.Laps {
        c.Finished = true
    }
    s.repo.UpdateCompetitor(c)
}

func (s *CompetitorService) MarkCannotContinue(c *entity.Competitor, reason string) {
    c.CannotContinue = true
    c.Reason = reason
    s.repo.UpdateCompetitor(c)
}
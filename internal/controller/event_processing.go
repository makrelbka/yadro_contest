package controller

import (
    "fmt"
    "strings"
    "time"
    "yadro/internal/entity"
    "yadro/internal/usecase/library"
)

const parseTimeConst = "15:04:05.000"

type EventProcessor struct {
    service *library.CompetitorService
}

func NewEventProcessor(service *library.CompetitorService) *EventProcessor {
    return &EventProcessor{service: service}
}

func (p *EventProcessor) Process(e entity.Event, logs *[]string) {
    c, exists := p.service.GetCompetitor(e.CompetitorID)
    if !exists {
        c = &entity.Competitor{ID: e.CompetitorID}
        p.service.RegisterCompetitor(c)
    }

    if c.Disqualified || c.CannotContinue || (!c.Started && e.ID > 4) {
        return
    }

    switch e.ID {
    case 1:
        p.service.RegisterCompetitor(c)
        *logs = append(*logs, fmt.Sprintf("[%s] The competitor(%s) registered", e.Time.Format(parseTimeConst), c.ID))
    case 2:
        t, err := time.Parse(parseTimeConst, e.Extra[0])
        if err != nil {
            *logs = append(*logs, fmt.Sprintf("[%s] Invalid start time for competitor(%s)", e.Time.Format(parseTimeConst), c.ID))
            return
        }
        p.service.SetStartTime(c, t)
        *logs = append(*logs, fmt.Sprintf("[%s] The start time for the competitor(%s) was set by a draw to %s", e.Time.Format(parseTimeConst), c.ID, t.Format(parseTimeConst)))
    case 3:
        *logs = append(*logs, fmt.Sprintf("[%s] The competitor(%s) is on the start line", e.Time.Format(parseTimeConst), c.ID))
    case 4:
        err := p.service.StartCompetitor(c, e.Time)
        if err != nil {
            *logs = append(*logs, fmt.Sprintf("[%s] %v", e.Time.Format(parseTimeConst), err))
            return
        }
        *logs = append(*logs, fmt.Sprintf("[%s] The competitor(%s) has started", e.Time.Format(parseTimeConst), c.ID))
    case 5:
        p.service.EnterFiringRange(c, e.Extra[0])
        *logs = append(*logs, fmt.Sprintf("[%s] The competitor(%s) is on the firing range(%s)", e.Time.Format(parseTimeConst), c.ID, e.Extra[0]))
    case 6:
        p.service.HitTarget(c)
        *logs = append(*logs, fmt.Sprintf("[%s] The target(%s) has been hit by competitor(%s)", e.Time.Format(parseTimeConst), e.Extra[0], c.ID))
    case 7:
        p.service.LeaveFiringRange(c)
        *logs = append(*logs, fmt.Sprintf("[%s] The competitor(%s) left the firing range", e.Time.Format(parseTimeConst), c.ID))
    case 8:
        p.service.EnterPenaltyLaps(c, e.Time)
        *logs = append(*logs, fmt.Sprintf("[%s] The competitor(%s) entered the penalty laps", e.Time.Format(parseTimeConst), c.ID))
    case 9:
        p.service.LeavePenaltyLaps(c, e.Time)
        *logs = append(*logs, fmt.Sprintf("[%s] The competitor(%s) left the penalty laps", e.Time.Format(parseTimeConst), c.ID))
    case 10:
        p.service.EndLap(c, e.Time)
        *logs = append(*logs, fmt.Sprintf("[%s] The competitor(%s) ended the main lap", e.Time.Format(parseTimeConst), c.ID))
    case 11:
        p.service.MarkCannotContinue(c, strings.Join(e.Extra, " "))
        *logs = append(*logs, fmt.Sprintf("[%s] The competitor(%s) can`t continue: %s", e.Time.Format(parseTimeConst), c.ID, strings.Join(e.Extra, " ")))
    }
}
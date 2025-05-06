package controller

import (
    "fmt"
    "sort"
    "time"
    "yadro/config"
    "yadro/internal/entity"
)

func formatDuration(d time.Duration) string {
    hours := int(d.Hours())
    minutes := int(d.Minutes()) % 60
    seconds := d.Seconds() - float64(hours*3600+minutes*60)
    return fmt.Sprintf("%02d:%02d:%06.3f", hours, minutes, seconds)
}

func calculateTotalTime(c *entity.Competitor) time.Duration {
    total := c.StartActual.Sub(c.StartPlanned)
    for _, lap := range c.LapTimes {
        total += lap
    }
    total += c.PenaltyTime
    return total
}

func formatLaps(c *entity.Competitor, cfg config.Config) string {
    lapsStr := "["
    for i, lap := range c.LapTimes {
        if i > 0 {
            lapsStr += ", "
        }
        speed := float64(cfg.LapLen) / (lap.Seconds())
        lapsStr += fmt.Sprintf("{%s, %.3f}", formatDuration(lap), speed)
    }
    for i := len(c.LapTimes); i < cfg.Laps; i++ {
        if i > 0 || len(c.LapTimes) > 0 {
            lapsStr += ", "
        }
        lapsStr += "{,}"
    }
    lapsStr += "]"
    return lapsStr
}

func formatPenalty(c *entity.Competitor, cfg config.Config) string {
    penaltyCount := c.Shots - c.Hits 
    penaltySpeed := 0.0
    if c.PenaltyTime > 0 && penaltyCount > 0 {
        penaltySpeed = float64(cfg.PenaltyLen*penaltyCount) / c.PenaltyTime.Seconds()
    }
    return fmt.Sprintf("{%s, %.3f}", formatDuration(c.PenaltyTime), penaltySpeed)
}

func GenerateFinalReport(cfg config.Config, competitors map[string]*entity.Competitor) []string {
    compList := make([]*entity.Competitor, 0, len(competitors))
    for _, c := range competitors {
        compList = append(compList, c)
    }

    sort.Slice(compList, func(i, j int) bool {
        return compareCompetitors(compList[i], compList[j])
    })

    var report []string
    for _, c := range compList {
        status := determineStatus(c)
        if status == "Finished" {
            status = formatDuration(calculateTotalTime(c))
        }

        lapsStr := formatLaps(c, cfg)
        penaltyStr := formatPenalty(c, cfg)

        report = append(report, fmt.Sprintf("[%s] %s %s %s %d/%d", status, c.ID, lapsStr, penaltyStr, c.Hits, c.Shots))
    }
    return report
}

func compareCompetitors(ci, cj *entity.Competitor) bool {
    if ci.Disqualified != cj.Disqualified {
        return !ci.Disqualified
    }
    if ci.CannotContinue != cj.CannotContinue {
        return !ci.CannotContinue
    }
    if ci.Finished != cj.Finished {
        return ci.Finished
    }
    return calculateTotalTime(ci) < calculateTotalTime(cj)
}

func determineStatus(c *entity.Competitor) string {
    if c.Disqualified {
        return "NotStarted"
    }
    if c.CannotContinue {
        return "NotFinished"
    }
    if c.Finished {
        return "Finished"
    }
    return "Running"
}
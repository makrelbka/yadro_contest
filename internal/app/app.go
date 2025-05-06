package app

import (
    "fmt"

    "yadro/config"
    processing "yadro/internal/controller"
    "yadro/internal/usecase/library"
    "yadro/internal/usecase/repository"
)

func Run() error {
    cfg, err := config.Load("data/config.json")
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    fmt.Println(cfg)

    repo := inmemory.NewInMemoryRepository(cfg)

    service := library.NewCompetitorService(repo)

    processor := processing.NewEventProcessor(service)

    var logs []string

    err = processing.ProcessEvents("data/events.txt", &logs, processor)
    if err != nil {
        return fmt.Errorf("failed to parse events: %w", err)
    }

    for _, log := range logs {
        fmt.Println(log)
    }

    final := processing.GenerateFinalReport(cfg, repo.Competitors)
    fmt.Println("\nFinal Report:")
    for _, line := range final {
        fmt.Println(line)
    }

    return nil
}
package inmemory

import (
	"errors"
	"yadro/config"
	"yadro/internal/entity"
)

type InMemoryRepository struct {
    Competitors map[string]*entity.Competitor
    Cfg         config.Config
}

func NewInMemoryRepository(cfg config.Config) *InMemoryRepository {
    return &InMemoryRepository{
        Competitors: make(map[string]*entity.Competitor),
        Cfg:         cfg,
    }
}

func (r *InMemoryRepository) CreateCompetitor(c *entity.Competitor) error {
    if _, exists := r.Competitors[c.ID]; exists {
        return errors.New("competitor already exists")
    }
    r.Competitors[c.ID] = c
    return nil
}

func (r *InMemoryRepository) UpdateCompetitor(c *entity.Competitor) error {
    if _, exists := r.Competitors[c.ID]; !exists {
        return errors.New("competitor does not exist")
    }
    r.Competitors[c.ID] = c
    return nil
}

func (r *InMemoryRepository) GetCompetitor(id string) (*entity.Competitor, bool) {
    c, ok := r.Competitors[id]
    return c, ok
}
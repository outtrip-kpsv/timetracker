package bl

import (
  "timetracker/internal/bl/repo"
  "timetracker/internal/db"
)

type BL struct {
  People repo.IPeopleBL
  Task   repo.ITaskBL
}

func New(db *db.DbRepo) *BL {
  return &BL{
    People: repo.NewPeopleBL(db),
    Task:   repo.NewTaskBL(db),
  }
}

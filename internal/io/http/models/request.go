package models

import "timetracker/internal/bl/dto"

type TaskCreate struct {
  IdPerson string
  Name     string `json:"task_name"`
}

func (t *TaskCreate) ToDto() *dto.Task {
  if t == nil {
    return nil
  }

  return &dto.Task{
    IdPerson: t.IdPerson,
    TaskName: t.Name}
}

type DateStartEnd struct {
  Start string `json:"start"`
  End   string `json:"end"`
}

package models

import (
  "fmt"
  "timetracker/internal/bl/dto"
  "timetracker/internal/utils/var/endpoint"
)

type ErrorResponse struct {
  Error string `json:"error"`
}

type Ok struct {
  Msg    string `json:"msg"`
  Status uint   `json:"status"`
}

type TaskResp struct {
  IdTask     string  `json:"id_task"`
  IdPerson   string  `json:"id_person"`
  TaskName   string  `json:"task_name,omitempty"`
  TaskStatus string  `json:"task_status"`
  Urls       UrlTask `json:"urls"`
}

type UrlTask struct {
  Start    string `json:"start_task"`
  Pause    string `json:"pause_task"`
  Complete string `json:"complete_task"`
}

func TaskFromDto(model *dto.Task) *TaskResp {
  if model == nil {
    return nil
  }
  Urls := UrlTask{
    Start:    fmt.Sprintf("http://%s/people/%s/%s/start", endpoint.SRV, model.IdPerson, model.IdTask),
    Pause:    fmt.Sprintf("http://%s/people/%s/%s/pause", endpoint.SRV, model.IdPerson, model.IdTask),
    Complete: fmt.Sprintf("http://%s/people/%s/%s/complete", endpoint.SRV, model.IdPerson, model.IdTask),
  }
  return &TaskResp{
    IdTask:     model.IdTask,
    IdPerson:   model.IdPerson,
    TaskName:   model.TaskName,
    TaskStatus: model.TaskStatus,
    Urls:       Urls,
  }
}

type PeopleResp struct {
  Total  int          `json:"total"`
  Limit  int          `json:"limit"`
  Offset int          `json:"offset"`
  People []dto.Person `json:"people"`
}

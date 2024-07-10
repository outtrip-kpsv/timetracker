package dto

type Task struct {
  IdTask     string `json:"id_task"`
  IdPerson   string `json:"id_person"`
  TaskName   string `json:"task_name"`
  TaskStatus string `json:"task_status"`
}

type TaskTimeResult struct {
  IDTask    string `json:"idtask"`
  TaskName  string `json:"task_name,omitempty"`
  TotalTime string `json:"total_time"`
}

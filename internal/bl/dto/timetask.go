package dto

import "time"

type TimeTask struct {
  ID        int        `json:"id"`
  IDTask    string     `json:"id_task"`
  StartTime time.Time  `json:"start_time"`
  EndTime   *time.Time `json:"end_time"`
}

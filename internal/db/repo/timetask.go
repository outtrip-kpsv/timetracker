package repo

import (
  "context"
  "fmt"
  "github.com/jmoiron/sqlx"
  "time"
  "timetracker/internal/bl/dto"
)

type TimeTask struct {
  ID        int        `json:"id" db:"id"`
  IDTask    string     `json:"id_task" db:"idtask"`
  StartTime time.Time  `json:"start_time" db:"start_time"`
  EndTime   *time.Time `json:"end_time" db:"end_time"`
}

func (t *TimeTask) toDTO() *dto.TimeTask {
  if t == nil {
    return nil
  }

  return &dto.TimeTask{
    ID:        t.ID,
    IDTask:    t.IDTask,
    StartTime: t.StartTime,
    EndTime:   t.EndTime,
  }
}

func (t *TimeTask) fromDTO(model *dto.TimeTask) *TimeTask {
  if model == nil {
    return nil
  }
  t.ID = model.ID
  t.IDTask = model.IDTask
  t.StartTime = model.StartTime
  t.EndTime = model.EndTime
  return t
}

type timeTaskRepo struct {
  db *sqlx.DB
}

type ITimeTaskRepo interface {
  StartTimer(ctx context.Context, id string) error
  StopTimer(ctx context.Context, id string) error
}

func NewTimeTaskRepo(db *sqlx.DB) ITimeTaskRepo {
  return &timeTaskRepo{db: db}
}

func (t *timeTaskRepo) StartTimer(ctx context.Context, id string) error {
  tx, ok := ctx.Value("tx").(*sqlx.Tx)

  query := `INSERT INTO timetask (idtask, start_time) VALUES ($1, NOW())`
  var err error

  if ok {
    _, err = tx.ExecContext(ctx, query, id)
  } else {
    _, err = t.db.ExecContext(ctx, query, id)
  }

  if err != nil {
    return fmt.Errorf("ошибка вставки данных в таблицу timetask: %v", err)
  }

  return nil
}

func (t *timeTaskRepo) StopTimer(ctx context.Context, id string) error {
  query := `UPDATE timetask SET end_time = NOW() WHERE idtask = $1 AND end_time IS NULL`
  result, err := t.db.ExecContext(ctx, query, id)
  if err != nil {
    return fmt.Errorf("ошибка обновления данных в таблице timetask: %v", err)
  }
  rowsAffected, err := result.RowsAffected()
  if err != nil {
    return fmt.Errorf("ошибка получения количества затронутых строк: %v", err)
  }
  if rowsAffected == 0 {
    return fmt.Errorf("нет активного таймера для задачи с id %s", id)
  }
  return nil
}

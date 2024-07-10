package repo

import (
  "context"
  "database/sql"
  "fmt"
  "github.com/jmoiron/sqlx"
  "timetracker/internal/bl/dto"
  "timetracker/internal/utils/const/status"
)

type Task struct {
  IdTask     string `json:"id_task" db:"idtask"`
  IdPerson   string `json:"id_person" db:"idperson"`
  TaskName   string `json:"task_name" db:"task_name"`
  TaskStatus string `json:"task_status" db:"task_status"`
}

func (t *Task) toDTO() *dto.Task {
  if t == nil {
    return nil
  }

  return &dto.Task{
    IdTask:     t.IdTask,
    IdPerson:   t.IdPerson,
    TaskName:   t.TaskName,
    TaskStatus: t.TaskStatus,
  }
}

func (t *Task) fromDTO(model *dto.Task) *Task {
  if model == nil {
    return nil
  }

  t.IdTask = model.IdTask
  t.IdPerson = model.IdPerson
  t.TaskName = model.TaskName
  t.TaskStatus = model.TaskStatus

  return t
}

type taskRepo struct {
  db *sqlx.DB
}

type ITaskRepo interface {
  CreateTask(ctx context.Context, task *dto.Task) (*dto.Task, error)
  GetTaskStatus(ctx context.Context, id string) (string, error)
  UpdateStatus(ctx context.Context, id, st string) error
  TaskTimes(ctx context.Context, id, start, end string) ([]dto.TaskTimeResult, error)
}

func NewTaskRepo(db *sqlx.DB) ITaskRepo {
  return &taskRepo{db: db}
}

func (t *taskRepo) CreateTask(ctx context.Context, task *dto.Task) (*dto.Task, error) {
  query := `INSERT INTO tasks ( idperson, task_name, task_status) VALUES ( :idperson, :task_name, :task_status) RETURNING idtask`
  taskModel := new(Task).fromDTO(task)
  taskModel.TaskStatus = status.New
  fmt.Println(taskModel.IdPerson)

  rows, err := t.db.NamedQueryContext(ctx, query, taskModel)
  if err != nil {
    return nil, fmt.Errorf("ошибка вставки данных в базу: %v", err)
  }
  defer rows.Close()

  if rows.Next() {
    err := rows.Scan(&taskModel.IdTask)
    if err != nil {
      return nil, fmt.Errorf("ошибка сканирования результата: %v", err)
    }
  } else {
    return nil, fmt.Errorf("не удалось вставить задачу")
  }

  return taskModel.toDTO(), nil
}

func (t *taskRepo) GetTaskStatus(ctx context.Context, id string) (string, error) {
  tx, ok := ctx.Value("tx").(*sqlx.Tx)

  var taskStatus string
  query := `SELECT task_status FROM tasks WHERE idtask = $1`
  var err error

  if ok {
    err = tx.GetContext(ctx, &taskStatus, query, id)
  } else {
    err = t.db.GetContext(ctx, &taskStatus, query, id)
  }

  if err == sql.ErrNoRows {
    return "", fmt.Errorf("задача с id %s не найдена", id)
  } else if err != nil {
    return "", fmt.Errorf("ошибка получения статуса задачи из базы данных: %v", err)
  }

  return taskStatus, nil
}
func (t *taskRepo) UpdateStatus(ctx context.Context, id, st string) error {
  tx, ok := ctx.Value("tx").(*sqlx.Tx)
  query := `UPDATE tasks SET task_status = $1 WHERE idtask = $2`
  var result sql.Result
  var err error

  if ok {
    result, err = tx.ExecContext(ctx, query, st, id)
  } else {
    result, err = t.db.ExecContext(ctx, query, st, id)
  }

  if err != nil {
    return fmt.Errorf("ошибка обновления статуса задачи в базе данных: %v", err)
  }

  rowsAffected, err := result.RowsAffected()
  if err != nil {
    return fmt.Errorf("ошибка получения количества затронутых строк: %v", err)
  }

  if rowsAffected == 0 {
    return fmt.Errorf("задача с id %s не найдена", id)
  }

  return nil
}

type TaskTimeCollect []TaskTimeResult

type TaskTimeResult struct {
  IDTask    string `json:"idtask" db:"idtask"`
  TaskName  string `json:"task_name" db:"task_name"`
  TotalTime string `json:"total_time" db:"total_time"`
}

func (tc TaskTimeCollect) toDTO() []dto.TaskTimeResult {
  var res []dto.TaskTimeResult
  for _, result := range tc {

    res = append(res, dto.TaskTimeResult{
      IDTask:    result.IDTask,
      TaskName:  result.TaskName,
      TotalTime: result.TotalTime,
    })
  }
  return res
}

func (t *taskRepo) TaskTimes(ctx context.Context, id, start, end string) ([]dto.TaskTimeResult, error) {
  query := `WITH filtered_times AS (
    SELECT
        t.idtask,
        task.task_name,
        CASE
            WHEN t.start_time < $2 THEN $2
            ELSE t.start_time
        END AS adjusted_start_time,
        CASE
            WHEN t.end_time > $3 THEN $3
            ELSE t.end_time
        END AS adjusted_end_time
    FROM
        timetask t
    JOIN
        tasks task ON t.idtask = task.idtask
    WHERE
        task.idperson = $1
        AND t.end_time IS NOT NULL
        AND (
            (t.start_time >= $2 AND t.start_time <= $3)
            OR
            (t.end_time >= $2 AND t.end_time <= $3)
            OR
            (t.start_time < $2 AND t.end_time > $3)
        )
)
SELECT
    idtask,
    task_name,
    SUM(adjusted_end_time - adjusted_start_time) AS total_time
FROM
    filtered_times
GROUP BY
    idtask,
    task_name
ORDER BY
    total_time DESC;
`

  var results TaskTimeCollect
  err := t.db.SelectContext(ctx, &results, query, id, start, end)
  if err != nil {
    return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
  }
  return results.toDTO(), nil
}

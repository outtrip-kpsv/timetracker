package repo

import (
  "context"
  "fmt"
  "timetracker/internal/bl/dto"
  "timetracker/internal/db"
  "timetracker/internal/utils/const/status"
)

type ITaskBL interface {
  CreateTask(ctx context.Context, task *dto.Task) (*dto.Task, error)
  StartTask(ctx context.Context, idP, idT string) error
  PauseTask(ctx context.Context, idP, idT string) error
  CompleteTask(ctx context.Context, idP, idT string) error
  TimeTasks(ctx context.Context, id, start, end string) ([]dto.TaskTimeResult, error)
}

type taskBL struct {
  db *db.DbRepo
}

func NewTaskBL(db *db.DbRepo) ITaskBL {
  return &taskBL{db: db}
}

func (t *taskBL) CreateTask(ctx context.Context, task *dto.Task) (*dto.Task, error) {

  _, err := t.db.People.GetByUUID(ctx, task.IdPerson)
  if err != nil {
    return nil, err
  }

  createTask, err := t.db.Task.CreateTask(ctx, task)
  if err != nil {
    return nil, err
  }

  return createTask, nil
}

func (t *taskBL) StartTask(ctx context.Context, idP, idT string) error {
  ctx, err := t.db.Begin(ctx)
  if err != nil {
    return err
  }
  defer func() {
    t.db.End(ctx, err)
  }()

  st, err := t.db.Task.GetTaskStatus(ctx, idT)
  if err != nil {
    return err
  }
  if st == status.New || st == status.Pause {
    err = t.db.TimeTask.StartTimer(ctx, idT)
    if err != nil {
      return err
    }
  } else {
    err = fmt.Errorf("задача в статусе %s", st)
    fmt.Println(err)
    return err
  }

  err = t.db.Task.UpdateStatus(ctx, idT, status.Work)
  if err != nil {
    return err
  }
  return nil
}

func (t *taskBL) PauseTask(ctx context.Context, idP, idT string) error {
  ctx, err := t.db.Begin(ctx)
  if err != nil {
    return err
  }
  defer func() {
    t.db.End(ctx, err)
  }()

  st, err := t.db.Task.GetTaskStatus(ctx, idT)
  if err != nil {
    return err
  }

  if st != status.Work {
    err = fmt.Errorf("задача в статусе %s", st)
    return err
  }

  err = t.db.TimeTask.StopTimer(ctx, idT)
  if err != nil {
    return err
  }

  err = t.db.Task.UpdateStatus(ctx, idT, status.Pause)
  if err != nil {
    return err
  }

  return nil
}

func (t *taskBL) CompleteTask(ctx context.Context, idP, idT string) error {
  ctx, err := t.db.Begin(ctx)
  if err != nil {
    return err
  }
  defer func() {
    t.db.End(ctx, err)
  }()

  st, err := t.db.Task.GetTaskStatus(ctx, idT)
  if err != nil {
    return err
  }

  if st != status.Work {
    err = fmt.Errorf("задача в статусе %s", st)
    return err
  }

  err = t.db.TimeTask.StopTimer(ctx, idT)
  if err != nil {
    return err
  }

  err = t.db.Task.UpdateStatus(ctx, idT, status.Complete)
  if err != nil {
    return err
  }

  return nil
}

func (t *taskBL) TimeTasks(ctx context.Context, id, start, end string) ([]dto.TaskTimeResult, error) {
  times, err := t.db.Task.TaskTimes(ctx, id, start, end)
  if err != nil {
    return nil, err
  }
  return times, nil
}

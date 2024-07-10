package db

import (
  "context"
  "github.com/jmoiron/sqlx"
  "timetracker/internal/db/repo"
)

type DbRepo struct {
  db       *sqlx.DB
  People   repo.IPeopleRepo
  Task     repo.ITaskRepo
  TimeTask repo.ITimeTaskRepo
}

func New(connStr string) *DbRepo {
  res := DbRepo{db: newDb(connStr)}
  res.People = repo.NewPeopleRepo(res.db)
  res.Task = repo.NewTaskRepo(res.db)
  res.TimeTask = repo.NewTimeTaskRepo(res.db)
  return &res
}

func (d *DbRepo) Begin(ctx context.Context) (context.Context, error) {
  tx, err := d.db.BeginTxx(ctx, nil)
  if err != nil {
    return ctx, err
  }
  ctx = context.WithValue(ctx, "tx", tx)

  return ctx, nil
}

func (d *DbRepo) End(ctx context.Context, err error) {
  tx, ok := ctx.Value("tx").(*sqlx.Tx)
  if !ok {
    return
  }
  if err != nil {
    err := tx.Rollback()
    if err != nil {
      return
    }
  } else {
    err := tx.Commit()
    if err != nil {
      return
    }
  }
}

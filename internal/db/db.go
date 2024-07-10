package db

import (
  "fmt"
  "github.com/golang-migrate/migrate/v4"
  _ "github.com/golang-migrate/migrate/v4/database/postgres"
  _ "github.com/golang-migrate/migrate/v4/source/file"
  "github.com/jmoiron/sqlx"
  "path/filepath"

  "log"
)

func newDb(connStr string) *sqlx.DB {
  conn, err := sqlx.Open("postgres", connStr)
  if err != nil {
    log.Fatal(err)
  }
  err = runMigrations(connStr)
  if err != nil {
    log.Println(err.Error())
    return nil
  }
  return conn
}

func runMigrations(dbURL string) error {
  dir, err := filepath.Abs("internal/db/migrations")
  fmt.Println(dir)
  if err != nil {
    return fmt.Errorf("could not get absolute path: %w", err)
  }

  m, err := migrate.New(
    "file://"+dir,
    dbURL)
  if err != nil {
    return err
  }
  if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    return err
  }
  return nil
}

package main

import (
  "context"
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "timetracker/internal/bl"
  "timetracker/internal/config"
  "timetracker/internal/config/logger"
  "timetracker/internal/db"
  "timetracker/internal/io/http"
)

func main() {
  conf, _ := config.InitConfServ()

  lg := logger.New(conf.Options.Log)

  ctx, cancel := context.WithCancel(context.Background())
  fin := make(chan struct{})
  done := make(chan os.Signal, 1)
  signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

  d := db.New(conf.Options.DbString())
  blRepo := bl.New(d)
  fmt.Println(conf.Options.DbString())
  serv := http.New(conf.Options.ServStr(), lg, blRepo, fin)

  serv.Run()

  lg.Info("Server Started")

  <-done

  lg.Info("Выключение сервера...")

  defer cancel()
  serv.Stop(ctx)

}

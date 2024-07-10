package logger

import (
  "log/slog"
  "os"
)

func New(mode string) *slog.Logger {
  var opts *slog.HandlerOptions
  if mode == "debug" {
    opts = &slog.HandlerOptions{
      Level: slog.LevelDebug,
    }
  }

  log := slog.New(slog.NewJSONHandler(os.Stdout, opts))

  return log
}

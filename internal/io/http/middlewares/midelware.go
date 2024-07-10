package middlewares

import (
  "context"
  "log/slog"
  "net/http"
  "time"
)

type Mw struct {
  l *slog.Logger
}

func New(log *slog.Logger) *Mw {
  return &Mw{l: log}
}

func (m *Mw) WithLogger(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()

    atr := slog.GroupValue(
      slog.String("method", r.Method),
      slog.String("url", r.Host+r.URL.String()),
      slog.String("userAgent", r.Header.Values("User-Agent")[0]),
    )
    log := m.l.With("request", atr, slog.Any("startTime", startTime))
    log.Debug("requestStart")
    ctx := context.WithValue(r.Context(), "logger", log)
    ctx = context.WithValue(ctx, "startTime", startTime)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

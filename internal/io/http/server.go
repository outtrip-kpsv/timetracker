package http

import (
  "context"
  "log/slog"
  "net/http"
  "timetracker/internal/bl"
)

type serv struct {
  l   *slog.Logger
  srv *http.Server
  bl  *bl.BL

  fin chan struct{}
}

func New(address string, log *slog.Logger, bl *bl.BL, fin chan struct{}) *serv {
  srv := &http.Server{
    Addr:    address,
    Handler: InitRoutes(bl, log),
  }
  return &serv{
    l:   log.With(slog.String("layer", "serv")),
    srv: srv,
    bl:  bl,
    fin: fin}
}

func (s *serv) Run() {
  go func() {
    if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      s.l.Error("listen: %s\n", err)
    }
  }()

}

func (s *serv) Stop(ctx context.Context) {
  if err := s.srv.Shutdown(ctx); err != nil {
    s.l.Error("Ошибка выключения сервера:", err)
  }
  s.l.Info("Сервер успешно выключен")

}

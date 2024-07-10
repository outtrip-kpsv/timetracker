package http

import (
  "log/slog"
  "net/http"
  "timetracker/internal/bl"
  "timetracker/internal/io/http/handlers"
  "timetracker/internal/io/http/middlewares"
)

type router struct {
  logger      *slog.Logger
  router      *http.ServeMux
  middlewares *middlewares.Mw
}

func InitRoutes(bl *bl.BL, logger *slog.Logger) http.Handler {
  r := &router{
    logger:      logger,
    router:      http.NewServeMux(),
    middlewares: middlewares.New(logger),
  }
  controller := handlers.NewController(bl, r.logger)
  r.logger.Debug("init handler")

  r.router.HandleFunc("GET /people", r.wrapHandler(controller.GetPeoples))
  r.router.HandleFunc("POST /people", r.wrapHandler(controller.CreatePeople))
  r.router.HandleFunc("DELETE /people/{uuid}", r.wrapHandler(controller.DeletePeople))
  r.router.HandleFunc("GET /people/{uuid}", r.wrapHandler(controller.GetPeopleByUUID))
  r.router.HandleFunc("PATCH /people/{uuid}", r.wrapHandler(controller.UpdatePeopleByUUID))

  r.router.HandleFunc("POST /people/{uuid}/create-task", r.wrapHandler(controller.CreateTask))

  r.router.HandleFunc("GET /people/{uuidP}/{uuidT}/start", r.wrapHandler(controller.StartTimer))
  r.router.HandleFunc("GET /people/{uuidP}/{uuidT}/pause", r.wrapHandler(controller.PauseTimer))
  r.router.HandleFunc("GET /people/{uuidP}/{uuidT}/complete", r.wrapHandler(controller.CompleteTask))

  r.router.HandleFunc("POST /people/{uuid}/worktime", r.wrapHandler(controller.WorkTime))

  r.router.HandleFunc("GET /info", r.wrapHandler(controller.InfoPeople))

  r.router.HandleFunc("/", r.wrapHandler(controller.NotFound))
  muxN := use(r.router, r.middlewares.WithLogger)

  return muxN
}

func use(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
  var s http.Handler
  s = r

  for _, mw := range middlewares {
    s = mw(s)
  }

  return s
}

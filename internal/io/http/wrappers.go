package http

import (
  "encoding/json"
  "log/slog"
  "net/http"
  "time"
  "timetracker/internal/io/http/models"
)

func (r *router) wrapHandler(handler func(w http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error)) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {

    res, status, attr, err := handler(w, req)
    ctxLogger := req.Context().Value("logger").(*slog.Logger)
    startTime := req.Context().Value("startTime").(time.Time)

    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(status)
    if err != nil {
      attr = slog.Group("error", slog.String("msg", err.Error()), slog.Any("status", status), attr)
      res = models.ErrorResponse{
        Error: err.Error(),
      }
    }
    resBytes, _ := json.Marshal(res)
    endTime := time.Now()
    allTime := endTime.Sub(startTime).String()
    ctxLogger.Info("reqDone", slog.Any("endTime", endTime), slog.Any("duration", allTime), attr)
    _, err = w.Write(resBytes)
    if err != nil {
      return
    }
  }
}

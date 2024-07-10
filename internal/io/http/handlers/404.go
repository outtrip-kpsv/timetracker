package handlers

import (
  "errors"
  "log/slog"
  "net/http"
)

func (c *Controller) NotFound(_ http.ResponseWriter, _ *http.Request) (interface{}, int, slog.Attr, error) {
  return nil, http.StatusNotFound, slog.Attr{}, errors.New("handler not found")
}

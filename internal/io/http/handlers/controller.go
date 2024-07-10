package handlers

import (
  "log/slog"
  "timetracker/internal/bl"
)

//import (
//  //"auth-ms/internal/bl"
//)

type Controller struct {
  bl *bl.BL
  l  *slog.Logger
}

func NewController(bl *bl.BL, log *slog.Logger) *Controller {
  return &Controller{bl: bl, l: log}
}

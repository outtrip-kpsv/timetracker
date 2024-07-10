package handlers

import (
  "fmt"
  "log/slog"
  "net/http"
  "timetracker/internal/io/http/models"
  "timetracker/internal/utils"
)

// StartTimer начинает таймер для задачи указанного человека по UUID человека и UUID задачи
// @Summary Начало таймера для задачи
// @Description Начинает таймер для задачи указанного человека по его UUID и UUID задачи
// @Tags tasks
// @Accept json
// @Produce json
// @Param uuidP path string true "UUID человека"
// @Param uuidT path string true "UUID задачи"
// @Success 200 {object} models.Ok "Успешное начало таймера для задачи"
// @Failure 400 {object} models.ErrorResponse "Некорректные параметры запроса"
// @Failure 404 {object} models.ErrorResponse "Человек или задача с указанным UUID не найдены"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /people/{uuidP}/{uuidT}/start [get]
func (c *Controller) StartTimer(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
  idPerson := req.PathValue("uuidP")
  if !utils.IsValidUUID(idPerson) {
    attr := slog.String("not uuidP", idPerson)
    return nil, http.StatusBadRequest, attr, fmt.Errorf("uuidP %s не валидный", idPerson)
  }
  idTask := req.PathValue("uuidT")

  if !utils.IsValidUUID(idTask) {
    attr := slog.String("not uuidT", idTask)
    return nil, http.StatusBadRequest, attr, fmt.Errorf("uuidT %s не валидный", idTask)
  }

  err := c.bl.Task.StartTask(req.Context(), idPerson, idTask)
  if err != nil {
    return nil, http.StatusBadRequest, slog.Attr{}, err
  }
  return models.Ok{
    Msg:    "задача в работе",
    Status: http.StatusOK,
  }, http.StatusOK, slog.Attr{}, nil
}

// PauseTimer приостанавливает таймер для задачи указанного человека по UUID человека и UUID задачи
// @Summary Приостановка таймера для задачи
// @Description Приостанавливает таймер для задачи указанного человека по его UUID и UUID задачи
// @Tags tasks
// @Accept json
// @Produce json
// @Param uuidP path string true "UUID человека"
// @Param uuidT path string true "UUID задачи"
// @Success 200 {object} models.Ok "Успешная приостановка таймера для задачи"
// @Failure 400 {object} models.ErrorResponse "Некорректные параметры запроса"
// @Failure 404 {object} models.ErrorResponse "Человек или задача с указанным UUID не найдены"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /people/{uuidP}/{uuidT}/pause [get]
func (c *Controller) PauseTimer(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
  idPerson := req.PathValue("uuidP")
  if !utils.IsValidUUID(idPerson) {
    attr := slog.String("not uuid", idPerson)
    return nil, http.StatusBadRequest, attr, fmt.Errorf("uuid %s не валидный", idPerson)
  }
  idTask := req.PathValue("uuidT")

  if !utils.IsValidUUID(idTask) {
    attr := slog.String("not uuid", idTask)
    return nil, http.StatusBadRequest, attr, fmt.Errorf("uuid %s не валидный", idTask)
  }

  err := c.bl.Task.PauseTask(req.Context(), idPerson, idTask)
  if err != nil {
    return nil, http.StatusBadRequest, slog.Attr{}, err
  }
  return models.Ok{
    Msg:    "задача на паузе",
    Status: http.StatusOK,
  }, http.StatusOK, slog.Attr{}, nil
}

// WorkTime возвращает список задач с временем работы в указанном диапазоне для указанного человека по UUID
// @Summary Получение списка задач с временем работы
// @Description Возвращает список задач с временем работы в указанном диапазоне для указанного человека по его UUID
// @Tags tasks
// @Accept json
// @Produce json
// @Param uuid path string true "UUID человека"
// @Param body body models.DateStartEnd true "Временной диапазон"
// @Success 200 {object} []dto.TaskTimeResult "Список задач с временем работы"
// @Failure 400 {object} models.ErrorResponse "Некорректные параметры запроса"
// @Failure 404 {object} models.ErrorResponse "Человек с указанным UUID не найден или нет задач в указанном диапазоне"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /people/{uuid}/worktime [post]
func (c *Controller) WorkTime(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
  id := req.PathValue("uuid")
  if !utils.IsValidUUID(id) {
    attr := slog.String("not uuid", id)
    return nil, http.StatusBadRequest, attr, fmt.Errorf("uuid %s не валидный", id)
  }

  var tm models.DateStartEnd
  body, err := utils.DecodeRequestBody(req, &tm)
  if err != nil {
    attr := slog.Group("body", slog.String("reqBody", body))
    return nil, http.StatusBadRequest, attr, err
  }

  if !utils.IsValidDateTime(tm.Start) ||
      !utils.IsValidDateTime(tm.End) ||
      len(tm.Start) == 0 ||
      len(tm.End) == 0 {
    attr := slog.Group("body", slog.String("start", tm.Start), slog.String("end", tm.End))
    return nil, http.StatusBadRequest, attr, fmt.Errorf("проблемма входных данных")
  }

  if !utils.IsValidTimeRange(tm.Start, tm.End) {
    attr := slog.Group("body", slog.String("start", tm.Start), slog.String("end", tm.End))
    return nil, http.StatusBadRequest, attr, fmt.Errorf("дата конца диапозона раньше чем начало")

  }
  tasks, err := c.bl.Task.TimeTasks(req.Context(), id, tm.Start, tm.End)
  if err != nil {
    return nil, http.StatusBadRequest, slog.Attr{}, err
  }
  return tasks, http.StatusOK, slog.Attr{}, nil
}

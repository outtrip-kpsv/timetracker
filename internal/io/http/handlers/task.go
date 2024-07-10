package handlers

import (
  "fmt"
  "log/slog"
  "net/http"
  "timetracker/internal/io/http/models"
  "timetracker/internal/utils"
)

// CreateTask создает новую задачу для указанного человека по его UUID
// @Summary Создание новой задачи для человека
// @Description Создает новую задачу для человека по его UUID
// @Tags tasks
// @Accept json
// @Produce json
// @Param uuid path string true "UUID человека"
// @Param task body models.TaskCreate true "Данные для создания задачи"
// @Success 200 {object} models.Ok "Созданная задача"
// @Failure 400 {object} models.ErrorResponse "Некорректные параметры запроса"
// @Failure 404 {object} models.ErrorResponse "Человек с указанным UUID не найден"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /people/{uuid}/create-task [post]
func (c *Controller) CreateTask(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
  id := req.PathValue("uuid")
  if !utils.IsValidUUID(id) {
    attr := slog.String("not uuid", id)
    return nil, http.StatusBadRequest, attr, fmt.Errorf("uuid %s не валидный", id)
  }
  var task models.TaskCreate
  utils.DecodeRequestBody(req, &task)
  task.IdPerson = id
  createTask, err := c.bl.Task.CreateTask(req.Context(), task.ToDto())
  if err != nil {
    return nil, http.StatusBadRequest, slog.Attr{}, err
  }
  resp := models.TaskFromDto(createTask)
  return resp, http.StatusOK, slog.Attr{}, nil
}

// CompleteTask завершает задачу для указанного человека по UUID человека и UUID задачи
// @Summary Завершение задачи для человека
// @Description Завершает задачу для человека по его UUID и UUID задачи
// @Tags tasks
// @Accept json
// @Produce json
// @Param uuidP path string true "UUID человека"
// @Param uuidT path string true "UUID задачи"
// @Success 200 {object} models.Ok "Успешное завершение задачи"
// @Failure 400 {object} models.ErrorResponse "Некорректные параметры запроса"
// @Failure 404 {object} models.ErrorResponse "Человек или задача с указанным UUID не найдены"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /people/{uuidP}/{uuidT}/complete [get]
func (c *Controller) CompleteTask(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
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

  err := c.bl.Task.CompleteTask(req.Context(), idPerson, idTask)
  if err != nil {
    return nil, http.StatusBadRequest, slog.Attr{}, err
  }

  return models.Ok{
    Msg:    "задача завершена",
    Status: http.StatusOK,
  }, http.StatusOK, slog.Attr{}, nil
}

package handlers

import (
  "errors"
  "fmt"
  "log/slog"
  "net/http"
  "strconv"
  "timetracker/internal/bl/dto"
  "timetracker/internal/io/http/models"
  "timetracker/internal/utils"
)

// CreatePeople создает новую запись о человеке
// @Summary Создание человека
// @Description Создание новой записи о человеке по номеру паспорта
// @Tags people
// @Accept json
// @Produce json
// @Param passport body dto.Passport true "Паспортные данные"
// @Success 200 {object} dto.Person
// @Failure 400 {object} models.ErrorResponse "Неверный запрос"
// @Failure 409 {object} models.ErrorResponse "Конфликт данных"
// @Router /people [post]
func (c *Controller) CreatePeople(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
  var passport dto.Passport
  body, err := utils.DecodeRequestBody(req, &passport)
  if err != nil {
    attr := slog.Group("body", slog.String("reqBody", body))
    return nil, http.StatusBadRequest, attr, err
  }
  if err = utils.PassportValidate(passport.PassportNumber); err != nil {
    attr := slog.Group("body", slog.String("passportNumber", passport.PassportNumber))
    return nil, http.StatusBadRequest, attr, err
  }

  people, err := c.bl.People.CreatePeople(req.Context(), passport)
  if err != nil {
    return nil, http.StatusConflict, slog.Attr{}, err
  }

  return people, http.StatusOK, slog.Attr{}, nil
}

// DeletePeople удаляет запись о человеке по UUID
// @Summary Удаление человека
// @Description Удаление записи о человеке по UUID
// @Tags people
// @Accept json
// @Produce json
// @Param uuid path string true "UUID человека для удаления"
// @Success 200 {object} models.Ok
// @Failure 400 {object} models.ErrorResponse "Неверный UUID"
// @Failure 404 {object} models.ErrorResponse "Человек не найден"
// @Router /people/{uuid} [delete]
func (c *Controller) DeletePeople(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
  id := req.PathValue("uuid")
  if !utils.IsValidUUID(id) {
    attr := slog.String("not uuid", id)
    return nil, http.StatusBadRequest, attr, fmt.Errorf("uuid %s не валидный", id)
  }

  err := c.bl.People.DeletePeople(req.Context(), id)
  if err != nil {
    return nil, http.StatusNotFound, slog.Attr{}, err
  }
  return models.Ok{
    Msg:    fmt.Sprintf("id: %s удален", id),
    Status: http.StatusOK,
  }, http.StatusOK, slog.Attr{}, nil
}

// InfoPeople возвращает информацию о человеке по серии и номеру паспорта
// @Summary Получение информации о человеке
// @Description Получение информации о человеке по серии и номеру паспорта
// @Tags people
// @Accept json
// @Produce json
// @Param passportSerie query string true "Серия паспорта (только цифры)"
// @Param passportNumber query string true "Номер паспорта (только цифры)"
// @Success 200 {object} dto.Person "Информация о человеке"
// @Failure 400 {object} models.ErrorResponse "Неверный формат серии или номера паспорта"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /info [get]
func (c *Controller) InfoPeople(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
  passportSerie := req.URL.Query().Get("passportSerie")
  passportNumber := req.URL.Query().Get("passportNumber")
  attr := slog.Group("getQuery",
    slog.String("passportSerie", passportSerie),
    slog.String("passportNumber", passportNumber))

  if !utils.OnlyDigit(passportSerie) || !utils.SeriesValid(passportSerie) {
    return nil, http.StatusBadRequest, attr, errors.New("не правильный формат серии паспорта: пример 'passportSerie=1234'")
  }

  if !utils.OnlyDigit(passportNumber) || !utils.NoValid(passportNumber) {
    return nil, http.StatusBadRequest, attr, errors.New("не правильный формат номера паспорта: пример 'passportSerie=1234'")
  }
  people, err := c.bl.People.FakePeople(req.Context(), passportSerie, passportNumber)
  if err != nil {
    attr := slog.Group("people",
      slog.Any("all", people))
    return nil, http.StatusInternalServerError, attr, err
  }
  return people, http.StatusOK, attr, nil
}

// GetPeopleByUUID возвращает информацию о человеке по UUID
// @Summary Получение информации о человеке по UUID
// @Description Получение информации о человеке по его уникальному идентификатору
// @Tags people
// @Accept json
// @Produce json
// @Param uuid path string true "Уникальный идентификатор человека (UUID)"
// @Success 200 {object} dto.Person "Информация о человеке"
// @Failure 400 {object} models.ErrorResponse "Неверный формат UUID"
// @Failure 404 {object} models.ErrorResponse "Человек с указанным UUID не найден"
// @Router /people/{uuid} [get]
func (c *Controller) GetPeopleByUUID(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
  id := req.PathValue("uuid")
  if !utils.IsValidUUID(id) {
    attr := slog.String("not uuid", id)
    return nil, http.StatusBadRequest, attr, fmt.Errorf("uuid %s не валидный", id)
  }

  people, err := c.bl.People.GetPeopleUUID(req.Context(), id)
  if err != nil {
    return nil, http.StatusNotFound, slog.Attr{}, err
  }
  return people, http.StatusOK, slog.Attr{}, nil
}

// UpdatePeopleByUUID обновляет информацию о человеке по его UUID
// @Summary Обновление информации о человеке
// @Description Обновляет информацию о человеке по его UUID
// @Tags people
// @Accept json
// @Produce json
// @Param uuid path string true "UUID человека"
// @Param body body dto.Person true "Информация о человеке для обновления"
// @Success 200 {object} dto.Person "Обновленная информация о человеке"
// @Failure 400 {object} models.ErrorResponse "Некорректные параметры запроса"
// @Failure 404 {object} models.ErrorResponse "Человек с указанным UUID не найден"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /people/{uuid} [patch]
func (c *Controller) UpdatePeopleByUUID(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
  id := req.PathValue("uuid")
  if !utils.IsValidUUID(id) {
    attr := slog.String("not uuid", id)
    return nil, http.StatusBadRequest, attr, fmt.Errorf("uuid %s не валидный", id)
  }

  var people dto.Person
  body, err := utils.DecodeRequestBody(req, &people)
  if err != nil {
    attr := slog.Group("body", slog.String("reqBody", body))
    return nil, http.StatusBadRequest, attr, err
  }

  people.ID = id
  person, err := c.bl.People.UpdatePeople(req.Context(), people)
  if err != nil {
    return nil, http.StatusNotFound, slog.Attr{}, err
  }
  return person, http.StatusOK, slog.Attr{}, nil
}

// GetPeoples возвращает список людей с возможностью фильтрации и пагинации
// @Summary Получение списка людей с фильтрацией и пагинацией
// @Description Получение списка людей с возможностью фильтрации по различным полям и пагинацией
// @Tags people
// @Accept json
// @Produce json
// @Param surname query string false "Фамилия человека"
// @Param name query string false "Имя человека"
// @Param patronymic query string false "Отчество человека"
// @Param address query string false "Адрес человека"
// @Param passport_number query string false "Номер паспорта человека"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество записей на странице (по умолчанию 10)"
// @Success 200 {object} models.PeopleResp "Список людей"
// @Failure 400 {object} models.ErrorResponse "Некорректные параметры запроса"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /people [get]
func (c *Controller) GetPeoples(_ http.ResponseWriter, req *http.Request) (interface{}, int, slog.Attr, error) {
  queryParams := req.URL.Query()
  filter := &dto.Person{
    People: dto.People{
      Surname:    queryParams.Get("surname"),
      Name:       queryParams.Get("name"),
      Patronymic: queryParams.Get("patronymic"),
      Address:    queryParams.Get("address"),
    },
    Passport: dto.Passport{
      PassportNumber: queryParams.Get("passport_number"),
    },
  }
  pageStr := queryParams.Get("page")
  if pageStr == "" {
    pageStr = "1"
  }
  page, err := strconv.Atoi(pageStr)
  if err != nil || page < 1 {
    return nil, http.StatusBadRequest, slog.Attr{}, errors.New("некорректное значение параметра page")
  }
  limitStr := queryParams.Get("limit")
  if limitStr == "" {
    limitStr = "10"
  }
  limit, err := strconv.Atoi(limitStr)
  if err != nil || limit < 1 {
    return nil, http.StatusBadRequest, slog.Attr{}, errors.New("некорректное значение параметра limit")
  }
  offset := (page - 1) * limit

  people, total, err := c.bl.People.GetPeople(req.Context(), filter, offset, limit)
  if err != nil {
    return nil, http.StatusInternalServerError, slog.Attr{}, fmt.Errorf("ошибка при получении данных пользователей: %s", err.Error())
  }
  return models.PeopleResp{
    Total:  total,
    Limit:  limit,
    Offset: offset,
    People: people,
  }, http.StatusOK, slog.Int("total", total), nil
}

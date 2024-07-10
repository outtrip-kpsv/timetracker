package repo

import (
  "context"
  "fmt"
  "github.com/brianvoe/gofakeit/v7"
  "log/slog"
  "strings"
  "timetracker/internal/bl/dto"
  "timetracker/internal/db"
  "timetracker/internal/utils"
  "timetracker/internal/utils/var/endpoint"
)

type IPeopleBL interface {
  CreatePeople(ctx context.Context, passport dto.Passport) (*dto.Person, error)
  FakePeople(ctx context.Context, series, number string) (dto.People, error)
  DeletePeople(ctx context.Context, uuid string) error
  GetPeopleUUID(ctx context.Context, uuid string) (*dto.Person, error)
  UpdatePeople(ctx context.Context, people dto.Person) (*dto.Person, error)
  GetPeople(ctx context.Context, filter *dto.Person, offset, limit int) ([]dto.Person, int, error)
}

type peopleBL struct {
  db *db.DbRepo
}

func NewPeopleBL(db *db.DbRepo) IPeopleBL {
  return &peopleBL{
    db: db,
  }
}

func (p peopleBL) CreatePeople(ctx context.Context, passport dto.Passport) (*dto.Person, error) {
  ctxLogger := ctx.Value("logger").(*slog.Logger)
  var people dto.People
  byPassport, err := p.db.People.GetByPassport(ctx, passport.PassportNumber)
  if byPassport != nil {
    return byPassport, fmt.Errorf("человек с паспортом: %s, уже добавлен", passport)
  }

  if err != nil {
    ctxLogger.Debug(err.Error())
    return nil, err
  }

  passportSlice := strings.Split(passport.PassportNumber, " ")
  getQuery := fmt.Sprintf("?passportSerie=%s&passportNumber=%s", passportSlice[0], passportSlice[1])
  err = utils.GetAPIRequest(ctx, endpoint.GetInfoUrl+getQuery, &people)
  if err != nil {
    ctxLogger.Debug("/info error")
    people, err = p.FakePeople(ctx, passportSlice[0], passportSlice[1])
    if err != nil {
      return nil, err
    }
  }
  person := &dto.Person{
    People:   people,
    Passport: passport,
  }

  person.ID, err = p.db.People.CreatePerson(ctx, person)
  if err != nil {
    ctxLogger.Debug(err.Error())
    return nil, err
  }
  return person, err
}

func (p peopleBL) FakePeople(_ context.Context, series, number string) (dto.People, error) {
  faker := gofakeit.New(utils.HashStringToInt(series + number))
  people := dto.People{
    Surname:    faker.LastName(),
    Name:       faker.FirstName(),
    Patronymic: faker.MiddleName(),
    Address:    faker.Address().Address,
  }
  return people, nil
}

func (p peopleBL) DeletePeople(ctx context.Context, uuid string) error {
  err := p.db.People.DeleteByPerson(ctx, uuid)
  if err != nil {
    return err
  }
  return nil
}

func (p peopleBL) GetPeopleUUID(ctx context.Context, uuid string) (*dto.Person, error) {
  people, err := p.db.People.GetByUUID(ctx, uuid)
  if err != nil {
    return people, err
  }
  return people, nil
}

func (p peopleBL) UpdatePeople(ctx context.Context, people dto.Person) (*dto.Person, error) {
  oldPerson, err := p.db.People.GetByUUID(ctx, people.ID)
  if err != nil {
    return nil, err
  }
  oldPerson.Name = UpdateField(people.Name, oldPerson.Name)
  oldPerson.Surname = UpdateField(people.Surname, oldPerson.Surname)
  oldPerson.Address = UpdateField(people.Address, oldPerson.Address)
  oldPerson.Patronymic = UpdateField(people.Patronymic, oldPerson.Patronymic)
  person, err := p.db.People.UpdatePerson(ctx, oldPerson)
  if err != nil {
    return nil, err
  }
  return person, nil
}

func (p peopleBL) GetPeople(ctx context.Context, filter *dto.Person, offset, limit int) ([]dto.Person, int, error) {
  persons, total, err := p.db.People.GetPeople(ctx, filter, offset, limit)
  if err != nil {
    return nil, 0, err
  }
  return persons, total, nil
}

func UpdateField(new, old string) string {
  if len(new) != 0 {
    return new
  }
  return old
}

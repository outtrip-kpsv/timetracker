package repo

import (
  "context"
  "database/sql"
  "fmt"
  "github.com/jmoiron/sqlx"
  "log/slog"
  "timetracker/internal/bl/dto"
)

type Person struct {
  Id             string `json:"id" db:"id"`
  Surname        string `json:"surname" db:"surname"`
  Name           string `json:"name" db:"name"`
  Patronymic     string `json:"patronymic" db:"patronymic"`
  Address        string `json:"address" db:"address"`
  PassportNumber string `json:"passportNumber" db:"passport_number"`
}

func (p *Person) toDTO() *dto.Person {
  if p == nil {
    return nil
  }
  return &dto.Person{
    ID: p.Id,
    People: dto.People{
      Surname:    p.Surname,
      Name:       p.Name,
      Patronymic: p.Patronymic,
      Address:    p.Address,
    },
    Passport: dto.Passport{
      PassportNumber: p.PassportNumber,
    },
  }

}

func (p *Person) fromDTO(model *dto.Person) *Person {
  if model == nil {
    return nil
  }
  p.Id = model.ID
  p.Surname = model.Surname
  p.Address = model.Address
  p.Name = model.Name
  p.Patronymic = model.Patronymic
  p.PassportNumber = model.PassportNumber
  return p
}

type IPeopleRepo interface {
  GetByPassport(ctx context.Context, passport string) (*dto.Person, error)
  GetByUUID(ctx context.Context, uuid string) (*dto.Person, error)
  CreatePerson(ctx context.Context, person *dto.Person) (string, error)
  DeleteByPerson(ctx context.Context, uuid string) error
  UpdatePerson(ctx context.Context, person *dto.Person) (*dto.Person, error)
  GetPeople(ctx context.Context, filter *dto.Person, offset, limit int) ([]dto.Person, int, error)
}

type peopleRepo struct {
  db *sqlx.DB
}

func NewPeopleRepo(db *sqlx.DB) IPeopleRepo {
  return &peopleRepo{db: db}
}

// GetByPassport ищет человека по номеру паспорта
func (p *peopleRepo) GetByPassport(ctx context.Context, passport string) (*dto.Person, error) {
  query := "SELECT surname, name, patronymic, address, passport_number FROM person WHERE passport_number = $1"
  var person Person
  fmt.Println(ctx)
  err := p.db.GetContext(ctx, &person, query, passport)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("ошибка бд: %v", err)
  }
  fmt.Println(person.Id)
  return person.toDTO(), nil
}

// CreatePerson создает новую запись о человеке в базе данных
func (p *peopleRepo) CreatePerson(ctx context.Context, person *dto.Person) (string, error) {
  query := `INSERT INTO person (surname, name, patronymic, address, passport_number) 
	          VALUES (:surname, :name, :patronymic, :address, :passport_number) 
	          RETURNING id`

  var id string
  var per Person

  rows, err := p.db.NamedQueryContext(ctx, query, per.fromDTO(person))
  if err != nil {
    return "", fmt.Errorf("ошибка вставки данных в базу: %v", err)
  }
  defer rows.Close()

  for rows.Next() {
    if err := rows.Scan(&id); err != nil {
      return "", fmt.Errorf("ошибка при получении ID: %v", err)
    }
  }

  return id, nil
}

// DeleteByPerson удаляет человека из таблицы
func (p *peopleRepo) DeleteByPerson(ctx context.Context, uuid string) error {
  query := "DELETE FROM person WHERE id = $1"

  result, err := p.db.ExecContext(ctx, query, uuid)
  if err != nil {
    return fmt.Errorf("ошибка удаления человека: %v", err)
  }

  rowsAffected, err := result.RowsAffected()
  if err != nil {
    return fmt.Errorf("ошибка получения количества затронутых строк: %v", err)
  }

  if rowsAffected == 0 {
    return fmt.Errorf("человек с UUID %s не найден", uuid)
  }

  return nil
}

// GetByUUID ищет человека по UUID
func (p *peopleRepo) GetByUUID(ctx context.Context, uuid string) (*dto.Person, error) {
  ctxLogger := ctx.Value("logger").(*slog.Logger)
  ctxLogger.Debug("db getbyuuid")
  query := "SELECT id, surname, name, patronymic, address, passport_number FROM person WHERE id = $1"
  var person Person
  err := p.db.GetContext(ctx, &person, query, uuid)
  if err != nil {
    if err == sql.ErrNoRows {
      ctxLogger.Error(err.Error())
      return nil, fmt.Errorf("человек с uuid %s не найден", uuid)
    }
    ctxLogger.Error(err.Error())
    return nil, fmt.Errorf("ошибка получения данных из базы: %v", err)
  }

  return person.toDTO(), nil
}

// UpdatePerson обновляет запись о человеке в базе данных
func (p *peopleRepo) UpdatePerson(ctx context.Context, person *dto.Person) (*dto.Person, error) {
  query := `UPDATE person
	          SET surname = :surname,
	              name = :name,
	              patronymic = :patronymic,
	              address = :address,
	              passport_number = :passport_number
	          WHERE id = :id
	          RETURNING id, surname, name, patronymic, address, passport_number`

  personDb := Person{}

  rows, err := p.db.NamedQueryContext(ctx, query, personDb.fromDTO(person))
  if err != nil {
    return nil, fmt.Errorf("ошибка обновления данных в базе: %v", err)
  }
  defer rows.Close()

  if rows.Next() {
    var updatedPerson Person
    if err := rows.StructScan(&updatedPerson); err != nil {
      return nil, fmt.Errorf("ошибка сканирования обновленных данных: %v", err)
    }
    return updatedPerson.toDTO(), nil
  }

  return nil, fmt.Errorf("человек с UUID %s не найден", person.ID)
}

func (p *peopleRepo) GetPeople(ctx context.Context, filter *dto.Person, offset, limit int) ([]dto.Person, int, error) {
  query := `SELECT id, surname, name, patronymic, address, passport_number
              FROM person
              WHERE (:surname = '' OR surname = :surname)
                AND (:name = '' OR name = :name)
                AND (:patronymic = '' OR patronymic = :patronymic)
                AND (:address = '' OR address = :address)
                AND (:passport_number = '' OR passport_number = :passport_number)
              ORDER BY surname, name
              LIMIT :limit OFFSET :offset`

  filterValues := map[string]interface{}{
    "surname":         filter.Surname,
    "name":            filter.Name,
    "patronymic":      filter.Patronymic,
    "address":         filter.Address,
    "passport_number": filter.PassportNumber,
    "limit":           limit,
    "offset":          offset,
  }
  rows, err := p.db.NamedQueryContext(ctx, query, filterValues)
  if err != nil {
    return nil, 0, fmt.Errorf("ошибка выполнения запроса: %v", err)
  }
  defer rows.Close()

  var people []dto.Person
  for rows.Next() {
    var person Person
    if err := rows.StructScan(&person); err != nil {
      return nil, 0, fmt.Errorf("ошибка сканирования данных: %v", err)
    }
    people = append(people, *person.toDTO())
  }

  countQuery := `SELECT COUNT(*)
                   FROM person
                   WHERE (:surname = '' OR surname = :surname)
                     AND (:name = '' OR name = :name)
                     AND (:patronymic = '' OR patronymic = :patronymic)
                     AND (:address = '' OR address = :address)
                     AND (:passport_number = '' OR passport_number = :passport_number)`

  nstmt, args, err := p.db.BindNamed(countQuery, filterValues)
  if err != nil {
    return nil, 0, fmt.Errorf("ошибка биндинга именованных параметров: %v", err)
  }

  var totalCount int
  err = p.db.GetContext(ctx, &totalCount, nstmt, args...)
  if err != nil {
    return nil, 0, fmt.Errorf("ошибка получения общего количества записей: %v", err)
  }

  return people, totalCount, nil

}

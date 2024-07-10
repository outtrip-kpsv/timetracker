package dto

type Person struct {
  ID string
  People
  Passport
}

type People struct {
  Surname    string `json:"surname"`
  Name       string `json:"name"`
  Patronymic string `json:"patronymic"`
  Address    string `json:"address"`
}

type Passport struct {
  PassportNumber string `json:"passportNumber"`
}

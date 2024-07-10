package utils

import (
  "errors"
  "github.com/google/uuid"
  "strings"
  "time"
  "unicode"
)

func PassportValidate(passport string) error {
  passport = strings.TrimSpace(passport)
  if len(passport) == 0 {
    return errors.New("поле не заполнено")
  }

  if !OnlyDigit(passport) {
    return errors.New("недопустимые символы во входных данных: пример '1234 567890'")
  }
  for _, r := range passport {
    if !unicode.IsDigit(r) && !unicode.IsSpace(r) {
      return errors.New("недопустимые символы во входных данных: пример '1234 567890'")
    }
  }
  pSlice := strings.Split(passport, " ")
  if len(pSlice) != 2 || !SeriesValid(pSlice[0]) || !NoValid(pSlice[1]) {
    return errors.New("не правильный формат: пример '1234 567890'")
  }

  return nil
}

func OnlyDigit(str string) bool {
  for _, r := range str {
    if !unicode.IsDigit(r) && !unicode.IsSpace(r) {
      return false
    }
  }
  return true
}

func SeriesValid(series string) bool {
  if len(series) != 4 {
    return false
  }
  return true
}

func NoValid(no string) bool {
  if len(no) != 6 {
    return false
  }
  return true
}

func IsValidUUID(u string) bool {
  _, err := uuid.Parse(u)
  return err == nil
}

func IsValidDateTime(dateTimeStr string) bool {
  _, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
  _, err2 := time.Parse("2006-01-02", dateTimeStr)

  return err == nil || err2 == nil
}

func IsValidTimeRange(startTimeStr, endTimeStr string) bool {
  var startTime, endTime time.Time
  var err1, err2 error

  if len(strings.Split(startTimeStr, " ")) == 2 {
    startTime, err1 = time.Parse("2006-01-02 15:04:05", startTimeStr)
  } else {
    startTime, err1 = time.Parse("2006-01-02", startTimeStr)
  }

  if len(strings.Split(endTimeStr, " ")) == 2 {
    endTime, err2 = time.Parse("2006-01-02 15:04:05", endTimeStr)
  } else {
    endTime, err2 = time.Parse("2006-01-02", endTimeStr)
  }

  if err1 != nil || err2 != nil {
    return false
  }

  return startTime.Before(endTime)
}

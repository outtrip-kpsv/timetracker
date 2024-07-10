package utils

import (
  "bytes"
  "encoding/json"
  "io"
  "net/http"
)

func DecodeRequestBody(req *http.Request, res interface{}) (string, error) {
  body, err := io.ReadAll(req.Body)
  if err != nil {
    return "", err
  }
  err = json.Unmarshal(body, res)
  req.Body = io.NopCloser(bytes.NewReader(body))
  return string(body), err
}

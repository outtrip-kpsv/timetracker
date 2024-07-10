package utils

import (
  "context"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "time"
)

func GetAPIRequest(ctx context.Context, url string, response interface{}) error {
  req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
  if err != nil {
    return fmt.Errorf("failed to create request: %w", err)
  }
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{Timeout: 3 * time.Second}

  resp, err := client.Do(req)
  if err != nil {
    return fmt.Errorf("request failed: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    bodyBytes, _ := io.ReadAll(resp.Body)
    return fmt.Errorf("received non-OK response: %s", string(bodyBytes))
  }

  bodyBytes, err := io.ReadAll(resp.Body)
  if err != nil {
    return fmt.Errorf("failed to read response body: %w", err)
  }

  if err := json.Unmarshal(bodyBytes, &response); err != nil {
    return fmt.Errorf("failed to unmarshal response data: %w", err)
  }

  return nil
}

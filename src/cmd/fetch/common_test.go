package fetch

import (
  "encoding/json"
  "os"
  "path"
  "testing"
)

var (
  inputFile = path.Join(os.TempDir(), "test-input.json")
  ts = fetchTestSuite{
    out: os.TempDir(),
  }
)

type fetchTestSuite struct {
  id   string
  out  string
  test *testing.T
}

func (ts *fetchTestSuite) checkError(kind string, err error, preErr func()) {
  if err != nil {
    if preErr != nil {
      preErr()
    }
    switch kind {
    case "fetch":
      ts.test.Fatalf("Expected fetch success, got fetch error: %v", err)

    case "stat":
      ts.test.Fatalf("Failed to get file info: %v", err)
    }
  }
}

func (ts *fetchTestSuite) prepareInputFile(source []int) {
  json, err := json.Marshal(source)
  if err != nil {
    ts.test.Fatalf("Failed to prepare input file: %v", err)
  }

  err = os.WriteFile(inputFile, json, 0644)
  if err != nil {
    ts.test.Fatalf("Failed to prepare input file: %v", err)
  }
}

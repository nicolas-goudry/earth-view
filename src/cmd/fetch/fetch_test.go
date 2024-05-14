package fetch

import (
  "os"
  "path"
  "regexp"
  "testing"
)

const testId = "1003"

func TestFetchSuccess(t *testing.T) {
  ts.test = t
  ts.id = testId

  filePath, err := runFetchCmd(ts.id, ts.out, true)
  ts.checkError("fetch", err, nil)
  os.Remove(filePath)
}

func TestFetchFail(t *testing.T) {
  ts.test = t
  ts.id = "999"

  filePath, err := runFetchCmd(ts.id, ts.out, true)
  want := regexp.MustCompile("not found")
  if err == nil {
    t.Fatalf("Expected error, got success. File is at %s", filePath)
  } else if !want.MatchString(err.Error()) {
    t.Fatalf("Received unexpected error: %v", err)
  }
}

func TestFetchCustomOutput(t *testing.T) {
  ts.test = t
  ts.id = testId
  ts.out = path.Join(os.TempDir(), "custom-out.jpeg")
  os.Remove(ts.out)

  filePath, err := runFetchCmd(ts.id, ts.out, false)
  ts.checkError("fetch", err, nil)

  if filePath != ts.out {
    os.Remove(filePath)
    t.Fatalf("Expected file to be saved at %s, got %s", ts.out, filePath)
  }

  os.Remove(filePath)
}

func TestFetchOverwrite(t *testing.T) {
  ts.test = t
  ts.id = testId

  filePath, err := runFetchCmd(ts.id, ts.out, true)
  ts.checkError("fetch", err, nil)

  clean := func() { os.Remove(filePath) }

  stat, err := os.Stat(filePath)
  ts.checkError("stat", err, clean)
  initialModTime := stat.ModTime()

  _, err = runFetchCmd(ts.id, ts.out, true)
  ts.checkError("fetch", err, clean)

  stat, err = os.Stat(filePath)
  ts.checkError("stat", err, clean)
  newModTime := stat.ModTime()

  if initialModTime == newModTime {
    clean()
    t.Fatal("Expected file to be modified")
  }

  clean()
}

func TestFetchNoOverwrite(t *testing.T) {
  ts.test = t
  ts.id = testId

  filePath, err := runFetchCmd(ts.id, ts.out, true)
  ts.checkError("fetch", err, nil)

  clean := func() { os.Remove(filePath) }

  stat, err := os.Stat(filePath)
  ts.checkError("stat", err, clean)
  initialModTime := stat.ModTime()

  _, err = runFetchCmd(ts.id, ts.out, false)
  ts.checkError("fetch", err, clean)

  stat, err = os.Stat(filePath)
  ts.checkError("stat", err, clean)
  newModTime := stat.ModTime()

  if initialModTime != newModTime {
    clean()
    t.Fatalf("Expected file to be untouched but modification time changed. Initial: %v, New: %v", initialModTime, newModTime)
  }

  clean()
}

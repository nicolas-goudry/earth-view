package fetch

import (
  "os"
  "path"
  "slices"
  "testing"

  "earth-view/lib"
)

var (
  inputIds = []int{ 1003, 1004, 1006, 1007 }
  altInputIds = []int{ 1548, 1550, 1551, 1553 }
)

func TestPickRandomIdInRange(t *testing.T) {
  for i := 0; i < 1000; i++ {
    randomId, _ := pickRandomId("")

    if randomId < lib.KnownIdLowerBoundary || randomId > lib.KnownIdUpperBoundary {
      t.Fatalf("Picked an invalid random id: %d", randomId)
    }
  }
}

func TestPickRandomIdFromFile(t *testing.T) {
  ts.test = t
  ts.prepareInputFile(inputIds)

  for i := 0; i < 20; i++ {
    randomId, err := pickRandomId(inputFile)
    if err != nil {
      os.Remove(inputFile)
      t.Fatalf("Got error while picking random id: %v", err)
    }

    if !slices.Contains(inputIds, randomId) {
      os.Remove(inputFile)
      t.Fatalf("Picked an invalid random id: %d", randomId)
    }
  }

  os.Remove(inputFile)
}

func TestFetchRandomSuccessFromFile(t *testing.T) {
  ts.test = t
  ts.prepareInputFile(inputIds)

  filePath, err := runFetchRandomCmd(inputFile, ts.out, true)
  ts.checkError("fetch", err, func() { os.Remove(inputFile) })
  os.Remove(filePath)
  os.Remove(inputFile)
}

func TestFetchRandomSuccessFromRange(t *testing.T) {
  ts.test = t

  filePath, err := runFetchRandomCmd("", ts.out, true)
  ts.checkError("fetch", err, nil)
  os.Remove(filePath)
}

func TestFetchRandomFailNoInputFile(t *testing.T) {
  ts.test = t

  filePath, err := runFetchRandomCmd(inputFile, ts.out, true)
  if err == nil {
    t.Fatal("Expected error, got success")
  }

  os.Remove(filePath)
}

func TestFetchRandomCustomOutput(t *testing.T) {
  ts.test = t
  ts.out = path.Join(os.TempDir(), "custom-out.jpeg")
  os.Remove(ts.out)

  filePath, err := runFetchRandomCmd("", ts.out, true)
  ts.checkError("fetch", err, nil)

  if filePath != ts.out {
    os.Remove(filePath)
    os.Remove(inputFile)
    t.Fatalf("Expected file to be saved at %s, got %s", ts.out, filePath)
  }

  os.Remove(filePath)
  os.Remove(inputFile)
}

func TestFetchRandomOverwrite(t *testing.T) {
  ts.test = t
  ts.out = path.Join(os.TempDir(), "custom-out.jpeg")
  ts.prepareInputFile(inputIds)

  filePath, err := runFetchRandomCmd(inputFile, ts.out, true)
  ts.checkError("fetch", err, nil)

  clean := func() {
    os.Remove(filePath)
    os.Remove(inputFile)
  }

  stat, err := os.Stat(filePath)
  ts.checkError("stat", err, clean)
  initialModTime := stat.ModTime()

  // Prepare an alternate input file to avoid picking the same random id
  ts.prepareInputFile(altInputIds)
  _, err = runFetchRandomCmd(inputFile, ts.out, true)
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

func TestFetchRandomNoOverwrite(t *testing.T) {
  ts.test = t
  ts.out = path.Join(os.TempDir(), "custom-out.jpeg")

  filePath, err := runFetchRandomCmd("", ts.out, true)
  ts.checkError("fetch", err, nil)

  clean := func() { os.Remove(filePath) }

  stat, err := os.Stat(filePath)
  ts.checkError("stat", err, clean)
  initialModTime := stat.ModTime()

  _, err = runFetchRandomCmd("", ts.out, false)
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

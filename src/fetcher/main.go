package main

import (
  "encoding/base64"
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "path/filepath"
  "strings"
)

func fetch(url string) ([]byte, error) {
  response, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer response.Body.Close()

  if response.StatusCode != http.StatusOK {
    return nil, errors.New("status code is " + string(response.StatusCode))
  }

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return nil, err
  }

  return body, nil
}

func toJson(content []byte) (map[string]interface{}, error) {
  var data map[string]interface{}

  if err := json.Unmarshal(content, &data); err != nil {
    return nil, err
  }

  return data, nil
}

func decode(content string) ([]byte, error) {
  imgWithMetadata := strings.Split(content, ",")
  encodedImg := imgWithMetadata[len(imgWithMetadata) - 1]

  decoded, err := base64.StdEncoding.DecodeString(encodedImg)
  if err != nil {
    return nil, err
  }

  return decoded, nil
}

func fileExists(path string) bool {
  _, err := os.Stat(path)
  return err == nil
}

// Function to extract the id from the URL
func extractId(url string) string {
  paths := strings.Split(url, "/")
  basename := paths[len(paths) - 1]
  idStr := strings.TrimSuffix(basename, ".json")

  return idStr
}

func writeFile(data []byte, path string) error {
  file, err := os.Create(path)
  if err != nil {
    return err
  }
  defer file.Close()

  _, err = file.Write(data)
  if err != nil {
    return err
  }

  return nil
}

const errorStringFormat = "error: %v\n"

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintf(os.Stderr, "URL to fetch is required\n")
    os.Exit(1)
  }

  if len(os.Args) < 3 {
    fmt.Fprintf(os.Stderr, "Destination directory is required\n")
    os.Exit(1)
  }

  url := os.Args[1]
  outDir := os.Args[2]
  outFile := filepath.Clean(outDir + "/" + extractId(url) + ".jpeg")

  if fileExists(outFile) {
    fmt.Println(outFile)
    os.Exit(0)
  }

  body, err := fetch(url)
  if err != nil {
    fmt.Fprintf(os.Stderr, errorStringFormat, err)
    os.Exit(1)
  }

  json, err := toJson(body)
  if err != nil {
    fmt.Fprintf(os.Stderr, errorStringFormat, err)
    os.Exit(1)
  }

  dataUri, err := decode(fmt.Sprint(json["dataUri"]))
  if err != nil {
    fmt.Fprintf(os.Stderr, errorStringFormat, err)
    os.Exit(1)
  }

  if err := writeFile(dataUri, outFile); err != nil {
    fmt.Fprintf(os.Stderr, errorStringFormat, err)
    os.Exit(1)
  }

  fmt.Println(outFile)
}

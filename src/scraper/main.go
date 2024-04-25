package main

import (
  "fmt"
  "net/http"
  "os"
  "sort"
  "strconv"
  "strings"
  "sync"
)

const urlPrefix = "https://www.gstatic.com/prettyearth/assets/data/v3/"

// Function to check if a given URL is valid
func checkUrl(url string, wg *sync.WaitGroup, ch chan<- string) {
  defer wg.Done()
  fmt.Sprint("Processing %s", url)

  response, err := http.Get(url)
  if err != nil {
    return
  }
  defer response.Body.Close()

  if response.StatusCode >= 400 {
    if response.StatusCode != 404 {
      fmt.Println("Error status:", response.StatusCode)
    }
    return
  }

  ch <- url
}

// Function to extract the id from the URL
func extractId(url string) int {
  paths := strings.Split(url, "/")
  basename := paths[len(paths) - 1]
  idStr := strings.TrimSuffix(basename, ".json")

  id, err := strconv.Atoi(idStr)
  if err != nil {
    panic(err)
  }

  return id
}

func main() {
  var wg sync.WaitGroup
  ch := make(chan string, 14000)

  // Check URLs from 1000 to 15000 in parallel
  for i := 1000; i <= 15000; i++ {
    wg.Add(1)
    go checkUrl(urlPrefix + strconv.Itoa(i) + ".json", &wg, ch)
  }

  wg.Wait()
  close(ch)

  // Prepare output file
  cwd, err := os.Getwd()
  if err != nil {
    fmt.Println("Error getting current working directory:", err)
    return
  }

  file, err := os.Create(cwd + "/_earthview.txt")
  if err != nil {
    fmt.Println("Error creating output file:", err)
    return
  }
  defer file.Close()

  // Sort URLs numerically by their last path
  var urls []string
  for str := range ch {
    urls = append(urls, str)
  }

  sort.Slice(urls, func(i, j int) bool {
    id1 := extractId(urls[i])
    id2 := extractId(urls[j])
    return id1 < id2
  })

  // Write each URL to output file on single line
  for _, url := range urls {
    _, err := file.WriteString(url + "\n")
    if err != nil {
      continue
    }
  }

  fmt.Println("Data written to file successfully.")
}

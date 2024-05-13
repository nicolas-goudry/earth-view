/*
Copyright © 2024 Nicolas Goudry <goudry.nicolas@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package list

import (
  "encoding/json"
  "fmt"
  "os"
  "sort"
  "sync"

  "earth-view/lib"
  "github.com/charmbracelet/bubbles/progress"
  tea "github.com/charmbracelet/bubbletea"
)

var program *tea.Program

type fetcher struct {
  done            bool
  errors          []error
  results         []int
  onFetchProgress func(fetchProgress)
}
type fetchProgress struct {
  percent float64
  results []result
}
type result struct {
  id    int
  error error
}

func (f *fetcher) Start() {
  var results []result
  idsCount := lib.KnownIdUpperBoundary - lib.KnownIdLowerBoundary

  for i := lib.KnownIdLowerBoundary; i < lib.KnownIdUpperBoundary; i += batchSize {
    chunkSize := batchSize
    chunkEnd := i + batchSize
    if chunkEnd > lib.KnownIdUpperBoundary {
      chunkSize = lib.KnownIdUpperBoundary - i
    }

    chunk := make([]int, chunkSize)
    for j := 0; j < chunkSize; j++ {
      chunk[j] = i + j
    }

    var wg sync.WaitGroup
    wg.Add(1)

    chunkResults := make(chan result, chunkSize)
    go func(ids []int) {
      defer wg.Done()
      f.fetchChunk(ids, chunkResults)
    }(chunk)

    wg.Wait()
    close(chunkResults)

    for result := range chunkResults {
      results = append(results, result)
      if result.error != nil {
        f.errors = append(f.errors, result.error)
      }
    }

    f.onFetchProgress(fetchProgress{
      percent: float64(len(results)) / float64(idsCount),
      results: results,
    })
  }

  for _, result := range results {
    if result.error == nil {
      f.results = append(f.results, result.id)
    }
  }

  f.done = true
}

func (f *fetcher) fetchChunk(ids []int, ch chan<- result) {
  for _, id := range ids {
    asset := lib.Asset{ Id: id }
    if err := asset.Fetch(retry); err != nil {
      ch <- result{ id: id, error: err }
    } else {
      ch <- result{ id: id }
    }
  }
}

func generateJSONContent(assetIds []int) ([]byte, error) {
  sort.Slice(assetIds, func(i, j int) bool {
    return assetIds[i] < assetIds[j]
  })

  json, err := json.Marshal(assetIds)
  if err != nil {
    return nil, err
  }

  return json, nil
}

func main() {
  program = tea.NewProgram(model{
    progress: progress.New(progress.WithDefaultGradient()),
  })

  f := &fetcher{
    onFetchProgress: func (progress fetchProgress) {
      program.Send(progressMsg(progress))
    },
  }

  go f.Start()

  if _, err := program.Run(); err != nil {
    fmt.Println("error running program:", err)
    os.Exit(1)
  }

  if f.done == true {
    json, err := generateJSONContent(f.results)
    if err != nil {
      if quiet == false {
        fmt.Println(err)
      }

      os.Exit(1)
    }

    if output == "" {
      fmt.Println(string(json))
    } else {
      outPath, err := lib.WriteFile(json, output, "earth-view.json")
      if err != nil {
        if quiet == false {
          fmt.Println(err)
        }

        os.Exit(1)
      }

      if quiet == false {
        fmt.Printf("List is available at %s\n", outPath)
      }
    }
  }
}

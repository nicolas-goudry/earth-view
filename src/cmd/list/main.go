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
  "strings"
  "sync"

  "earth-view/lib"
  "github.com/charmbracelet/bubbles/progress"
  tea "github.com/charmbracelet/bubbletea"
)

var program *tea.Program

// fetcher struct is used to track the state of the fetching process
type fetcher struct {
  done            bool
  errors          []error
  results         []int
  onFetchProgress func(fetchProgress)
}

// fetchProgress struct is used to report the fetch progress to the TUI program
type fetchProgress struct {
  percent float64
  results []result
}

// result struct is used to hold fetch result information
type result struct {
  id    int
  error error
}

// Start fetching all assets
func (f *fetcher) Start() {
  var results []result
  totalIds := lib.KnownIdUpperBoundary - lib.KnownIdLowerBoundary

  // Loop over all known ids in batches
  for i := lib.KnownIdLowerBoundary; i < lib.KnownIdUpperBoundary; i += batchSize {
    // Compute the real size of this batch
    chunkSize := batchSize
    chunkEnd := i + batchSize
    if chunkEnd > lib.KnownIdUpperBoundary {
      chunkSize = lib.KnownIdUpperBoundary - i
    }

    // Create a slice of ids to fetch in this batch
    chunk := make([]int, chunkSize)
    for j := 0; j < chunkSize; j++ {
      chunk[j] = i + j
    }

    // Start a WaitGroup to process the batch elements concurrently before next batch
    var wg sync.WaitGroup
    wg.Add(1)

    // Process the batch elements
    chunkCh := make(chan result, chunkSize)
    go func(ids []int) {
      defer wg.Done()
      f.fetchChunk(ids, chunkCh)
    }(chunk)

    // Wait for batch process results
    wg.Wait()
    close(chunkCh)

    // Process batch results
    var chunkResults []result
    for result := range chunkCh {
      // chunkResults is used to report the current batch results to the TUI program
      chunkResults = append(chunkResults, result)

      // Save result for later use when all batches are complete
      results = append(results, result)
    }

    // Report the fetch progress to the TUI program
    f.onFetchProgress(fetchProgress{
      percent: float64(len(results)) / float64(totalIds),
      results: chunkResults,
    })
  }

  // All batches have been processed, split results in actual results and errors to be reported to user
  for _, result := range results {
    if result.error == nil {
      f.results = append(f.results, result.id)
    } else if ! strings.Contains(result.error.Error(), "not found") {
      f.errors = append(f.errors, result.error)
    }
  }

  // Mark fetching as done
  // This is needed to avoid outputting partial results/errors to user in case Ctrl+C was pressed in TUI program
  f.done = true
}

// Fetch assets and report results in channel
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

// Generate a JSON array of given asset ids
func generateJSONContent(assetIds []int) ([]byte, error) {
  // Sort assets ascending by their id
  sort.Slice(assetIds, func(i, j int) bool {
    return assetIds[i] < assetIds[j]
  })

  json, err := json.Marshal(assetIds)
  if err != nil {
    return nil, err
  }

  return json, nil
}

// Execute program
func main() {
  // Create tea program with initial model
  program = tea.NewProgram(model{
    progress: progress.New(progress.WithDefaultGradient()),
  })

  // Create a fetcher instance
  f := &fetcher{
    onFetchProgress: func (progress fetchProgress) {
      // Send a progressMsg with actual progress
      program.Send(progressMsg(progress))
    },
  }

  // Start fetching assets
  go f.Start()

  // Start tea program
  if _, err := program.Run(); err != nil {
    fmt.Fprintf(os.Stderr, "error running program: %v\n", err)
    os.Exit(1)
  }

  // Handle results if fetch is done
  if f.done {
    // Report errors if any and not quiet
    if len(f.errors) > 0 && quiet == false {
      fmt.Fprintln(os.Stderr, "Encountered the following errors:")

      for _, err := range f.errors {
        fmt.Fprintf(os.Stderr, "  %v\n", err)
      }
      fmt.Fprintln(os.Stderr, "")
    }

    // Handle empty results slice
    if len(f.results) == 0 {
      if quiet == false {
        fmt.Fprintln(os.Stderr, "No results to save")
      }

      os.Exit(1)
    }

    // Generate JSON array from results
    json, err := generateJSONContent(f.results)
    if err != nil {
      if quiet == false {
        fmt.Fprintln(os.Stderr, err)
      }

      os.Exit(1)
    }

    // Save results to stdout or given file
    if output == "" {
      fmt.Println(string(json))
    } else {
      filePath, err := lib.ResolveAbsFilePath(output, "earth-view.json")
      if err != nil {
        if quiet == false {
          fmt.Fprintln(os.Stderr, err)
        }

        os.Exit(1)
      }

      err = lib.WriteFile(json, filePath)
      if err != nil {
        if quiet == false {
          fmt.Fprintln(os.Stderr, err)
        }

        os.Exit(1)
      }

      // Report location of file containing results
      if quiet == false {
        fmt.Printf("Results saved to %s\n", filePath)
      }
    }
  } else {
    // Something went wrong, fetching is not done, exit with non-zero status code
    os.Exit(1)
  }
}

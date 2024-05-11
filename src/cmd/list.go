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
package cmd

import (
  "encoding/json"
  "fmt"
  "os"
  "sort"
  "strings"
  "sync"

  "earth-view/lib"
	"github.com/spf13/cobra"
)

var (
  batchSize int
  retry     int
  output    string

  listCmd = &cobra.Command{
    Use:   "list",
    Aliases: []string{"ls"},
    Short: "Get images list",
    Long: `Get a list of Google Earth View images

Description:
  This command will try to fetch images from gstatic.com using a known range of
  possible identifiers and generate a JSON array of valid identifiers for images.

  If the fetch succeeds, the image is added to the list.
  If the fetch fails with a 404 HTTP status code, the image is skipped.
  If the fetch fails with a non 200 HTTP status code, the error is reported and
  the image is skipped.

  By default, the generated list is output to the standard output. This
  behaviour can be changed by using the '--output' flag.`,
    DisableFlagsInUseLine: true,
    SilenceUsage: true,
    Args: cobra.MaximumNArgs(0),
    Run: func(_ *cobra.Command, _ []string) {
      assets := fetchAssets()
      json, err := generateJSONContent(assets)
      cobra.CheckErr(err)

      if output == "" {
        fmt.Println(string(json))
      } else {
        outPath, err := lib.WriteFile(json, output)
        cobra.CheckErr(err)

        fmt.Printf("List is available at %s\n", outPath)
      }
    },
  }
)

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().IntVarP(&batchSize, "batch-size", "b", 20, `number of parallel calls to gstatic.com
Using a high value may result in potentially wrong failures to fetch images`)
  listCmd.Flags().IntVarP(&retry, "retry", "r", 3, "number of retries before skipping an image in case of non 200 HTTP status code")
  listCmd.Flags().StringVarP(&output, "output", "o", "", "write to file instead of stdout")
}

func generateJSONContent(assets []lib.Asset) ([]byte, error) {
  var content []int
  for _, asset := range assets {
    content = append(content, asset.Id)
  }

  sort.Slice(content, func(i, j int) bool {
    return content[i] < content[j]
  })

  json, err := json.Marshal(content)
  if err != nil {
    return nil, err
  }

  return json, nil
}

func fetchAssets() []lib.Asset {
  var wg sync.WaitGroup
  ch := make(chan lib.Asset, lib.KnownIdUpperBoundary - lib.KnownIdLowerBoundary)

  for i := lib.KnownIdLowerBoundary; i < lib.KnownIdUpperBoundary; i += batchSize {
    chunkEnd := i + batchSize
    if chunkEnd > lib.KnownIdUpperBoundary {
      chunkEnd = lib.KnownIdUpperBoundary
    }

    chunk := make([]int, batchSize)
    for j := 0; j < batchSize; j++ {
      chunk[j] = i + j
    }

    wg.Add(1)
    go fetchChunk(chunk, &wg, ch)
  }

  go func() {
    wg.Wait()
    close(ch)
  }()

  assets := make([]lib.Asset, len(ch))
  for asset := range ch {
    assets = append(assets, asset)
  }

  return assets
}

func fetchChunk(chunk []int, wg *sync.WaitGroup, ch chan<- lib.Asset) {
  defer wg.Done()

  for _, id := range chunk {
    asset := lib.Asset{ Id: id }
    if err := asset.Fetch(retry); err != nil {
      if ! strings.Contains(err.Error(), "not found") {
        fmt.Fprintln(os.Stderr, err)
      }

      continue
    }

    ch <- asset
  }
}

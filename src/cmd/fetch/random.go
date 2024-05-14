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
package fetch

import (
  "encoding/json"
  "fmt"
  "math/rand"
  "os"
  "strconv"

  "earth-view/lib"
  "github.com/spf13/cobra"
)

var (
  input string

  randomCmd = &cobra.Command{
    Use: "random",
    Aliases: []string{"rnd", "rand"},
    Short: "Fetch random images",
    Long: fmt.Sprintf(`Download a random Google Earth View image.

Description:
  This command will download a random image using either the known possible
  image identifiers or an input file containing an array of identifiers in JSON
  format.

%s

  When '--input' flag is provided, the command expects it to be fed with the
  output of the 'list' command, either via standard input or the path to a file.

  When '--input' flag is not provided, a random image identifier will be chosen
  from the known range of possible identifiers. If the selected identifier is
  not valid, another one will be chosen, until a valid identifier is found.

%s`, helpText.process, helpText.output),
    DisableFlagsInUseLine: true,
    SilenceUsage: true,
    Args: cobra.MaximumNArgs(0),
    Run: func(cmd *cobra.Command, args []string) {
      asset := lib.Asset{}
      filePath, err := fetchRandomAsset(&asset)
      cobra.CheckErr(err)

      // Only write file if it does not yet exist or if overwrite is set
      if lib.FileExists(filePath) == false || overwrite {
        err := lib.WriteFile(asset.Content, filePath)
        cobra.CheckErr(err)
      }

      fmt.Println(filePath)
    },
  }
)

func init() {
  fetchCmd.AddCommand(randomCmd)

  randomCmd.Flags().StringVarP(&input, "input", "i", "", "input file to choose an image from, or standard input if not specified")
  addCommonFlags(randomCmd.Flags())
}

func pickRandomId() (int, error) {
  if input == "" {
    return lib.KnownIdLowerBoundary + rand.Intn(lib.KnownIdUpperBoundary - lib.KnownIdLowerBoundary + 1), nil
  }

  ids, err := readInput()
  if err != nil {
    return -1, err
  }

  return ids[rand.Intn(len(ids))], nil
}

func fetchRandomAsset(asset *lib.Asset) (string, error) {
  randomId, err := pickRandomId()
  if err != nil {
    return "", err
  }

  filePath, err := lib.ResolveAbsFilePath(output, strconv.Itoa(randomId) + ".jpeg")
  if err != nil {
    return "", err
  }

  // Only fetch file if it does not yet exist or if overwrite is set
  if lib.FileExists(filePath) == false || overwrite {
    asset.Id = randomId
    if _, err := asset.GetContent(); err != nil {
      return fetchRandomAsset(asset)
    }
  }

  return filePath, nil
}

func readInput() ([]int, error) {
  content, err := os.ReadFile(input)
  if err != nil {
    return nil, err
  }

  var ids []int
  if err := json.Unmarshal(content, &ids); err != nil {
    return nil, err
  }

  return ids, nil
}

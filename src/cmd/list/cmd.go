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
  "fmt"

  "earth-view/cmd"
	"github.com/spf13/cobra"
)

var (
  batchSize int
  output    string
  quiet     bool
  retry     int

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
  behaviour can be changed by using the '--output' flag. If the provided value
  is a directory, the file will be named 'earth-view.json'.`,
    DisableFlagsInUseLine: true,
    SilenceUsage: true,
    Args: cobra.MaximumNArgs(0),
    PreRunE: func(_ *cobra.Command, _ []string) error {
      if output == "" && quiet {
        return fmt.Errorf("--quiet cannot be provided when --output is not set")
      }

      return nil
    },
    Run: func(_ *cobra.Command, _ []string) {
      main()
    },
  }
)

func init() {
	cmd.RootCmd.AddCommand(listCmd)

	listCmd.Flags().IntVarP(&batchSize, "batch-size", "b", 20, `number of parallel calls to gstatic.com
Using a high value may result in potentially wrong failures to fetch images`)
  listCmd.Flags().StringVarP(&output, "output", "o", "", "write to file instead of stdout")
  listCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "do not output anything")
  listCmd.Flags().IntVarP(&retry, "retry", "r", 3, "number of retries before skipping an image in case of non 200 HTTP status code")
}

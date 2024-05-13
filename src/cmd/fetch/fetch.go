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
	"fmt"
  "strconv"

  "earth-view/cmd"
  "earth-view/lib"
	"github.com/spf13/cobra"
)

var (
  idNumeric int

  fetchCmd = &cobra.Command{
    Use:   "fetch identifier",
    Aliases: []string{"get", "download", "dl"},
    Short: "Fetch images",
    Long: fmt.Sprintf(`Download a Google Earth View image by its identifier.

Description:
%s

%s`, helpText.process, helpText.output),
    DisableFlagsInUseLine: true,
    SilenceUsage: true,
    Args: func(cmd *cobra.Command, args []string) error {
      if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
        return fmt.Errorf("missing required argument 'identifier'")
      }

      id, err := strconv.Atoi(args[0])
      if err != nil {
        return fmt.Errorf("invalid identifier provided: %s. Identifier must be a number", args[0])
      }

      idNumeric = id

      return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
      defaultFilename := args[0] + ".jpeg"

      filePath, err := lib.ResolveAbsFilePath(output, defaultFilename)
      cobra.CheckErr(err)

      // Only fetch and write file if it does not yet exist or if overwrite is set
      if lib.FileExists(filePath) == false || overwrite {
        assetContent, err := fetchAssetContent(idNumeric)
        cobra.CheckErr(err)

        err = lib.WriteFile(assetContent, filePath)
        cobra.CheckErr(err)
      }

      fmt.Println(filePath)
    },
  }
)

func init() {
  cmd.RootCmd.AddCommand(fetchCmd)

  addCommonFlags(fetchCmd.Flags())
}

func fetchAssetContent(id int) ([]byte, error) {
  asset := lib.Asset{ Id: id }

  content, err := asset.GetContent()
  if err != nil {
    return nil, err
  }

  return content, nil
}

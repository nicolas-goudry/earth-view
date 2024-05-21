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
    Use: "fetch identifier",
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

      return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
      filePath, err := runFetchCmd(args[0], output, overwrite)
      cobra.CheckErr(err)
      fmt.Println(filePath)
    },
  }
)

func init() {
  cmd.RootCmd.AddCommand(fetchCmd)

  addCommonFlags(fetchCmd.Flags())
}

func runFetchCmd(id string, output string, overwrite bool) (string, error) {
  idNumeric, err := strconv.Atoi(id)
  if err != nil {
    return "", fmt.Errorf("invalid identifier provided: %s. Identifier must be a number", id)
  }

  filePath, err := lib.ResolveAbsFilePath(output, id + ".jpeg")
  if err != nil {
    return "", err
  }

  // Only fetch and write file if it does not yet exist or if overwrite is set
  if lib.FileExists(filePath) == false || overwrite {
    asset := lib.Asset{ Id: idNumeric }
    content, err := asset.GetContent()
    if err != nil {
      return "", err
    }

    err = lib.WriteFile(content, filePath)
    if err != nil {
      return "", err
    }
  }

  return filePath, nil
}

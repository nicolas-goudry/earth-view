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
package lib

import (
  "os"
  "path"
  "path/filepath"
)

// WriteFile writes some content to a given path.
// If the path is a directory, the file is written to this directory with the defaultFilename as its filename.
// If the path is empty, the file is written to the current working directory with the defaultFilename as its filename.
// It returns the absolute path to the file written.
func WriteFile(content []byte, outPath string, defaultFilename string) (string, error) {
  if stat, err := os.Stat(outPath); err == nil {
    if stat.IsDir() {
      outPath = path.Join(outPath, defaultFilename)
    }
  }

  if outPath == "" {
    outPath = defaultFilename
  }

  if err := os.WriteFile(outPath, content, 0644); err != nil {
    return "", err
  }

  finalPath, err := filepath.Abs(outPath)
  if err != nil {
    return "", err
  }

  return finalPath, nil
}

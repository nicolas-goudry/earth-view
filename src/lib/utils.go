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

func IsPathValid(path string) bool {
  if _, err := os.Stat(path); err == nil {
    return true
  }

  var d []byte
  if err := os.WriteFile(path, d, 0644); err == nil {
    os.Remove(path)
    return true
  }

  return false
}

func WriteFile(content []byte, outPath string, defaultFilename string) (string, error) {
  stat, err := os.Stat(outPath)
  if err != nil {
    return "", err
  }

  if stat.IsDir() {
    outPath = path.Join(outPath, defaultFilename)
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

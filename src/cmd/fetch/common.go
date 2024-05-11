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
  "github.com/spf13/pflag"
)

var (
  output string

  helpText = struct {
    process string
    output string
  }{
    process: `  The image metadata is first retrieved from gstatic.com, the server hosting the
  images assets, then the image is decoded before being saved on the filesystem.`,
    output: `  By default, the image is saved in the current working directory and its
  identifier is used as the filename. This behaviour can be changed by using the
  '--output' flag. If the provided value is a directory, the file is saved into
  it and its identifier is used as the filename.`,
  }
)

func addOutputFlag(f *pflag.FlagSet) {
  f.StringVarP(&output, "output", "o", "", "write to given file path instead of current working directory")
}

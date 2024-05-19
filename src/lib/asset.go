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
  "encoding/base64"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "strconv"
  "strings"
  "time"
)

// Asset represents a Google Earth View asset
type Asset struct {
  Id       int
  raw      []byte
  Metadata map[string]interface{}
  Content  []byte
}

// Fetch tries to fetch an asset and save the response content in the raw field
func (a *Asset) Fetch(retry int) error {
  client := http.Client{
    Timeout: 5 * time.Second,
  }
  response, err := client.Get(baseUrl + "/" + strconv.Itoa(a.Id) + ".json")
  if err != nil {
    return err
  }
  defer response.Body.Close()

  if response.StatusCode == http.StatusNotFound {
    return fmt.Errorf("[%d] fetch failed: asset not found", a.Id)
  }

  if response.StatusCode != http.StatusOK {
    if retry > 0 {
      return a.Fetch(retry - 1)
    }

    return fmt.Errorf("[%d] fetch failed: received HTTP %d", a.Id, response.StatusCode)
  }

  body, err := io.ReadAll(response.Body)
  if err != nil {
    return fmt.Errorf("[%d] fetch failed: error while parsing response body: %s", a.Id, err)
  }

  a.raw = body

  return nil
}

// GetMetadata parses the asset JSON content and stores it in the Metadata field
// If no data is available yet for the asset, it will call Fetch by itself
func (a *Asset) GetMetadata() (map[string]interface{}, error) {
  if a.raw == nil {
    if err := a.Fetch(0); err != nil {
      return nil, err
    }
  }

  var metadata map[string]interface{}

  if err := json.Unmarshal(a.raw, &metadata); err != nil {
    return nil, err
  }

  a.Metadata = metadata

  return a.Metadata, nil
}

// GetContent parses and decode the actual asset image from its metadata and stores it in the Content field
// If no metadata is available yet for the asset, it will call GetMetadata by itself
func (a *Asset) GetContent() ([]byte, error) {
  if a.Metadata == nil {
    if _, err := a.GetMetadata(); err != nil {
      return nil, err
    }
  }

  dataUri := fmt.Sprint(a.Metadata["dataUri"])
  if dataUri == "" {
    return nil, fmt.Errorf("missing 'dataUri' field on asset metadata")
  }

  dataUriSplit := strings.Split(dataUri, ",")
  encodedImg := dataUriSplit[len(dataUriSplit) -1]
  if encodedImg == "" {
    return nil, fmt.Errorf("failed to decode image")
  }

  decodedImg, err := base64.StdEncoding.DecodeString(encodedImg)
  if err != nil {
    return nil, fmt.Errorf("failed to decode image: %s", err)
  }

  a.Content = decodedImg

  return a.Content, nil
}

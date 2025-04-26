// internal/fetchpage.go
package internal

import (
  "context"
  "fmt"
  "net/http"
  "strings"

  "github.com/PuerkitoBio/goquery"
  "github.com/pontus-devoteam/agent-sdk-go/pkg/tool"
)

func NewFetchPageTool() tool.Tool {
  return tool.NewFunctionTool(
    "fetch_page",
    "Download a URL and return the main text content",
    func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
      urlStr, ok := params["url"].(string)
      if !ok {
        return nil, fmt.Errorf("parameter 'url' must be a string")
      }
      resp, err := http.Get(urlStr)
      if err != nil {
        return nil, err
      }
      defer resp.Body.Close()

      doc, err := goquery.NewDocumentFromReader(resp.Body)
      if err != nil {
        return nil, err
      }

      // Extract all paragraphs as text
      var sb strings.Builder
      doc.Find("p").Each(func(i int, s *goquery.Selection) {
        text := strings.TrimSpace(s.Text())
        if text != "" {
          sb.WriteString(text + "\n\n")
        }
      })
      out := sb.String()
      if out == "" {
        out = "⚠️ No <p> tags found—page may be non-HTML or JS-rendered."
      }
      return out, nil
    },
  ).WithSchema(map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
      "url": map[string]interface{}{
        "type":        "string",
        "description": "The URL to fetch and extract text from",
      },
    },
    "required": []string{"url"},
  })
}

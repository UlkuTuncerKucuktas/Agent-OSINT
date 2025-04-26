package internal

import (
  "context"
  "fmt"
  "net/http"
  "net/url"
  "sync"
  "time"

  "github.com/PuerkitoBio/goquery"
  "github.com/pontus-devoteam/agent-sdk-go/pkg/tool"
)

// Result is one search hit.
type Result struct {
  Source  string `json:"source"`
  Title   string `json:"title"`
  URL     string `json:"url"`
  Snippet string `json:"snippet"`
}


func NewMultiSearchTool() tool.Tool {
  return tool.NewFunctionTool(
    "multi_search",
    "Search Google, DuckDuckGo, and Bing concurrently, aggregate and dedupe results",
    func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
      rawQuery, ok := params["query"].(string)
      if !ok {
        return nil, fmt.Errorf("parameter 'query' must be a string")
      }
      // Prepare endpoints
      q := url.QueryEscape(rawQuery)
      engines := []struct {
        name string
        url  string
        ua   string
      }{
        {
          name: "google",
          url:  "https://www.google.com/search?q=" + q + "&gl=tr&hl=tr&num=10",
          ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
        },
        {
          name: "duckduckgo",
          url:  "https://html.duckduckgo.com/html/?q=" + q,
          ua:   "Mozilla/5.0",
        },
        {
          name: "bing",
          url:  "https://www.bing.com/search?q=" + q + "&setlang=tr",
          ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
        },
      }

      // helper to fetch & parse one engine
      fetch := func(name, urlStr, ua string, ch chan<- Result, wg *sync.WaitGroup) {
        defer wg.Done()
        client := http.Client{Timeout: 5 * time.Second}
        req, err := http.NewRequest("GET", urlStr, nil)
        if err != nil {
          return
        }
        req.Header.Set("User-Agent", ua)
        resp, err := client.Do(req)
        if err != nil {
          return
        }
        defer resp.Body.Close()

        doc, err := goquery.NewDocumentFromReader(resp.Body)
        if err != nil {
          return
        }

        switch name {
        case "google":
          doc.Find("div.yuRUbf").Each(func(i int, s *goquery.Selection) {
            if i >= 10 {
              return
            }
            a := s.Find("a")
            href, _ := a.Attr("href")
            title := a.Text()
            snippet := s.Parent().Find("div.IsZvec span").First().Text()
            ch <- Result{name, title, href, snippet}
          })
        case "duckduckgo":
          doc.Find("div.result__body").Each(func(i int, s *goquery.Selection) {
            if i >= 10 {
              return
            }
            a := s.Find("a.result__a")
            href, _ := a.Attr("href")
            title := a.Text()
            snippet := s.Find(".result__snippet").Text()
            ch <- Result{name, title, href, snippet}
          })
        case "bing":
          doc.Find("li.b_algo").Each(func(i int, s *goquery.Selection) {
            if i >= 10 {
              return
            }
            a := s.Find("h2 > a")
            href, _ := a.Attr("href")
            title := a.Text()
            snippet := s.Find(".b_caption p").Text()
            ch <- Result{name, title, href, snippet}
          })
        }
      }

      // run them all
      ch := make(chan Result, 30)
      var wg sync.WaitGroup
      for _, e := range engines {
        wg.Add(1)
        go fetch(e.name, e.url, e.ua, ch, &wg)
      }
      go func() {
        wg.Wait()
        close(ch)
      }()

      // collect & dedupe
      seen := map[string]bool{}
      var out []Result
      for r := range ch {
        if r.URL == "" || seen[r.URL] {
          continue
        }
        seen[r.URL] = true
        out = append(out, r)
      }
      return out, nil
    },
  ).WithSchema(map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
      "query": map[string]interface{}{
        "type":        "string",
        "description": "What to search for",
      },
    },
    "required": []string{"query"},
  })
}

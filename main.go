package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/cache/diskcache"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"github.com/geziyor/geziyor/metrics"
)

type Course struct {
	Title       string
	Description string
	Creator     string
	Level       string
	URL         string
	Language    string
	Commitment  string
	Rating      string
}

func main() {

	// base := "https://boardgamegeek.com"
	url := "https://boardgamegeek.com/browse/boardgame"

	func TestAllLinksWithRender(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}
	
		geziyor.NewGeziyor(&geziyor.Options{
			AllowedDomains: []string{"boardgamegeek.com"},
			StartURLs:      []string{"https://boardgamegeek.com/browse/boardgame"},
			ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
				g.Exports <- []string{r.Request.URL.String()}
	
				r.HTMLDoc.Find("div.game-header-body").Each(func(i int, s *goquery.Selection) {
					fmt.Println("JUEGO encontrado")
					if strings.TrimSpace(s.Find("h1 > a").Text()) != "" {
	
						//fmt.Println("Juego: ", strings.TrimSpace(s.Find("h1 > a").Text()))
						//fmt.Println("Description: ", s.Find("div.game-header-title-info > p").Text())
						// fmt.Println("Short description", s.Find(""))
						var p, m string
						s.Find("ul.gameplay > li").Each(func(i int, s *goquery.Selection) {
							txt := s.Find("div").Text()
							switch {
							case strings.Contains(txt, "Players"):
								p = strings.ReplaceAll(strings.TrimSpace(txt), `	`, "")
							case strings.Contains(txt, "Min"):
								m = strings.TrimSpace(txt)
							}
						})
	
						g.Exports <- map[string]interface{}{
							"number":      i,
							"juego":       strings.TrimSpace(s.Find("h1 > a").Text()),
							"Description": s.Find("div.game-header-title-info > p").Text(),
							"Players":     p,
							"Min":         m,
						}
	
						// TODO add big description
					}
				})
	
				r.HTMLDoc.Find("a.primary").Each(func(i int, s *goquery.Selection) {
					if href, ok := s.Attr("href"); ok {
						absoluteURL, _ := r.Request.URL.Parse(href)
						switch {
						case strings.Contains(absoluteURL.String(), "/browse/"):
							fmt.Println("url -> ", absoluteURL.String())
							//g.Get(absoluteURL.String(), g.Opt.ParseFunc)
						case strings.Contains(absoluteURL.String(), "/boardgame/"):
							fmt.Println("game url -> ", absoluteURL.String())
							g.GetRendered(absoluteURL.String(), g.Opt.ParseFunc)
						}
	
					}
				})
			},
			Exporters:       []export.Exporter{&export.JSONLine{FileName: "1.json"}, &export.JSON{FileName: "2.json"}},
			Cache:           diskcache.New(".new"),
			BrowserEndpoint: "ws://localhost:9002",
			MetricsType:     metrics.Prometheus,
		}).Start()
	}
	
	//c.Visit(url)

}

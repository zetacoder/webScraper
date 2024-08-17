package main

import (
	"encoding/json"

	"github.com/gocolly/colly/v2"
	"github.com/zetacoder/webScraper/scraper"
)

var domain string = "https://news.ycombinator.com/"
var maxEntries int = 30

func main() {
	c := colly.NewCollector()

	s := scraper.NewScraper(c, domain, maxEntries)

	s.Scrape()

	for _, p := range s.Posts {
		jsoned, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			panic(err)
		}

		println(string(jsoned))
		println("------------------------------------------------")
	}

}

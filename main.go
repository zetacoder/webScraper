package main

import (
	"log"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
	"github.com/zetacoder/webScraper/database"
	"github.com/zetacoder/webScraper/scraper"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("cant load env variables: ", err.Error())
	}

	if err := database.InitDB(); err != nil {
		log.Fatal("cant initilize db: ", err.Error())
	}

	if err := database.AutoMigrateTables(&scraper.Post{}, &scraper.UsageData{}); err != nil {
		log.Fatal("error migrating models into db tables")
	}
}

func main() {
	c := colly.NewCollector()

	s := scraper.NewScraper(c, scraper.Domain, scraper.MaxEntries, database.DB)

	if err := s.Scrape(); err != nil {
		log.Fatal("error scraping: ", err.Error())
	}

	if err := s.FilterPosts(); err != nil {
		log.Fatal("error filtering posts: ", err.Error())
	}

	if err := s.SavePostsInDB(); err != nil {
		log.Fatal("error saving entries in db: ", err.Error())
	}

	if err := s.SaveUsageData(); err != nil {
		log.Fatal("error saving usage data in db: ", err.Error())
	}
}

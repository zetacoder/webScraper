package scraper

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Post struct {
	Title                  string
	IsTitleLenghtMoreThan5 bool
	Points                 int
	Comments               int
	Source                 string
	SourceDomain           string
}

type Scraper struct {
	Collector  *colly.Collector
	Domain     string
	Posts      []Post
	MaxEntries int
}

func NewScraper(c *colly.Collector, domain string, maxEntries int) Scraper {
	return Scraper{
		Collector: c,
		Domain:    domain,
	}
}

func (s *Scraper) Scrape() {
	s.Posts = []Post{}

	s.Collector.OnHTML(`tbody`, func(e *colly.HTMLElement) {
		e.ForEach(`tr.athing`, func(idx int, el *colly.HTMLElement) {
			title := s.getTitle(el)
			source := s.getSource(el)
			sourceDomain := s.getSourceDomain(el)
			points := s.getPoints(el)
			comments := s.getComments(el)
			isTitleLenghtMoreThan5 := countWords(title) > 5

			s.Posts = append(s.Posts, Post{
				Title:                  title,
				Source:                 source,
				SourceDomain:           sourceDomain,
				Points:                 points,
				Comments:               comments,
				IsTitleLenghtMoreThan5: isTitleLenghtMoreThan5,
			})
		})
	})

	s.Collector.Visit(s.Domain)
}

func (s *Scraper) getTitle(el *colly.HTMLElement) string {
	title := ""
	el.ForEach(`td.title span.titleline a`, func(idx int, el *colly.HTMLElement) {
		if idx == 0 {
			title = el.Text
		}
	})

	return title
}

func (s *Scraper) getSource(el *colly.HTMLElement) string {
	return el.ChildAttr(`td.title span.titleline a`, "href")
}

func (s *Scraper) getSourceDomain(el *colly.HTMLElement) string {
	sourceDomain := ""
	el.ForEach(`td.title span.titleline a`, func(idx int, el *colly.HTMLElement) {
		if idx == 1 {
			sourceDomain = el.Text
		}
	})
	return sourceDomain
}

func (s *Scraper) getPoints(el *colly.HTMLElement) int {
	points := 0
	pointsText := el.DOM.Next().Find(`td.subtext span.score`).Text()
	if pointsText != "" {
		points, _ = strconv.Atoi(strings.TrimSpace(strings.Split(pointsText, "points")[0]))
	}
	return points
}

func (s *Scraper) getComments(el *colly.HTMLElement) int {
	comments := 0
	commentsText := el.DOM.Next().Find(`td.subtext a`).Last().Text()
	if strings.Contains(commentsText, "comments") {
		comments, _ = strconv.Atoi(strings.TrimSpace(strings.Split(commentsText, "comments")[0]))
	}
	return comments
}

func countWords(input string) int {
	// Define a regular expression to match words, allowing periods, hyphens, and apostrophes within them
	re := regexp.MustCompile(`[a-zA-Z0-9]+(?:[-'\.][a-zA-Z0-9]+)*`)

	// Find all matches of the regular expression in the input string
	matches := re.FindAllString(input, -1)

	// Return the count of matched words
	return len(matches)
}

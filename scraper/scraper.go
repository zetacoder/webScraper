package scraper

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"gorm.io/gorm"
)

const Domain string = "https://news.ycombinator.com/"

const MaxEntries int = 30

type Scraper struct {
	Collector   *colly.Collector
	Domain      string
	Posts       []Post
	MaxEntries  int
	Database    *gorm.DB
	UniquePosts map[string]bool
	UsageData   UsageData
}

type Post struct {
	gorm.Model
	Title         string `gorm:"unique"`
	Points        int
	Comments      int
	Source        string
	SourceDomain  string
	AppliedFilter string
}

// NewScraper creates a new Scraper instance
func NewScraper(c *colly.Collector, domain string, maxEntries int, db *gorm.DB) Scraper {
	return Scraper{
		Collector:  c,
		Domain:     domain,
		Database:   db,
		MaxEntries: maxEntries,
		UsageData:  UsageData{},
	}
}

// Scrape scrapes the website and collects the posts
func (s *Scraper) Scrape() error {
	startRequest := time.Now()
	s.UsageData.StartedAt = startRequest

	s.Posts = []Post{}
	s.UniquePosts = map[string]bool{}

	s.Collector.OnHTML(`tbody`, func(e *colly.HTMLElement) {

		// get bytes scraped
		s.UsageData.TotalBytesScraped = len(e.Response.Body)

		e.ForEach(`tr.athing`, func(idx int, el *colly.HTMLElement) {
			title := s.getTitle(el)

			if s.isInDB(title) {
				return
			}

			if s.UniquePosts[title] {
				return
			}

			source := s.getSource(el)
			sourceDomain := s.getSourceDomain(el)
			points := s.getPoints(el)
			comments := s.getComments(el)

			post := Post{
				Title:        title,
				Source:       source,
				SourceDomain: sourceDomain,
				Points:       points,
				Comments:     comments,
			}

			s.Posts = append(s.Posts, post)

			s.UniquePosts[post.Title] = true

			if len(s.Posts) >= s.MaxEntries {
				return
			}
		})
	})

	err := s.Collector.Visit(s.Domain)
	if err != nil {
		return err
	}

	return nil
}

// SavePostsInDB saves the posts in the database
func (s *Scraper) SavePostsInDB() error {
	if len(s.Posts) == 0 {
		return nil
	}

	tx := s.Database.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, post := range s.Posts {
		if err := tx.Create(&post).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	s.UsageData.FinishedAt = time.Now()
	s.UsageData.Duration = time.Since(s.UsageData.StartedAt)
	s.UsageData.TotalPosts = len(s.Posts)
	s.UsageData.AverageTimePerPost = s.UsageData.Duration / time.Duration(s.UsageData.TotalPosts)

	return tx.Commit().Error
}

func (s *Scraper) SaveUsageData() error {
	tx := s.Database.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(&s.UsageData).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// FilterPosts filters the posts depending on the title length
func (s *Scraper) FilterPosts() error {
	moreThanFiveWords := s.filterMoreThanFiveWords(s.Posts)

	fiveWordsOrLess := s.filterFiveWordsOrLess(s.Posts)

	s.Posts = []Post{}

	s.Posts = append(s.Posts, moreThanFiveWords...)
	s.Posts = append(s.Posts, fiveWordsOrLess...)

	return nil
}

// isInDB checks if the post is already in the database
func (s *Scraper) isInDB(title string) bool {
	var post Post
	tx := s.Database.Where("title = ?", title).First(&post)

	return tx.Error != gorm.ErrRecordNotFound
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

func (s *Scraper) CheckPageStatus() int {
	statusCode := 0
	s.Collector.OnResponse(func(r *colly.Response) {
		statusCode = r.StatusCode
	})

	s.Collector.Visit(s.Domain)

	if statusCode != 200 {
		return statusCode
	}

	// close the collector
	s.Collector.Wait()

	return statusCode

}

// Method to count words
func (s *Scraper) countWords(input string) int {
	re := regexp.MustCompile(`[a-zA-Z0-9]+(?:[-'\.][a-zA-Z0-9]+)*`)
	matches := re.FindAllString(input, -1)
	return len(matches)
}

// Method to filter posts with more than five words in the title
func (s *Scraper) filterMoreThanFiveWords(posts []Post) []Post {
	filtered := []Post{}
	for _, post := range posts {
		if s.countWords(post.Title) > 5 {
			post.AppliedFilter = "more_than_five_words"
			filtered = append(filtered, post)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Comments > filtered[j].Comments
	})

	s.UsageData.TotalPostsWithMoreThanFiveWords = len(filtered)

	return filtered
}

// Method to filter posts with five words or less in the title
func (s *Scraper) filterFiveWordsOrLess(posts []Post) []Post {
	filtered := []Post{}
	for _, post := range posts {
		if s.countWords(post.Title) <= 5 {
			post.AppliedFilter = "five_words_or_less"
			filtered = append(filtered, post)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Points > filtered[j].Points
	})

	s.UsageData.TotalPostsWithFiveOrFewerWords = len(filtered)

	return filtered
}

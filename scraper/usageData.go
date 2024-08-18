package scraper

import (
	"time"

	"gorm.io/gorm"
)

type UsageData struct {
	gorm.Model
	StartedAt                       time.Time
	FinishedAt                      time.Time
	Duration                        time.Duration
	TotalPosts                      int
	TotalBytesScraped               int
	TotalPostsWithMoreThanFiveWords int
	TotalPostsWithFiveOrFewerWords  int
	AverageTimePerPost              time.Duration
}

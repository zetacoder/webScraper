package scraper

import (
	"testing"

	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
)

func TestCountWords(t *testing.T) {
	s := Scraper{}

	tests := []struct {
		input    string
		expected int
	}{
		{"This is - a self-explained example", 5},
		{"Another example", 2},
		{"", 0},
		{"This is a very long example with more than 10 words", 11},
		{"This is a very long example with more than 10 words and some special characters like !@#$%^&*()", 16},
		{"Example with a number 1234567890", 5},
		{`Discrete Mathematics – An Open Introduction, 4th edition`, 7},
		{"Interviewing the Interviewer: Questions to Uncover a Company's True Culture", 10},
		{"This is &a very $%long example with more than 10 words^*()  and some!@# special! characters like and -", 17},
		{"This is a example with more than 10 words and some special characters like don't", 15},
		{"The company's culture", 3},
		{"Charlie's Angels", 2},
	}

	for _, test := range tests {
		if got := s.countWords(test.input); got != test.expected {
			t.Errorf("countWords(%s) = %d; want %d", test.input, got, test.expected)
		}
	}
}

func TestFilterPosts(t *testing.T) {
	posts := []Post{
		{Title: "two words", Comments: 999, Points: 999},
		{Title: "three words title", Comments: 888, Points: 888},
		{Title: "discrete mathematics – an open introduction, 4th edition", Comments: 333, Points: 333},
		{Title: "this is a very long title with more than ten words in it", Comments: 666, Points: 666},
		{Title: "example with a number 1234567890", Comments: 444, Points: 444},
		{Title: "with more than five words in the title", Comments: 777, Points: 777},
		{Title: "this is a very long title with more than ten words in it and some special characters like !@#$%^&*()", Comments: 555, Points: 555},
	}

	expected := []Post{
		// Posts with more than five words, sorted by comments (descending)
		{Title: "with more than five words in the title", Comments: 777, Points: 777},
		{Title: "this is a very long title with more than ten words in it", Comments: 666, Points: 666}, // Highest Comments
		{Title: "this is a very long title with more than ten words in it and some special characters like !@#$%^&*()", Comments: 555, Points: 555},
		{Title: "discrete mathematics – an open introduction, 4th edition", Comments: 333, Points: 333},

		// Posts with five or fewer words, sorted by points (descending)
		{Title: "two words", Comments: 999, Points: 999}, // Highest Points
		{Title: "three words title", Comments: 888, Points: 888},
		{Title: "example with a number 1234567890", Comments: 444, Points: 444},
	}

	s := Scraper{
		Posts: posts,
	}

	err := s.FilterPosts()
	assert.NoError(t, err, "Unexpected error in FilterPosts")

	assert.Equal(t, len(expected), len(s.Posts), "Number of filtered posts is incorrect")

	for i, post := range expected {
		assert.Equal(t, post.Title, s.Posts[i].Title, "Post title mismatch at index %d", i)
		assert.Equal(t, post.Comments, s.Posts[i].Comments, "Post comments mismatch at index %d", i)
		assert.Equal(t, post.Points, s.Posts[i].Points, "Post points mismatch at index %d", i)
	}
}

func TestFilterMoreThanFiveWords(t *testing.T) {

	s := Scraper{}

	posts := []Post{
		{Title: "two words", Comments: 999},
		{Title: "three words title", Comments: 888},
		{Title: "discrete mathematics – an open introduction, 4th edition", Comments: 333},
		{Title: "this is a very long title with more than ten words in it", Comments: 666},
		{Title: "example with a number 1234567890", Comments: 444},
		{Title: "with more than five words in the title", Comments: 777},
		{Title: "this is a very long title with more than ten words in it and some special characters like !@#$%^&*()", Comments: 555},
	}

	expected := []Post{
		{Title: "with more than five words in the title", Comments: 777},
		{Title: "this is a very long title with more than ten words in it", Comments: 666},
		{Title: "this is a very long title with more than ten words in it and some special characters like !@#$%^&*()", Comments: 555},
		{Title: "discrete mathematics – an open introduction, 4th edition", Comments: 333},
	}

	result := s.filterMoreThanFiveWords(posts)

	assert.Equal(t, len(expected), len(result))

	for i, post := range expected {
		assert.Equal(t, post.Title, result[i].Title, "Post title mismatch")
		assert.Equal(t, post.Comments, result[i].Comments, "Post comments mismatch")
	}
}

func TestFilterFiveWordsOrLess(t *testing.T) {

	s := Scraper{}

	posts := []Post{
		{Title: "two words", Points: 999},
		{Title: "three words title", Points: 888},
		{Title: "discrete mathematics – an open introduction, 4th edition", Points: 333},
		{Title: "this is a very long title with more than ten words in it", Points: 666},
		{Title: "example with a number 1234567890", Points: 444},
		{Title: "with more than five words in the title", Points: 777},
		{Title: "this is a very long title with more than ten words in it and some special characters like !@#$%^&*()", Points: 555},
	}

	expected := []Post{
		{Title: "two words", Points: 999},
		{Title: "three words title", Points: 888},
		{Title: "example with a number 1234567890", Points: 444},
	}

	result := s.filterFiveWordsOrLess(posts)

	assert.Equal(t, len(expected), len(result))

	for i, post := range expected {
		assert.Equal(t, post.Title, result[i].Title, "Post title mismatch")
		assert.Equal(t, post.Points, result[i].Points, "Post points mismatch")
	}
}

func TestPageConnection(t *testing.T) {
	s := NewScraper(colly.NewCollector(), "https://news.ycombinator.com/", 0, nil)

	if s.CheckPageStatus() != 200 {
		t.Errorf("Expected status code 200, got %d", s.CheckPageStatus())
	}
}

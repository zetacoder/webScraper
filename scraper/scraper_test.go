package scraper

import "testing"

func TestCountWords(t *testing.T) {
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
		// example with special characters in between words
		{"This is &a very $%long example with more than 10 words^*()  and some!@# special! characters like and -", 17},
	}

	for _, test := range tests {
		if got := countWords(test.input); got != test.expected {
			t.Errorf("countWords(%s) = %d; want %d", test.input, got, test.expected)
		}
	}
}

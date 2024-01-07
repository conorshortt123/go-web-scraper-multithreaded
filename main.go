package main

import (
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	websites := []string{
		"https://en.wikipedia.org/wiki/Main_Page",
		"https://www.imdb.com/",
		"https://github.com/trending",
		"https://news.ycombinator.com/",
		"https://www.reddit.com/r/programming/",
		"https://openweathermap.org/",
		"https://www.amazon.com/Best-Sellers/zgbs",
		"https://www.goodreads.com/",
		"https://stackoverflow.com/questions",
		"https://www.cnn.com/",
	}

	// Loop over websites and scrape
	for _, website := range websites {
		scrape(website)
	}
}

// Scrapes for occurences of a word (case-insensitive)
func scrape(website string) []string {
	// Slice to store occurrences
	var occurrences []string

	// Create a new collector
	c := colly.NewCollector()

	// Set up a callback for when a visited HTML element is found
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Extract text content from the body
		bodyText := e.Text

		// Check for occurrences of the word "ChatGPT" (case-insensitive)
		if strings.Contains(strings.ToLower(bodyText), "chatgpt") {
			occurrences = append(occurrences, e.Request.URL.String())
		}
	})

	// Set up error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Visit the initial URL
	err := c.Visit(website)
	if err != nil {
		log.Fatal(err)
	}

	return occurrences
}

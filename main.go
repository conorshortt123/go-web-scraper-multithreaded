package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	// List of websites to search
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

	// Prompt the user for input
	fmt.Print("Enter the keyword you want to search websites for: ")

	// Declare a variable to store user input
	var keyword string

	// Read user input from the console
	_, err := fmt.Scan(&keyword)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Loop over websites and scrape
	for _, website := range websites {
		occurrences := scrape(website, keyword)

		for _, occurence := range occurrences {
			log.Printf("Found occurence of %s! \n%s", keyword, occurence)
		}
	}
}

// Scrapes for occurences of a word
func scrape(website, searchWord string) []string {
	// Slice to store occurrences
	var occurrences []string

	// Create a new collector
	c := colly.NewCollector()

	// Set up a callback for when a visited HTML element is found
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Extract text content from the body
		captured := captureText(e.Text, searchWord)

		if !isEmpty(captured) {
			occurrences = append(occurrences, "Website: "+e.Request.URL.String()+"\nCaptured text : "+captured)
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

func captureText(body, searchWord string) string {
	beforeCount := 50
	afterCount := 50
	surroundingText := ""

	// Find the target word in the body text (case-insensitive)
	index := strings.Index(strings.ToLower(body), strings.ToLower(searchWord))
	if index != -1 {
		// Calculate the start and end indices for the captured substring
		startIndex := index - beforeCount
		if startIndex < 0 {
			startIndex = 0
		}

		endIndex := index + len(searchWord) + afterCount
		if endIndex > len(body) {
			endIndex = len(body)
		}

		// Capture the surrounding substring
		surroundingText = body[startIndex:endIndex]
	}

	// Remove excess whitespace within the captured text using a regular expression
	re := regexp.MustCompile(`\s+`)
	cleanedText := re.ReplaceAllString(surroundingText, " ")

	return cleanedText
}

func isEmpty(s string) bool {
	return len(s) == 0
}

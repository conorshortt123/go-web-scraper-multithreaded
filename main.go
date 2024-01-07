package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

type ScrapedResult struct {
	WebsiteURL   string
	CapturedText string
	ErrorMessage string
}

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

	// Channel to receive results of scrape from goroutines i.e results of multi threaded execution
	results := make(chan ScrapedResult, len(websites))

	// WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup

	// Increment the counter for each goroutine
	wg.Add(len(websites))

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
	for _, url := range websites {
		go scrape(url, keyword, &wg, results)
	}

	// Close the channel once all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results from the channel
	for result := range results {
		if result.ErrorMessage != "" {
			fmt.Printf("Error scraping %s: %s\n", result.WebsiteURL, result.ErrorMessage)
		} else {
			fmt.Printf("Website: %s : Captured keyword %s from text: %s\n", result.WebsiteURL, keyword, result.CapturedText)
		}
	}
}

// Scrapes for occurences of a word
func scrape(url string, searchWord string, wg *sync.WaitGroup, ch chan<- ScrapedResult) {
	defer wg.Done() // Decrement the counter when the function completes

	// Create a new collector
	c := colly.NewCollector()

	// Initialize the result structure
	result := ScrapedResult{
		WebsiteURL:   url,
		CapturedText: "",
	}

	// Set up a callback for when a visited HTML element is found
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Extract text content from the body
		captured := captureText(e.Text, searchWord)

		if !isEmpty(captured) {
			// Assign the captured text to the result structure
			result.CapturedText = captured
		}
	})

	// Set up error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)

		// Assign the error message to the result structure
		result.ErrorMessage = err.Error()
	})

	// Start scraping from the provided URL
	err := c.Visit(url)
	if err != nil {
		log.Printf("Error scraping %s: %v\n", url, err)
	}

	// Send the result to the channel if text was captured
	if !isEmpty(result.CapturedText) {
		ch <- result
	}
}

func captureText(body string, searchWord string) string {
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

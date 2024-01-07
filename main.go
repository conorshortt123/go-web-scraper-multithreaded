package main

import (
	"fmt"
	"log"

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

func scrape(website string) {
	// Create a new collector
	c := colly.NewCollector()

	// Set up a callback for when a visited HTML element is found
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(link)
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
}

package main

import (
	"fmt"
	"time"

	"github.com/asciimoo/colly"
)

func main() {

	url := "https://en.wikipedia.org"

	// Instantiate default collector
	c := colly.NewCollector()

	// Visit only domains
	c.AllowedDomains = []string{"wikipedia.org", "en.wikipedia.org"}

	// Limit the maximum parallelism to 10
	// when visiting links which domains' matches "*wikipedia.*" glob
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*wikipedia.*", Parallelism: 10})

	// MaxDepth is 2, so only the links on the scraped page
	// and links on those pages are visited
	c.MaxDepth = 2

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Println(link)
		// Visit link found on page on a new thread
		go e.Request.Visit(link)
	})

	// Before making a request print "Starting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting", r.URL, time.Now())
	})

	// After making a request print "Finished ..."
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL, time.Now())
	})

	// Start scraping on https://en.wikipedia.org
	c.Visit(url)
	// Wait until threads are finished
	c.Wait()
}

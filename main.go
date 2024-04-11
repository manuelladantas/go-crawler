package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	crawl()
}

func crawl() {
	var arr []string
	count := 0

	c := colly.NewCollector(
		colly.AllowedDomains("rottentomatoes.com", "www.rottentomatoes.com"),
	)
	infoCollector := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	c.OnHTML("span.p--small", func(h *colly.HTMLElement) {
		name := h.Text

		arr = append(arr, name)

		fmt.Printf("Name: ", name)
	})

	c.OnHTML("div.discovery__actions", func(e *colly.HTMLElement) {
		test := e.ChildAttr("button", "class")
		fmt.Println(test)
		count++
		url := fmt.Sprintf("https://www.rottentomatoes.com/browse/movies_in_theaters/?page=%v ", count)
		fmt.Println(count)
		c.Visit(url)
	})

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Profile URL: ", r.URL.String())
	})

	c.Visit("https://www.rottentomatoes.com/browse/movies_in_theaters/")
}

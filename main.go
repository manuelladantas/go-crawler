package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/manuelladantas/go-crawler/database"
)

type Movie struct {
	Name        string
	Synopse     string
	Rating      string
	Genre       string
	Director    string
	Producer    string
	Writer      string
	ReleaseDate string
	Runtime     string
	Distributor string
}

const CARD_SELECTOR string = "#main-page-content > div.discovery > div.discovery-grids-container > div > div.discovery-tiles__wrap > div > div > tile-dynamic > a[href]"

func main() {
	database.Ping()
	database.Disconnect()
	// crawl()
}

func crawl() {
	movies := make([]Movie, 0)
	// count := 0

	c := colly.NewCollector(
		colly.AllowedDomains("rottentomatoes.com", "www.rottentomatoes.com"),
		// colly.MaxDepth(3),
		colly.MaxBodySize(0),
		// colly.Async(true),
	)
	// TODO: Analisar a necessidade disso
	// c.Limit(&colly.LimitRule{
	// 	DomainGlob:  "*rottentomatoes.*",
	// 	Parallelism: 2,
	// 	Delay:       5 * time.Second,
	// })

	// c.SetRequestTimeout(120 * time.Second)

	infoCollector := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	c.OnHTML(CARD_SELECTOR, func(he *colly.HTMLElement) {
		link := he.Attr("href")
		link = he.Request.AbsoluteURL(link)
		fmt.Println("Movie Link", link)

		infoCollector.Visit(link)
	})

	// c.OnHTML("div.discovery__actions", func(e *colly.HTMLElement) {
	// 	test := e.ChildAttr("button", "class")
	// 	fmt.Println(test)
	// 	count++
	// 	url := fmt.Sprintf("%sbrowse/movies_in_theaters/?page=%v ", URL, count)
	// 	fmt.Println(count)
	// 	c.Visit(url)
	// })

	// c.OnHTML("")

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Movie URL: ", r.URL.String())
	})

	infoCollector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	infoCollector.OnHTML("#main", func(e *colly.HTMLElement) {
		movie := Movie{}

		movie.Name = e.ChildText("#scoreboard > h1")
		movie.Synopse = e.ChildText("#movie-info > div > div > drawer-more > p")

		infoMap := map[string]*string{
			"Distributor:":              &movie.Distributor,
			"Rating:":                   &movie.Rating,
			"Director:":                 &movie.Director,
			"Producer:":                 &movie.Producer,
			"Writer:":                   &movie.Writer,
			"Release Date (Streaming):": &movie.ReleaseDate,
			"Runtime:":                  &movie.Runtime,
		}

		e.ForEach("#info > li", func(_ int, h *colly.HTMLElement) {
			label := h.ChildText("p > b")
			value := h.ChildText("p > span")

			if fieldPointer, ok := infoMap[label]; ok {
				*fieldPointer = strings.ReplaceAll(value, "\n", "")
			}
		})

		movies = append(movies, movie)
	})

	c.OnResponse(func(r *colly.Response) {
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		js, err := json.MarshalIndent(movies, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("movies.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})
	c.Visit("https://www.rottentomatoes.com/browse/movies_at_home/")
	// c.Wait()
	// infoCollector.Wait()
}

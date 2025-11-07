package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Book struct {
	title string
	url   string
}

func main() {

	currentlyReading := make([]Book, 0, 200)

	c := colly.NewCollector(
		colly.AllowedDomains("www.goodreads.com"))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Accessing goodreads profile")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("❌ Request error:", err)
		if r != nil {
			fmt.Println("Status Code:", r.StatusCode)
			fmt.Println("URL:", r.Request.URL)
			fmt.Println("Body (if any):", string(r.Body))
		}
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Received response!", r.StatusCode)
	})

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		if h.Attr("class") != "bookTitle" {
			return
		}

		link := h.Attr("href")
		title := h.Text

		currentlyReading = append(currentlyReading, Book{title, link})
	})

	c.Visit("https://www.goodreads.com/user/show/38718384-sophie-marshall-unitt")

	// reviewCollector := c.Clone()

	// reviews := make([]Review, 0, 200)

	// e.Request.Visit(link)

	// reviewCollector.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("runninggg")
	// })

	// reviewCollector.OnResponse(func(r *colly.Response) {
	// 	fmt.Println("Visited", r.Request.URL)
	// })
	// reviewCollector.Visit("https://thestorygraph.com")

	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", "  ")

	// // Dump json to the standard output
	// enc.Encode(courses)
}

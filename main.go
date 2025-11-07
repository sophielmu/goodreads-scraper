package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gocolly/colly"
)

type Book struct {
	title string
	url   string
}

type Shelf struct {
	name string
	url  string
}

func main() {
	defaultShelfNames := []string{"to-read", "currently-reading", "read", "favorites"}
	currentlyReading := make([]Book, 0, 200)
	shelves := make([]Shelf, 0, 200)

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

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		if h.Attr("class") != "actionLinkLite userShowPageShelfListItem" {
			return
		}

		splitShelfName := strings.Split(h.Text, "\u200E")
		cleaned := strings.TrimLeft(splitShelfName[0], "\r\n ")
		shelf := Shelf{cleaned, h.Attr("href")}
		shelves = append(shelves, shelf)

		fmt.Println(shelf.name)
		if !slices.Contains(defaultShelfNames, shelf.name) {
			fmt.Println("Returning, can't find shelf")
			return
		}

		switch foundShelfName := shelf.name; foundShelfName {
		case defaultShelfNames[1]:
			fmt.Println("currently-reading")
		case defaultShelfNames[2]:
			fmt.Println("read")
		}

	})

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		if h.Attr("class") != "userPagePhoto" {
			return
		}

		userPhoto := h.Attr("href")
		fmt.Println("User photo link - https://www.goodreads.com" + userPhoto)
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

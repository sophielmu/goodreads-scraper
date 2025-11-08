package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type UserProfile struct {
	photoUrl string
	shelves  []Shelf
}

type Book struct {
	coverUrl  string
	title     string
	author    string
	isbn      string
	numPages  string
	rating    int
	dateAdded string
	dateRead  string
}

type Shelf struct {
	name string
	url  string
}

func getBooksFromShelf(url string) []Book {
	result := make([]Book, 0, 10)
	collector := colly.NewCollector()

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Grabbing data from " + url)
	})

	collector.OnHTML("tr.bookalike.review", func(h *colly.HTMLElement) {
		book := Book{
			title:    strings.TrimSpace(h.ChildText("td.field.title .value a")),
			author:   strings.TrimSpace(h.ChildText("td.field.author .value a")),
			coverUrl: h.ChildAttr("td.field.cover img", "src"),
		}
		result = append(result, book)
	})

	collector.Visit(url)
	return result
}

func getUserProfile(url string, c *colly.Collector) UserProfile {
	var profile UserProfile

	c.OnHTML("a.userPagePhoto[href]", func(h *colly.HTMLElement) {

		userPhotoLink := "https://www.goodreads.com" + h.Attr("href")
		fmt.Println("User photo link - ", userPhotoLink)
	})

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		if h.Attr("class") != "actionLinkLite userShowPageShelfListItem" {
			return
		}

		splitShelfName := strings.Split(h.Text, "\u200E")
		cleaned := strings.TrimLeft(splitShelfName[0], "\r\n ")
		shelf := Shelf{cleaned, "https://www.goodreads.com" + h.Attr("href")}
		profile.shelves = append(profile.shelves, shelf)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished scraping:", url)
	})

	c.Visit(url)

	return profile
}

func main() {

	shelvesToRead := []string{"currently-reading", "read"}
	profileUrl := "https://www.goodreads.com/user/show/38718384-sophie-marshall-unitt"

	c := colly.NewCollector(
		colly.AllowedDomains("www.goodreads.com"))

	user := getUserProfile(profileUrl, c)

	// make a map collection to find the item we want
	shelfMap := make(map[string]Shelf)
	for _, myShelf := range user.shelves {
		shelfMap[myShelf.name] = myShelf
	}

	if found, ok := shelfMap[shelvesToRead[0]]; ok {
		currentlyReading := getBooksFromShelf(found.url)
		fmt.Println(currentlyReading)
	}

	if found, ok := shelfMap[shelvesToRead[1]]; ok {
		recentlyRead := getBooksFromShelf(found.url)
		fmt.Println(recentlyRead)
	}
}

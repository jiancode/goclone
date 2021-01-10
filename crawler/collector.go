package crawler

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func afterStr(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	vl := len(value)
	if adjustedPos >= vl {
		return ""
	}
	return value[adjustedPos:vl]
}

// Collector searches for css, js, and images within a given link
// TODO improve for better performance
func Collector(urlStr string, projectPath string) {
	// create a new collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
		// asynchronous boolean
		colly.Async(false),
	)

	// search for all link tags that have a rel attribute that is equal to stylesheet - CSS
	c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		// hyperlink reference
		link := e.Attr("href")
		// print css file was found
		fmt.Println("Css found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// search for all script tags with src attribute -- JS
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		// Print link
		fmt.Println("Js found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// search for all img tags with src attribute -- Images
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		// Print link
		fmt.Println("Img found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// recursive internal link
	c.OnHTML("meta[http-equiv]", func(e *colly.HTMLElement) {

		var link string

		// http-equiv=refresh attribute
		if e.Attr("http-equiv") == "refresh" {
			c := e.Attr("content")
			link = afterStr(c, "url=")
		} else {
			return
		}

		if !strings.HasPrefix(link, "javascript") {
			// Print link
			fmt.Printf("\n>>>>>> reflush: %s\n", link)
			sublink := e.Request.AbsoluteURL(link)
			fmt.Printf("Visiting %s\n", sublink)
			Collector(sublink, projectPath)
		}
	})

	// recursive internal link
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("href")
		if !strings.HasPrefix(link, "http") && !strings.HasPrefix(link, "javascript") && !strings.HasPrefix(link, "#") {
			// Print link
			fmt.Printf("\n>>>>>> Sublink: %s\n", link)
			sublink := e.Request.AbsoluteURL(link)
			fmt.Printf("Visiting %s\n", sublink)
			Collector(sublink, projectPath)

		}
	})

	//Before making a request
	c.OnRequest(func(r *colly.Request) {
		link := r.URL.String()
		if urlStr == link {
			HTMLExtractor(link, projectPath)
		}
	})

	// Visit each url and wait for stuff to load :)
	c.Visit(urlStr)
	c.Wait()
}

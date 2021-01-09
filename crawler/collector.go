package crawler

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

// Collector searches for css, js, and images within a given link
// TODO improve for better performance
func Collector(urlStr string, projectPath string) {
	// create a new collector
	c := colly.NewCollector(
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
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// src attribute
		link := strings.TrimSpace(e.Attr("href"))
		if strings.HasPrefix(link, "/") || strings.HasPrefix(link, ".") {
			// Print link
			fmt.Printf("\n\n>>>========================<<<\nSublink href: %s\n", link)
			sublink := e.Request.AbsoluteURL(link)
			fmt.Printf("Downloading %s\n", sublink)

			// extraction
			u, err := url.Parse(sublink)
			dirPath, base := filepath.Split(u.Path)
			if base == "" {
				base = "index.html"
			}

			if filepath.Ext(base) == "" {
				dirPath = u.Path
				base = "index.html"
			}
			fileDir := filepath.Join(projectPath, dirPath)
			fileName := filepath.Join(fileDir, base)
			// Check if page has downloaded
			_, err = os.Stat(fileName)
			if err != nil {
				Collector(sublink, projectPath)
			}

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

package crawler

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
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

// Link2FileName check if page is downloaded
// Create file directory if its necessary
func Link2FileName(link, projectPath string) (fileName string, newPage bool) {

	var fileDir, base, dirPath string
	newPage = true

	u, err := url.Parse(link)
	p := strings.TrimSpace(u.Path)
	if p == "" || p == "/" {
		fileName = filepath.Join(projectPath, "/", "index.html")
	} else {
		dirPath, base = filepath.Split(p)
		if base == "" {
			base = "index.html"
		}
		fileExt := filepath.Ext(base)
		if fileExt == "" || fileExt == ".php" {
			if u.RawQuery != "" {
				dirPath = p + "?" + u.RawQuery
			} else {
				dirPath = p
			}
			base = "index.html"
		}

		fileDir = filepath.Join(projectPath, dirPath)
		fileName = filepath.Join(fileDir, base)
		// Check if page has downloaded
		_, err = os.Stat(fileName)
		if err == nil {
			newPage = false
			return fileName, newPage
		}
		err = os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			fmt.Println("Mkdir error:", err)
		}
	}

	// Check if page has downloaded
	_, err = os.Stat(fileName)
	if err == nil {
		newPage = false
	}

	return fileName, newPage
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
			_, newPage := Link2FileName(link, projectPath)
			if newPage {
				fmt.Printf("\n>>>>>> reflush: %s\n", link)
				sublink := e.Request.AbsoluteURL(link)
				fmt.Printf("Visiting %s\n", sublink)
				Collector(sublink, projectPath)
			}
		}
	})

	// recursive internal link
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("href")
		if !strings.HasPrefix(link, "http") && !strings.HasPrefix(link, "javascript") && !strings.HasPrefix(link, "#") {
			_, newPage := Link2FileName(link, projectPath)
			if newPage {
				fmt.Printf("\n>>>>>> Sublink: %s\n", link)
				sublink := e.Request.AbsoluteURL(link)
				fmt.Printf("Visiting %s\n", sublink)
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

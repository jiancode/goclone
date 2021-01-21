package crawler

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func parseHTML(a, s, p string) {
	d, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		fmt.Println(err)
	}
	d.Find("link").Each(func(index int, e *goquery.Selection) {
		h, found := e.Attr("href")
		if found {
			fmt.Println("Check link in js:", path.Join(a, h))
			Extractor(path.Join(a, h), p)
		}
	})

	d.Find("img").Each(func(index int, e *goquery.Selection) {
		h, found := e.Attr("src")
		fmt.Println("video src", h)
		if found {
			Extractor(path.Join(a, h), p)
		}
		h, found = e.Attr("pic-src")
		if found {
			Extractor(path.Join(a, h), p)
		}
		h, found = e.Attr("webppic-src")
		if found {
			Extractor(path.Join(a, h), p)
		}
	})

	d.Find("video").Each(func(index int, e *goquery.Selection) {
		h, found := e.Attr("src")
		fmt.Println("video src", h)
		if found {
			Extractor(path.Join(a, h), p)
		}
		h, found = e.Attr("image")
		if found {
			Extractor(path.Join(a, h), p)
		}
		h, found = e.Attr("file")
		if found {
			Extractor(path.Join(a, h), p)
		}
		h, found = e.Attr("poster")
		if found {
			Extractor(path.Join(a, h), p)
		}
	})

	d.Find("embed").Each(func(index int, e *goquery.Selection) {
		h, found := e.Attr("src")
		if found {
			Extractor(path.Join(a, h), p)
		}
		h, found = e.Attr("image")
		if found {
			Extractor(path.Join(a, h), p)
		}
		h, found = e.Attr("file")
		if found {
			Extractor(path.Join(a, h), p)
		}
		h, found = e.Attr("poster")
		if found {
			Extractor(path.Join(a, h), p)
		}
	})

}

func jsLink(absURL, pStr, pPath string) {
	gMatch := regexp.MustCompile(`document.writeln\(\'(.*)\'\)`)
	hstr := gMatch.FindAllStringSubmatch(pStr, -1)
	for _, s := range hstr {
		//fmt.Println("Find js link:", s[1])
		parseHTML(absURL, s[1], pPath)
	}
}

// Collector searches for css, js, and images within a given link
// TODO improve for better performance
func Collector(urlStr string, projectPath string) {
	// create a new collector
	u, _ := url.Parse(urlStr)
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
		// asynchronous boolean
		colly.Async(false),
		colly.AllowedDomains(u.Hostname()),
	)

	// search for all link tags that have a rel attribute that is equal to stylesheet - CSS
	c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		// hyperlink reference
		link := e.Attr("href")
		// print css file was found
		//fmt.Println("Css found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// search for all script tags with src attribute -- JS
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// search for all link, image and video in javascript
	c.OnHTML("script", func(e *colly.HTMLElement) {
		pageStr := e.DOM.Contents().Text()
		absURL := fmt.Sprintf("%s://%s", e.Request.URL.Scheme, e.Request.URL.Host)
		//fmt.Println("script text:", pageStr)
		jsLink(absURL, pageStr, projectPath)
	})

	// search for all img tags with src attribute -- Images
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		// Print link
		//fmt.Println("Img found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// search for all img tags with pic-src attribute -- Images
	c.OnHTML("img[pic-src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("pic-src")
		// Print link
		//fmt.Println("Img found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// search for all img tags with src attribute -- Images
	c.OnHTML("img[webppic-src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("webppic-src")
		// Print link
		//fmt.Println("Img found", "-->", link)
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
			sublink := e.Request.AbsoluteURL(link)
			// Print link
			_, newPage := Link2FileName(sublink, projectPath)
			if newPage {
				fmt.Printf("\n>>>>>> reflush: %s\n", sublink)
				Collector(sublink, projectPath)
			}
		}
	})

	// recursive internal link
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("href")
		if !strings.HasPrefix(link, "http") && !strings.HasPrefix(link, "javascript") && !strings.HasPrefix(link, "#") {
			sublink := e.Request.AbsoluteURL(link)
			_, newPage := Link2FileName(sublink, projectPath)
			if newPage {
				fmt.Printf("\n>>>>>> Sublink: %s\n", sublink)
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

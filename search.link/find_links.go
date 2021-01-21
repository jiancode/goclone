// find_links_in_page.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// This will get called for each HTML element found
func modHref(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	href, exists := element.Attr("href")
	if exists {
		h := strings.TrimSpace(href)
		u, _ := url.Parse(h)
		ext := filepath.Ext(u.Path)
		if ext == "" || ext == ".php" {
			newhref := h + ".html"
			element.SetAttr("href", newhref)
			fmt.Printf(">>> %s \n", newhref)
		}
	}
}

func parseHTML(s string) {
	d, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	d.Find("link").Each(func(index int, e *goquery.Selection) {
		//log.Println(e.Contents().Text())
		h, _ := e.Attr("href")
		fmt.Println("link", h)
	})
	d.Find("video").Each(func(index int, e *goquery.Selection) {
		h, found := e.Attr("src")
		if found {
			fmt.Println("video src", h)
		}
		h, found = e.Attr("image")
		if found {
			fmt.Println("video image", h)
		}
		h, found = e.Attr("file")
		if found {
			fmt.Println("video file", h)
		}
		h, found = e.Attr("poster")
		if found {
			fmt.Println("video poster", h)
		}

	})
	d.Find("embed").Each(func(index int, e *goquery.Selection) {
		h, found := e.Attr("src")
		if found {
			fmt.Println("embed", h)
		}
		h, found = e.Attr("image")
		if found {
			fmt.Println("embed image", h)
		}
		h, found = e.Attr("file")
		if found {
			fmt.Println("embed file", h)
		}
		h, found = e.Attr("poster")
		if found {
			fmt.Println("embed poster", h)
		}

	})

}

var filename string

func main() {
	if len(os.Args) != 2 {
		fmt.Println("search.link path")
		os.Exit(0)
	}
	var document *goquery.Document
	err := filepath.Walk(os.Args[1],
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return nil
				}
				filename = path
				// Create a goquery document from the HTTP response
				document, err = goquery.NewDocumentFromReader(bufio.NewReader(file))
				if err != nil {
					log.Fatal("Error loading HTTP response body. ", err)
				}

				// Find all links and process them with the function
				// defined earlier
				//document.Find("a").Each(modHref)
				document.Find("meta").Each(func(index int, e *goquery.Selection) {
					c, _ := e.Attr("content")
					nc := strings.ReplaceAll(c, ".php", ".php.html")
					e.SetAttr("content", nc)
				})
				document.Find("script").Each(func(index int, e *goquery.Selection) {
					//log.Println(e.Contents().Text())
					gstr := e.Contents().Text()
					gMatch := regexp.MustCompile(`document.writeln\(\'(.*)\'\)`)
					hstr := gMatch.FindAllStringSubmatch(gstr, -1)
					for _, s := range hstr {
						//fmt.Println(s)
						parseHTML(s[1])
					}
					//log.Println(src)
				})
			}
			return nil
		})

	//fmt.Println(document.Html())
	if err != nil {
		log.Println(err)
	}
}

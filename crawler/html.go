package crawler

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
		} else {
			// Add .html extention to php pages
			fileExt := filepath.Ext(base)
			if fileExt == "" || fileExt == ".php" {
				if u.RawQuery != "" {
					base = base + "_" + pageFileName(u.RawQuery)
				}
				base = base + ".html"
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
	}

	// Check if page has downloaded
	_, err = os.Stat(fileName)
	if err == nil {
		newPage = false
	}

	return fileName, newPage
}

// This will get called for each HTML element found
func modHref(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	href, exists := element.Attr("href")
	if exists {
		if strings.HasPrefix(href, "#") {
			return
		}
		h := strings.TrimSpace(href)
		u, _ := url.Parse(h)
		ext := filepath.Ext(u.Path)
		if ext == "" || ext == ".php" {
			newhref := pageFileName(h) + ".html"
			element.SetAttr("href", newhref)
		}
	}
	return
}

func pageFileName(f string) string {
	f = strings.ReplaceAll(f, "&", "_")
	f = strings.ReplaceAll(f, "?", "_")
	f = strings.ReplaceAll(f, "=", "_")
	return f
}

// HTMLExtractor ...
func HTMLExtractor(link string, projectPath string) {

	fileName, newPage := Link2FileName(link, projectPath)
	if !newPage {
		return
	}
	// get the project name and path we use the path to
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	fmt.Printf("Extracting HTML %s --> %s\n", link, fileName)
	// get the html body
	resp, err := HTTPGet(link)
	// Close the body once everything else is compled
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		f.WriteString("<html><body><p>Download page error!</p></body></html>")
		return
	}
	// Modify internel link
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error loading HTTP response body. ", err)
		f.WriteString("<html><body><p>Download page error!</p></body></html>")
	} else {
		// Find all links and process
		doc.Find("a").Each(modHref)
		// Modify refresh page
		doc.Find("meta").Each(func(index int, e *goquery.Selection) {
			he, found := e.Attr("http-equiv")
			if found && he == "refresh" {
				c, _ := e.Attr("content")
				nc := strings.ReplaceAll(c, ".php", ".php.html")
				e.SetAttr("content", nc)
			}
		})
		htmlData, err := doc.Html()
		if err == nil {
			f.WriteString(htmlData)
		} else {
			fmt.Println(err)
		}
	}
	return
}

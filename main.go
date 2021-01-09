package main

import (
	"fmt"

	"github.com/jianconnect/goclone/crawler"
	"github.com/jianconnect/goclone/file"
	"github.com/jianconnect/goclone/parser"
)

func main() {
	url := args[0]

	if Serve == true {
		// grab the url from the
		if !parser.ValidateURL(url) && !parser.ValidateDomain(url) {
			fmt.Println("goclone <url>")
		} else {
			name := url
			// CreateProject
			projectpath := file.CreateProject(name)
			// create the url
			validURL := parser.CreateURL(name)
			// Crawler
			crawler.Crawl(validURL, projectpath)
		}
	}
}

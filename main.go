package main

import (
	"fmt"
	"os"

	"github.com/jianconnect/goclone/crawler"
	"github.com/jianconnect/goclone/file"
	"github.com/jianconnect/goclone/parser"
)

func main() {
	url := os.Args[1]

	if !parser.ValidateURL(url) && !parser.ValidateDomain(url) {
		fmt.Println("goclone <url>")
	} else {
		name := url
		// CreateProject
		projectpath := file.CreateProject(name)
		// create the url
		//validURL := parser.CreateURL(name)
		// Crawler
		crawler.Crawl(url, projectpath)
	}

}

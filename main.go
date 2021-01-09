package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jianconnect/goclone/crawler"
	"github.com/jianconnect/goclone/file"
)

func main() {
	urlStr := os.Args[1]

	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Printf("goclone <url> \n Error url:%s", urlStr)
	} else {
		name := u.Hostname()
		// CreateProject
		projectpath := file.CreateProject(name)
		// create the url
		//validURL := parser.CreateURL(name)
		// Crawler
		crawler.Crawl(urlStr, projectpath)
	}

}

package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/jianconnect/goclone/crawler"
)

func createHomeFolder(projectName string) (projectPath string, err error) {
	curPath, _ := os.Getwd()
	// define project path
	projectPath = filepath.Join(curPath, "/", projectName)

	// create base directory
	err = os.MkdirAll(projectPath, 0755)
	return projectPath, err
}

func main() {
	urlStr := os.Args[1]

	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Printf("goclone <url> \n Error url:%s", urlStr)
	} else {
		name := u.Hostname()
		// CreateProject
		projectpath, err := createHomeFolder(name)
		if err != nil {
			fmt.Printf("Setup project home folder err: %v", err)
		} else {
			crawler.Crawl(urlStr, projectpath)
		}
	}

}

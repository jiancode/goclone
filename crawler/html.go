package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// HTMLExtractor ...
func HTMLExtractor(link string, projectPath string) {
	fmt.Println("Extracting --> ", link)

	// get the html body
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}

	// Close the body once everything else is compled
	defer resp.Body.Close()

	u, err := url.Parse(link)
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
	os.MkdirAll(fileDir, os.ModePerm)

	fmt.Printf("path:%s, name:%s", fileDir, fileName)
	// get the project name and path we use the path to
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	htmlData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	f.Write(htmlData)

}

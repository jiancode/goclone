package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// HTMLExtractor ...
func HTMLExtractor(link string, projectPath string) {
	var fileName, fileDir string
	u, err := url.Parse(link)
	p := strings.TrimSpace(u.Path)
	if p == "" || p == "/" {
		fileName = filepath.Join(projectPath, "/", "index.html")
	} else {
		dirPath, base := filepath.Split(p)
		if base == "" {
			base = "index.html"
		}
		if filepath.Ext(base) == "" {
			dirPath = u.Path
			base = "index.html"
		}

		fileDir = filepath.Join(projectPath, dirPath)
		fileName = filepath.Join(fileDir, base)
		// Check if page has downloaded
		_, err = os.Stat(fileName)
		if err == nil {
			return
		}
		os.MkdirAll(fileDir, os.ModePerm)
	}

	fmt.Printf("Extracting %s --> %s", link, fileName)
	// get the html body
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	// Close the body once everything else is compled
	defer resp.Body.Close()

	//fmt.Printf("path:%s, name:%s", fileDir, fileName)
	// get the project name and path we use the path to
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
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

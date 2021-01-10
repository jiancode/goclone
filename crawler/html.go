package crawler

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// HTMLExtractor ...
func HTMLExtractor(link string, projectPath string) {
	var fileName, fileDir, base, dirPath string
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
			return
		}
		err = os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			fmt.Println("Mkdir error:", err)
		}
	}

	// Check if page has downloaded
	_, err = os.Stat(fileName)
	if err == nil {
		return
	}

	fmt.Printf("Extracting HTML %s --> %s\n", link, fileName)
	// get the html body
	resp, err := HTTPGet(link)
	if err != nil {
		fmt.Println(err)
	}
	// Close the body once everything else is compled
	defer resp.Body.Close()

	//fmt.Printf("path:%s, name:%s", fileDir, fileName)
	// get the project name and path we use the path to
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		f.Write(htmlData)
	} else {
		fmt.Println(err)
	}

}

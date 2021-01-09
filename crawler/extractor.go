package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// Extractor visits a link determines if its a page or sublink
// downloads the contents to a correct directory in project folder
// TODO add functionality for determining if page or sublink
func Extractor(link string, projectPath string) {
	fmt.Println("Extracting new --> ", link)

	// get the html body
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println("Ignore bad link --> ", link)
		//panic(err)
		return
	}

	// Closure
	defer resp.Body.Close()

	u, err := url.Parse(link)
	dirPath, base := filepath.Split(u.Path)
	//fmt.Printf("Download file: %s %s", dirPath, base)
	writeFileToPath(projectPath, base, dirPath, resp)
}

func writeFileToPath(projectPath, base, dirPath string, resp *http.Response) {
	//var name = base[0 : len(base)-len(oldFileExt)]
	//document := name + newFileExt
	fileDir := filepath.Join(projectPath, dirPath)

	os.MkdirAll(fileDir, os.ModePerm)
	fileName := filepath.Join(fileDir, base)

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

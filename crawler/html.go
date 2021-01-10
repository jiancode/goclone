package crawler

import (
	"fmt"
	"io/ioutil"
	"os"
)

// HTMLExtractor ...
func HTMLExtractor(link string, projectPath string) {

	fileName, newPage := Link2FileName(link, projectPath)
	if !newPage {
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

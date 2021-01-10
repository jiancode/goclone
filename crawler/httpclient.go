package crawler

import (
	"fmt"
	"net/http"
	"time"
)

// HTTPGet --- http GET client with custom properties
func HTTPGet(link string) (resp *http.Response, err error) {

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 Firefox")
		resp, err = client.Do(req)
	}
	return resp, err

}

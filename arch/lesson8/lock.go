package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type WebSite struct {
	Url  string
	Size int64
}

var websites []*WebSite

func getsize(url string) {
	resp, err := http.Head(url)
	if err != nil {
		log.Print(err)
		return
	}

	defer resp.Body.Close()
	fmt.Println(resp.Status, resp.ContentLength)
	site := &WebSite{
		Url:  url,
		Size: resp.ContentLength,
	}
	websites = append(websites, site)
}

func main() {
	urls := []string{"http://www.baidu.com", "http://www.51reboot.com", "http://godoc.org"}
	for _, url := range urls {
		go getsize(url)
	}

	time.Sleep(5 * time.Second)
	for _, site := range websites {
		fmt.Println(site.Url, site.Size)
	}
}

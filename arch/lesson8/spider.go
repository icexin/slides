package main

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type TaskPool struct {
	queue   chan func()
	wait    *sync.WaitGroup
	workers int
}

func NewTaskPool(workers int) *TaskPool {
	return &TaskPool{
		queue:   make(chan func()),
		wait:    new(sync.WaitGroup),
		workers: workers,
	}
}

func (t *TaskPool) work() {
	defer t.wait.Done()
	for task := range t.queue {
		task()
	}
}

func (t *TaskPool) Start() {
	for i := 0; i < t.workers; i++ {
		t.wait.Add(1)
		go t.work()
	}
}

func (t *TaskPool) Submit(task func()) {
	t.queue <- task
}

func (t *TaskPool) Stop() {
	close(t.queue)
}

func (t *TaskPool) Wait() {
	t.wait.Wait()
}

func makelink(uri *url.URL, link string) string {
	switch {
	case strings.HasPrefix(link, "https") ||
		strings.HasPrefix(link, "http"):
		return link
	case strings.HasPrefix(link, "//"):
		return uri.Scheme + ":" + link
	case strings.HasPrefix(link, "/"):
		return fmt.Sprintf("%s://%s%s", uri.Scheme, uri.Host, link)
	default:
		return fmt.Sprintf("%s://%s/%s/%s", uri.Scheme, uri.Host, uri.Path, link)
	}
}

func fetch(target string) ([]string, error) {
	uri, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	var urls []string
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		link, ok := s.Attr("src")
		if !ok {
			return
		}
		urls = append(urls, makelink(uri, link))
	})
	return urls, nil
}

func saveimg(dir, target string) error {
	log.Print(target)
	uri, err := url.Parse(target)
	if err != nil {
		return err
	}

	resp, err := http.Get(target)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	name := path.Base(uri.Path)
	fullpath := filepath.Join(dir, name)
	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(f, resp.Body)
	return nil
}

func saveimgs(dir string, urls []string) error {
	pool := NewTaskPool(5)
	pool.Start()
	for _, url := range urls {
		url := url
		pool.Submit(func() {
			if err := saveimg(dir, url); err != nil {
				log.Print(err)
			}
		})
	}
	pool.Stop()
	pool.Wait()
	return nil
}

func maketar(dir string, w io.Writer) error {
	base := filepath.Base(dir)
	tr := tar.NewWriter(w)
	defer tr.Close()
	return filepath.Walk(dir, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()

		h, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		p, _ := filepath.Rel(dir, name)
		h.Name = filepath.Join(base, p)
		if err = tr.WriteHeader(h); err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			io.Copy(tr, f)
		}
		return nil
	})
}

func main() {
	url := os.Args[1]
	urls, err := fetch(url)
	if err != nil {
		log.Panic(err)
	}

	dir, err := ioutil.TempDir("", "img")
	if err != nil {
		log.Panic(err)
	}

	defer os.RemoveAll(dir)

	err = saveimgs(dir, urls)
	if err != nil {
		log.Panic(err)
	}

	err = maketar(dir, os.Stdout)
	if err != nil {
		log.Panic(err)
	}
}

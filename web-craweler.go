package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type SafeCache struct {
	cache map[string]bool
	mux   sync.Mutex
}

func (c *SafeCache) Set(s string) {
	c.mux.Lock()
	c.cache[s] = true
	c.mux.Unlock()
}

func (c *SafeCache) Get(s string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.cache[s]
}

var (
	sc = SafeCache{cache: make(map[string]bool)}
	errs, ress = make(chan error), make(chan string)
)

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}

	var (
		body string
		err error
		urls []string
	)

	if ok := sc.Get(url); !ok {
		sc.Set(url)
		body, urls, err = fetcher.Fetch(url)
	} else {
		err = fmt.Errorf("Already fetched: %s", url)
	}

	if err != nil {
		errs <- err
		return
	}

	ress <- fmt.Sprintf("found: %s %q\n", url, body)
	for _, u := range urls {
		go Crawl(u, depth-1, fetcher)
	}
	return
}

func main() {
	go Crawl("http://golang.org/", 4, fetcher)
	for {
		select {
		case res, ok := <-ress:
			fmt.Println(res)
			if !ok {
				break
			}
		case err, ok := <-errs:
			fmt.Println(err)
			if !ok {
				break
			}
		}
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

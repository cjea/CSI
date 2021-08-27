package main

import (
	"fmt"
	"math"
	"time"

	"github.com/valyala/fasthttp"
)

/*
                                             |               |
                            ||_|_|_|_|_|_|_|_|               |__|__|__|__|__|__|__|__|__|   <------------\
                          /                 \                                         \    URL QUEUE      \
             ____________/_                  \____________                             \_____________      \
             | subscriber |                  | publisher |                              |           |      |
                                                                                        |           |      |
WritePage :: (Page, Directory) -> ()  |     FetchPage :: URL -> Page
DequeuePage :: Queue -> Page          |     ParsePage :: Page -> HTMLNode
NewQueue :: () -> Queue               |     GetAnchorTags :: HTMLNode -> []AnchorTag
GetDir :: Page -> Directory           |     GetURLs :: []AnchorTag -> []URL
                                      |     PublishPage :: (Page, Queue) -> ()
                                      |     Dequeue :: Queue -> URL
                                      |     Enqueue :: Queue -> URL
                      Page :: { URL, Body }

Indexer flow                                     Fetcher flow
============                                     ============
NewQueue |> Dequeue |> GetDir |> WritePage       NewQueue |> Dequeue |> FetchPage |> ParsePage |> GetAnchorTags |> GetURLs |> Publish(URLQueue)
                                                                     |
                                                                     |> Publish

*/

func main() {
	fmt.Println("Running")
	done := make(chan int)
	pageQueue := NewPageQueue(100)

	go Publish(pageQueue, done)
	go Subscribe(pageQueue, done)
}

func Subscribe(pageQueue PageQueue, quit chan int) {
	for {
		done := make(chan int)
		go WritePage(DequeuePage(pageQueue), done)
		select {
		case <-done:
			continue
		case <-quit:
			return
		}
	}
}

func Publish(pageQueue PageQueue, quit chan int) {
	urlQueue := NewURLQueue(100)
	page := FetchURL(DequeueURL(urlQueue))
	if page.URL.Value == "" {
		return
	}

	go EnqueuePage(pageQueue, page)
	go Crawl(urlQueue, page)
}

func Crawl(urlQueue URLQueue, page Page) {
	for _, href := range GetURLs(GetAnchorTags(ParsePage(page))) {
		PublishURL(urlQueue, href)
	}
}

func WritePage(page Page, done chan int) {
}

const maxWaitSeconds = 3

func FetchURL(url URL) Page {
	status, body, err := fasthttp.Get(nil, url.Value)
	if err != nil {
		return Page{}
	}
	if status > 299 {
		for retries := float64(0); retries < 4 && status > 299; retries++ {
			sleepTime := math.Min(float64(math.Pow(2, retries)), maxWaitSeconds)
			time.Sleep(time.Duration(sleepTime) * time.Second)
			status, body, err = fasthttp.Get(nil, url.Value)
		}
	}
	return Page{
		URL:  url,
		Data: body,
	}
}

func ParsePage(Page) HTMLNode {
	return HTMLNode{}
}

func GetAnchorTags(HTMLNode) []AnchorTag {
	return []AnchorTag{}

}

func NewPageQueue(slots int) PageQueue {
	return make(chan Page, slots)
}

func NewURLQueue(slots int) URLQueue {
	return make(chan URL, slots)
}

func GetURLs([]AnchorTag) []URL {
	return []URL{}
}

func GetDir(Page) Directory {
	return Directory{}
}

func EnqueuePage(pageQueue PageQueue, page Page) {
	select {
	case pageQueue <- page:
	case <-time.After(5 * time.Second):
		fmt.Println("Couldn't publish page for more than 5 seconds")
	}
}

func DequeuePage(queue PageQueue) Page {
	page := Page{}
	select {
	case page = <-queue:
		return page
	case <-time.After(5 * time.Second):
		fmt.Println("No message after 5 seconds")
		return Page{}
	}
}

func DequeueURL(queue URLQueue) URL {
	url := URL{}
	select {
	case url = <-queue:
		return url
	case <-time.After(5 * time.Second):
		fmt.Println("No message after 5 seconds")
		return URL{}
	}
}

func PublishURL(queue URLQueue, url URL) {
	select {
	case queue <- url:
	case <-time.After(5 * time.Second):
		fmt.Println("Couldn't publish URL for more than 5 seconds")
	}
}

type Page struct {
	URL  URL
	Data []byte
}
type PageQueue chan Page
type URLQueue chan URL
type AnchorTag struct{}
type Directory struct{}
type URL struct {
	Value string
}
type HTMLNode struct{}
type Subscription struct{}

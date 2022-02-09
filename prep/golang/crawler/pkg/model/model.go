package model

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var (
	ErrTimeout = errors.New("timed out")
)

type Page struct {
	URL  URL
	Data []byte
}

func (page Page) ParseHTML() (*html.Node, error) {
	doc, err := html.Parse(bytes.NewReader(page.Data))
	if err != nil {
		fmt.Printf("Could not parse html: %v\n", err)
		return nil, err
	}

	return doc, nil
}

type PageQueue struct {
	seen *sync.Map
	Q    chan Page
}

type URLQueue struct {
	seen *sync.Map
	Q    chan URL
}

func (pageQueue PageQueue) Enqueue(page Page) {
	if pageQueue.HasSeen(page.URL) {
		return
	}
	select {
	case pageQueue.Q <- page:
		pageQueue.MarkSeen(page.URL)
		return
	case <-time.After(5 * time.Second):
		fmt.Println("Couldn't publish page for more than 5 seconds")
	}
}

func (urlQueue URLQueue) Enqueue(u URL) {
	parsed, err := url.Parse(u.Value)
	if err != nil {
		return
	}
	if parsed.Host == "" {

		parsed.Host = "www.ebay.com"
	}
	if urlQueue.HasSeen(u) {
		return
	}

	select {
	case urlQueue.Q <- u:
		urlQueue.MarkSeen(u)
		return
	case <-time.After(5 * time.Second):
		fmt.Println("Couldn't publish URL for more than 5 seconds")
	}
}

func (queue PageQueue) Dequeue() (Page, error) {
	select {
	case page := <-queue.Q:
		return page, nil
	case <-time.After(5 * time.Second):
		fmt.Println("No message after 5 seconds")
		return Page{}, ErrTimeout
	}
}

func (queue URLQueue) Dequeue() (URL, error) {
	select {
	case url := <-queue.Q:
		return url, nil
	case <-time.After(5 * time.Second):
		fmt.Println("No message after 5 seconds")
		return URL{}, ErrTimeout
	}
}

type AnchorTag struct{ Value string }
type Directory struct{ Path string }
type URL struct {
	Value string
}
type HTMLNode struct {
	Node *html.Node
}

func (uq URLQueue) HasSeen(u URL) bool {
	_, ok := uq.seen.Load(u.Value)
	return ok
}
func ApplyAllHTMLNode(n *html.Node, f func(*html.Node)) {
	queue := []*html.Node{n}
	for len(queue) > 0 {
		el := queue[0]
		queue = queue[1:]
		if el == nil {
			continue
		}
		for ; el != nil; el = el.NextSibling {
			f(el)
			queue = append(queue, el.FirstChild)
		}
	}
}

func NewPageQueue(slots int) PageQueue {
	return PageQueue{
		seen: &sync.Map{},
		Q:    make(chan Page, slots),
	}
}

func (pq PageQueue) MarkSeen(u URL) {
	pq.seen.Store(u.Value, struct{}{})
}

func (pq PageQueue) HasSeen(u URL) bool {
	_, ok := pq.seen.Load(u.Value)
	return ok
}

func NewURLQueue(slots int) URLQueue {
	return URLQueue{
		seen: &sync.Map{},
		Q:    make(chan URL, slots),
	}
}

func (uq URLQueue) MarkSeen(u URL) {
	uq.seen.Store(u.Value, struct{}{})
}

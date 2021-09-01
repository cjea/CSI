package main

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"golang.org/x/net/html"
)

/*
                                             |               |
                            ||_|_|_|_|_|_|_|_|               |__|__|__|__|__|__|__|__|__|   <------------\
                          /                 \                                         \    URL QUEUE      \
             ____________/_                  \____________                             \_____________      \
             | subscriber |                  | publisher |                              |           |      |
                                                                                        |           |      |
WritePage :: (Page, Directory) -> ()  |     FetchPage :: URL -> Page
DequeuePage :: Queue -> Page          |     ParseHTML :: Page -> HTMLNode
NewQueue :: () -> Queue               |     GetAnchorTags :: HTMLNode -> []AnchorTag
GetDir :: Page -> Directory           |     GetURLs :: []AnchorTag -> []URL
                                      |     PublishPage :: (Page, Queue) -> ()
                                      |     Dequeue :: Queue -> URL
                                      |     Enqueue :: (URL, Queue) -> ()
                      Page :: { URL, Body }

Indexer flow                                     Fetcher flow
============                                     ============
NewQueue |> Dequeue |> GetDir |> WritePage       NewQueue |> Dequeue |> FetchPage |> ParseHTML |> GetAnchorTags |> GetURLs |> Publish(URLQueue)
                                                                     |
                                                                     |> Publish

*/

func main() {
	fmt.Println("Running")
	done := make(chan int)
	pageQueue := NewPageQueue(100)
	urlQueue := NewURLQueue(100)

	go EnqueueURL(urlQueue, URL{Value: "https://www.ebay.com"})
	go Publish(pageQueue, urlQueue, done)
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

func Publish(pageQueue PageQueue, urlQueue URLQueue, quit chan int) {
	page := FetchURL(DequeueURL(urlQueue))
	if page.URL.Value == "" {
		return
	}

	go EnqueuePage(pageQueue, page)
	go Crawl(urlQueue, page)
}

func Crawl(urlQueue URLQueue, page Page) {
	for _, href := range GetURLs(GetAnchorTags(ParseHTML(page))) {
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
		fmt.Printf("Failed to fetch %s\n", url)
		return Page{}
	}
	return Page{
		URL:  url,
		Data: body,
	}
}

func ParseHTML(page Page) HTMLNode {
	doc, err := html.Parse(bytes.NewReader(page.Data))
	if err != nil {
		fmt.Printf("Could not parse html: %v\n", err)
		return HTMLNode{}
	}

	return HTMLNode{Node: doc}
}

func GetAnchorTags(n HTMLNode) []AnchorTag {
	ret := []AnchorTag{}
	n.ApplyAll(func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					ret = append(ret, AnchorTag{Value: attr.Val})
				}
			}
		}
	})
	return ret

}

func GetURLs(tags []AnchorTag) []URL {
	ret := []URL{}
	for _, tag := range tags {
		ret = append(ret, URL{tag.Value})
	}
	return ret
}

func NewPageQueue(slots int) PageQueue {
	return PageQueue{
		Seen: map[string]struct{}{},
		Q:    make(chan Page, slots),
	}
}

func NewURLQueue(slots int) URLQueue {
	return URLQueue{
		Seen: map[string]struct{}{},
		Q:    make(chan URL, slots),
	}
}

func GetDir(page Page) Directory {
	// u := page.URL.Value
	// if !(strings.HasPrefix(u, "http") || strings.HasPrefix(u, "https")) {
	// 	if strings.Index(u, ".") < strings.Index(u, "")
	// 	u = "//" + u
	// }

	return Directory{Path: fmt.Sprintf(`./%s/`, page.URL.Value)}
}

func EnqueuePage(pageQueue PageQueue, page Page) {
	if _, ok := pageQueue.Seen[page.URL.Value]; ok {
		return
	}

	select {
	case pageQueue.Q <- page:
		pageQueue.Seen[page.URL.Value] = struct{}{}
		return
	case <-time.After(5 * time.Second):
		fmt.Println("Couldn't publish page for more than 5 seconds")
	}
}

func EnqueueURL(urlQueue URLQueue, u URL) {
	if _, ok := urlQueue.Seen[u.Value]; ok {
		return
	}

	select {
	case urlQueue.Q <- u:
		urlQueue.Seen[u.Value] = struct{}{}
		return
	case <-time.After(5 * time.Second):
		fmt.Println("Couldn't publish URL for more than 5 seconds")
	}
}

var (
	chars              = `[^\/]+`
	dot                = `\.`
	slash              = `\/`
	blankDotBlank      = regexp.MustCompile(chars + dot + chars)
	blankDotBlankSlash = regexp.MustCompile(chars + dot + chars + slash)
)

// TODO(cjea): figure out how to do this sanely.
// Problems e.g. how do you know if "index.html" is a host, or a path?
func ensureFullURL(u string, host string) string {
	hasScheme := strings.HasPrefix(u, "http") || strings.HasPrefix(u, "https") ||
		strings.HasPrefix(u, "||")
	bs := []byte(u)
	if u[0] == '/' {
		return "//" + host + u
	}
	if matches := blankDotBlankSlash.FindSubmatch(bs); len(matches) > 0 {
		if hasScheme {
			return u
		}
		return "//" + u
	} else if matches := blankDotBlank.FindSubmatch(bs); len(matches) > 0 {
		if string(matches[0]) == host {
			if hasScheme {
				return u
			}
			return "//" + u
		}
		return "//" + host + "/" + u
	} else if !hasScheme {
		return "//" + host + "/" + u
	} else {
		// can probably never get here
		return u
	}
}

func DequeuePage(queue PageQueue) Page {
	page := Page{}
	select {
	case page = <-queue.Q:
		return page
	case <-time.After(5 * time.Second):
		fmt.Println("No message after 5 seconds")
		return page
	}
}

func DequeueURL(queue URLQueue) URL {
	select {
	case url := <-queue.Q:
		return url
	case <-time.After(5 * time.Second):
		fmt.Println("No message after 5 seconds")
		return URL{}
	}
}

func rewriteHost(url URL, host string) {

}

func PublishURL(queue URLQueue, url URL) {
	select {
	case queue.Q <- url:
	case <-time.After(5 * time.Second):
		fmt.Println("Couldn't publish URL for more than 5 seconds")
	}
}

type Page struct {
	URL  URL
	Data []byte
}
type PageQueue struct {
	Seen map[string]struct{}
	Q    chan Page
}

type URLQueue struct {
	Seen map[string]struct{}
	Q    chan URL
}

type AnchorTag struct{ Value string }
type Directory struct{ Path string }
type URL struct {
	Value string
}
type HTMLNode struct {
	Node *html.Node
}

func (n HTMLNode) ApplyAll(f func(*html.Node)) {
	queue := []*html.Node{n.Node}
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

type Subscription struct{}

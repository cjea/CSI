package pub

import (
	"csi_prep_golang/pkg/model"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"golang.org/x/net/html"
)

func Publish(pageQueue model.PageQueue, urlQueue model.URLQueue, done chan error) {
	u, err := urlQueue.Dequeue()
	if err != nil {
		done <- err
		return
	}
	page, err := FetchURL(u)
	if err != nil || page.URL.Value == "" {
		done <- err
		return
	}
	grp := sync.WaitGroup{}
	grp.Add(2)
	go func() {
		pageQueue.Enqueue(page)
		grp.Done()
	}()
	go func() {
		parsed, err := page.ParseHTML()
		if err != nil {
			grp.Done()
			return
		}
		for _, u := range GetURLs(GetAnchorTags(parsed)) {
			urlQueue.Enqueue(u)
		}
		grp.Done()
	}()
	grp.Wait()
	done <- nil
}

func Crawl(urlQueue model.URLQueue, n *html.Node) error {
	for _, u := range GetURLs(GetAnchorTags(n)) {
		urlQueue.Enqueue(u)
	}
	return nil
}

const maxWaitSeconds = 3

func FetchURL(url model.URL) (model.Page, error) {
	status, body, err := fasthttp.Get(nil, url.Value)
	if err != nil {
		return model.Page{}, err
	}
	if status > 299 {
		for retries := float64(0); retries < 4 && status > 299; retries++ {
			sleepTime := math.Min(float64(math.Pow(2, retries)), maxWaitSeconds)
			time.Sleep(time.Duration(sleepTime) * time.Second)
			status, body, err = fasthttp.Get(nil, url.Value)
		}
		return model.Page{}, fmt.Errorf(
			"failed to fetch '%s' (status %d)",
			url, status,
		)
	}
	return model.Page{
		URL:  url,
		Data: body,
	}, nil
}

func GetAnchorTags(n *html.Node) []model.AnchorTag {
	ret := []model.AnchorTag{}
	model.ApplyAllHTMLNode(n, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					ret = append(ret, model.AnchorTag{Value: attr.Val})
				}
			}
		}
	})
	return ret

}

func GetURLs(tags []model.AnchorTag) []model.URL {
	ret := []model.URL{}
	for _, tag := range tags {
		ret = append(ret, model.URL{Value: tag.Value})
	}
	return ret
}

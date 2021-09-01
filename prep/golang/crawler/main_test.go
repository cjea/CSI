package main_test

import (
	"fmt"
	"regexp"
	"time"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
	"golang.org/x/net/html"

	main "csi_prep_golang"
)

var (
	asyncTestTimeout = float64(time.Second * 5)
)

var (
	whitespace = regexp.MustCompile(`\s+`)
)

var _ = gk.Describe(".FetchURL", func() {
	gk.It("fetches a URL", func() {
		url := main.URL{Value: "https://www.ebay.com"}
		page := main.FetchURL(url)
		str := string(whitespace.ReplaceAll(page.Data[0:20], nil))
		gm.Expect(str).To(gm.HavePrefix("<!DOCTYPEhtml>"))
		gm.Expect(page.URL.Value).To(gm.Equal("https://www.ebay.com"))
	})
})

var _ = gk.Describe(".ParseHTML", func() {
	gk.It("parses HTML", func() {
		var key, val string
		page := main.Page{Data: []byte(`<a href="fake.com">fake body</a>`)}
		parsed := main.ParseHTML(page)
		parsed.ApplyAll(func(n *html.Node) {
			if n.Data == "a" {
				hrefAttr := n.Attr[0]
				key, val = hrefAttr.Key, hrefAttr.Val
				return
			}
		})
		gm.Expect(key).To(gm.Equal("href"))
		gm.Expect(val).To(gm.Equal("fake.com"))
	})
})

var _ = gk.Describe(".GetAnchorTags", func() {
	gk.It("returns a list of AnchorTags", func() {
		page := main.Page{Data: []byte(`
<a href="fake.com">fake body</a>
<p href="bad.com">bad</p>
<a href="fake2.com">fake body 2</a>
`,
		)}
		tags := main.GetAnchorTags(main.ParseHTML(page))
		gm.Expect(tags).To(gm.ConsistOf(
			main.AnchorTag{Value: "fake.com"},
			main.AnchorTag{Value: "fake2.com"},
		))
	})
})

var _ = gk.Describe("PageQueue", func() {
	gk.Describe("EnqueuePage", func() {
		gk.It("enqueues Pages", func() {
			queue := main.NewPageQueue(5)
			main.EnqueuePage(queue, main.Page{URL: main.URL{Value: "url-1"}})
			main.EnqueuePage(queue, main.Page{URL: main.URL{Value: "url-2"}})

			gm.Expect(queue.Q).To(gm.HaveLen(2))
		}, asyncTestTimeout)

		gk.It("does not enqueue the same URL twice", func() {
			queue := main.NewPageQueue(5)
			main.EnqueuePage(queue, main.Page{URL: main.URL{Value: "url-1"}})
			main.EnqueuePage(queue, main.Page{URL: main.URL{Value: "url-1"}})

			gm.Expect(queue.Q).To(gm.HaveLen(1))
		})
	})

	gk.Describe("DequeuePage", func() {
		gk.It("dequeues Pages asynchronously", func() {
			queue := main.NewPageQueue(5)
			urls := []string{}
			for i := 0; i < 5; i++ {
				u := main.URL{Value: fmt.Sprintf("url-%d", i)}
				go main.EnqueuePage(queue, main.Page{URL: u})
			}
			for i := 0; i < 5; i++ {
				p := main.DequeuePage(queue)
				urls = append(urls, p.URL.Value)
			}
			gm.Expect(urls).To(gm.ContainElements([]string{
				"url-0",
				"url-1",
				"url-2",
				"url-3",
				"url-4",
			}))
		}, asyncTestTimeout)
	})
})

var _ = gk.Describe("URLQueue", func() {
	gk.Describe("EnqueueURL", func() {
		gk.It("enqueues URLs", func() {
			queue := main.NewURLQueue(5)
			main.EnqueueURL(queue, main.URL{Value: "url-1"})
			main.EnqueueURL(queue, main.URL{Value: "url-2"})

			gm.Expect(queue.Q).To(gm.HaveLen(2))
		}, asyncTestTimeout)

		gk.It("does not enqueue the same URL twice", func() {
			queue := main.NewURLQueue(5)
			main.EnqueueURL(queue, main.URL{Value: "url-1"})
			main.EnqueueURL(queue, main.URL{Value: "url-1"})

			gm.Expect(queue.Q).To(gm.HaveLen(1))
		})
	})

	gk.Describe("DequeueURL", func() {
		gk.It("dequeues URLs asynchronously", func() {
			queue := main.NewURLQueue(5)
			urls := []string{}
			for i := 0; i < 5; i++ {
				go main.EnqueueURL(queue, main.URL{Value: fmt.Sprintf("url-%d", i)})
			}
			for i := 0; i < 5; i++ {
				u := main.DequeueURL(queue)
				urls = append(urls, u.Value)
			}
			gm.Expect(urls).To(gm.ContainElements([]string{
				"url-0",
				"url-1",
				"url-2",
				"url-3",
				"url-4",
			}))
		}, asyncTestTimeout)
	})
})

var _ = gk.Describe(".Page{}", func() {
	// gk.Describe("#GetDir", func() {
	// 	gk.It("returns the page's URL under current directory", func() {
	// 		p := main.Page{URL: main.URL{Value: "fake.com/foo"}}
	// 		gm.Expect(main.GetDir(p)).To(gm.Equal("fake.com"))
	// 	})
	// })
})

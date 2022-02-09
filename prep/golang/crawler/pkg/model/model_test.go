package model_test

import (
	"fmt"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
	"golang.org/x/net/html"

	"csi_prep_golang/pkg/model"
)

var _ = gk.Describe("Model", func() {

	var _ = gk.Describe("Page", func() {
		gk.Describe("#ParseHTML", func() {
			gk.It("parses HTML", func() {
				var key, val string
				page := model.Page{Data: []byte(`<a href="fake.com">fake body</a>`)}
				n, err := page.ParseHTML()
				noErr(err)
				model.ApplyAllHTMLNode(n, func(n *html.Node) {
					if n.Data == "a" {
						hrefAttr := n.Attr[0]
						key, val = hrefAttr.Key, hrefAttr.Val
						return
					}
				})
				gm.Expect(key).To(gm.Equal("href"))
				gm.Expect(val).To(gm.Equal("fake.com"))
			})

			// gk.It("has a nil Node for invalid HTML", func() {
			// 	page := model.Page{Data: []byte(`<not html`)}
			// 	x, err := page.ParseHTML()
			// 	fmt.Printf("\n*****\n%#v\n", x.FirstChild.Data)
			// 	gm.Expect(err).To(gm.HaveOccurred())

			// })
		})
	})

	var _ = gk.Describe("PageQueue", func() {
		gk.Describe("#Enqueue", func() {
			gk.It("enqueues Pages", func() {
				queue := model.NewPageQueue(5)
				queue.Enqueue(model.Page{URL: model.URL{Value: "url-1"}})
				queue.Enqueue(model.Page{URL: model.URL{Value: "url-2"}})

				gm.Expect(queue.Q).To(gm.HaveLen(2))
			})

			gk.It("does not enqueue the same URL twice", func() {
				queue := model.NewPageQueue(5)
				queue.Enqueue(model.Page{URL: model.URL{Value: "url-1"}})
				queue.Enqueue(model.Page{URL: model.URL{Value: "url-1"}})

				gm.Expect(queue.Dequeue().URL.Value).To(gm.Equal("url-1"))
				gm.Expect(queue.Q).To(gm.BeEmpty())
			})
		})

		gk.Describe("DequeuePage", func() {
			gk.It("dequeues Pages asynchronously", func() {
				queue := model.NewPageQueue(5)
				urls := []string{}
				go func() {
					for i := 0; i < 5; i++ {
						u := model.URL{Value: fmt.Sprintf("url-%d", i)}
						go queue.Enqueue(model.Page{URL: u})
					}
				}()
				for i := 0; i < 5; i++ {
					p := queue.Dequeue()
					urls = append(urls, p.URL.Value)
				}
				gm.Expect(urls).To(gm.ContainElements([]string{
					"url-0",
					"url-1",
					"url-2",
					"url-3",
					"url-4",
				}))
			})
		})
	})

	var _ = gk.Describe("URLQueue", func() {
		gk.Describe("Enqueue", func() {
			gk.It("enqueues URLs", func() {
				queue := model.NewURLQueue(5)
				queue.Enqueue(model.URL{Value: "url-1"})
				queue.Enqueue(model.URL{Value: "url-2"})

				gm.Expect(queue.Q).To(gm.HaveLen(2))
			})

			gk.It("does not enqueue the same URL twice", func() {
				queue := model.NewURLQueue(5)
				queue.Enqueue(model.URL{Value: "url-1"})
				queue.Enqueue(model.URL{Value: "url-1"})

				gm.Expect(queue.Dequeue().Value).To(gm.Equal("url-1"))
				gm.Expect(queue.Q).To(gm.BeEmpty())
			})
		})

		gk.Describe("DequeueURL", func() {
			gk.It("dequeues URLs asynchronously", func() {
				queue := model.NewURLQueue(5)
				urls := []string{}
				for i := 0; i < 5; i++ {
					go queue.Enqueue(model.URL{Value: fmt.Sprintf("url-%d", i)})
				}
				for i := 0; i < 5; i++ {
					u := queue.Dequeue()
					urls = append(urls, u.Value)
				}
				gm.Expect(urls).To(gm.ContainElements([]string{
					"url-0",
					"url-1",
					"url-2",
					"url-3",
					"url-4",
				}))
			})
		})
	})
})

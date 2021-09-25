package pub_test

import (
	"bytes"
	"net"
	"regexp"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
	"golang.org/x/net/html"

	"csi_prep_golang/pkg/model"
	"csi_prep_golang/pkg/pub"
)

var (
	whitespace = regexp.MustCompile(`\s+`)
)

var _ = gk.Describe("Pub", func() {
	var _ = gk.Describe(".FetchURL", func() {
		gk.It("fetches a URL", func() {
			url := model.URL{Value: "https://www.ebay.com"}
			page, err := pub.FetchURL(url)
			noErr(err)
			gm.Expect(len(page.Data)).To(gm.BeNumerically(">", 0))

			gm.Expect(page.URL.Value).To(gm.Equal("https://www.ebay.com"))
			gm.Expect(string(page.Data)).To(gm.HavePrefix("<!DOCTYPE html>"))
		})

		gk.It("returns an empty page when fetch fails", func() {
			url := model.URL{Value: "bad-url"}
			page, err := pub.FetchURL(url)
			_, ok := err.(*net.DNSError)

			gm.Expect(ok).To(gm.BeTrue())
			gm.Expect(page.URL.Value).To(gm.Equal(""))
			gm.Expect(page.Data).To(gm.BeEmpty())
		})
	})

	var _ = gk.Describe(".GetAnchorTags", func() {
		gk.It("returns a list of AnchorTags", func() {
			d := []byte(
				`
<a href="fake.com">fake body</a>
<p href="bad.com">bad</p>
<a href="fake2.com">fake body 2</a>
`,
			)
			n, err := html.Parse(bytes.NewReader(d))
			noErr(err)

			tags := pub.GetAnchorTags(n)
			gm.Expect(tags).To(gm.ConsistOf(
				model.AnchorTag{Value: "fake.com"},
				model.AnchorTag{Value: "fake2.com"},
			))
		})
	})
})

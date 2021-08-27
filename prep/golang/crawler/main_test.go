package main_test

import (
	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"

	main "csi_prep_golang"
)

var _ = gk.Describe(".FetchURL", func() {
	gk.It("fetches a URL", func() {
		url := main.URL{Value: "https://www.ebay.com"}
		page := main.FetchURL(url)
		gm.Expect(page.URL.Value).To(gm.Equal("https://www.ebay.com"))
		gm.Expect(len(page.Data)).To(gm.BeNumerically(">", 0))
	})
})

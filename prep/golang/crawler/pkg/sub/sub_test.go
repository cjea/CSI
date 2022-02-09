package sub_test

import (
	"io/ioutil"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"

	"csi_prep_golang/pkg/model"
	"csi_prep_golang/pkg/sub"
)

var _ = gk.Describe("Sub", func() {
	var _ = gk.Describe(".WritePage", func() {
		gk.It("writes an HTML page", func() {
			c := make(chan int)
			go sub.WritePage(
				model.Page{
					URL:  model.URL{Value: "a.com/"},
					Data: []byte(`<a href="supdude.com">linky</a>`),
				},
				"test/fixtures",
				c,
			)
			gm.Expect(<-c).To(gm.Equal(0))
			data, err := ioutil.ReadFile("test/fixtures/a.com")
			noErr(err)
			gm.Expect(string(data)).To(gm.Equal(`<a href="supdude.com">linky</a>`))
		})

		// gk.It("fails if the page is not valid HTML", func() {
		// 	c := make(chan int)
		// 	go sub.WritePage(
		// 		model.Page{
		// 			URL:  model.URL{Value: "a.com/"},
		// 			Data: []byte(`not html`),
		// 		},
		// 		"test/fixtures",
		// 		c,
		// 	)
		// 	gm.Expect(<-c).To(gm.Equal(1))

		// })
	})
})

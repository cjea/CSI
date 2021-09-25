package bin_den_test

import (
	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"

	"float-to-den/pkg/bin_den"
)

var _ = gk.Describe("Bin_den", func() {
	gk.Describe(".Btod", func() {
		gk.It("returns the steps for converting bits to denary", func() {
			steps := bin_den.Btod("00001010")
			gm.Expect(steps.Steps).To(gm.HaveLen(8))
		})

		gk.It("gets the right result", func() {
			steps := bin_den.Btod("00001010")
			gm.Expect(steps.Calc()).To(gm.Equal(float64(10)))
		})
	})

	gk.Describe(".BtodFraction", func() {
		gk.It("returns the steps for converting fractional bits to denary", func() {
			steps := bin_den.BtodFraction("00001010")
			gm.Expect(steps.Steps).To(gm.HaveLen(8))
		})

		gk.It("gets the right result", func() {
			steps := bin_den.BtodFraction("11001")
			gm.Expect(steps.Calc()).To(gm.Equal(float64(0.78125)))
		})
	})
})

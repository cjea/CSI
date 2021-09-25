package floatbits_test

import (
	"strings"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"

	"float-to-den/pkg/floatbits"
)

var _ = gk.Describe("Floatbits", func() {
	gk.Describe(".New", func() {
		gk.It("returns a 32-bit float", func() {
			bits := "00000001001000000000010000000000"
			f, err := floatbits.New(bits)
			noErr(err)
			gm.Expect(f).
				To(gm.Equal(floatbits.FloatBits("00000001001000000000010000000000")))
		})

		gk.It("returns an error if not enough bits", func() {
			_, err := floatbits.New("1")
			gm.Expect(err).To(gm.HaveOccurred())
		})

		gk.It("returns an error if too many bits", func() {
			_, err := floatbits.New(strings.Repeat("0", 33))
			gm.Expect(err).To(gm.HaveOccurred())
		})

		gk.It("returns an error if not valid bits", func() {
			_, err := floatbits.New(strings.Repeat("a", 32))
			gm.Expect(err).To(gm.HaveOccurred())
		})
	})

	gk.Describe("FloatBits", func() {
		gk.Describe("#Sign", func() {
			gk.It("treats bits with leading-0 as positive", func() {
				bits, err := floatbits.New(strings.Repeat("0", 32))
				noErr(err)
				gm.Expect(bits.Sign()).To(gm.Equal(floatbits.SIGN_POSITIVE))
			})

			gk.It("treats bits with leading-1 as negative", func() {
				bits, err := floatbits.New("1" + strings.Repeat("0", 31))
				noErr(err)
				gm.Expect(bits.Sign()).To(gm.Equal(floatbits.SIGN_NEGATIVE))
			})
		})

		gk.Describe("#BiasedExp", func() {
			gk.It("returns the first eight bits after the sign bit", func() {
				bits, err := floatbits.New("010101010" + strings.Repeat("0", 23))
				noErr(err)
				gm.Expect(bits.BiasedExp()).To(gm.Equal("10101010"))
			})
		})

		gk.Describe("#Mantissa", func() {
			gk.It("returns the last 23 bits", func() {
				bits, err := floatbits.New("000000000" + "10101010101010101010101")
				noErr(err)
				gm.Expect(bits.Mantissa()).To(gm.Equal("10101010101010101010101"))
			})
		})
	})
})

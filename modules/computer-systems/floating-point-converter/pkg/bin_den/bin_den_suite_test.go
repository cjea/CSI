package bin_den_test

import (
	"testing"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
)

func TestFloatbits(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "Bin_den Suite")
}

func noErr(err error) {
	gm.Expect(err).NotTo(gm.HaveOccurred())
}

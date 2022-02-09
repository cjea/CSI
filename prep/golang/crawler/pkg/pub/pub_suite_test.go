package pub_test

import (
	"testing"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
)

func TestPub(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "Pub Suite")
}

func noErr(err error) {
	gm.Expect(err).NotTo(gm.HaveOccurred())
}

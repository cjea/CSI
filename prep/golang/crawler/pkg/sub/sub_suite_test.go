package sub_test

import (
	"testing"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
)

func TestSub(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "Sub Suite")
}

func noErr(err error) {
	gm.Expect(err).NotTo(gm.HaveOccurred())
}

package model_test

import (
	"testing"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
)

func TestModel(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "Model Suite")
}

func noErr(err error) {
	gm.Expect(err).NotTo(gm.HaveOccurred())
}

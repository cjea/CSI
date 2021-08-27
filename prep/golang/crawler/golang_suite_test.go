package main_test

import (
	"testing"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
)

func TestGolang(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "site crawler tests")
}

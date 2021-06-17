package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLocal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Test")
}

var _ = Describe("Build Workflows Test", func() {
	Context("Get a specific workflow", func() {
		It("get docker build image", func() {
			_, err := GetSnippet("docker")
			Expect(err).Should(BeNil())
		})
	})
})

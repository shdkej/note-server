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
		files, _ := GetFile("/home/sh/workspace/golang/src/github.com/shdkej/note-server/snippet", ".yml")

		It("get docker build image", func() {
			d, err := GetSnippet("docker")
			Expect(d).Should(Equal(""))
			Expect(err).Should(BeNil())
		})
	})
})

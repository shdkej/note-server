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

var _ = Describe("Parser Test", func() {
	Context("One async post test", func() {
		d := getSnippet()
		c := Counter{}

		It("sync once", func() {
			prevCount := ConfirmComplete()
			ch := make(chan bool, 1)

			c.SendPostRequest(d[0], ch)
			endCount := ConfirmComplete()
			Expect(endCount - prevCount).NotTo(Equal(prevCount))
		})
		It("async once", func() {
			count := c.runningAsyncronizly(d)
			Expect(count).Should(Equal(len(d)))
		})
		It("delete test", func() {
			CleanUp(d)
		})
	})
})

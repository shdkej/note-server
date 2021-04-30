package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestLocal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Test")
}

var _ = Describe("Parser Test", func() {
	Context("One async post test", func() {
		d := getSnippet()

		It("sync once", func() {
			prevCount := ConfirmComplete()
			ch := make(chan bool, len(d))

			SendPostRequest(d[0], ch)
			endCount := ConfirmComplete()
			Expect(endCount - prevCount).Should(Equal(len(d)))
		})
		It("async once", func() {
			prevCount := ConfirmComplete()
			runningAsyncronizly(d)
			endCount := ConfirmComplete()
			Expect(endCount - prevCount).Should(Equal(len(d)))
		})

		It("delete test", func() {
			//CleanUp(d)
		})
	})
})

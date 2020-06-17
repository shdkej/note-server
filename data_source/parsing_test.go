package data_source

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parsing file", func() {
	Context("Test Parsing", func() {
		taglines, err := getTagAll()
		It("get tagline all", func() {
			Expect(taglines).NotTo(BeNil())
			Expect(err).Should(BeNil())
		})
		It("get Tag List", func() {
			tagList, err := getTagList()
			Expect(tagList).NotTo(BeNil())
			Expect(err).Should(BeNil())
		})
	})
})

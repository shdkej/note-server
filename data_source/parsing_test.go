package data_source

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parsing file", func() {
	Context("Test CSV", func() {
		tags, err := getTagAll()
		if err != nil {
			return
		}
		csv := toCSV(tags)

		It("parse success", func() {
			Expect(csv).Should(BeNil())
		})
	})
	Context("Test Parsing", func() {
		taglines, err := getTagAll()
		It("get tagline all", func() {
			Expect(taglines).NotTo(BeNil())
			Expect(err).Should(BeNil())
		})
	})
})

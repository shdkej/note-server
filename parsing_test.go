package main

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
		taglines, err := getTaglineAll()
		It("get tagline all", func() {
			Expect(taglines[0]).NotTo(BeNil())
			Expect(err).Should(BeNil())
		})
		/*
			filename := "/home/sh/vimwiki/Architecture.md"
			result, _ := makeTagSet(filename)
			It("get tagline set", func() {
				Expect(result).NotTo(BeNil())
			})
		*/
	})
})

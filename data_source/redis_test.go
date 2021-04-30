package data_source

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Redis", func() {
	pool := Redis{}
	pool.Init()

	Context("Test ping", func() {
		It("pong", func() {
			Expect(pool.Ping()).Should(BeNil())
			Expect(pool.Hits("test")).NotTo(BeNil())
		})
	})

	Context("Test sets", func() {
		tag := Note{
			FileName: "main.md",
			Tag:      "Good",
			TagLine:  "Good Enough",
		}

		It("set Sets", func() {
			Expect(pool.SetStruct(tagPrefix, tag)).Should(BeNil())
		})
		It("get Sets", func() {
			Expect(pool.GetStruct(tagPrefix, tag.Tag)).Should(Equal(tag))
		})
		It("get empty Sets", func() {
			Expect(pool.GetStruct(tagPrefix, "empty")).Should(Equal(Note{}))
		})
		It("delete Sets", func() {
			Expect(pool.Delete(tag)).Should(BeNil())
		})
	})

	Context("Test Misc Function", func() {
		tag := Note{
			FileName: "main.md",
			Tag:      "Good",
			TagLine:  "Good Enough",
		}
		pool.SetStruct(tagPrefix, tag)
		tags, err := pool.GetTags(tagPrefix)
		It("get scan body", func() {
			Expect(tags).NotTo(BeNil())
		})
		It("check error", func() {
			Expect(err).Should(BeNil())
		})

		tagParagraph, err := pool.GetSet(tag.Tag)
		It("check tag paragraph value", func() {
			Expect(tagParagraph).NotTo(BeNil())
		})
		It("check error", func() {
			Expect(err).Should(BeNil())
		})
		It("delete sets", func() {
			Expect(pool.Delete(tag)).Should(BeNil())
		})
	})
})

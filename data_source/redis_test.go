package data_source

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRedis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redis Test")
}

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
		tag := Tag{
			FileName: "main.md",
			Tag:      "Good",
			TagLine:  "Good Enough",
		}

		It("set Sets", func() {
			Expect(pool.SetStruct(tag)).Should(BeNil())
		})
		It("get Sets", func() {
			Expect(pool.GetStruct(tag.Tag)).Should(Equal(tag))
		})
		It("get empty Sets", func() {
			Expect(pool.GetStruct("empty")).Should(Equal(Tag{}))
		})
		It("delete Sets", func() {
			Expect(pool.Delete(tag)).Should(BeNil())
		})
	})

	Context("Test Misc Function", func() {
		tag := Tag{
			FileName: "main.md",
			Tag:      "Good",
			TagLine:  "Good Enough",
		}
		pool.SetStruct(tag)
		tags, err := pool.GetTagList(tag.Tag)
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

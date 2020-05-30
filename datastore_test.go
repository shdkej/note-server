package main

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
	pool := Client{}
	pool.Init()

	Context("Test ping", func() {
		It("pong", func() {
			Expect(pool.ping()).Should(BeNil())
		})
	})

	Context("Test strings", func() {
		key := "test"
		value := "value"
		It("set strings", func() {
			Expect(pool.set(key, value)).Should(BeNil())
		})
		It("get strings", func() {
			Expect(pool.get(key)).Should(Equal("value"))
		})
		It("get strings empty", func() {
			Expect(pool.get("empty")).Should(BeEmpty())
		})
	})

	Context("Test sets", func() {
		tag := Tag{
			FileName:    "main.md",
			FileContent: "0",
			Tag:         "Good",
			TagLine:     "Good Enough",
		}

		It("set sets", func() {
			Expect(pool.setStruct(tag)).Should(BeNil())
		})
		It("get sets", func() {
			Expect(pool.getStruct(tag.Tag)).Should(Equal(tag))
		})
		It("get sets empty", func() {
			Expect(pool.getStruct("empty")).Should(Equal(Article{}))
		})
	})

	/*
		Context("Test initial content", func() {
			pool.setInitial()
			It("initial Dir", func() {
				Expect(os.Getenv("VIMWIKI")).Should(Equal("/home/sh/vimwiki"))
			})
			article, err := pool.getStruct("2020-04-06-WEEK14.md")
			It("get content1", func() {
				Expect(article.Filename).Should(Equal("2020-04-06-WEEK14.md"))
			})
			It("error check", func() {
				Expect(err).Should(BeNil())
			})
			It("get ##Need Component", func() {
				Expect(pool.get("##Need Component")).NotTo(BeNil())
			})
			It("keys space bar test", func() {
				Expect(pool.getTagParagraph("Need Component")).NotTo(BeNil())
			})
		})
	*/
})

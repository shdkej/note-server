package main

import (
	"os"
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
	pool.NewClient()

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
		article := Article{
			Title:    "test",
			Category: "test@gmail.com",
			Content:  "kim",
		}

		It("set sets", func() {
			Expect(pool.setStruct(article)).Should(BeNil())
		})
		It("get sets", func() {
			Expect(pool.getStruct(article.Title)).Should(Equal(article))
		})
		It("get sets empty", func() {
			Expect(pool.getStruct("empty")).Should(Equal(Article{}))
		})
	})

	Context("Test initial content", func() {
		pool.setInitial()
		//It("initial", func() {
		//})
		It("initial Dir", func() {
			Expect(os.Getenv("VIMWIKI")).Should(Equal("/home/sh/vimwiki"))
		})
		article, err := pool.getStruct("2020-04-06-WEEK14.md")
		It("get content1", func() {
			Expect(article.Title).Should(Equal("2020-04-06-WEEK14.md"))
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
})

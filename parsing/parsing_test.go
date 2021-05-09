package parsing

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"testing"
)

func TestLocal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Local file parsing Test")
}

var _ = Describe("Parsing file", func() {
	Context("Test Parsing", func() {
		It("get tag by tagline", func() {
			tagline := "#### tag1\n tagline is too ling"
			Expect(GetTagByTagline(tagline)).Should(Equal("#### tag1"))
		})

		taglines, err := GetTagAll()
		It("get tagline all", func() {
			Expect(len(taglines)).NotTo(BeZero())
			Expect(err).Should(BeNil())
		})
	})

	Context("Test Write File", func() {
		file := "../data_source/synonyms.txt"
		text := "love is text"

		It("write first line", func() {
			Expect(WriteToFile(file, text)).Should(BeNil())
			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(text))
		})

		It("write to specific line", func() {
			keyword := "love"
			text = "text"
			expected := "love is text,text"

			err := AppendToDictionary(file, keyword, text)

			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(expected))
		})

		It("write to new keyword", func() {
			keyword := "newkeyword"
			text = "text"
			expected := "love is text,text\nnewkeyword,text"

			err := AppendToDictionary(file, keyword, text)

			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(expected))
		})

		It("write to specific new line", func() {
			keyword := "newkeyword"
			text = "new_text"
			expected := "love is text,text\nnewkeyword,text,new_text"

			err := AppendToDictionary(file, keyword, text)

			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(expected))
		})

		It("write to specific original line", func() {
			keyword := "love"
			text = "new_text"
			expected := "love is text,text,new_text\nnewkeyword,text,new_text"

			err := AppendToDictionary(file, keyword, text)

			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(expected))
		})

		It("write to specific multiple line", func() {
			keyword := "newkeyword2"
			text = "new_text"
			expected := "love is text,text\nnewkeyword,text,new_text\nnewkeyword2,new_text"

			err := AppendToDictionary(file, keyword, text)

			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(expected))
		})
	})

	Context("Test Get Snippet", func() {
		It("get tagline all", func() {
			taglines, err := GetSnippet("../snippet", ".yml")
			Expect(len(taglines)).NotTo(BeZero())
			Expect(err).Should(BeNil())
		})
	})
})

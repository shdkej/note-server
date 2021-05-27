package parsing

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	Context("Test Get Snippet", func() {
		It("get tagline all", func() {
			taglines, err := GetSnippet("../snippet", ".yml")
			Expect(len(taglines)).NotTo(BeZero())
			Expect(err).Should(BeNil())
		})
	})
})

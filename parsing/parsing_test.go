package parsing

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
})

package data_source

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	parsing "github.com/shdkej/note-server/parsing"
	"testing"
)

func TestLocal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Load from Data source Test")
}

var _ = Describe("Test Load Data", func() {
	Context("Test Initial", func() {
		taglines, err := parsing.GetTagAll()
		It("get tagline all", func() {
			Expect(len(taglines)).NotTo(BeZero())
			Expect(err).Should(BeNil())
		})
		It("convert map to Tag", func() {
			tags, _ := ListToTag(taglines)
			isTag := func(t interface{}) bool {
				switch t.(type) {
				case Tag:
					return true
				default:
					return false
				}
			}(tags[0])
			Expect(tags).NotTo(BeNil())
			Expect(len(tags)).NotTo(BeZero())
			Expect(isTag).Should(BeTrue())
		})
	})

	Context("Test with Redis", func() {
		db := Redis{}
		err := db.Init()
		It("Test initial", func() {
			Expect(err).Should(BeNil())
		})
		It("Test is exist initial content", func() {
			value, err := db.GetStruct("#### kubernetes")
			Expect(value).NotTo(BeNil())
			Expect(err).Should(BeNil())
		})
	})
})

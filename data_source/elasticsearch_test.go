package data_source

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Elasticsearch", func() {
	Context("Test CRUD", func() {
		c := Elastic{}
		tag := Note{
			FileName: "test.md",
			Tag:      "## test",
			TagLine:  "this is file content",
		}
		It("Init", func() {
			Expect(c.Init()).Should(BeNil())
		})
		It("Put", func() {
			Expect(c.SetStruct(tag)).Should(BeNil())
		})
		It("Get", func() {
			Expect(c.GetStruct(tag.Tag)).Should(BeNil())
		})
	})
})

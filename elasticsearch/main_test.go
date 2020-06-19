package elasticsearch

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	db "github.com/shdkej/note-server/data_source"
)

var _ = Describe("Running Elasticsearch", func() {
	Context("Test CRUD", func() {
		c := Elastic{}
		tag := db.Tag{
			FileName:    "test.md",
			FileContent: "## test\n this is file content",
			Tag:         "## test",
			TagLine:     "this is file content",
		}
		It("Init", func() {
			Expect(c.Init()).Should(BeNil())
		})
		It("Put", func() {
			Expect(c.Put(tag)).Should(BeNil())
		})
		It("Get", func() {
			Expect(c.Get()).Should(BeNil())
		})
		It("Update", func() {
			Expect(c.Update()).Should(BeNil())
		})
	})
})

package data_source

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Dynamodb", func() {
	Context("Test CRUD", func() {
		conn := Dynamodb{}
		tableName := "myBlog"
		tag := Tag{
			FileName: "main.md",
			Tag:      "Good",
			TagLine:  "Good Enough",
		}
		It("Init", func() {
			Expect(conn.Init()).Should(BeNil())
		})
		It("Get Table", func() {
			Expect(conn.getTable()).Should(BeNil())
			Expect(conn.TableName).Should(Equal(tableName))
		})
		It("Create Item", func() {
			Expect(conn.SetStruct(tag)).Should(BeNil())
		})
		It("Get Item", func() {
			Expect(conn.GetStruct(tag.Tag)).Should(Equal(tag))
		})
		It("Scan Item", func() {
			Expect(conn.Scan(tag.Tag)).Should(Equal(tag))
		})
		It("Delete Item", func() {
			Expect(conn.Delete(tag)).Should(BeNil())
		})
	})
})

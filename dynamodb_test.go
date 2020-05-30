package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Dynamodb", func() {
	Context("Test CRUD", func() {
		conn := Dynamodb{}
		tableName := "myBlog"
		tag := Tag{
			FileName:    "main.md",
			FileContent: "0",
			Tag:         "Good",
			TagLine:     "Good Enough",
		}
		It("Init", func() {
			Expect(conn.Init()).Should(BeNil())
		})
		It("Get Table", func() {
			Expect(conn.getTable()).Should(Equal(tableName))
		})
		It("Create Item", func() {
			Expect(conn.put(tableName, tag)).Should(BeNil())
		})
		It("Get Item", func() {
			Expect(conn.get(tableName, tag.Tag)).Should(Equal(tag))
		})
		It("Delete Item", func() {
			Expect(conn.deleteItem(tableName, tag)).Should(BeNil())
		})
		//It("Load Json", func() {
		//	Expect(conn.loadData(tableName, "moviedata.json")).Should(BeNil())
		//})
	})
})

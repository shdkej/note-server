package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Dynamodb", func() {
	Context("test ping", func() {
		conn := Dynamodb{}
		It("pong", func() {
			Expect(conn.initDB()).Should(BeNil())
		})
	})

	Context("Test CRUD", func() {
		conn := Dynamodb{}
		tableName := "myBlog"
		item := Item{
			Year:   2013,
			Title:  "The Monster",
			Plot:   "BONG",
			Rating: 5.0,
		}
		It("Init", func() {
			Expect(conn.initDB()).Should(BeNil())
		})
		It("Get Table", func() {
			Expect(conn.getTable()).Should(Equal(tableName))
		})
		It("Create Item", func() {
			Expect(conn.putItem(tableName, item)).Should(BeNil())
		})
		It("Get Item", func() {
			Expect(conn.getItem(tableName, item.Title, "5.0")).Should(Equal(item))
		})
		It("Delete Item", func() {
			Expect(conn.deleteItem(tableName, item)).Should(BeNil())
		})
		//It("Load Json", func() {
		//	Expect(conn.loadData(tableName, "moviedata.json")).Should(BeNil())
		//})
	})
})

package main

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Get AWS", func() {
	Context("Test get telegram", func() {
		It("pong", func() {
			//Expect(getTelegram()).Should(BeNil())
		})
	})
	Context("Test call SQS", func() {
		It("get queue", func() {
			//Expect(sendSqs()).Should(BeNil())
		})
	})
})

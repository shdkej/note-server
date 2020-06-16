package messenger

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"io/ioutil"
	"testing"
)

func TestAWS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parsing AWS")
}

var _ = Describe("Get AWS", func() {
	/*
		filepath := "../recommend.txt"
		Context("Test call SQS", func() {
			message, err := ioutil.ReadFile(filepath)
			It("get queue", func() {
				Expect(sendSqs(string(message))).Should(BeNil())
				Expect(err).Should(BeNil())
			})
		})
		Context("Test call SNS", func() {
			It("send sns", func() {
				Expect(sendSNS("3")).Should(BeNil())
			})
		})
	*/
})

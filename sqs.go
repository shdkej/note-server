package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func sendSqs() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("MyQueue"),
	})
	if err != nil {
		return err
	}
	qURL := *result.QueueUrl
	fmt.Println("Success", qURL)

	sendMessage, err := svc.SendMessage(&sqs.SendMessageInput{
		//DelaySeconds: aws.Int64(10),
		//MessageAttributes: map[string]*sqs.MessageAttributeValue{
		//	"Message": &sqs.MessageAttributeValue{
		//		DataType:    aws.String("String"),
		//		StringValue: aws.String("gogogogogogogo"),
		//	},
		//},
		MessageBody: aws.String("Information about current NY Times"),
		QueueUrl:    &qURL,
	})
	if err != nil {
		return err
	}

	fmt.Println("Send Success", *sendMessage.MessageId)

	return nil
}

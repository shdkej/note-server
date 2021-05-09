package messenger

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

func getSQSMessage() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	queue, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("MyQueue"),
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	qURL := *queue.QueueUrl

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(20),
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return
	}

	resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &qURL,
		ReceiptHandle: result.Messages[0].ReceiptHandle,
	})

	if err != nil {
		fmt.Println("Delete Error", err)
	}

	fmt.Println("Message Deleted", resultDelete)
}

func sendSqs(message string) error {
	if message == "" {
		message = "message is empty"
	}
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
	log.Println("Success", qURL)

	sendMessage, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    &qURL,
	})
	if err != nil {
		return err
	}

	log.Println("Send Success", *sendMessage.MessageId)

	return nil
}

func sendSNS(message string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	svc := sns.New(sess)

	topic, err := svc.ListTopics(nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var topicArn string
	for _, t := range topic.Topics {
		topicArn = *t.TopicArn
	}

	result, err := svc.Publish(&sns.PublishInput{
		Message: aws.String(message),
		//		TopicArn: aws.String("arn:aws:sns:eu-central-1:917213086376:sns-sqs-upload-topic"),
		TopicArn: aws.String(topicArn),
	})

	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("SNS Send Success", *result.MessageId)
	return nil
}

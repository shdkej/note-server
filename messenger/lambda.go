package messenger

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"net/http"
	"os"
)

func getTelegram() error {
	url := "https://j1pzc4zmw9.execute-api.eu-central-1.amazonaws.com/dev/send-telegram"
	var data = []byte(`{"message": {"text":"testda!"}}`)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func uploadS3(filepath string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		log.Fatal(err)
		return err
	}
	svc := s3manager.NewUploader(sess)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	_, err = svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String("my-note-0.0.1"),
		Key:    aws.String("note/tags.csv"),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
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

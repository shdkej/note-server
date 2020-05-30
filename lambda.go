package main

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

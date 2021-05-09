package grpcserver

import (
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
)

func AddHandler(body string) string {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := NewTagManagerClient(conn)
	message := Message{
		Body: body,
	}

	stream, err := c.GetTags(context.Background(), &message)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan int)

	go func() {
		for {
			tags, err := stream.Recv()
			if err == io.EOF {
				done <- 1
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Println(tags)
		}
	}()

	<-done
	return "finish"
}

func GetFromGRPC(body string) string {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := NewTagManagerClient(conn)
	message := Message{
		Body: body,
	}

	response, err := c.GetTag(context.Background(), &message)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("response")
	log.Println(response)
	return response.String()
}

package server

import (
	"context"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type tagManagerServer struct {
}

type GRPCServer struct {
}

func (server GRPCServer) RunServer() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	var opts []grpc.ServerOption
	s := tagManagerServer{}

	grpcServer := grpc.NewServer(opts...)
	RegisterTagManagerServer(grpcServer, &s)

	log.Println("Listening gRPC server port 9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
}

func (server GRPCServer) HealthCheck() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := NewTagManagerClient(conn)

	message := Message{
		Body: "Hello This is Client!",
	}

	response, err := c.GetTag(context.Background(), &message)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(response)
}

func gRPCStreamHandler() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := NewTagManagerClient(conn)
	message := Message{
		Body: "Hello This is Client!",
	}

	stream, err := c.GetTags(context.Background(), &message)
	if err != nil {
		log.Fatal(err)
	}
	for {
		tags, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tags)
	}
}

func (s *tagManagerServer) GetTag(ctx context.Context, message *Message) (*Note, error) {
	return &Note{Tag: "## It is Tag"}, nil
}

func (s *tagManagerServer) GetTags(message *Message, stream TagManager_GetTagsServer) error {
	list := getTag("##")
	for _, tags := range list {
		tag := &Note{Tag: tags}
		if err := stream.Send(tag); err != nil {
			return err
		}
	}
	return nil
}

package data_source

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type tagManagerServer struct {
}

func (s *tagManagerServer) GetTag(ctx context.Context, message *Message) (*Note, error) {
	return &Note{Tag: "## It is Tag"}, nil
}

func (s *tagManagerServer) GetTags(message *Message, stream TagManager_GetTagsServer) error {
	list, err := getTagList()
	if err != nil {
		log.Fatal(err)
	}
	for _, tags := range list {
		tag := &Note{Tag: tags}
		if err := stream.Send(tag); err != nil {
			return err
		}
	}
	return nil
}

func RungRPC() {
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

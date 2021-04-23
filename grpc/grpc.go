//go:generate protoc --go_out=. --go_opt=paths=source_relative tag.proto
package grpcserver

import (
	"context"
	db "github.com/shdkej/note-server/data_source"
	"google.golang.org/grpc"
	"log"
	"net"
)

type tagManagerServer struct {
	datasource db.DB
	UnimplementedTagManagerServer
}

func ListenGRPC(datasource db.DB, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption
	s := tagManagerServer{datasource: datasource}

	grpcServer := grpc.NewServer(opts...)
	RegisterTagManagerServer(grpcServer, &s)

	log.Println("Listening gRPC server port ", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %s: %v", port, err)
	}
}

func (s *tagManagerServer) GetTag(ctx context.Context, message *Message) (*GetTagResponse, error) {
	return &GetTagResponse{Note: &Note{Tag: "## It is Tag"}}, nil
}

func (s *tagManagerServer) GetTags(message *Message, stream TagManager_GetTagsServer) error {
	list := [3]string{"1", "2", "3"}
	for _, tags := range list {
		tag := &Note{Tag: tags}
		if err := stream.Send(tag); err != nil {
			return err
		}
	}
	return nil
}

func (s *tagManagerServer) PutTag(ctx context.Context, note *Note) (*Message, error) {
	n := db.Note{
		Tag:      note.Tag,
		FileName: note.Filename,
		TagLine:  note.Tagline,
	}
	err := s.datasource.PutTag(n)
	if err != nil {
		log.Fatal(err)
	}
	return &Message{Body: "put is done"}, nil
}

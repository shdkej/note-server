package server

import (
	"log"
	"strings"

	db "github.com/shdkej/note-server/data_source"
)

type Server struct {
	Handler    Protocol
	Datasource *db.DB
}

type Protocol interface {
	Init()
	RunServer()
	AddHandler(string, func() string)
}

func (s *Server) RunServer() {
	s.Handler.Init()
	s.Handler.AddHandler("/", s.HealthCheck)
	s.Handler.AddHandler("/test", s.HTTP2)
	s.Handler.AddHandler("/health", s.HealthCheck)
	s.Handler.AddHandler("/tag", s.GetTag)
	s.Handler.RunServer()
}

func (s *Server) SetProtocol(p Protocol) {
	s.Handler = p
}

func (s *Server) SetDatasource(ds *db.DB) {
	s.Datasource = ds
}

func (s *Server) HTTP2() string {
	return "http2"
}

func (s *Server) GetTag() string {
	tags, err := s.Datasource.GetAllKey("####")
	if err != nil {
		log.Fatal(err)
	}
	var taglines string
	taglines += "<h1>Tags</h1>"

	for _, tag := range tags {
		taglines += "<a href='/tag/" +
			strings.Trim(tag, "# ") + "'>" +
			tag + "</a><br/><p>" + tag + "</p>"
	}
	if err != nil {
		log.Fatal(err)
	}

	return taglines
}

func (s *Server) HealthCheck() string {
	return `{"alive": true}`
}

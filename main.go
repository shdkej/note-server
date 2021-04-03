package main

import (
	"flag"

	db "github.com/shdkej/note-server/data_source"
	server "github.com/shdkej/note-server/server"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", "./app", "directory to serve")
)

func main() {
	flag.Parse()

	// Redis, Dynamodb, File
	c := &db.Redis{}
	httpserver := server.HTTPServer{}
	grpcserver := server.GRPCServer{}

	c.Init()
	go httpserver.RunServer(*listen)
	go grpcserver.RunServer()
}

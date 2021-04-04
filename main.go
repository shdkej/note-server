package main

import (
	"flag"
	"log"

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
	data := db.DB{Store: c}
	c.Init()
	log.Println(data)

	httpserver := &server.HTTPServer{}
	srv := server.Server{Handler: httpserver, Datasource: &data}
	srv.RunServer()
	/*
		grpcserver := &server.GRPCServer{}

		grpcserver.RunServer()
	*/
}

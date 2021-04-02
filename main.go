package main

import (
	"flag"
	db "github.com/shdkej/note-server/data_source"
	"log"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", "./app", "directory to serve")
)

type DataServer struct {
	Source db.DataSource
}

func main() {
	flag.Parse()

	// Redis, Dynamodb, File
	c := &db.Redis{}

	data := DataServer{c}
	err := data.Source.Init()
	if err != nil {
		log.Fatal(err)
	}

	go db.RungRPC()
	go data.httpServer()
}

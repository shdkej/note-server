package main

import (
	"flag"
	"log"

	db "github.com/shdkej/note-server/data_source"
	"github.com/shdkej/note-server/parsing"
	server "github.com/shdkej/note-server/server"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
)

func main() {
	flag.Parse()

	// Redis, Dynamodb, File
	/*
		c2 := &db.Dynamodb{}

		c2.Init()

		grpcserver := &server.GRPCServer{}
		gsrv := server.Server{Handler: grpcserver, Datasource: &data}
		go gsrv.RunServer()
	*/

	c := &db.Redis{}
	c.Init()
	data := db.DB{Store: c}

	s := server.NewServer()

	s.HandleFunc("GET", "/tag", func(c *server.Context) {
		tags, err := parsing.GetTagAll()
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(tags)
	})

	s.HandleFunc("GET", "/tag/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		t, err := data.GetTag(parameter)
		if err != nil {
			log.Println(err)
		}
		data.Hits(parameter)
		c.RenderJson(t)
	})

	s.HandleFunc("GET", "/hits/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		t := data.Hits(parameter)
		c.RenderJson(t)
	})

	s.Run(*listen)
}

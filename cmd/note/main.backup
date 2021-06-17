package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	db "github.com/shdkej/database"
	grpcserver "github.com/shdkej/note-server/grpc"
	server "github.com/shdkej/note-server/server"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
)

func main() {
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Redis, Dynamodb, File
	/*
		c2 := &db.Dynamodb{}

		c2.Init()

		grpcserver := &server.GRPCServer{}
		gsrv := server.Server{Handler: grpcserver, Datasource: &data}
		go gsrv.RunServer()
	*/

	c := &db.Redis{}
	data := db.DB{Store: c}
	data.Init()

	go grpcserver.ListenGRPC(data, ":9000")

	s := server.NewServer()

	s.HandleFunc("GET", "/", func(c *server.Context) {
		c.RenderJson("{'health':'ok'}")
	})

	s.HandleFunc("GET", "/tag", func(c *server.Context) {
		tags, err := data.GetTags()
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

	s.HandleFunc("POST", "/tag/:tag", func(c *server.Context) {
		// 파라미터를 받고
		// 메인로직을 부르고
		// 응답을 보낸다
		parameter := c.Params["tag"].(string)
		t := data.PutTag(db.Note{Tag: parameter})
		c.RenderJson(t)
	})

	s.HandleFunc("POST", "/tag/:tag/:value", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		value := c.Params["value"].(string)
		t := data.PutTagForSearch(parameter, value)
		c.RenderJson(t)
	})

	s.HandleFunc("GET", "/hits/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		t := data.Hits(parameter)
		c.RenderJson(t)
	})

	s.HandleFunc("GET", "/dict", func(c *server.Context) {
		tags, err := data.GetEverything("dict")
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(tags)
	})

	s.HandleFunc("POST", "/dict", func(c *server.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
		}
		note := db.Note{}
		err = json.Unmarshal(body, &note)
		t := data.Put("dict", note)
		c.RenderJson(t)
	})

	s.HandleFunc("GET", "/grpc", func(c *server.Context) {
		message := grpcserver.GetFromGRPC("test")
		log.Println("grpc")
		c.RenderJson(message)
	})

	s.HandleFunc("GET", "/grpcstream", func(c *server.Context) {
		finish := grpcserver.AddHandler("test")
		c.RenderJson(finish)
	})

	s.Run(*listen)
}

package main

import (
	"flag"
	"io/ioutil"
	"log"

	db "github.com/shdkej/note-server/data_source"
	grpcserver "github.com/shdkej/note-server/grpc"
	server "github.com/shdkej/note-server/server"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
)

func main() {
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	c := &db.Redis{}
	data := db.DB{Store: c}
	data.Init()

	go grpcserver.ListenGRPC(data, ":9000")

	s := server.NewServer()

	// health check endpoint
	s.HandleFunc("GET", "/health", func(c *server.Context) {
		c.RenderJson("{'health':'ok'}")
	})

	// it can make more various data acceptable
	s.HandleFunc("GET", "/set/:table", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		data.SetPrefix(parameter)
		c.RenderJson("{'table change':'ok', 'table':'" + parameter + "'}")
	})

	// get all list of data
	s.HandleFunc("GET", "/", func(c *server.Context) {
		tags, err := data.GetEverything()
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(tags)
	})

	// get one item
	s.HandleFunc("GET", "/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		t, err := data.Get(parameter)
		if err != nil {
			log.Println(err)
		}
		data.Hits(parameter)
		c.RenderJson(t)
	})

	// create one item
	s.HandleFunc("POST", "/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		t := data.Put(db.Note{Tag: parameter})
		c.RenderJson(t)
	})

	// update one item
	s.HandleFunc("PUT", "/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		body := c.Request.Body
		value, err := ioutil.ReadAll(body)
		if err != nil {
			log.Println(err)
		}

		err = data.Update(parameter, string(value))
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(err)
	})

	// delete one item
	s.HandleFunc("DELETE", "/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)

		err := data.Delete(parameter)
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(err)
	})

	s.Run(*listen)
}

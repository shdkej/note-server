package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	elastic "github.com/shdkej/note-server/elastic"
	parsing "github.com/shdkej/note-server/parsing"
	server "github.com/shdkej/note-server/server"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
)

func main() {
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	e := elastic.Elastic{}
	e.Init()

	s := server.NewServer()

	s.HandleFunc("GET", "/health", func(c *server.Context) {
		c.RenderJson("{'health':'ok'}")
	})

	s.HandleFunc("GET", "/", func(c *server.Context) {
		tags, err := e.GetAll()
		if err != nil {
			log.Println("get error", err)
		}
		c.RenderJson(tags)
	})

	s.HandleFunc("GET", "/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		t, err := e.GetSynonym(parameter)
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(t)
	})

	s.HandleFunc("POST", "/", func(c *server.Context) {
		parameter := c.Request.Body
		body, err := ioutil.ReadAll(parameter)
		note := elastic.Tag{}

		err = json.Unmarshal(body, &note)
		if err != nil {
			log.Println(err)
		}

		log.Println("Insert data", note.Name)
		err = e.SetStruct(note)
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(err)
	})

	s.HandleFunc("PUT", "/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		body := c.Request.Body
		value, err := ioutil.ReadAll(body)
		if err != nil {
			log.Println(err)
		}

		err = e.Update(parameter, string(value))
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(err)
	})

	s.HandleFunc("DELETE", "/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)

		err := e.Delete(parameter)
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(err)
	})

	// Append new keyword to Dictionary
	s.HandleFunc("POST", "/:tag", func(c *server.Context) {
		keyword := c.Params["tag"].(string)
		parameter := c.Request.Body
		body, err := ioutil.ReadAll(parameter)
		if err != nil {
			log.Println(err)
		}

		data := string(body)

		file := "elastic/synonyms.txt"
		err = parsing.AppendToDictionary(file, keyword, data)
		if err != nil {
			log.Println(err)
		}

		err = e.Update(keyword, data)
		if err != nil {
			log.Println(err)
		}

		c.RenderJson(err)
	})

	s.Run(*listen)
}

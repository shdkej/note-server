package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	elastic "github.com/shdkej/note-server/data_source"
	parsing "github.com/shdkej/note-server/parsing"
	server "github.com/shdkej/note-server/server"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
)

func main() {
	flag.Parse()

	e := elastic.Elastic{}
	e.Init()

	s := server.NewServer()

	s.HandleFunc("GET", "/", func(c *server.Context) {
		c.RenderJson("{'health':'ok'}")
	})

	s.HandleFunc("GET", "/search", func(c *server.Context) {
		tags, err := e.GetAll()
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(tags)
	})

	s.HandleFunc("GET", "/search/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)
		t, err := e.GetSynonym(parameter)
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(t)
	})

	s.HandleFunc("POST", "/search", func(c *server.Context) {
		parameter := c.Request.Body
		body, err := ioutil.ReadAll(parameter)
		note := elastic.Note{}

		err = json.Unmarshal(body, &note)
		if err != nil {
			log.Println(err)
		}

		log.Println("Insert data", note.Tag)
		err = e.SetStruct(note)
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(err)
	})

	s.HandleFunc("PUT", "/search/:tag", func(c *server.Context) {
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

	s.HandleFunc("DELETE", "/search/:tag", func(c *server.Context) {
		parameter := c.Params["tag"].(string)

		err := e.Delete(parameter)
		if err != nil {
			log.Println(err)
		}
		c.RenderJson(err)
	})

	// Append new keyword to Dictionary
	s.HandleFunc("POST", "/search/:tag", func(c *server.Context) {
		keyword := c.Params["tag"].(string)
		parameter := c.Request.Body
		body, err := ioutil.ReadAll(parameter)
		if err != nil {
			log.Println(err)
		}

		data := string(body)

		file := "data_source/synonyms.txt"
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

func decisionIndex(e elastic.Elastic, data interface{}) error {
	index := "analyze"

	switch data.(type) {
	case string:
		index = "analyze"
	case elastic.Note:
		index = "note"
	}

	e.ChangeIndex(index)
	return nil
}

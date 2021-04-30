package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	datasource "github.com/shdkej/note-server/data_source"
	parsing "github.com/shdkej/note-server/parsing"
)

var (
	source      = flag.String("src", "/home/sh/workspace/note-reminder", "file path for parsing")
	destination = flag.String("dest", "http://localhost:8081/search", "database endpoint to push data")
	db          = flag.String("db", "redis", "database to push data")
	dataType    = flag.String("t", "Note", "data type for json")
)

func main() {
	flag.Parse()
	data1, err := parsing.GetTagAll()
	if err != nil {
		log.Fatalf("Error when parsing tags %s", err)
	}
	d1, _ := datasource.ListToNote(data1)
	d2 := getSnippet()

	d := append(d1, d2...)
	runningAsyncronizly(d)
}

func runningAsyncronizly(d []datasource.Note) {
	ch := make(chan bool, len(d))

	for _, s := range d {
		go SendPostRequest(s, ch)
	}
	<-ch
}

func getSnippet() []datasource.Note {
	data, err := parsing.GetSnippet(*source, ".yml")
	if err != nil {
		log.Fatalf("Error when parsing snippets %s", err)
	}
	d, err := datasource.ListToNote(data)
	if err != nil {
		log.Fatalf("Error when Convert to interface %s", err)
	}
	return d
}

func SendPostRequest(s datasource.Note, ch chan bool) {
	j, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("Error when Marshal to json %s", err)
	}
	params := bytes.NewBuffer(j)
	_, err = http.Post(*destination, "application/json", params)
	if err != nil {
		log.Fatalf("Error when post api %s", err)
	}
	log.Println("done")
	ch <- true
}

func ConfirmComplete() int {
	resp, err := http.Get(*destination)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	note := []datasource.Note{}
	json.Unmarshal(data, &note)
	return len(note)
}

func CleanUp(values []datasource.Note) {
	for _, value := range values {
		client := &http.Client{}
		url := *destination + "/" + value.Tag
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("cleaning")
}

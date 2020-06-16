package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
	db "github.com/shdkej/go-wasm/data_source"
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
	c.Init()

	data := DataServer{c}

	err := c.Ping()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/tag", data.TagHandler)
	r.HandleFunc("/tag/{key}", data.TagOneHandler)
	r.HandleFunc("/random", data.TagHandler)
	r.HandleFunc("/Initial", data.InitialHandler)

	r.HandleFunc("/health", HealthCheckHandler)
	r.HandleFunc("/test", HTTP2TestHandler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("/src/app"))))

	log.Printf("listening on %q...", *listen)
	srv := &http.Server{Addr: *listen, Handler: r}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
	//log.Fatal(srv.ListenAndServe())
}

func (d DataServer) InitialHandler(w http.ResponseWriter, r *http.Request) {
	d.Source.SetInitial()
}

func (d DataServer) TagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hits, err := d.Source.Hits(vars["key"])
	if err != nil {
		log.Fatal(err)
	}
	tags, err := d.Source.GetAllKey("test")
	if err != nil {
		log.Fatal(err)
	}
	var taglines string
	taglines += "<h1>Random</h1>"

	for _, tag := range tags {
		taglines += "<a href='/tag/" + strings.Trim(tag, "# ") + "'>" + tag + "</a><br/><p>" + tag + "</p>"
	}
	if err != nil {
		log.Fatal(err)
	}
	err = RenderOutput(w, taglines)
	if err != nil {
		log.Fatal(err)
	}
	val, err := ioutil.ReadFile("app/recommend.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Hits: %v\n <h1>Content Based</h1>%v", hits, string(val))
}

func (d DataServer) TagOneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hits, err := d.Source.Hits(vars["key"])
	if err != nil {
		log.Fatal(err)
	}

	content, err := d.Source.GetTagParagraph(vars["key"])
	if err != nil {
		log.Fatal(err)
	}

	err = RenderArrayOutput(w, content)
	fmt.Fprintf(w, "Hits: %v\n", hits)
}

func (d DataServer) getRandomContent() []string {
	tags, err := d.Source.GetSet("## markdown")
	if err != nil {
		log.Fatal(err)
	}

	//random := getRandom(len(tags))

	var randomtags []string
	for _, v := range tags {
		randomtags = append(randomtags, string(v))
	}
	return randomtags
}

type Book struct {
	Title   string
	Content template.HTML
}

func RenderOutput(w http.ResponseWriter, values string) error {
	path, err := os.Getwd()
	fp := path + "/app/index.html"
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}

	output := Book{"article", template.HTML(values)}
	if err := tmpl.Execute(w, output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func RenderArrayOutput(w http.ResponseWriter, values []string) error {
	path, err := os.Getwd()
	fp := path + "/app/index.html"
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatal(err)
	}

	var html string
	for _, value := range values {
		output := blackfriday.Run([]byte(value))
		html += "<div style='padding:10px;'>" + string(output) + "</div>"
	}

	output := Book{"article", template.HTML(html)}
	if err := tmpl.Execute(w, output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func HTTP2TestHandler(w http.ResponseWriter, r *http.Request) {
	if pusher, ok := w.(http.Pusher); ok {
		options := &http.PushOptions{
			Header: http.Header{
				"Accept-Encoding": r.Header["Accept-Encoding"],
			},
		}
		if err := pusher.Push("/styles.css", options); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	w.Write([]byte("Hello"))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}

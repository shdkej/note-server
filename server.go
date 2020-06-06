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
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", "./app", "directory to serve")
)

var wikiDir = os.Getenv("VIMWIKI")

func main() {
	flag.Parse()

	// Redis, Dynamodb, File
	c := &Client{}
	c.Init()

	//
	db := Database{c}

	err := c.ping()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/category", db.CategoryHandler)
	r.HandleFunc("/category/{key}", db.ArticleHandler)
	r.HandleFunc("/tag", db.TagHandler)
	r.HandleFunc("/tag/{key}", db.TagOneHandler)
	r.HandleFunc("/random", db.TagHandler)

	r.HandleFunc("/health", HealthCheckHandler)
	r.HandleFunc("/test", HTTP2TestHandler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("/src/app"))))

	log.Printf("listening on %q...", *listen)
	srv := &http.Server{Addr: *listen, Handler: r}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
	//log.Fatal(srv.ListenAndServe())
}

func (c Database) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	files, err := getFileAll()
	var filename string
	for _, file := range files {
		filename += "<a href='/data/" + file.Name() + "'>" + file.Name() + "</a><br/> "
	}
	err = RenderOutput(w, filename)
	c.source.get("test")
	if err != nil {
		log.Fatal(err)
	}
}

func (c Database) TagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hits, err := c.source.hits(vars["key"])
	if err != nil {
		log.Fatal(err)
	}

	tags := c.getRandomContent()
	var taglines string
	taglines += "<h1>Random</h1>"

	for _, tag := range tags {
		taglines += "<a href='/tag/" + strings.Trim(getTag(tag), "# ") + "'>" + getTag(tag) + "</a><br/><p>" + tag + "</p>"
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

func (c Database) TagOneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hits, err := c.source.hits(vars["key"])
	if err != nil {
		log.Fatal(err)
	}

	content, err := c.source.getTagParagraph(vars["key"])
	if err != nil {
		log.Fatal(err)
	}

	err = RenderArrayOutput(w, content)
	fmt.Fprintf(w, "Hits: %v\n", hits)
}

func (c Database) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	filename := wikiDir + vars["key"] + ".md"
	val, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	output := blackfriday.Run(val)
	err = RenderOutput(w, string(output))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Key: %v\n", vars["key"])
}

func (c Database) getRandomContent() []string {
	tags, err := c.source.getSet("## markdown")
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

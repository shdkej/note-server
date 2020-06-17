package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
	db "github.com/shdkej/go-wasm/data_source"
	"google.golang.org/grpc"
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

	r := mux.NewRouter()

	r.HandleFunc("/tag", data.TagHandler)
	r.HandleFunc("/tag/{key}", data.TagOneHandler)
	r.HandleFunc("/random", data.RandomHandler)
	r.HandleFunc("/Initial", data.InitialHandler)

	r.HandleFunc("/grpc", gRPCHandler)
	r.HandleFunc("/stream", gRPCStreamHandler)
	r.HandleFunc("/health", HealthCheckHandler)
	r.HandleFunc("/test", HTTP2TestHandler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("/src/app"))))

	log.Printf("listening on %q...", *listen)
	srv := &http.Server{Addr: *listen, Handler: r}
	//log.Fatal(srv.ListenAndServe())
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

func (d DataServer) InitialHandler(w http.ResponseWriter, r *http.Request) {
	d.Source.PutTags()
}

func (d DataServer) RandomHandler(w http.ResponseWriter, r *http.Request) {
	random := d.getRandomContent()
	fmt.Fprintf(w, "<h1>Random</h1>%v", random)
}

func (d DataServer) TagHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := d.Source.GetAllKey("####")
	if err != nil {
		log.Fatal(err)
	}
	var taglines string
	taglines += "<h1>Tags</h1>"

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

func gRPCHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := db.NewTagManagerClient(conn)

	message := db.Message{
		Body: "Hello This is Client!",
	}

	response, err := c.GetTag(context.Background(), &message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Response: %v", response)
}

func gRPCStreamHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := db.NewTagManagerClient(conn)
	message := db.Message{
		Body: "Hello This is Client!",
	}

	stream, err := c.GetTags(context.Background(), &message)
	if err != nil {
		log.Fatal(err)
	}
	for {
		tags, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Response: %v", tags)
	}
}

func (d DataServer) getRandomContent() db.Tag {
	tags, err := d.Source.GetStruct("##")
	if err != nil {
		log.Fatal(err)
	}
	return tags
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

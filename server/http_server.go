package server

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

type HTTPServer struct {
}

func (server HTTPServer) RunServer(listen string) {
	r := mux.NewRouter()

	r.HandleFunc("/tag", func(w http.ResponseWriter, r *http.Request) {
		taglines := getTags()
		err := RenderOutput(w, taglines)
		if err != nil {
			log.Println(err)
		}
	})

	r.HandleFunc("/tag/{key}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Fprintf(w, "Hits: %v\n", vars)
		taglines := getTag(vars["key"])
		err := RenderArrayOutput(w, taglines)
		if err != nil {
			log.Println(err)
		}
	})

	r.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		random := RandomHandler()
		fmt.Fprintf(w, "<h1>Random</h1>%v", random)
	})

	r.HandleFunc("/health", server.HealthCheck)
	r.HandleFunc("/test", HTTP2TestHandler)

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("/src/app"))))

	log.Printf("listening on %q...", listen)
	srv := &http.Server{Addr: listen, Handler: r}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

func (s HTTPServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
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

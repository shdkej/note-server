package server

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

type HTTPServer struct {
	srv *http.Server
	r   *mux.Router
}

func (s *HTTPServer) Init() {
	s.r = mux.NewRouter()
	s.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("/src/app"))))
}

func (s *HTTPServer) RunServer() {
	log.Printf("listening on %q...", ":8080")
	s.srv = &http.Server{Addr: ":8080", Handler: s.r}
	log.Fatal(s.srv.ListenAndServe())
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.srv.ListenAndServeTLS(path+"/server.crt", path+"/server.key"))
}

func (s *HTTPServer) AddHandler(url string, handler func() string) {
	s.r.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		value := handler()
		RenderOutput(w, value)
	})
}

func (s *HTTPServer) restful(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var value string
	switch r.Method {
	case "GET":
		value = "get"
	case "POST":
		value = "post"
	case "PUT":
		value = "put"
	case "DELETE":
		value = "delete"
	default:
		w.WriteHeader(http.StatusNotFound)
	}
	value = `{"message":` + value + `}`
	w.Write([]byte(value))
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

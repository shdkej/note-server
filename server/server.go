package server

import (
	"log"
	"net/http"
)

type Protocol interface {
	Init()
	RunServer()
	AddHandler(string, func() string)
}

type Server struct {
	*router
	middlewares  []Middleware
	startHandler HandlerFunc
}

func NewServer() *Server {
	r := &router{make(map[string]map[string]HandlerFunc)}
	s := &Server{router: r}
	s.middlewares = []Middleware{
		logHandler,
		recoverHandler,
	}
	return s
}
func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Run(addr string) {
	s.startHandler = s.router.handler()

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.startHandler = s.middlewares[i](s.startHandler)
	}

	log.Println("running server...", addr)
	if err := http.ListenAndServe(addr, s); err != nil {
		log.Println("server bump")
		panic(err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{
		Params:         make(map[string]interface{}),
		ResponseWriter: w,
		Request:        r,
	}

	for k, v := range r.URL.Query() {
		c.Params[k] = v[0]
	}
	c.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	c.ResponseWriter.Header().Add("Access-Control-Allow-Methods", "*")
	c.ResponseWriter.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	s.startHandler(c)
}

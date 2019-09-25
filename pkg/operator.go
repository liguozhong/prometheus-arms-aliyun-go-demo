package pkg

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
)

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Run() error {
	port := ":" + strconv.Itoa(s.port)
	path := "/metrics"
	http.Handle(path, promhttp.Handler()) //初始一个http handler
	fmt.Println("http.url: http://localhost"+port+path)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		return err
	}
	return nil
}

package pkg

import (
	"fmt"
	stat "github.com/liguozhong/prometheus-arms-aliyun-go-demo/pkg/opentelemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
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
	service := "/buy"
	http.Handle(path, promhttp.Handler()) //初始一个http handler
	http.HandleFunc(service, func(writer http.ResponseWriter, request *http.Request) {
		content, err := stat.DoBuy()
		if err != nil {
			io.WriteString(writer, err.Error())
			return
		}
		io.WriteString(writer, content)
	})
	stat.InitMeter("buy2", true)
	fmt.Println("http.url: http://localhost" + port + path)
	fmt.Println("service.url: http://localhost" + port + service)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		return err
	}
	return nil
}

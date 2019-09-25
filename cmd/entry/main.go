package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	gs "github.com/liguozhong/prometheus-arms-aliyun-go-demo/pkg"
	"os"

	_ "net/http/pprof"
)

const (
	logLevelAll   = "all"
	logLevelDebug = "debug"
	logLevelInfo  = "info"
	logLevelWarn  = "warn"
	logLevelError = "error"
	logLevelNone  = "none"
)

func main() {
	fmt.Println("开始启动 prometheus-arms-aliyun-go-demo -v0")
	logger := initLogger()
	runServer(logger)
	fmt.Println("结束启动 prometheus-arms-aliyun-go-demo module", )
	<-make(chan bool)
}

//部署在用户的k8s内，抓取exporter内的数据，写入SLS
func runServer(logger log.Logger) {
	server := gs.NewServer(8077)
	err := server.Run()
	if err != nil {
		logger.Log("server.Run() err:", err)
	}
}

func initLogger() log.Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

	logLevel := logLevelWarn
	switch logLevel {
	case logLevelAll:
		logger = level.NewFilter(logger, level.AllowAll())
	case logLevelDebug:
		logger = level.NewFilter(logger, level.AllowDebug())
	case logLevelInfo:
		logger = level.NewFilter(logger, level.AllowInfo())
	case logLevelWarn:
		logger = level.NewFilter(logger, level.AllowWarn())
	case logLevelError:
		logger = level.NewFilter(logger, level.AllowError())
	case logLevelNone:
		logger = level.NewFilter(logger, level.AllowNone())
	default:
		fmt.Fprintf(os.Stderr, "log level %v unknown, %v are possible values", logLevel, logLevelInfo)
		return nil
	}
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	logger.Log("msg", fmt.Sprintf("Starting ARMS Prometheus Operator version '%v'.", "0.31.0"))
	return logger
}

module github.com/liguozhong/prometheus-arms-aliyun-go-demo

go 1.12

require (
	github.com/go-kit/kit v0.9.0
	github.com/prometheus/client_golang v1.7.1
	go.opentelemetry.io/contrib/exporters/metric/cortex v0.15.0
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.15.0
	go.opentelemetry.io/otel/exporters/otlp v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
	golang.org/x/text v0.3.8 // indirect
)

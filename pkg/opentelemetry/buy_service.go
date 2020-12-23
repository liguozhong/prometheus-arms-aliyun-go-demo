package stat

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"time"
)

var buyCounter metric.Int64Counter

func init() {
	fmt.Println(time.Now(), " - initMetrics start......")
	meter := otel.GetMeterProvider().Meter("github.com/liguozhong/prometheus-arms-aliyun-go-demo")
	buyCounter = metric.Must(meter).NewInt64Counter(
		"buy_total",
		metric.WithDescription("Measures  buy"),
	)
}

func DoBuy() (string, error) {
	buyCounter.Add(context.Background(), 1)
	return "buy success", nil
}

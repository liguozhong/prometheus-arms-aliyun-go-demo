package stat

import (
	"context"
	prometheusPushExporter "go.opentelemetry.io/contrib/exporters/metric/cortex"
	prometheusExporter "go.opentelemetry.io/otel/exporters/metric/prometheus"

	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"net/http"
	"time"

	_ "net/http/pprof" // 为了诊断问题 shichun.fsc 20190712
)

func InitMeter(app string, push bool) error {
	fmt.Println(time.Now(), " - initMeter start......")
	if push {
		fmt.Println(time.Now(), " - initMeter opentelemetry push......")
		remoteUrl := "http://region.arms.aliyuncs.com/prometheus/../../../../api/v3/write"
		ak := "ak"
		sk := "sk"
		return initPushMeter(app, remoteUrl, ak, sk)
	}
	fmt.Println(time.Now(), " - initMeter opentelemetry pull......")
	return initPullMeter(app)
}

func initPushMeter(regionId string, remoteWriteUrl string, ak string, sk string) error {
	fmt.Println(time.Now(), " - initPushMeter start......")
	var validatedStandardConfig = prometheusPushExporter.Config{
		Endpoint:      remoteWriteUrl,
		Name:          "AliyunConfig",
		RemoteTimeout: 30 * time.Second,
		PushInterval:  10 * time.Second,
		Quantiles:     []float64{0.5, 0.9, 0.95, 0.99},
		BasicAuth: map[string]string{
			"username": ak,
			"password": sk,
		},
	}
	if validatedStandardConfig.Endpoint == "" {
		return errors.New(" validatedStandardConfig.Endpoint==empty.regionId:" + regionId)
	}
	fmt.Println("Success: Created Config struct")
	r, err := resource.New(context.Background(),
		resource.WithAttributes(
			label.String("cluster", "test-otel"),
			label.String("app", "buy")))
	if err != nil {
		fmt.Println("resource Error:", err)
	}
	pusher, err := prometheusPushExporter.InstallNewPipeline(validatedStandardConfig,
		push.WithPeriod(30*time.Second), push.WithResource(r))
	if err != nil {
		fmt.Println("InstallNewPipeline Error:", err)
	}
	otel.SetMeterProvider(pusher.MeterProvider())
	return nil
}

func initPullMeter(app string) error {
	fmt.Println(time.Now(), " - initPullMeter start......")
	r, err := resource.New(context.Background(),
		resource.WithAttributes(
			label.String("cluster", "test-otel"),
			label.String("app", app)))
	if err != nil {
		fmt.Println("resource Error:", err)
	}
	exporter, err := prometheusExporter.NewExportPipeline(
		prometheusExporter.Config{
			DefaultHistogramBoundaries: []float64{-0.5, 1},
		},
		pull.WithCachePeriod(0),
		pull.WithResource(r),
	)
	if err != nil {
		return err
	}
	http.HandleFunc("/opentelemetry", exporter.ServeHTTP)
	otel.SetMeterProvider(exporter.MeterProvider())
	return nil
}

func initOtlpProvider(regionId string) (*push.Controller, error) {
	exporter, err := otlp.NewExporter(
		context.Background(),
		otlp.WithInsecure(),
		otlp.WithAddress(regionId+"-intranet.arms.aliyuncs.com:8000"),
	)
	if err != nil {
		return nil, err
	}

	pusher := push.New(
		basic.New(
			simple.NewWithExactDistribution(),
			exporter,
		),
		exporter,
		push.WithPeriod(30*time.Second),
	)

	otel.SetMeterProvider(pusher.MeterProvider())
	pusher.Start()

	return pusher, err
}

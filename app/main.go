package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tkanos/api_trace_metrics_demo/app/config"
	"github.com/tkanos/api_trace_metrics_demo/app/logger"
	"github.com/tkanos/api_trace_metrics_demo/app/metrics"
	"github.com/tkanos/api_trace_metrics_demo/app/pixel"
	"go.uber.org/zap"
)

var tracer stdopentracing.Tracer

func init() {
	config.Init()
	logger.Init()
}

func main() {

	//Zipkin Connection
	if config.ZipkinServer != "" {
		collector, err := zipkin.NewHTTPCollector(config.ZipkinServer + "/api/v1/spans")

		if err != nil {
			logger.Debug("collector_error", zap.Error(err))
		}
		defer collector.Close()

		recorder := zipkin.NewRecorder(collector, false, "localhost", "pixel-api")
		tracer, err = zipkin.NewTracer(
			recorder,
			zipkin.ClientServerSameSpan(false),
			zipkin.TraceID128Bit(true),
		)

		if err != nil {
			logger.Debug("recorder_error", zap.Error(err))
		}

		// explicitly set our tracer to be the default tracer.
		stdopentracing.SetGlobalTracer(tracer)
	} else {
		tracer = stdopentracing.NoopTracer{}
		logger.Debug("noop tracer")

	}
	logger.Debug("tracing", zap.String("ZipkinServer", config.ZipkinServer))

	// Errors channel
	errc := make(chan error)

	// Interrupt handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP Transport
	go func() {
		httpAddr := ":" + strconv.Itoa(config.Port)
		mux := http.NewServeMux()

		mux.Handle("/pixel/", pixel.MakeHTTPHandler(pixelHandler(), tracer))

		mux.HandleFunc("/", indexHandler)
		mux.HandleFunc("/healthz", healthzHandler)
		mux.Handle("/metrics", promhttp.Handler())

		httpServer := &http.Server{
			Addr:    httpAddr,
			Handler: mux,
		}

		logger.Info(fmt.Sprintf("the Pixel API is starting on port %v", config.Port), logger.IntField("port", config.Port))
		errc <- httpServer.ListenAndServe()
	}()

	logger.Error("exit", <-errc)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to the Pixel API!")
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func pixelHandler() pixel.Endpoints {
	var pxEndpoints pixel.Endpoints
	{
		pxRepository := pixel.NewRepository()
		pxRepository = pixel.NewTracingRepository(pxRepository)

		pxService := pixel.NewService(pxRepository)
		pxService = pixel.NewpxTracing(pxService)

		getByIDEndpoint := pixel.MakeGetByIDEndpoint(pxService)
		getByIDEndpoint = opentracing.TraceServer(tracer, "pxService::GetpxByID")(getByIDEndpoint)

		pxEndpoints = pixel.Endpoints{
			GetByID: metrics.NewMetricsMiddleware("pixel", getByIDEndpoint),
		}
	}

	return pxEndpoints
}

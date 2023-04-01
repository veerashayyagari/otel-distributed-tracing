package tracer

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	otrace "go.opentelemetry.io/otel/trace"
)

type TraceConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	ExportURI      string
}

// RouterHandler represents the func type that httprouter handles on any given route match
type RouterHandler func(w http.ResponseWriter, r *http.Request, p httprouter.Params)

// NewTraceProvider sets up a global trace provider for the service and configures trace data to be published to exportURI
func NewTraceProvider(cfg *TraceConfig) (otrace.Tracer, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("creating new resource: %w", err)
	}

	exp, err := newExporter(cfg.ExportURI)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
		// sample half of the traces
		trace.WithSampler(trace.TraceIDRatioBased(0.5)),
		trace.WithBatcher(exp,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(r),
	)

	otel.SetTracerProvider(tp)
	return tp.Tracer(cfg.ServiceName), nil
}

// newExporter configures zipkin exporter when a valid export URI is passed,
// else will configure a stdouttrace exporter that writes to traces.txt file
func newExporter(exportURI string) (trace.SpanExporter, error) {
	// if the supplied export uri is invalid, write to stdouttrace
	if _, err := url.Parse(exportURI); err != nil {
		if f, e := os.Create("traces.txt"); e != nil {

		} else {
			return stdouttrace.New(
				// use traces.txt for writing out traces
				stdouttrace.WithWriter(f),

				// use human-readable output.
				stdouttrace.WithPrettyPrint(),
			)
		}

	}

	// if we have a valid export URI , export logs to the service running zipkin
	return zipkin.New(exportURI)
}

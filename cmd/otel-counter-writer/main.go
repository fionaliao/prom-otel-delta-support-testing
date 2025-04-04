package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func main() {
	// Set env variable to convert exponential histograms to Prometheus native histograms
	os.Setenv("OTEL_EXPORTER_OTLP_METRICS_DEFAULT_HISTOGRAM_AGGREGATION", "base2_exponential_bucket_histogram")

	ctx := context.Background()

	delta, err := createMeterSetup(ctx, metricdata.DeltaTemporality)
	if err != nil {
		log.Fatal(err)
	}

	cumulative, err := createMeterSetup(ctx, metricdata.CumulativeTemporality)
	if err != nil {
		log.Fatal(err)
	}

	count := int64(0)
	log.Println("Starting to send metrics. Press Ctrl+C to stop...")

	for {
		count++
		delta.counter.Add(ctx, 1, metric.WithAttributes(delta.labels...))
		cumulative.counter.Add(ctx, 1, metric.WithAttributes(cumulative.labels...))
		delta.histogram.Record(ctx, float64(count%100), metric.WithAttributes(delta.labels...))
		cumulative.histogram.Record(ctx, float64(count%100), metric.WithAttributes(cumulative.labels...))

		// Only increment sparse counters every 20 iterations
		if count%20 == 0 {
			delta.sparseCounter.Add(ctx, 1, metric.WithAttributes(delta.labels...))
			cumulative.sparseCounter.Add(ctx, 1, metric.WithAttributes(cumulative.labels...))
			delta.sparseHistogram.Record(ctx, float64(count%100), metric.WithAttributes(delta.labels...))
			cumulative.sparseHistogram.Record(ctx, float64(count%100), metric.WithAttributes(cumulative.labels...))
			log.Printf("Added sparse metric values (count: %d)", count)
		}

		log.Printf("Added metric value %d (count: %d)", 1, count)
		time.Sleep(1 * time.Second)
	}
}

type meterSetup struct {
	provider        *sdkmetric.MeterProvider
	meter           metric.Meter
	counter         metric.Int64Counter
	sparseCounter   metric.Int64Counter
	histogram       metric.Float64Histogram
	sparseHistogram metric.Float64Histogram
	labels          []attribute.KeyValue
}

func createMeterSetup(ctx context.Context, temporality metricdata.Temporality) (*meterSetup, error) {
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("localhost:9090"),
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithURLPath("/api/v1/otlp/v1/metrics"),
		otlpmetrichttp.WithTemporalitySelector(func(_ sdkmetric.InstrumentKind) metricdata.Temporality {
			return temporality
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s exporter: %v", temporality, err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(10*time.Second),
		)),
	)

	meter := provider.Meter(temporality.String() + "-metrics")

	counter, err := meter.Int64Counter("test_counter",
		metric.WithDescription("Test "+temporality.String()+" counter example"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s counter: %v", temporality.String(), err)
	}

	sparseCounter, err := meter.Int64Counter("test_sparse_counter",
		metric.WithDescription("Test sparse "+temporality.String()+" counter example"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create sparse %s counter: %v", temporality.String(), err)
	}

	histogram, err := meter.Float64Histogram("test_histogram",
		metric.WithDescription("Test "+temporality.String()+" histogram example"),
		metric.WithUnit("1"))

	if err != nil {
		return nil, fmt.Errorf("failed to create %s histogram: %v", temporality.String(), err)
	}

	sparseHistogram, err := meter.Float64Histogram("test_sparse_histogram",
		metric.WithDescription("Test sparse "+temporality.String()+" histogram example"),
		metric.WithUnit("1"))

	if err != nil {
		return nil, fmt.Errorf("failed to create sparse %s histogram: %v", temporality.String(), err)
	}

	return &meterSetup{
		provider:        provider,
		meter:           meter,
		counter:         counter,
		sparseCounter:   sparseCounter,
		histogram:       histogram,
		sparseHistogram: sparseHistogram,
		labels:          createLabels(temporality.String()),
	}, nil
}

func createLabels(temporality string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("service", "test-service"),
		attribute.String("environment", "dev"),
		attribute.String("otel_temporality", temporality),
	}
}

package nsq

import (
	"context"

	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/pkg/nsq"
	"github.com/jaegertracing/jaeger/plugin/storage/kafka"
)

type spanWriterMetrics struct {
	SpansWrittenSuccess metrics.Counter
	SpansWrittenFailure metrics.Counter
}

// SpanWriter writes spans to NSQ. Implements spanstore.Writer
type SpanWriter struct {
	metrics    *spanWriterMetrics
	producer   *nsq.NsqProducer
	marshaller kafka.Marshaller
}

// NewSpanWriter initiates and returns a new NSQ spanwriter
func NewSpanWriter(
	options Options,
	marshaller kafka.Marshaller,
	factory *metrics.Factory,
	logger *zap.Logger,
) *SpanWriter {
	metricsFactory := (*factory)
	writeMetrics := &spanWriterMetrics{
		SpansWrittenSuccess: metricsFactory.Counter(metrics.Options{Name: "nsq_spans_written", Tags: map[string]string{"status": "success"}}),
		SpansWrittenFailure: metricsFactory.Counter(metrics.Options{Name: "nsq_spans_written", Tags: map[string]string{"status": "failure"}}),
	}

	p, e := nsq.NewProducer(options.Brokers, options.Topic, logger)
	if e != nil {
		logger.Error("Failed to create NSQ Producer", zap.Error(e))
		return nil
	}

	return &SpanWriter{
		producer:   p,
		marshaller: marshaller,
		metrics:    writeMetrics,
	}
}

// WriteSpan writes the span to NSQ.
func (s *SpanWriter) WriteSpan(ctx context.Context, span *model.Span) error {
	spanBytes, err := s.marshaller.Marshal(span)
	if err != nil {
		s.metrics.SpansWrittenFailure.Inc(1)
		return err
	}

	s.producer.WriteSpan(spanBytes)
	if err != nil {
		s.metrics.SpansWrittenFailure.Inc(1)
		return err
	}

	s.metrics.SpansWrittenSuccess.Inc(1)

	return nil
}

// Close closes SpanWriter by closing producer
func (s *SpanWriter) Close() {
	s.producer.Destroy()
}

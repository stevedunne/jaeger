package nsq

import (
	"context"
	"errors"

	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/pkg/nsq"
	"github.com/jaegertracing/jaeger/plugin/storage/kafka"
)

type spanReaderMetrics struct {
	SpansReadSuccess metrics.Counter
	SpansReadFailure metrics.Counter
}

// SpanWriter writes spans to NSQ. Implements spanstore.Writer
type SpanReader struct {
	metrics      *spanReaderMetrics
	consumer     *nsq.NsqConsumer
	unmarshaller kafka.Unmarshaller
}

// NewSpanWriter initiates and returns a new NSQ spanwriter
func NewSpanReader(
	options Options,
	readerOptions ReaderOptions,
	unmarshaller kafka.Unmarshaller,
	factory *metrics.Factory,
	logger *zap.Logger,
) *SpanReader {
	metricsFactory := (*factory)
	readMetrics := &spanReaderMetrics{
		SpansReadSuccess: metricsFactory.Counter(metrics.Options{Name: "nsq_spans_read", Tags: map[string]string{"status": "success"}}),
		SpansReadFailure: metricsFactory.Counter(metrics.Options{Name: "nsq_spans_read", Tags: map[string]string{"status": "failure"}}),
	}

	c, e := nsq.NewConsumer(options.Topic, "jaeger-consumer", []string{options.Brokers}, []string{readerOptions.Lookups}, logger)
	if e != nil {
		logger.Error("Failed to create NSQ Producer", zap.Error(e))
	}

	return &SpanReader{
		consumer:     c,
		unmarshaller: unmarshaller,
		metrics:      readMetrics,
	}
}

// WriteSpan writes the span to NSQ.
func (s *SpanReader) ReadSpan(ctx context.Context, span *model.Span) error {
	// spanBytes, err := s.marshaller.Marshal(span)
	// if err != nil {
	// 	s.metrics.SpansWrittenFailure.Inc(1)
	// 	return err
	// }

	// s.producer.WriteSpan(spanBytes)
	// if err != nil {
	// 	s.metrics.SpansWrittenFailure.Inc(1)
	// 	return err
	// }

	// s.metrics.SpansWrittenSuccess.Inc(1)

	return errors.New("Not implemented")
}

// Start the Spanreader
func (s *SpanReader) Start() {
	s.consumer.HandleMessages(func() {})
}

// Close closes SpanReader by closing consumer
func (s *SpanReader) Close() error {
	s.consumer.Destroy()

	return nil
}

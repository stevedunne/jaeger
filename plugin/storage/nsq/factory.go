package nsq

import (
	"errors"
	"flag"
	"sync"

	"github.com/jaegertracing/jaeger/plugin/storage/kafka"
	"github.com/jaegertracing/jaeger/storage/dependencystore"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type Factory struct {
	options       Options
	readerOptions ReaderOptions
	logger        *zap.Logger

	writer         *SpanWriter
	reader         *SpanReader
	marshaller     kafka.Marshaller
	unmarshaller   kafka.Unmarshaller
	metricsFactory *metrics.Factory

	mu_read  sync.Mutex
	mu_write sync.Mutex
}

// NewFactory creates a new Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// AddFlags implements plugin.Configurable
func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	f.options.AddFlags(flagSet)
}

// AddFlags implements plugin.Configurable
func (f *Factory) AddReaderFlags(flagSet *flag.FlagSet) {
	f.readerOptions.AddReaderFlags(flagSet)
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitReaderFromViper(v *viper.Viper) {
	f.readerOptions.InitReaderFromViper(v)
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitFromViper(v *viper.Viper) {
	f.options.InitFromViper(v)
}

// InitFromOptions initializes factory from options.
func (f *Factory) InitFromOptions(o Options) {
	f.options = o
}

// InitFromOptions initializes factory from options.
func (f *Factory) InitFromReaderOptions(o ReaderOptions) {
	f.readerOptions = o
}

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.logger = logger

	logger.Info("NSQ factory",
		zap.Any("brokers", f.options.Brokers),
		zap.Any("topic", f.options.Topic),
		zap.Any("encoding", f.options.Encoding),
	)

	switch f.options.Encoding {
	case EncodingProto:
		f.marshaller = kafka.NewProtobufMarshaller()
		f.unmarshaller = kafka.NewProtobufUnmarshaller()
	case EncodingJSON:
		f.marshaller = kafka.NewJSONMarshaller()
		f.unmarshaller = kafka.NewJSONUnmarshaller()
	default:
		return errors.New("Nsq encoding is not one of '" + EncodingJSON + "' or '" + EncodingProto + "'")
	}

	f.metricsFactory = &metricsFactory

	return nil
}

// CreateSpanReader creates a spanstore.Reader.
func (f *Factory) CreateSpanReader() (spanstore.Reader, error) {
	if f.reader == nil {
		f.mu_read.Lock()
		f.reader = NewSpanReader(f.options, f.readerOptions, f.unmarshaller, f.metricsFactory, f.logger)
		f.mu_read.Unlock()
	}
	return nil, errors.New("nsq CreateSpanReader is not implemented")
}

// CreateSpanWriter creates a spanstore.Writer.
func (f *Factory) CreateSpanWriter() (spanstore.Writer, error) {

	if f.writer == nil {
		f.mu_write.Lock()
		opts := Options{Brokers: f.options.Brokers, Topic: f.options.Topic, Encoding: f.options.Encoding}
		f.writer = NewSpanWriter(opts, f.marshaller, f.metricsFactory, f.logger)
		f.mu_write.Unlock()
	}

	return f.writer, nil
}

// CreateDependencyReader creates a dependencystore.Reader.
func (f *Factory) CreateDependencyReader() (dependencystore.Reader, error) {
	return nil, errors.New("nsq CreateDependencyReader is not implemented")
}

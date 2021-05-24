package iowriter

import (
	"context"

	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter/grpc"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/model"
)

type IOWriterFactory struct {
}

func NewIOWriterFactory() (*IOWriterFactory, error) {
	w := &IOWriterFactory{}
	return w, nil
}

type IOWriter struct {
	reporter grpc.Reporter
	logger   *zap.Logger
}

func (IOWriterFactory) NewIOWriter(r grpc.Reporter, l *zap.Logger) (*IOWriter, error) {
	w := &IOWriter{reporter: r, logger: l}
	return w, nil
}

func (i IOWriter) WriteSpan(ctx context.Context, span *model.Span) error {

	err := i.reporter.EmitProtobufBatch(context.TODO(), span)
	if err != nil {
		i.logger.Error("Error writing span to reporter", zap.Error(err))
		return err
	}
	i.logger.Debug("Span Sent " + span.OperationName)
	return nil
}

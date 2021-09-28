// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/cmd/collector/app/handler"
	"github.com/jaegertracing/jaeger/cmd/collector/app/processor"
	zs "github.com/jaegertracing/jaeger/cmd/collector/app/sanitizer/zipkin"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

// SpanHandlerBuilder holds configuration required for handlers
type SpanHandlerBuilder struct {
	SpanWriter     spanstore.Writer
	CollectorOpts  CollectorOptions
	Logger         *zap.Logger
	MetricsFactory metrics.Factory
}

// SpanHandlers holds instances to the span handlers built by the SpanHandlerBuilder
type SpanHandlers struct {
	ZipkinSpansHandler   handler.ZipkinSpansHandler
	JaegerBatchesHandler handler.JaegerBatchesHandler
	GRPCHandler          *handler.GRPCHandler
}

// BuildSpanProcessor builds the span processor to be used with the handlers
func (b *SpanHandlerBuilder) BuildSpanProcessor(additional ...ProcessSpan) processor.SpanProcessor {
	hostname, _ := os.Hostname()
	svcMetrics := b.metricsFactory()
	hostMetrics := svcMetrics.Namespace(metrics.NSOptions{Tags: map[string]string{"host": hostname}})
	logger := b.logger()
	preprocessor := &Preprocessor{Logger: logger}

	return NewSpanProcessor(
		b.SpanWriter,
		additional,
		Options.ServiceMetrics(svcMetrics),
		Options.HostMetrics(hostMetrics),
		Options.Logger(b.logger()),
		Options.SpanFilter(customSpanFilter),
		Options.NumWorkers(b.CollectorOpts.NumWorkers),
		Options.QueueSize(b.CollectorOpts.QueueSize),
		Options.CollectorTags(b.CollectorOpts.CollectorTags),
		Options.DynQueueSizeWarmup(uint(b.CollectorOpts.QueueSize)), // same as queue size for now
		Options.DynQueueSizeMemory(b.CollectorOpts.DynQueueSizeMemory),
		Options.PreProcessSpans(preprocessor.ProcessSpans),
	)
}

// BuildHandlers builds span handlers (Zipkin, Jaeger)
func (b *SpanHandlerBuilder) BuildHandlers(spanProcessor processor.SpanProcessor) *SpanHandlers {
	return &SpanHandlers{
		handler.NewZipkinSpanHandler(b.Logger, spanProcessor, zs.NewChainedSanitizer(zs.StandardSanitizers...)),
		handler.NewJaegerSpanHandler(b.Logger, spanProcessor),
		handler.NewGRPCHandler(b.Logger, spanProcessor),
	}
}

func customSpanFilter(span *model.Span) bool {

	if span != nil {
		if span.Process != nil && span.Process.ServiceName == "demo-service" {
			return false
		}

		if span.OperationName == "Deserialize" || span.OperationName == "Serialize" {
			return false
		}
		for i, v := range span.Tags {
			if strings.ToLower(v.Key) == "http.url" {
				lstr := strings.ToLower(span.Tags[i].VStr)
				if strings.HasSuffix(lstr, "ping") ||
					strings.HasSuffix(lstr, "health") ||
					strings.HasSuffix(lstr, "hc") ||
					strings.Contains(lstr, ".newrelic.") {
					return false
					// strings.HasSuffix(span.Tags[i].VStr, "heartbeat") ||
				}
			}
		}
	}

	return true
}

func (b *SpanHandlerBuilder) logger() *zap.Logger {
	if b.Logger == nil {
		return zap.NewNop()
	}
	return b.Logger
}

func (b *SpanHandlerBuilder) metricsFactory() metrics.Factory {
	if b.MetricsFactory == nil {
		return metrics.NullFactory
	}
	return b.MetricsFactory
}

type Preprocessor struct {
	Logger *zap.Logger
}

func (f *Preprocessor) ProcessSpans(spans []*model.Span) {
	// Only for a transient data gathering test
	// please dont use this ever
	for _, s := range spans {
		if s.Duration < 0 {
			s.Duration *= -1
			f.Logger.Debug(fmt.Sprintf("Corrected negative duration for %s %s", s.TraceID.String(), s.SpanID.String()))
			s.Tags = append(s.Tags, model.KeyValue{Key: "duration-adjusted", VStr: fmt.Sprintf("%v", s.Duration)})
		}
	}
}

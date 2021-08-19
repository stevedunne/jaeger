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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/cmd/flags"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/pkg/config"
	"github.com/jaegertracing/jaeger/plugin/storage/memory"
)

func TestNewSpanHandlerBuilder(t *testing.T) {
	v, command := config.Viperize(flags.AddFlags, AddFlags)

	require.NoError(t, command.ParseFlags([]string{}))
	cOpts := new(CollectorOptions).InitFromViper(v)

	spanWriter := memory.NewStore()

	builder := &SpanHandlerBuilder{
		SpanWriter:    spanWriter,
		CollectorOpts: *cOpts,
	}
	assert.NotNil(t, builder.logger())
	assert.NotNil(t, builder.metricsFactory())

	builder = &SpanHandlerBuilder{
		SpanWriter:     spanWriter,
		CollectorOpts:  *cOpts,
		Logger:         zap.NewNop(),
		MetricsFactory: metrics.NullFactory,
	}

	spanProcessor := builder.BuildSpanProcessor()
	spanHandlers := builder.BuildHandlers(spanProcessor)
	assert.NotNil(t, spanHandlers.ZipkinSpansHandler)
	assert.NotNil(t, spanHandlers.JaegerBatchesHandler)
	assert.NotNil(t, spanHandlers.GRPCHandler)
	assert.NotNil(t, spanProcessor)
}

func TestDefaultSpanFilter(t *testing.T) {
	assert.True(t, defaultSpanFilter(nil))
}

type FilterTestCase struct {
	key      string
	value    string
	expected bool
	message  string
}

func TestDefaultSpanFilter_RemovesUnwantedSpan(t *testing.T) {

	cases := []FilterTestCase{
		{key: "http.url", value: "http://localhost/ping", expected: false, message: "basic ping test"},
		{key: "http.url", value: "http://someserver/health", expected: false, message: "basic health test"},
		{key: "http.url", value: "http://localhost/hc", expected: false, message: "hc - health check test"},
		{key: "HTTP.url", value: "http://localhost/PING", expected: false, message: "upper case ping test"},
		{key: "http.url", value: "https://collector-003.newrelic.com/agent_listener/invoke_raw_method?method=connect&license_key=b8f18c786bebaacaac5344c65a437846cac79f82&marshal_format=json&protocol_version=16", expected: false, message: "basic ping test"},
	}

	for _, tc := range cases {

		span := model.Span{}
		span.Tags = append(span.Tags, model.KeyValue{Key: tc.key, VStr: tc.value})
		res := defaultSpanFilter(&span)
		assert.Equal(t, tc.expected, res, tc.message)

	}
}

// Copyright (c) 2018 The Jaeger Authors.
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

package builder

import (
	"fmt"
	"strings"

	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/plugin/storage/kafka"
	"github.com/jaegertracing/jaeger/plugin/storage/nsq"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

// CreateConsumer creates a new span consumer for the ingester
func CreateConsumer(logger *zap.Logger, metricsFactory metrics.Factory, spanWriter spanstore.Writer, options nsq.Options, readerOptions nsq.ReaderOptions) (*nsq.SpanReader, error) {
	var unmarshaller kafka.Unmarshaller
	switch options.Encoding {
	case kafka.EncodingJSON:
		unmarshaller = kafka.NewJSONUnmarshaller()
	case kafka.EncodingProto:
		unmarshaller = kafka.NewProtobufUnmarshaller()
		//  case kafka.EncodingZipkinThrift:
		//  	unmarshaller = kafka.NewZipkinThriftUnmarshaller()
	default:
		return nil, fmt.Errorf(`encoding '%s' not recognised, use one of ("%s")`,
			options.Encoding, strings.Join(kafka.AllEncodings, "\", \""))
	}

	spanProcessor := nsq.NewSpanReader(options, readerOptions, unmarshaller, &metricsFactory, logger)

	return spanProcessor, nil
}

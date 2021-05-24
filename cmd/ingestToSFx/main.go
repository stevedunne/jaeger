package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter/grpc"
	"github.com/jaegertracing/jaeger/cmd/flags"
	"github.com/jaegertracing/jaeger/cmd/ingestToSFx/iowriter"
	"github.com/jaegertracing/jaeger/cmd/ingester/app"
	"github.com/jaegertracing/jaeger/cmd/ingester/app/builder"
	"github.com/jaegertracing/jaeger/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	jexpvar "github.com/uber/jaeger-lib/metrics/expvar"
	"github.com/uber/jaeger-lib/metrics/fork"
	"go.uber.org/zap"
)

func main() {

	svc := flags.NewService(9999)

	spanWriterFactory, err := iowriter.NewIOWriterFactory()
	if err != nil {
		log.Fatalf("Cannot initialize iowriter factory: %v", err)
	}

	v := viper.New()
	command := &cobra.Command{
		Use:   "jaeger-ingest-signalfx",
		Short: "Jaeger ingester consumes from Kafka and writes to signalFx.",
		Long:  `Jaeger ingester consumes spans from a particular Kafka topic and writes them to a configured signalFx listener.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger // shortcut

			baseFactory := svc.MetricsFactory.Namespace(metrics.NSOptions{Name: "jaeger"})
			metricsFactory := baseFactory.Namespace(metrics.NSOptions{Name: "ingester"})

			options := app.Options{}
			options.InitFromViper(v)
			mFactory := fork.New("internal",
				jexpvar.NewFactory(10), // backend for internal opts
				metricsFactory)

			grpcBuilder := grpc.NewConnBuilder().InitFromViper(v)

			//create IoWriter
			grpcConn, err := grpcBuilder.CreateConnection(logger, mFactory)
			if err != nil {
				log.Fatalf("Cannot initialize GRPC connection to signalFx %v", err)
			}
			grpcReporter := grpc.NewReporter(grpcConn, map[string]string{}, logger)
			spanWriter, err := spanWriterFactory.NewIOWriter(*grpcReporter, logger)
			if err != nil {
				log.Fatalf("Cannot initialize iowriter: %v", err)
			}

			//Consumer
			consumer, err := builder.CreateConsumer(logger, metricsFactory, spanWriter, options)
			if err != nil {
				logger.Fatal("Unable to create consumer", zap.Error(err))
			}
			consumer.Start()

			//cleanup
			svc.RunAndThen(func() {
				if err := options.TLS.Close(); err != nil {
					logger.Error("Failed to close TLS certificates watcher", zap.Error(err))
				}
				if err = consumer.Close(); err != nil {
					logger.Error("Failed to close consumer", zap.Error(err))
				}
				if err = grpcConn.Close(); err != nil {
					logger.Error("Failed to close grpc reporter", zap.Error(err))
				}
			})
			return nil
		},
	}

	config.AddFlags(
		v,
		command,
		svc.AddFlags,
		app.AddFlags,
		reporter.AddFlags,
		grpc.AddFlags,
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

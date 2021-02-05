package binProcessor

import(
	"fmt"
	"os"
	"bufio"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/jaegertracing/jaeger/thrift-gen/agent"
	"github.com/jaegertracing/jaeger/cmd/agent/app/customtransport"

)

func GetBinProcessCommand() *cobra.Command {

	command := &cobra.Command{
		Use:   "processBinFile",
		Short: "Process binary files with dumps of failed data.",
		Long:  `Process binary files with dumps of failed data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
		
			logger, _ := zap.NewDevelopment()
			
			var name string
			if(len(args) < 1 ){
				//return fmt.Errorf("No file name specified")
				name = "jaeger-agent-err-2021-01-13T14-24-02-06-00.bin"
			} else {
				name = args[0]
			}

			name = "jaeger-agent-err-2021-01-15T18-05-04Z.bin"

			payload, err := readFileBytes(name)
			if err != nil {
				return err
			}
			logger.Info("Read File Contents")
			writeFileDebug(name, payload)

			cpf := thrift.NewTCompactProtocolFactory()
			trans := &customtransport.TBufferedReadTransport{}
			compactProtocol := cpf.GetProtocol(trans)
			compactProtocol.Transport().Write(payload)
			
			agt := agent.NewAgentEmitBatchArgs()

			if err = agt.Read(compactProtocol); err != nil {
				fmt.Println("Error reading batch - compact protocol", err)
			}

			// bpf := thrift.NewTBinaryProtocolFactoryDefault() 
			// trans1 := &customtransport.TBufferedReadTransport{}
			// binaryProtocol := bpf.GetProtocol(trans1)

			// binaryProtocol.Transport().Write(payload)

			// if err = agt.Read(binaryProtocol); err != nil {
			// 	fmt.Println("Error reading batch - binary protocol", err)
			// }
			
			return nil
		},
	}
	return command
}

func readFileBytes(name string) ([]byte, error) {
	
	file, err := os.Open("C:\\logs\\jaeger-agent-err\\" + name)

	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_,err = bufr.Read(bytes)
	
	return bytes, nil
}

func writeFileDebug(name string, payload []byte)  error {
	
	file, err := os.Create("C:\\logs\\jaeger-agent-err\\" + strings.ReplaceAll(name, ".bin", ".log") )
	if err != nil {
		return err
	}
	defer file.Close()

	var output string

//	s := string(payload) 
//	runes := []rune(s)
//	for index, item := range runes {
	for index, item := range payload {
		output += fmt.Sprintf("%d\t%#x\t%d\t%c\n", index, item, uint8(item), rune(item))
	}

	_,err = file.Write( []byte(output) )
	
	if err != nil {
		return err
	}

	return  nil
}


			// baseFactory := svc.MetricsFactory.
			// 	Namespace(metrics.NSOptions{Name: "jaeger"}).
			// 	Namespace(metrics.NSOptions{Name: "agent"})
			// mFactory := fork.New("internal",
			// 	jexpvar.NewFactory(10), // backend for internal opts
			// 	baseFactory)

			// rOpts := new(reporter.Options).InitFromViper(v, logger)
			// grpcBuilder := grpc.NewConnBuilder().InitFromViper(v)

			// builders := map[reporter.Type]app.CollectorProxyBuilder{
			// 	reporter.GRPC: app.GRPCCollectorProxyBuilder(grpcBuilder),
			// }
			// cp, err := app.CreateCollectorProxy(app.ProxyBuilderOptions{
			// 	Options: *rOpts,
			// 	Logger:  logger,
			// 	Metrics: mFactory,
			// }, builders)
			// if err != nil {
			// 	logger.Fatal("Could not create collector proxy", zap.Error(err))
			// }

			// builder := new(app.Builder).InitFromViper(v)
			// processor, err := builder.CreateBinaryCompactProcessor(cp, logger, mFactory)
			// if err != nil {
			// 	return fmt.Errorf("unable to initialize Jaeger processor: %w", err)
			// }
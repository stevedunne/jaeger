package nsq

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

// Options stores the configuration options for NSQ
// these options are shared between producers and consumers
// NOTE - consumers also need the ReaderOptions with the lookups list
type Options struct {
	Brokers  string `mapstructure:"brokers"`
	Topic    string `mapstructure:"topic"`
	Channel  string `mapstructure:"channel"`
	Encoding string `mapstructure:"encoding"`
}

type ReaderOptions struct {
	Lookups string `mapstructure:"brokers"`
}

const (
	// EncodingJSON is used for spans encoded as Protobuf-based JSON.
	EncodingJSON = "json"
	// EncodingProto is used for spans encoded as Protobuf.
	EncodingProto = "protobuf"
	//	// EncodingZipkinThrift is used for spans encoded as Zipkin Thrift.
	//	EncodingZipkinThrift = "zipkin-thrift"

	configPrefix         = "nsq"
	configConsumerPrefix = "nsq.consumer"

	suffixBrokers  = ".brokers"
	suffixLookups  = ".lookups"
	suffixTopic    = ".topic"
	suffixChannel  = ".channel"
	suffixEncoding = ".encoding"

	defaultBroker   = "127.0.0.1:4150"
	defaultTopic    = "jaeger-spans"
	defaultChannel  = "jaeger-ingester"
	defaultEncoding = EncodingProto

	defaultLookup = "127.0.0.1:4161"
)

// AddFlags adds flags for Options
func (opt *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		configPrefix+suffixBrokers,
		defaultBroker,
		"The comma-separated list of nsqd instances. i.e. '127.0.0.1:4150,0.0.0:1234'")
	flagSet.String(
		configPrefix+suffixTopic,
		defaultTopic,
		"The name of the NSQ topic")
	flagSet.String(
		configPrefix+suffixChannel,
		defaultChannel,
		"The name of the NSQ consumer channel")
	flagSet.String(
		configPrefix+suffixEncoding,
		defaultEncoding,
		fmt.Sprintf(`Encoding of spans ("%s" or "%s") sent to NSQ.`, EncodingJSON, EncodingProto),
	)
}

// InitFromViper initializes Options with properties from viper
func (opt *Options) InitFromViper(v *viper.Viper) {
	opt.Brokers = v.GetString(configPrefix + suffixBrokers)
	opt.Topic = v.GetString(configPrefix + suffixTopic)
	opt.Encoding = v.GetString(configPrefix + suffixEncoding)
}

// AddFlags adds flags for Options
func (opt *ReaderOptions) AddReaderFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		configConsumerPrefix+suffixLookups,
		defaultLookup,
		"The comma-separated list of nsq lookups. i.e. '127.0.0.1:4150,0.0.0:1234'")
}

// InitFromViper initializes Options with properties from viper
func (opt *ReaderOptions) InitReaderFromViper(v *viper.Viper) {
	opt.Lookups = v.GetString(configConsumerPrefix + suffixLookups)
}

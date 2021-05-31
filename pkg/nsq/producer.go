package nsq

import (
	"fmt"

	"github.com/jaegertracing/jaeger/pkg/version"
	"go.uber.org/zap"

	gonsq "github.com/nsqio/go-nsq"
)

type NsqProducer struct {
	producer *gonsq.Producer
	logger   *zap.Logger
	topic    string
}

func NewProducer(uri, topic string, l *zap.Logger) (*NsqProducer, error) {

	var producer = &NsqProducer{
		logger: l,
		topic:  topic,
	}

	cfg := gonsq.NewConfig()

	cfg.UserAgent = fmt.Sprintf("to_nsq/%s go-nsq/%s", version.Get().GitVersion, gonsq.VERSION)

	p, err := gonsq.NewProducer(uri, cfg)
	if err != nil {
		producer.logger.Error("failed to create nsq.Producer", zap.Error(err))
		return nil, err
	}

	producer.producer = p

	return producer, nil
}

func (n *NsqProducer) WriteSpan(bspan []byte) error {

	return n.producer.Publish(n.topic, bspan)

}

func (n *NsqProducer) Destroy() {
	n.producer.Stop()
}

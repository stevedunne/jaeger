package nsq

import (
	"log"

	"go.uber.org/zap"

	gonsq "github.com/nsqio/go-nsq"
)

type NsqConsumer struct {
	consumer *gonsq.Consumer
	logger   *zap.Logger
	options  *NsqConsumerOptions
}

type NsqConsumerOptions struct {
	Topic   string
	Channel string
	Nsqds   []string
	Lookups []string
}

func NewConsumer(topic, channel string, nsqds, lookups []string, l *zap.Logger) (*NsqConsumer, error) {

	var c = &NsqConsumer{
		logger: l,
		options: &NsqConsumerOptions{
			Topic:   topic,
			Channel: channel,
			Nsqds:   nsqds,
			Lookups: lookups,
		},
	}

	return c, nil
}

func (n *NsqConsumer) HandleMessages(handler func()) error {

	cfg := gonsq.NewConfig()

	cons, err := gonsq.NewConsumer(n.options.Topic, n.options.Channel, cfg)
	if err != nil {
		n.logger.Error("failed to create nsq.Consumer", zap.Error(err))
	}

	n.consumer = cons

	n.consumer.AddHandler(&MessageHandler{topicName: "jaeger-span", logger: n.logger})

	err = n.consumer.ConnectToNSQDs(n.options.Nsqds)
	if err != nil {
		log.Fatal(err)
	}

	err = n.consumer.ConnectToNSQLookupds(n.options.Lookups)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (n *NsqConsumer) Destroy() {
	n.consumer.Stop()
}

type MessageHandler struct {
	topicName     string
	totalMessages int
	logger        *zap.Logger
}

func (mh *MessageHandler) HandleMessage(m *gonsq.Message) error {
	mh.totalMessages++

	mh.logger.Info("Handled Message ",
		zap.Int("total", mh.totalMessages),
		zap.String("topicName", mh.topicName),
		zap.String("messageBody", string(m.Body)))

	return nil
}

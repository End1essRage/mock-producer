package publisher

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"gitlab.gitlab.bcs.ru/test-tools/mock-producer/api"
	"gitlab.gitlab.bcs.ru/test-tools/mock-producer/logger"
)

type Publisher struct {
	publisher amqp.Publisher
}

func New(connString, exchange string) *Publisher {
	p := &Publisher{}

	publisher, err := amqp.NewPublisher(initAmqpConfig(connString, exchange), logger.MsgLogger{})
	if err != nil {
		panic(err)
	}

	p.publisher = *publisher

	return p
}

func initAmqpConfig(connString, exchange string) amqp.Config {
	amqpConfig := amqp.Config{
		Connection: amqp.ConnectionConfig{
			AmqpURI: connString,
		},
		Exchange: amqp.ExchangeConfig{
			GenerateName: func(topic string) string {
				return exchange
			},
			Type:    "direct",
			Durable: true,
		},
		Queue: amqp.QueueConfig{
			GenerateName: func(topic string) string {
				return topic
			},
			Durable: true,
		},
		QueueBind: amqp.QueueBindConfig{
			GenerateRoutingKey: func(topic string) string {
				return topic
			},
		},
		Publish: amqp.PublishConfig{
			GenerateRoutingKey: func(topic string) string {
				return topic
			},
		},
		/*
			Consume: amqp.ConsumeConfig{
				Qos: amqp.QosConfig{
					PrefetchCount: 1,
				},
			},
		*/
		Marshaler:       amqp.DefaultMarshaler{},
		TopologyBuilder: &amqp.DefaultTopologyBuilder{},
	}

	return amqpConfig
}

func (p *Publisher) Publish(queue string, pattern []api.Pattern, delay int) error {
	for _, v := range pattern {
		logger.Log.WithField("caller", "publisher").Infof("pattern is %+v", v)
	}

	for item := range pattern {
		data, err := json.Marshal(item)
		if err != nil {
			logger.Log.WithField("caller", "publisher").Errorf("ошибка сериализации: %v", err)
			return err
		}

		msg := message.NewMessage(watermill.NewUUID(), data)

		if err := p.publisher.Publish(queue, msg); err != nil {
			logger.Log.WithField("caller", "publisher").Errorf("ошибка отправки сообщения: %v", err)
			return err
		}
	}

	return nil
}

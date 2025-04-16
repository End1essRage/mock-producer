package publisher

import (
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/api"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/logger"
)

type Publisher struct{}

func New() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(queue string, pattern []api.Pattern, delay int) {
	for _, v := range pattern {
		logger.Log.WithField("caller", "publisher").Infof("pattern is %+v", v)
	}

	logger.Log.WithField("caller", "Publisher").Warn("NOT IMPLEMENTED")
}

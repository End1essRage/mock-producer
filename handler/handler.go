package handler

import (
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/api"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/logger"
)

type Generator interface {
}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) HandlePattern(queue string, pattern api.Pattern, count int) error {
	logger.Log.WithField("caller", "HandlePattern").Warn("NOT IMPLEMENTED")
	return nil
}

func (h *Handler) HandleTemplate(template, queue string, override api.Pattern, count int) error {
	logger.Log.WithField("caller", "HandleTemplate").Warn("NOT IMPLEMENTED")
	return nil
}

func (h *Handler) GetTemplates() ([]string, error) {
	logger.Log.WithField("caller", "GetTemplates").Warn("NOT IMPLEMENTED")
	return nil, nil
}

func (h *Handler) GetTemplate(name string) (api.Pattern, error) {
	logger.Log.WithField("caller", "GetTemplate").Warn("NOT IMPLEMENTED")
	return nil, nil
}

func (h *Handler) AddUpdateTemplate(name string, pattern api.Pattern) error {
	logger.Log.WithField("caller", "AddUpdateTemplate").Warn("NOT IMPLEMENTED")
	return nil
}

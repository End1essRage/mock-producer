package handler

import (
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/api"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/logger"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/types"
)

type Generator interface {
	Generate(pattern api.Pattern) api.Pattern
}

type Buffer interface {
	GetAll() []string
	Get(name string) (api.Pattern, error)
	Set(name string, pattern api.Pattern) error
}

type Publisher interface {
	Publish(queue string, pattern []api.Pattern, delay int)
}

type Handler struct {
	buffer    Buffer
	generator Generator
	publisher Publisher
}

func New(b Buffer, p Publisher, g Generator) *Handler {
	return &Handler{buffer: b, publisher: p, generator: g}
}

func (h *Handler) HandlePattern(queue string, pattern api.Pattern, count, delay int) error {
	toSend := make([]api.Pattern, 0)

	for range count {
		toSend = append(toSend, h.generator.Generate(pattern))
	}

	go h.publisher.Publish(queue, toSend, delay)

	return nil
}

func (h *Handler) GenPattern(pattern api.Pattern) (api.Pattern, error) {
	return h.generator.Generate(pattern), nil
}

func (h *Handler) GenTemplate(template string, override api.Pattern) (api.Pattern, error) {
	tmpl, err := h.buffer.Get(template)
	if err != nil {
		logger.Log.WithField("caller", "HandleTemplate").Errorf("ошбика поулчения темплейтa %s из буфера: %v", template, err)
		return nil, types.NewTemplateNotFoundErr(err.Error())
	}

	//оверрайдим темплейт
	for k := range override {
		if override[k] != nil && tmpl[k] != nil {
			tmpl[k] = override[k]
		} else {
			logger.Log.WithField("caller", "HandleTemplate").Warnf("несуществующее поле %s в темплейте %s", k, template)
		}
	}

	toSend := h.generator.Generate(tmpl)

	return toSend, nil

}

func (h *Handler) HandleTemplate(template, queue string, override api.Pattern, count, delay int) error {
	toSend := make([]api.Pattern, 0)

	tmpl, err := h.buffer.Get(template)
	if err != nil {
		logger.Log.WithField("caller", "HandleTemplate").Errorf("ошбика поулчения темплейтa %s из буфера: %v", template, err)
		return types.NewTemplateNotFoundErr(err.Error())
	}

	//оверрайдим темплейт
	for k := range override {
		if override[k] != nil && tmpl[k] != nil {
			tmpl[k] = override[k]
		} else {
			logger.Log.WithField("caller", "HandleTemplate").Warnf("несуществующее поле %s в темплейте %s", k, template)
		}
	}

	for range count {
		toSend = append(toSend, h.generator.Generate(tmpl))
	}

	go h.publisher.Publish(queue, toSend, delay)

	return nil
}

func (h *Handler) GetTemplates() []string {

	return h.buffer.GetAll()
}

func (h *Handler) GetTemplate(name string) (api.Pattern, error) {
	pattern, err := h.buffer.Get(name)
	if err != nil {
		logger.Log.WithField("caller", "GetTemplate").Errorf("ошбика поулчения темплейтa %s из буфера: %v", name, err)
		return nil, err
	}

	return pattern, nil
}

func (h *Handler) AddUpdateTemplate(name string, pattern api.Pattern) error {
	if err := h.buffer.Set(name, pattern); err != nil {
		logger.Log.WithField("caller", "AddUpdateTemplate").Errorf("ошбика сохранения темплейтa %s в буфер: %v", name, err)
		return err
	}

	return nil
}

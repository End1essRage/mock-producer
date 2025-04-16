package handler

import (
	"gitlab.gitlab.bcs.ru/test-tools/mock-producer/api"
	"gitlab.gitlab.bcs.ru/test-tools/mock-producer/logger"
	"gitlab.gitlab.bcs.ru/test-tools/mock-producer/types"
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
	Publish(queue string, pattern []api.Pattern, delay int) error
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
		logger.Log.WithField("caller", "GenTemplate").Errorf("ошбика поулчения темплейтa %s из буфера: %v", template, err)
		return nil, types.NewTemplateNotFoundErr(err.Error())
	}

	logger.Log.WithField("caller", "GenTemplate").Debugf("поулчение темплейтa %s из буфера: %+v", template, tmpl)

	//оверрайдим темплейт
	tmpl = overrideTemplate(tmpl, override)

	toSend := h.generator.Generate(tmpl)

	return toSend, nil

}

func overrideTemplate(tmpl, override api.Pattern) api.Pattern {
	for key, overrideVal := range override {
		// Проверяем существование ключа в шаблоне
		tmplVal, existsInTemplate := tmpl[key]
		if !existsInTemplate {
			logger.Log.WithField("caller", "overrideTemplate").Warnf("несуществующее поле %s в темплейте", key)
			continue
		}

		// Обработка вложенных структур
		if overrideMap, ok := overrideVal.(map[string]interface{}); ok {
			if tmplMap, ok := tmplVal.(map[string]interface{}); ok {
				// Рекурсивный вызов для вложенной структуры
				tmpl[key] = overrideTemplate(tmplMap, overrideMap)
			} else {
				logger.Log.WithField("caller", "overrideTemplate").Warnf(
					"поле %s имеет разные типы (шаблон: %T, переопределение: %T)",
					key, tmplVal, overrideVal,
				)
			}
			continue
		}

		// Перезаписываем значение для примитивных типов
		tmpl[key] = overrideVal
	}
	return tmpl
}

func (h *Handler) HandleTemplate(template, queue string, override api.Pattern, count, delay int) error {
	toSend := make([]api.Pattern, 0)

	tmpl, err := h.buffer.Get(template)
	if err != nil {
		logger.Log.WithField("caller", "HandleTemplate").Errorf("ошбика поулчения темплейтa %s из буфера: %v", template, err)
		return types.NewTemplateNotFoundErr(err.Error())
	}

	//оверрайдим темплейт
	tmpl = overrideTemplate(tmpl, override)

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

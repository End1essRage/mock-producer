package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gitlab.gitlab.bcs.ru/test-tools/mock-producer/logger"
	"gitlab.gitlab.bcs.ru/test-tools/mock-producer/types"
)

type Pattern map[string]interface{}

type Request struct {
	Queue   string  `json:"queue,omitempty"`
	Pattern Pattern `json:"pattern"`
	Count   int     `json:"count,omitempty"`
	Delay   int     `json:"dekay,omitempty"`
}

type TemplateRequest struct {
	Request
	Template string `json:"template"`
}

type Handler interface {
	//отправка сообщения
	HandlePattern(queue string, pattern Pattern, count, delay int) error
	HandleTemplate(template, queue string, override Pattern, count, delay int) error

	//Генерация
	GenPattern(pattern Pattern) (Pattern, error)
	GenTemplate(template string, override Pattern) (Pattern, error)

	//управление темплейтами
	GetTemplates() []string
	GetTemplate(name string) (Pattern, error)
	AddUpdateTemplate(name string, pattern Pattern) error
}

type API struct {
	router  *chi.Mux
	server  *http.Server
	handler Handler
}

func New(h Handler) *API {
	r := chi.NewRouter()

	return &API{
		router:  r,
		handler: h,
	}
}

func (a *API) Start(addr string) error {
	// Регистрируем обработчики из конфига
	a.registerHandlers()

	a.server = &http.Server{
		Addr:    addr,
		Handler: a.router,
	}

	go func() {
		logrus.Infof("Starting API server on %s", addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("API server error: %v", err)
		}
	}()

	return nil
}

func (a *API) registerHandlers() {
	//сгенерировать 1 и получить по паттерну
	a.router.Post("/gen-pattern", func(w http.ResponseWriter, r *http.Request) {
		var body Request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.Log.WithField("endpoint", "/gen-pattern").Errorf("%v", err)

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		item, err := a.handler.GenPattern(body.Pattern)
		if err != nil {
			logger.Log.WithField("endpoint", "/gen-pattern").Errorf("%v", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		jsonData, err := json.Marshal(item)
		if err != nil {
			logger.Log.WithField("endpoint", "GET /gen-pattern").Errorf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")
		// Отправляем данные
		w.Write(jsonData)
	})

	//сгенерировать 1 и получить по темплейту
	a.router.Post("/gen-template", func(w http.ResponseWriter, r *http.Request) {
		var body TemplateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.Log.WithField("endpoint", "/gen-template").Errorf("%v", err)

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		if body.Template == "" {
			logger.Log.WithField("endpoint", "/gen-template").Error("template не может быть пустым")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("template не может быть пустым"))
			return
		}

		item, err := a.handler.GenTemplate(body.Template, body.Pattern)
		if err != nil {
			var templateNotFoundErr types.TemplateNotFoundErr
			switch {
			case errors.As(err, &templateNotFoundErr):
				logger.Log.WithField("endpoint", "/gen-template").Errorf("%v", err)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(fmt.Sprintf("темплейт с именем '%s' не найден", body.Template)))
				return
			default:
			}
			logger.Log.WithField("endpoint", "/gen-template").Errorf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		jsonData, err := json.Marshal(item)
		if err != nil {
			logger.Log.WithField("endpoint", "GET /gen-template").Errorf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")
		// Отправляем данные
		w.Write(jsonData)
	})

	//отправить сообщения, очередь, json'ка со значениями(паттерном),кол-во
	a.router.Post("/send-pattern", func(w http.ResponseWriter, r *http.Request) {
		var body Request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.Log.WithField("endpoint", "/send-pattern").Errorf("%v", err)

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		if body.Queue == "" {
			logger.Log.WithField("endpoint", "/send-pattern").Error("queue не может быть пустым")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("queue не может быть пустым"))
			return
		}

		if body.Count < 1 {
			logger.Log.WithField("endpoint", "/send-pattern").Error("кол-во сообщений не может быть < 1")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("кол-во сообщений не может быть < 1"))
			return
		}

		if err := a.handler.HandlePattern(body.Queue, body.Pattern, body.Count, body.Delay); err != nil {
			logger.Log.WithField("endpoint", "/send-pattern").Errorf("%v", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	//отправить сообщения по темплейту , типы map, кол-во
	a.router.Post("/send-template", func(w http.ResponseWriter, r *http.Request) {
		var body TemplateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.Log.WithField("endpoint", "/send-template").Errorf("%v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		if body.Queue == "" || body.Template == "" {
			logger.Log.WithField("endpoint", "/send-template").Error("queue и template не может быть пустым")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("queue и template не может быть пустым"))
			return
		}

		if body.Count < 1 {
			logger.Log.WithField("endpoint", "/send-template").Error("кол-во сообщений не может быть < 1")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("кол-во сообщений не может быть < 1"))
			return
		}

		if err := a.handler.HandleTemplate(body.Template, body.Queue, body.Pattern, body.Count, body.Delay); err != nil {
			var templateNotFoundErr types.TemplateNotFoundErr
			switch {
			case errors.As(err, &templateNotFoundErr):
				logger.Log.WithField("endpoint", "/send-template").Errorf("%v", err)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(fmt.Sprintf("темплейт с именем '%s' не найден", body.Template)))
				return
			default:
			}
			logger.Log.WithField("endpoint", "/send-template").Errorf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	//получить все шаблоны
	a.router.Get("/template", func(w http.ResponseWriter, r *http.Request) {
		templates := a.handler.GetTemplates()

		jsonData, err := json.Marshal(templates)
		if err != nil {
			logger.Log.WithField("endpoint", "GET /template").Errorf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")
		// Отправляем данные
		w.Write(jsonData)
	})

	//получить шаблон по имени
	a.router.Get("/template/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Имя шаблона обязательно"))
			return
		}

		template, err := a.handler.GetTemplate(name)
		if err != nil {
			var templateNotFoundErr types.TemplateNotFoundErr
			switch {
			case errors.As(err, &templateNotFoundErr):
				logger.Log.WithField("endpoint", fmt.Sprintf("GET /template/%s", name)).Errorf("%v", err)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(fmt.Sprintf("темплейт с именем '%s' не найден", name)))
				return
			default:
			}

			logger.Log.WithField("endpoint", fmt.Sprintf("GET /template/%s", name)).Errorf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		jsonData, err := json.Marshal(template)
		if err != nil {
			logger.Log.WithField("endpoint", fmt.Sprintf("GET /template/%s", name)).Errorf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")
		// Отправляем данные
		w.Write(jsonData)
	})

	//создать/обновить шаблон
	a.router.Post("/template/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Имя шаблона обязательно"))
			return
		}

		var pattern Pattern
		if err := json.NewDecoder(r.Body).Decode(&pattern); err != nil {
			logger.Log.WithField("endpoint", fmt.Sprintf("POST /template/%s", name)).Errorf("%v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		if err := a.handler.AddUpdateTemplate(name, pattern); err != nil {
			logger.Log.WithField("endpoint", fmt.Sprintf("POST /template/%s", name)).Errorf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		}

		w.WriteHeader(http.StatusOK)
	})

}

func (a *API) Stop() {
	if a.server != nil {
		if err := a.server.Shutdown(context.Background()); err != nil {
			logrus.Errorf("API server shutdown error: %v", err)
		}
	}
}

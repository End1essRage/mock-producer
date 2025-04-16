package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/api"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/buffer"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/config"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/generator"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/handler"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/logger"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/publisher"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/secrets"
)

var (
	Env   string
	Debug bool // 1, true, True
)

func init() {
	// Настройка формата вывода (JSON)
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
	})

	Env = os.Getenv("ENV")
	if Env == "" {
		if err := godotenv.Load(); err != nil {
			logger.Log.Warn("error while reading environment", err.Error())
		}
	}

	Env = os.Getenv("ENV")
	if Env == "" {
		logger.Log.Warn("cant set environment, setting to test by default")
		Env = config.ENV_TEST
	}

	eDebug := os.Getenv("DEBUG")
	Debug = eDebug == "1" || eDebug == "true" || eDebug == "True"

	//уровень логирования
	logrus.SetLevel(logrus.InfoLevel)
	if Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	//при локальном запуске логи пишутся в файл
	if Env == config.ENV_DEV {
		// Создаем папку для логов
		if err := os.MkdirAll("logs", 0755); err != nil {
			logger.Log.Fatalf("Failed to create logs directory: %v", err)
		}

		// Генерируем имя файла на основе времени запуска
		filename := time.Now().Format("2006-01-02_15-04-05") + ".log"
		filePath := filepath.Join("logs", filename)

		// Создаем/открываем файл логов
		logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.Log.Fatalf("Failed to open log file: %v", err)
		}

		// Настраиваем вывод в файл и консоль
		logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}
}

func main() {
	var cfg config.Config
	// parse with generics
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		panic(err)
	}

	logger.Log.Debugf("%+v", cfg)

	//загружаем темплейты в память
	b := buffer.New()
	b.FillFromFiles("./templates")

	//генерируем строку подключения к rmq
	login := ""
	pwd := ""
	if Env == config.ENV_DEV {
		login = cfg.RMQ_LOGIN
		pwd = cfg.RMQ_PWD
	} else {
		creds, err := secrets.NewVault(&cfg).GetCredentials()
		if err != nil {
			panic(err)
		}

		login = creds.RmqLogin
		pwd = creds.RmqPwd
	}

	connString := fmt.Sprintf("amqp://%s:%s@%s/", login, pwd, cfg.RMQ_ADDRESS)
	if cfg.RMQ_VIRTUAL_HOST != "/" {
		connString += cfg.RMQ_VIRTUAL_HOST
	}

	//запускаем приложение
	p := publisher.New(connString, cfg.RMQ_EXCHANGE)
	g := generator.New()

	handler := handler.New(b, p, g)
	api := api.New(handler)

	api.Start(":8080")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	api.Stop()

	logger.Log.Info("Server stopped")
}

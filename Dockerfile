# Используем образ golang для сборки
FROM artifactory.gitlab.bcs.ru/docker-local/golang:1.24 AS builder

ARG GO_LIBS_TOKEN

RUN git config --global url."https://oauth2:${GO_LIBS_TOKEN}@gitlab.gitlab.bcs.ru".insteadOf "https://gitlab.gitlab.bcs.ru"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/integr-bus ./main.go

FROM registry.gitlab.bcs.ru/devops/images/docker/alpine-3-18-with-certs:master

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранный бинарник из этапа сборки
COPY --from=builder /app/integr-bus .
COPY --from=builder /app/templates .

EXPOSE 80
EXPOSE 443

CMD ["./integr-bus"]


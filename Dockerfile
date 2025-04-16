# Используем образ golang для сборки

FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/mock-producer ./main.go

FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранный бинарник из этапа сборки
COPY --from=builder /app/mock-producer .
COPY --from=builder /app/templates ./templates

EXPOSE 80
EXPOSE 443

CMD ["./mock-producer"]


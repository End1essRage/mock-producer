# Mock Producer API

API для генерации и отправки сообщений с поддержкой шаблонов и динамических данных

## Основные возможности

- 🎨 Генерация данных по паттернам
- 📁 Работа с шаблонами сообщений
- ✨ Встроенные генераторы случайных данных
- ⏱ Контроль частоты отправки сообщений
- 🔄 Переопределение полей шаблонов

## Быстрый старт

### Переменные окружения
```env
ENV="ENV_DEV"
DEBUG="1"

#common
#rmq
RMQ_ADDRESS="rabbitmq.example.com"
RMQ_VIRTUAL_HOST="/"
RMQ_PORT="5672"        
RMQ_EXCHANGE="" 

SECRET_KEY_RMQ_LOGIN="rmq.login"
SECRET_KEY_RMQ_PWD="rmq.password"

#for dev
RMQ_LOGIN="login"
RMQ_PWD="pwd"

#vault    
VAULT_SERVER="https://vault.example.com:8200"
VAULT_SECRET_PATH="secret"
VAULT_MOUNT_POINT="secret"
VAULT_ROLE_NAME="app-role"
```

## Эндпоинты

у паттерн может быть любая глубина map[string]interface{}

1. Генерация данных
| Метод | URL | Описание |
|:------|:----|:---------|
| POST | {{ _.base_url }}/gen-pattern | Генерация по паттерну |
| POST | {{ _.base_url }}/gen-template | Генерация по шаблону |		

Пример запроса генерации:
```json
{
  "template": "templateName", // только для template
  "pattern": {
    "id": "randNum",
    "name": "randString",
    "struct": {
        "field" : "data"
    }
  }
}
```

2. Отправка сообщений
| Метод | URL | Описание |
|:------|:----|:---------|
| POST | {{ _.base_url }}/send-pattern | Отправка по паттерну |
| POST | {{ _.base_url }}/send-template | Отправка по шаблону |		


Параметры отправки:
```json
{
  "queue": "my_queue",
  "count": 5,
  "delay": 200
}
```

3. Управление шаблонами
| Метод | URL | Описание |
|:------|:----|:---------|
| GET | {{ _.base_url }}/template | Список шаблонов |
| GET | {{ _.base_url }}/template/:name | Получить шаблон |
| POST | {{ _.base_url }}/template/:name | Создать, обновить шаблон |		

## Генерация случайных данных

Используйте специальные ключи в полях паттерна:

| Ключ | Описание | Пример |
|:-----|:---------|:-------|
| randNum | Случайное целое число | 42 |
| randString | Случайная строка | "xYz3fG" |

Пример использования:
```json
{
  "user": {
    "id": "randNum",
    "token": "randString"
  }
}
```

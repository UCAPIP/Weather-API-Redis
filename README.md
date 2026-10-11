# Weather API Wrapper Service

A robust Go-based API wrapper that fetches weather data from a 3rd-party provider (Visual Crossing), caches responses in Redis with a 12-hour TTL, and protects the endpoint with rate limiting.

🇬🇧 [English](#english) | 🇷🇺 [Русский](#русский)

---

## English

### ✨ Features

- **3rd-Party API Integration**: Fetches real-time weather data from the Visual Crossing API.
- **Redis Caching**: Stores responses with a 12-hour Time-To-Live (TTL) to reduce external API calls and improve response time.
- **Rate Limiting**: Protects the endpoint from abuse using a token bucket algorithm.
- **Environment Configuration**: Securely manages the external API key via `.env` files.
- **Resilience**: HTTP client timeouts prevent the server from hanging if the external API is down.
- **Clean Architecture**: Clear separation of concerns between HTTP handlers, middleware, and external clients.

### 🛠 Tech Stack

- **Language**: Go (1.22+)
- **Cache**: Redis (`github.com/redis/go-redis/v9`)
- **Configuration**: `github.com/joho/godotenv`
- **Rate Limiting**: `golang.org/x/time/rate`
- **External API**: Visual Crossing Weather API

### 📁 Project Structure

```text
weather-api/
├── main.go              # Server entry point, routing, and initialization
├── handlers.go          # HTTP handlers and external API client logic
├── middleware.go        # Rate limiting middleware
├── .env                 # Local environment variables (ignored by Git)
├── .env.example         # Template for environment variables
├── go.mod               # Go module definition
└── go.sum               # Go module checksums
```

### 🚀 Getting Started

**Prerequisites:**
- Go 1.22 or newer
- A running local Redis instance (default: `localhost:6379`)
- A free API key from [Visual Crossing](https://www.visualcrossing.com/)

**1. Clone and install dependencies:**
```bash
git clone <your-repo-url>
cd weather-api
go mod tidy
```

**2. Configure environment variables:**
Copy the example file and fill in your Visual Crossing API key:
```bash
cp .env.example .env
```
Edit `.env`:
```env
API_KEY=your_visual_crossing_api_key_here
```
*(Note: The application is configured to connect to Redis at `localhost:6379` and serve on port `8080` by default).*

**3. Run the server:**
```bash
go run .
```
The server will start on `http://localhost:8080`.

### 📖 Usage

Make a GET request to the `/weather` endpoint with a `city` query parameter:

```bash
curl "http://localhost:8080/weather?city=London"
```

**Success Response (200 OK):**
```json
{
  "address": "London, United Kingdom",
  "currentConditions": {
    "temp": 15.5
  }
}
```

**Error Responses:**
- `400 Bad Request`: Missing `city` parameter or invalid input.
- `429 Too Many Requests`: Rate limit exceeded.
- `500 Internal Server Error`: External API failure or Redis connection issue.

### ⚙️ How It Works

1. A request arrives at `/weather?city=...`.
2. The **Rate Limiter** checks if the client has exceeded the allowed request quota.
3. The handler checks **Redis** for a cached response using the city name as the key.
4. **Cache Hit**: Returns the cached JSON immediately.
5. **Cache Miss**: The server makes an HTTP request to the Visual Crossing API, parses the response, saves it to Redis with a `12 * time.Hour` expiration, and returns it to the client.

### 🚧 Possible Improvements

- Move Redis connection string and server port to `.env` for easier deployment.
- Implement per-IP rate limiting instead of a global limiter.
- Add support for multiple weather providers with fallback logic.
- Implement graceful shutdown for the HTTP server and Redis connections.
- Add structured logging (e.g., `slog` or `zap`).
- Containerize the application and Redis using `docker-compose`.

---

## Русский

### ✨ Возможности

- **Интеграция со сторонним API**: Получение актуальных данных о погоде через Visual Crossing API.
- **Кэширование в Redis**: Сохранение ответов с временем жизни (TTL) 12 часов для снижения нагрузки на внешний API и ускорения ответов.
- **Ограничение частоты запросов (Rate Limiting)**: Защита эндпоинта от злоупотреблений с помощью алгоритма token bucket.
- **Конфигурация через окружение**: Безопасное хранение API-ключа внешнего сервиса через `.env` файл.
- **Отказоустойчивость**: Таймауты HTTP-клиента предотвращают зависание сервера при недоступности внешнего API.
- **Чистая архитектура**: Чёткое разделение ответственности между HTTP-обработчиками, middleware и клиентами внешних сервисов.

### 🛠 Технологии

- **Язык**: Go (1.22+)
- **Кэш**: Redis (`github.com/redis/go-redis/v9`)
- **Конфигурация**: `github.com/joho/godotenv`
- **Rate Limiting**: `golang.org/x/time/rate`
- **Внешний API**: Visual Crossing Weather API

### 📁 Структура проекта

```text
weather-api/
├── main.go              # Точка входа, маршруты и инициализация
├── handlers.go          # HTTP-обработчики и логика клиента внешнего API
├── middleware.go        # Middleware для ограничения частоты запросов
├── .env                 # Локальные переменные окружения (игнорируется Git)
├── .env.example         # Шаблон переменных окружения
├── go.mod               # Определение Go-модуля
└── go.sum               # Контрольные суммы зависимостей
```

### 🚀 Запуск

**Требования:**
- Go 1.22 или новее
- Запущенный локальный экземпляр Redis (по умолчанию: `localhost:6379`)
- Бесплатный API-ключ с сайта [Visual Crossing](https://www.visualcrossing.com/)

**1. Клонирование и установка зависимостей:**
```bash
git clone <your-repo-url>
cd weather-api
go mod tidy
```

**2. Настройка переменных окружения:**
Скопируйте шаблон и укажите свой API-ключ Visual Crossing:
```bash
cp .env.example .env
```
Отредактируйте `.env`:
```env
API_KEY=ваш_ключ_visual_crossing_здесь
```
*(Примечание: Приложение настроено на подключение к Redis по адресу `localhost:6379` и запуск сервера на порту `8080` по умолчанию).*

**3. Запуск сервера:**
```bash
go run .
```
Сервер запустится на `http://localhost:8080`.

### 📖 Использование

Отправьте GET-запрос на эндпоинт `/weather` с параметром `city`:

```bash
curl "http://localhost:8080/weather?city=Moscow"
```

**Успешный ответ (200 OK):**
```json
{
  "address": "Moscow, Russia",
  "currentConditions": {
    "temp": -2.5
  }
}
```

**Ответы с ошибками:**
- `400 Bad Request`: Отсутствует параметр `city` или некорректный ввод.
- `429 Too Many Requests`: Превышен лимит запросов.
- `500 Internal Server Error`: Ошибка внешнего API или проблемы с подключением к Redis.

### ⚙️ Как это работает

1. Запрос поступает на `/weather?city=...`.
2. **Rate Limiter** проверяет, не превысил ли клиент допустимую квоту запросов.
3. Обработчик проверяет **Redis** на наличие кэшированного ответа, используя название города в качестве ключа.
4. **Попадание в кэш (Cache Hit)**: Кэшированный JSON возвращается мгновенно.
5. **Промах кэша (Cache Miss)**: Сервер делает HTTP-запрос к Visual Crossing API, парсит ответ, сохраняет его в Redis с TTL `12 * time.Hour` и возвращает клиенту.

### 🚧 Возможные улучшения

- Вынести строку подключения к Redis и порт сервера в `.env` для упрощения деплоя.
- Реализовать ограничение запросов для каждого IP-адреса отдельно (per-IP rate limiting).
- Добавить поддержку нескольких провайдеров погоды с логикой fallback.
- Реализовать корректное завершение работы (graceful shutdown) для HTTP-сервера и соединений Redis.
- Добавить структурированное логирование (например, через `slog` или `zap`).
- Контейнеризировать приложение и Redis с помощью `docker-compose`.

---

## License / Лицензия

This project is created for learning purposes. / Проект создан в учебных целях.
# Структура проекта

## Обзор

Проект представляет собой высоконагруженный сервис на Go с базовой аналитикой,
детекцией аномалий и развёртыванием в Kubernetes с авто-масштабированием
и мониторингом.

---

## Структура директорий

```
Ршпрдщфв-final-project/
├── cmd/
│   └── service/
│       ├── main.go
│       └── service_test.go
│
├── internal/
│   ├── analytics/
│   │   ├── rolling_average.go
│   │   ├── anomaly_detector.go
│   │   ├── rolling_average_test.go
│   │   └── anomaly_detector_test.go
│   │
│   ├── cache/
│   │   ├── redis.go
│   │   └── cache.go
│   │
│   └── metrics/
│       └── prometheus.go
│
├── k8s/
│   ├── deployments/
│   │   ├── go-service-deployment.yaml
│   │   └── redis-deployment.yaml
│   │
│   ├── services/
│   │   ├── go-service-service.yaml
│   │   └── redis-service.yaml
│   │
│   ├── configmaps/
│   │   ├── app-config.yaml
│   │   └── redis-secret.yaml
│   │
│   ├── hpa/
│   │   └── go-service-hpa.yaml
│   │
│   ├── ingress/
│   │   └── go-service-ingress.yaml
│   │
│   └── monitoring/
│       ├── prometheus-config.yaml
│       └── service-monitor.yaml
│
├── scripts/
│   ├── build.sh
│   ├── deploy.sh
│   ├── setup-monitoring.sh
│   ├── load-test.sh
│   ├── test-ab.sh
│   ├── load_anomaly_test.sh
│   └── locustfile.py
│
├── xtemp/
│   ├── DEPLOYMENT_GUIDE.md
│   ├── TESTING_GUIDE.md
│   ├── GRAFANA_DASHBOARDS.md
│   └── PROJECT_REPORT_TEMPLATE.md
│
├── .gitignore
├── .dockerignore
├── Dockerfile
├── go.mod
├── go.sum
├── README.md
├── QUICKSTART.md
└── PROJECT_STRUCTURE.md
```


---

## Компоненты

### Go Service (`cmd/service/`)

Основной HTTP-сервис с эндпоинтами:
- `POST /metrics` — приём метрик
- `GET /analyze` — текущая аналитика и состояние детектора
- `GET /health` — health check
- `GET /metrics` — Prometheus метрики

**Особенности:**
- интеграция с Redis для кэширования
- rolling average для сглаживания нагрузки
- z-score детекция аномалий
- экспорт метрик для Prometheus
- потокобезопасная реализация

**Тестирование:**
- `service_test.go` — HTTP-тесты с in-memory кэшем (без Redis)

---

### Analytics (`internal/analytics/`)

**Rolling Average:**
- окно: 50 событий
- thread-safe реализация
- используется для сглаживания и анализа трендов

**Anomaly Detector:**
- метод: z-score
- threshold: 2σ
- окно: 50 событий
- поддержка warm-up фазы
- хранение последнего z-score и результата детекции

**Тестирование:**
- юнит-тесты для rolling average
- юнит-тесты для anomaly detector
- проверка детекции выбросов

---

### Cache (`internal/cache/`)

Redis-клиент для кэширования метрик:
- TTL: 5 минут
- JSON-сериализация
- connection pooling

Дополнительно:
- интерфейс `Cache` для повышения тестируемости
- возможность подмены Redis in-memory реализацией в тестах

---

### Metrics (`internal/metrics/`)

Prometheus метрики:
- `http_requests_total`
- `http_request_duration_seconds`
- `rps_rate`
- `anomalies_detected_total`
- `anomaly_rate_per_minute`
- `rolling_average_value`
- `cpu_usage_percent`

---

## Kubernetes манифесты

**Deployments:**
- Go-сервис: несколько реплик, ресурсы, probes
- Redis: пароль через Secret

**Services:**
- ClusterIP для внутренней коммуникации

**HPA:**
- масштабирование по CPU (70%)
- min: 2 реплики
- max: 5 реплик

**Ingress:**
- NGINX Ingress (опционально)
- маршруты `/metrics`, `/analyze`

**Monitoring:**
- ServiceMonitor для Prometheus
- автоматическое обнаружение метрик сервиса

---

## Скрипты

### `load_anomaly_test.sh`
Пользовательский скрипт воспроизводимой нагрузки:
- прогрев окна аналитики
- генерация аномальных значений
- получение `/analyze`
Используется для демонстрации детекции аномалий.

### `load-test.sh`
Нагрузочное тестирование с Locust:
- параметризуемый RPS
- генерация HTML-отчёта

### `test-ab.sh`
Нагрузочное тестирование с Apache Bench:
- высокая интенсивность запросов
- проверка пропускной способности и HPA

---

## Тестирование

### Юнит-тесты
- rolling average
- anomaly detector

### Интеграционные
- HTTP API (без Redis)
- проверка `/metrics` и `/analyze`

### Нагрузочные
- Locust (1000+ RPS)
- Apache Bench
- проверка авто-масштабирования HPA
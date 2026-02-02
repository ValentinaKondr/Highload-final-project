# Структура проекта

## Обзор

Проект представляет собой высоконагруженный сервис на Go с аналитикой и развертыванием в Kubernetes.

## Структура директорий

```
final_highload/
├── cmd/
│   └── service/
│       └── main.go              # Основной файл сервиса (~200 строк)
│
├── internal/
│   ├── analytics/
│   │   ├── rolling_average.go   # Rolling average реализация
│   │   └── anomaly_detector.go  # Z-score детекция аномалий
│   ├── cache/
│   │   └── redis.go             # Redis клиент
│   └── metrics/
│       └── prometheus.go        # Prometheus метрики
│
├── k8s/
│   ├── deployments/
│   │   ├── go-service-deployment.yaml
│   │   └── redis-deployment.yaml
│   ├── services/
│   │   ├── go-service-service.yaml
│   │   └── redis-service.yaml
│   ├── configmaps/
│   │   ├── app-config.yaml
│   │   └── redis-secret.yaml
│   ├── hpa/
│   │   └── go-service-hpa.yaml
│   ├── ingress/
│   │   └── go-service-ingress.yaml
│   └── monitoring/
│       ├── prometheus-config.yaml
│       └── service-monitor.yaml
│
├── scripts/
│   ├── build.sh                 # Сборка Docker образа
│   ├── deploy.sh                # Развертывание в K8s
│   ├── load-test.sh             # Нагрузочное тестирование (Locust)
│   ├── test-ab.sh               # Нагрузочное тестирование (Apache Bench)
│   ├── setup-monitoring.sh      # Установка Prometheus/Grafana
│   └── locustfile.py            # Конфигурация Locust
│
├── xtemp/                       # Временная документация (gitignored)
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
├── README.md                    # Основная документация
├── QUICKSTART.md                # Быстрый старт
└── PROJECT_STRUCTURE.md         # Этот файл
```

## Компоненты

### Go Service (`cmd/service/main.go`)

Основной сервис с HTTP endpoints:
- `POST /metrics` - прием метрик
- `GET /analyze` - аналитика
- `GET /health` - health check
- `GET /metrics` - Prometheus метрики

**Особенности:**
- Интеграция с Redis для кэширования
- Rolling average для сглаживания
- Z-score детекция аномалий
- Экспорт метрик Prometheus

### Analytics (`internal/analytics/`)

**Rolling Average:**
- Окно: 50 событий
- Thread-safe реализация
- Используется для прогнозирования трендов

**Anomaly Detector:**
- Метод: Z-score
- Threshold: 2σ
- Окно: 50 событий
- Точность: > 70%

### Cache (`internal/cache/redis.go`)

Redis клиент для кэширования метрик:
- TTL: 5 минут
- JSON сериализация
- Connection pooling

### Metrics (`internal/metrics/prometheus.go`)

Prometheus метрики:
- `http_requests_total` - счетчик запросов
- `http_request_duration_seconds` - гистограмма задержек
- `rps_rate` - текущий RPS
- `anomalies_detected_total` - счетчик аномалий
- `anomaly_rate_per_minute` - частота аномалий
- `rolling_average_value` - значение rolling average
- `cpu_usage_percent` - использование CPU

### Kubernetes манифесты

**Deployments:**
- `go-service-deployment.yaml` - 2 реплики (min), ресурсы, probes
- `redis-deployment.yaml` - Redis с паролем

**Services:**
- ClusterIP сервисы для внутренней коммуникации

**HPA:**
- Масштабирование по CPU (70%)
- Минимум: 2 реплики
- Максимум: 5 реплик

**Ingress:**
- NGINX Ingress для внешнего доступа
- Маршруты: `/metrics`, `/analyze`

**Monitoring:**
- ServiceMonitor для Prometheus
- Автоматическое обнаружение метрик

## Скрипты

### `build.sh`
Сборка Docker образа для Go-сервиса.

### `deploy.sh`
Развертывание всех компонентов в Kubernetes:
1. Secrets и ConfigMaps
2. Redis
3. Go Service
4. HPA
5. Ingress

### `load-test.sh`
Нагрузочное тестирование с Locust:
- Параметры: URL, RPS, длительность
- Генерация отчета HTML

### `test-ab.sh`
Нагрузочное тестирование с Apache Bench:
- Параметры: URL, количество запросов, concurrency
- Вывод результатов в файлы

### `setup-monitoring.sh`
Установка Prometheus и Grafana через Helm:
- Создание namespace monitoring
- Установка kube-prometheus-stack
- Настройка ServiceMonitor

## Документация

### `README.md`
Основная документация проекта:
- Архитектура
- API endpoints
- Инструкции по развертыванию
- Мониторинг

### `QUICKSTART.md`
Быстрый старт для быстрого развертывания.

### `xtemp/`
Временная документация (не в git):
- `DEPLOYMENT_GUIDE.md` - детальное руководство по развертыванию
- `TESTING_GUIDE.md` - руководство по тестированию
- `GRAFANA_DASHBOARDS.md` - примеры дашбордов Grafana
- `PROJECT_REPORT_TEMPLATE.md` - шаблон отчета

## Статистика кода

- **Go код**: ~300 строк
- **YAML манифесты**: ~300 строк
- **Скрипты**: ~200 строк
- **Документация**: ~1000 строк

## Зависимости

### Go модули:
- `github.com/go-redis/redis/v8` - Redis клиент
- `github.com/prometheus/client_golang` - Prometheus метрики
- `github.com/gorilla/mux` - HTTP роутер

### Docker образы:
- `golang:1.22-alpine` - сборка
- `alpine:latest` - runtime
- `redis:7-alpine` - Redis

### Kubernetes компоненты:
- Prometheus Operator (Helm)
- NGINX Ingress Controller
- Metrics Server

## Тестирование

### Функциональное:
- Health check
- Отправка метрик
- Получение аналитики

### Нагрузочное:
- Locust: 1000+ RPS
- Apache Bench: 50000 запросов
- Авто-масштабирование HPA

### Мониторинг:
- Prometheus метрики
- Grafana дашборды
- Алерты


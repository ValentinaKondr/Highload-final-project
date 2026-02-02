# High-Load Service with AI-Optimization

Высоконагруженный сервис на Go с аналитикой и развертыванием в Kubernetes.

## Архитектура

Сервис обрабатывает потоковые метрики от IoT-устройств, выполняет аналитику в реальном времени и детектирует аномалии.

### Компоненты

- **Go Service**: HTTP-сервис на Go для приема метрик
- **Redis**: Кэширование метрик
- **Analytics**: Rolling average и z-score детекция аномалий
- **Prometheus**: Мониторинг метрик
- **Kubernetes**: Оркестрация с HPA для авто-масштабирования

## Требования

- Go 1.22+
- Docker
- Kubernetes (Minikube/Kind) или доступ к облачному кластеру
- kubectl

## Быстрый старт

### 1. Локальная разработка

```bash
# Установка зависимостей
go mod download

# Запуск Redis локально
docker run -d -p 6379:6379 redis:7-alpine

# Запуск сервиса
export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=""
go run cmd/service/main.go
```

### 2. Развертывание в Kubernetes

#### Подготовка Minikube

```bash
# Запуск Minikube
minikube start --cpus=2 --memory=4g

# Включение метрик сервера для HPA
minikube addons enable metrics-server
```

#### Сборка и загрузка образа

```bash
# Сборка Docker образа
./scripts/build.sh

# Загрузка в Minikube
minikube image load go-service:latest
```

#### Развертывание

```bash
# Развертывание всех компонентов
./scripts/deploy.sh

# Проверка статуса
kubectl get pods,svc,hpa
```

### 3. Доступ к сервису

```bash
# Port-forward для доступа
kubectl port-forward service/go-service 8080:80

# Тестирование
curl http://localhost:8080/health
curl http://localhost:8080/analyze

# Отправка метрики
curl -X POST http://localhost:8080/metrics \
  -H "Content-Type: application/json" \
  -d '{"timestamp": 1234567890, "cpu": 45.5, "rps": 1500}'
```

## API Endpoints

- `POST /metrics` - Прием метрик (JSON: timestamp, cpu, rps)
- `GET /analyze` - Получение аналитики (rolling average, anomaly stats)
- `GET /health` - Health check
- `GET /metrics` - Prometheus метрики

## Мониторинг

### Prometheus

```bash
# Установка Prometheus через Helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack

# Доступ к Prometheus UI
kubectl port-forward svc/prometheus-kube-prometheus-prometheus 9090:9090
```

### Grafana

```bash
# Доступ к Grafana (после установки через Helm выше)
kubectl port-forward svc/prometheus-grafana 3000:80
# Логин: admin, пароль: prom-operator
```

## Нагрузочное тестирование

```bash
# Установка Locust
pip3 install locust

# Запуск теста (1000 RPS, 5 минут)
./scripts/load-test.sh http://localhost:8080 1000 300

# Или через kubectl port-forward
kubectl port-forward service/go-service 8080:80 &
./scripts/load-test.sh http://localhost:8080 1000 300
```

## Структура проекта

```
.
├── cmd/
│   └── service/          # Основной сервис
├── internal/
│   ├── analytics/        # Rolling average и anomaly detection
│   ├── cache/           # Redis клиент
│   └── metrics/         # Prometheus метрики
├── k8s/                 # Kubernetes манифесты
│   ├── deployments/
│   ├── services/
│   ├── configmaps/
│   ├── hpa/
│   └── ingress/
├── scripts/             # Скрипты развертывания
└── xtemp/              # Временная документация (gitignored)
```

## Аналитика

### Rolling Average
- Окно: 50 событий
- Используется для сглаживания нагрузки и прогнозирования

### Anomaly Detection
- Метод: Z-score
- Threshold: 2σ (стандартных отклонения)
- Окно: 50 событий

## Масштабирование

HPA настроен на:
- Минимум реплик: 2
- Максимум реплик: 5
- Целевая загрузка CPU: 70%

## Производительность

- Целевая нагрузка: 1000+ RPS
- Latency: < 50ms (p95)
- Точность детекции аномалий: > 70%

## Лицензия

MIT


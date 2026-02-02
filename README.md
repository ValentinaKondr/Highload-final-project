## Высоконагруженный Go-сервис с аналитикой и авто-масштабированием

Проект представляет собой высоконагруженный HTTP-сервис на Go для обработки
потоковых метрик (например, RPS и CPU), с простой статистической аналитикой
и детекцией аномалий.

Сервис развёртывается в Kubernetes и поддерживает горизонтальное
масштабирование (HPA) и мониторинг с помощью Prometheus и Grafana.

---

## Используемый стек

- **Go 1.22+**
- **Redis** — кэширование метрик
- **Docker**
- **Kubernetes (Minikube / Kind)**
- **Prometheus + Grafana**
- **Apache Bench и пользовательские скрипты нагрузки**

---

## Архитектура решения

- HTTP API для приёма и анализа метрик
- Аналитика нагрузки:
  - Rolling Average (скользящее среднее, окно 50 значений)
  - Детекция аномалий по z-score (порог 2.0)
- Redis используется как лёгкое хранилище
- Авто-масштабирование подов на основе загрузки CPU (HPA)
- Экспорт метрик для Prometheus

---

## HTTP API

| Endpoint | Метод | Описание |
|--------|-------|----------|
| `/health` | GET | Проверка работоспособности |
| `/metrics` | POST | Приём метрик (JSON) |
| `/analyze` | GET | Текущая аналитика и состояние детектора |
| `/metrics` | GET | Метрики Prometheus |

### Пример запроса `/metrics`

```json
{
  "timestamp": 1234567890,
  "cpu": 42.5,
  "rps": 1200
}

## Локальный запуск
### Требования
- Go 1.22+
- Redis (локально или через Docker)


```bash
go mod download

# Запуск Redis локально
docker run -d -p 6379:6379 redis:7-alpine

# Запуск сервиса
export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=""
go run cmd/service/main.go
```
---
```bash
# проверка
curl http://localhost:8081/health
curl http://localhost:8081/analyze
```
## Тестирование
В проекте реализованы юнит- и интеграционные тесты.
### Запуск тестов
```bash
go test ./... -v
```
### Запуск с детектором гонок
```bash
go test ./... -v
```
Тестами покрыты:
- скользящее среднее
- детекция аномалий (z-score)
- HTTP-эндпоинты (с in-memory кэшем без Redis)

## Нагрузочное тестирование
### Скрипт нагрузки с аномалиями
```bash
./scripts/load_anomaly_test.sh
```
Скрипт выполняет прогрев окна аналитики, генерацию аномальных значений, запрос /analyze для фиксации результата


```bash
# Установка Locust
pip3 install locust

# Запуск теста (1000 RPS, 5 минут)
./scripts/load-test.sh http://localhost:8080 1000 300

# Или через kubectl port-forward
kubectl port-forward service/go-service 8080:80 &
./scripts/load-test.sh http://localhost:8080 1000 300


### Apache Bench 
```bash 
ab -n 100000 -c 200 \
  -p payload.json \
  -T application/json \
  http://localhost:8081/metrics
```

## Развертывание в Kubernetes

### Подготовка Minikube
```bash
minikube start --cpus=2 --memory=4g

minikube addons enable metrics-server
```

#### Сборка и загрузка образа

```bash
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

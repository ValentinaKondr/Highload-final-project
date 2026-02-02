## Высоконагруженный Go-сервис с аналитикой и авто-масштабированием

Проект представляет собой высоконагруженный сервис анализа потоковых метрик,
разработанный на языке Go и развернутый в Kubernetes.
Сервис предназначен для обработки HTTP-запросов с метриками (CPU, RPS),
выполнения базовой статистической аналитики в реальном времени и автоматического
масштабирования в зависимости от нагрузки.

В качестве аналитики используются простые, но эффективные статистические методы:
rolling average для сглаживания нагрузки и z-score для детекции аномалий.
Решение ориентировано на работу в условиях высокой нагрузки и демонстрирует
применение облачно-ориентированных подходов (containerization, HPA, monitoring).

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

Архитектура сервиса включает следующие компоненты:

- Go-сервис, принимающий и обрабатывающий метрики по HTTP API
- Redis для кэширования входных данных
- Kubernetes для оркестрации и масштабирования
- Horizontal Pod Autoscaler (HPA) для автоматического масштабирования по CPU
- Prometheus для сбора метрик
- Grafana для визуализации и мониторинга

Поток данных:
1. Клиент отправляет метрики на эндпоинт POST /metrics
2. Данные кэшируются в Redis
3. Выполняется обновление rolling average и детекция аномалий
4. Метрики экспортируются в Prometheus
5. Grafana визуализирует состояние системы в реальном времени

---
## Конфигурация приложения

Конфигурация сервиса загружается один раз при старте из переменных окружения
(через ConfigMap в Kubernetes) и далее используется внутри приложения без повторного
чтения из окружения.

Основные параметры конфигурации:

- WINDOW_SIZE — размер окна для rolling average и z-score (по умолчанию 50)
- ANOMALY_THRESHOLD — порог детекции аномалий в сигмах (по умолчанию 2.0)
- REDIS_ADDR — адрес Redis
- REDIS_PASSWORD — пароль Redis (через Secret)

---
## HTTP API

| Endpoint | Метод | Описание |
|--------|-------|----------|
| `/health` | GET | Проверка работоспособности |
| `/metrics` | POST | Приём метрик (JSON) |
| `/analyze` | GET | Текущая аналитика и состояние детектора |
| `/metrics` | GET | Метрики Prometheus |

### Пример запроса `/health`
`/health`
пример ответа:
```json
{
  "status": "ok",
  "version": "v1.0.0"
}
```

### Пример запроса `/metrics`

```json
{
  "timestamp": 1700000000,
  "cpu": 42.5,
  "rps": 120
}
```

### Пример запроса `/analyze`

```json
{
  "rolling_average": 118,
  "anomaly_stats": {
    "mean": 115,
    "std_dev": 5,
    "threshold": 2,
    "window_size": 50,
    "data_points": 50
  },
  "timestamp": 1700000123
}
```

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
kubectl port-forward service/metrics-analyzer 8080:80 &
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
minikube image load metrics-analyzer:latest
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
kubectl port-forward service/metrics-analyzer 8080:80

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
# Доступ к Grafana (после установки через Helm)
kubectl port-forward svc/prometheus-grafana 3000:80
# Логин: admin, пароль: pwd1234
```
